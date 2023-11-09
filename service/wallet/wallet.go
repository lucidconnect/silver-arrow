package wallet

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lucidconnect/silver-arrow/erc20"

	// "github.com/lucidconnect/silver-arrow/erc4337"
	"github.com/lucidconnect/silver-arrow/api/graphql/wallet/graph/model"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/service/erc4337"
	"github.com/lucidconnect/silver-arrow/service/merchant"
	"github.com/lucidconnect/silver-arrow/service/turnkey"
	"github.com/pkg/errors"
	"github.com/rmanzoku/ethutils/ecrecover"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type WalletService struct {
	turnkey          *turnkey.TurnkeyService
	database         repository.Database
	bundlerService   *erc4337.AlchemyService
	validatorAddress string
}

func NewWalletService(r repository.Database, t *turnkey.TurnkeyService, chain int64) *WalletService {
	validatorAddress := os.Getenv("VALIDATOR_ADDRESS")
	var bundler *erc4337.AlchemyService
	var err error

	if chain != 0 {
		bundler, err = initialiseBundler(chain)
		if err != nil {
			return nil
		}
	}

	return &WalletService{
		turnkey:          t,
		database:         r,
		bundlerService:   bundler,
		validatorAddress: validatorAddress,
	}
}

func initialiseBundler(chain int64) (*erc4337.AlchemyService, error) {
	bundler, err := erc4337.NewAlchemyService(chain)
	if err != nil {
		err = errors.Wrap(err, "initialising bundler failed")
		log.Err(err).Caller().Send()
	}
	return bundler, nil
}

func (ws *WalletService) AddAccount(input model.Account) error {
	walletAddress := input.Address
	// Check if account exists
	_, err := ws.database.FetchAccountByAddress(walletAddress)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// create turnkey sub organization
			activityId, err := ws.turnkey.CreateSubOrganization("", walletAddress)
			if err != nil {
				log.Err(err).Caller().Send()
				return err
			}

			result, err := ws.turnkey.GetActivity("", activityId)
			fmt.Println("result: ", result)
			if err != nil {
				log.Err(err).Caller().Send()
				return err
			}

			orgId := turnkey.ExtractSubOrganizationIdFromResult(result)
			tag := fmt.Sprintf("key-tag-%s", input.Address)
			tagActivity, err := ws.turnkey.CreatePrivateKeyTag(orgId, tag)
			if err != nil {
				log.Err(err).Caller().Send()
				return err
			}

			tagResult, err := ws.turnkey.GetActivity(orgId, tagActivity)
			if err != nil {
				log.Err(err).Caller().Send()
				return err
			}

			tagId := turnkey.ExtractPrivateKeyTagIdFromResult(tagResult)
			wallet := &models.Wallet{
				WalletAddress:        walletAddress,
				SignerAddress:        *input.Signer,
				TurnkeySubOrgID:      orgId,
				TurnkeySubOrgName:    walletAddress,
				TurnkeyPrivateKeyTag: tagId,
			}

			err = ws.database.AddAccount(wallet)
			if err != nil {
				log.Err(err).Caller().Send()
				return err
			}
		} else {
			log.Err(err).Caller().Send()
			return err
		}
	}

	return nil
}

func (ws *WalletService) ValidateSubscription(userop map[string]any, chain int64) (*model.TransactionData, error) {
	opHash, err :=  ws.bundlerService.SendUserOperation(userop)
	if err != nil {
		log.Err(err).Msg("failed to send user op")
		return nil, err
	}
	fmt.Println("validating subscription with userop hash -", opHash)
	result, err := ws.database.FindSubscriptionByHash(opHash)
	if err != nil {
		log.Err(err).Msgf("failed to find subscription with hash %v", opHash)
		return nil, err
	}

	productId, err := merchant.Base64EncodeUUID(result.ProductID)
	if err != nil {
		log.Err(err).Msg("encoding product id failed")
		return nil, err
	}
	token := result.Token
	createdAt := result.CreatedAt.Format(time.RFC3339)
	amount := int(result.Amount)
	interval := int(result.Interval)
	subData := &model.TransactionData{
		Token:         token,
		Amount:        amount,
		Interval:      &interval,
		ProductID:     &productId,
		WalletAddress: result.WalletAddress,
		CreatedAt:     &createdAt,
	}
	fmt.Println("subscription result - ", result)

	transactionHash, err := ws.getTransactionHash(opHash)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	explorerUrl, err := erc20.GetChainExplorer(chain)
	if err != nil {
		log.Err(err).Msg("failed to get chain explorer url")
	}
	blockExplorerTx := fmt.Sprintf("%v/tx/%v", explorerUrl, transactionHash)

	update := map[string]interface{}{"active": true, "updated_at": time.Now(), "transaction_hash": transactionHash}
	err = ws.database.UpdateSubscription(result.ID, update)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}

	// if payment is due now, create a payment
	if isPaymentDue(result.NextChargeAt) {
		// create payment
		var sponsored bool
		switch os.Getenv("USE_PAYMASTER") {
		case "TRUE":
			sponsored = true
		default:
			sponsored = false
		}
		reference := uuid.New()
		payment := &models.Payment{
			Type:                  PaymentTypeRecurring.String(),
			Chain:                 result.Chain,
			Token:                 result.Token,
			Amount:                result.Amount,
			Source:                result.WalletAddress,
			WalletID:              result.WalletID,
			ProductID:             result.ProductID,
			Sponsored:             sponsored,
			Reference:             reference,
			Destination:           result.MerchantDepositAddress,
			SubscriptionID:        result.ID,
			SubscriptionPublicKey: result.Key.PublicKey,
		}

		userop, useropHash, err := ws.CreatePayment(payment)
		if err != nil {
			err = errors.Wrap(err, "creating payment operation failed")
			log.Err(err).Caller().Send()
			return nil, err
		}

		signature, err := ws.SignPaymentOperation(userop, useropHash)
		if err != nil {
			err = errors.Wrap(err, "signing payment operation failed")
			log.Err(err).Caller().Send()
			return nil, err
		}
		userop["signature"] = signature

		onchainTx, err := ws.ExecutePaymentOperation(userop, payment.Chain)
		if err != nil {
			log.Err(err).Send()
			return subData, err
		}
		nextChargeAt := time.Now().Add((time.Duration(result.Interval)))

		update := map[string]interface{}{
			"expires_at":     nextChargeAt,
			"next_charge_at": nextChargeAt,
		}
		err = ws.database.UpdateSubscription(result.ID, update)
		if err != nil {
			log.Err(err).Send()
		}
		subData.TransactionExplorer = &onchainTx
	} else {
		subData.TransactionExplorer = &blockExplorerTx
	}

	return subData, nil
}

func (ws *WalletService) AddSubscription(merchantId uuid.UUID, input model.NewSubscription, usePaymaster bool, index *big.Int, chain int64) (*model.ValidationData, map[string]any, error) {
	var nextChargeAt time.Time
	var initCode []byte
	var nonce, amount *big.Int

	// bundler := ws.bundlerService

	tagId, orgId, walletID, err := ws.database.GetWalletMetadata(input.WalletAddress)
	if err != nil {
		log.Err(err).Msgf("failed to fetch private key tag for wallet - %v", input.WalletAddress)
		return nil, nil, err
	}

	productId := merchant.ParseUUID(input.ProductID)

	product, err := ws.database.FetchProduct(productId)
	if err != nil {
		log.Err(err).Msg("failed to fetch product")
	}
	randomSalt := randKey(4)
	keyName := fmt.Sprintf("sub-%v-%v", randomSalt, productId)
	activityId, err := ws.turnkey.CreatePrivateKey(orgId, keyName, tagId)
	if err != nil {
		log.Err(err).Msg("failed to create subscription private key")
		return nil, nil, err
	}

	result, err := ws.turnkey.GetActivity(orgId, activityId)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, nil, err
	}
	privateKeyID, sessionKey, err := turnkey.GetPrivateKeyIdFromResult(result)
	// sessionKey, signingKey, err := CreateAccessKey()
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, nil, err
	}

	// bundler, err := erc4337.NewAlchemyService(chain)
	// if err != nil {
	// 	log.Err(err).Msg("failed to initialise bundler")
	// 	return nil, nil, err
	// }

	// supported token is still USDC, so minor factor is 1000000
	amount = big.NewInt(int64(input.Amount)) // This will cause a bug for amounts that are fractional
	interval := daysToNanoSeconds(int64(input.Interval))

	if input.NextChargeDate != nil {
		nextChargeAt = *input.NextChargeDate
	} else {
		nextChargeAt = time.Now().Add(interval)
	}

	isAccountDeployed := ws.isAccountDeployed(input.WalletAddress, chain)
	if !isAccountDeployed {
		initCode, err = GetContractInitCode(common.HexToAddress(input.OwnerAddress), index)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, nil, err
		}
		nonce = common.Big0
	} else {
		nonce, err = ws.bundlerService.GetAccountNonce(common.HexToAddress(input.WalletAddress))
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, nil, err
		}
	}

	callData, err := setValidatorExecutor(sessionKey, ws.validatorAddress, input.WalletAddress, int64(input.Chain))
	if err != nil {
		log.Err(err).Msg("failed to set a validator")
		return nil, nil, err
	}

	op, err := ws.bundlerService.CreateUnsignedUserOperation(input.WalletAddress, initCode, callData, nonce, usePaymaster, int64(input.Chain))
	if err != nil {
		log.Err(err).Msg("failed to create user operation")
		return nil, nil, err
	}

	entrypoint := erc4337.GetEntryPointAddress()
	operation, err := userop.New(op)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, nil, err
	}
	opHash := operation.GetUserOpHash(entrypoint, big.NewInt(int64(input.Chain)))

	key := &models.Key{
		WalletID:     walletID,
		PublicKey:    sessionKey,
		PrivateKeyId: privateKeyID,
	}

	tokenAddress := erc20.GetTokenAddress(input.Token, chain)
	sub := &models.Subscription{
		Token:                  input.Token,
		Amount:                 amount.Int64(),
		Active:                 false,
		Interval:               interval.Nanoseconds(),
		UserOpHash:             opHash.Hex(),
		MerchantId:             merchantId.String(),
		ProductID:              productId,
		MerchantDepositAddress: product.DepositAddress,
		NextChargeAt:           nextChargeAt,
		ExpiresAt:              nextChargeAt,
		WalletID:               walletID,
		WalletAddress:          input.WalletAddress,
		Chain:                  chain,
		Key:                    *key,
		TokenAddress:           tokenAddress,
	}

	err = ws.database.AddSubscription(sub, key)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, nil, err
	}

	// fmt.Println("New Subscription added", opHash.Hex())
	return &model.ValidationData{
		UserOpHash: opHash.Hex(),
	}, op, nil
}

func (w *WalletService) FetchSubscriptions(walletAddress string) ([]*model.SubscriptionData, error) {
	var subData []*model.SubscriptionData

	subs, err := w.database.FetchWalletSubscriptions(walletAddress)
	if err != nil {
		log.Err(err).Msgf("error while fetching subscriotions for %v", walletAddress)
		return nil, err
	}

	for _, v := range subs {
		interval := nanoSecondsToDay(v.Interval)
		createdAt := v.CreatedAt.Format("dd:mm:yyyy")
		sd := &model.SubscriptionData{
			ID:        v.ID.String(),
			Token:     v.Token,
			Amount:    int(v.Amount),
			Interval:  int(interval),
			ProductID: v.ProductID.String(),
			CreatedAt: &createdAt,
		}
		subData = append(subData, sd)
	}

	return subData, nil
}

func amountToWei(amount any) (*big.Int, error) {
	etherInWei := new(big.Int)
	etherInWei.SetString("1000000000000000000", 10)

	switch v := amount.(type) {
	case *big.Int:
		weiAmount := new(big.Int).Mul(v, etherInWei)
		return weiAmount, nil
	case *big.Float:
		weiAmount := new(big.Int)
		weiAmountFloat := new(big.Float).Mul(v, big.NewFloat(1e18))
		weiAmountFloat.Int(weiAmount)
		return weiAmount, nil
	default:
		return nil, fmt.Errorf("unsupported input type: %T", amount)
	}
}

func amountToMwei(amount int64) (*big.Int, error) {
	etherInMWei := new(big.Int)
	etherInMWei.SetString("1000000", 10)
	v := big.NewInt(int64(amount))
	mWeiAmount := new(big.Int).Mul(v, etherInMWei)
	return mWeiAmount, nil
	// switch v := amount.(type) {
	// case *big.Int:
	// 	mWeiAmount := new(big.Int).Mul(v, etherInMWei)
	// 	return mWeiAmount, nil
	// case *big.Float:
	// 	mWeiAmount := new(big.Int)
	// 	mWeiAmountFloat := new(big.Float).Mul(v, big.NewFloat(1e6))
	// 	mWeiAmountFloat.Int(mWeiAmount)
	// 	return mWeiAmount, nil
	// default:
	// 	return nil, fmt.Errorf("unsupported input type: %T", amount)
	// }
}

func mWeiToAmount(amt *big.Int) int64 {
	etherInMWei := new(big.Int)
	etherInMWei.SetString("1000000", 10)

	result := new(big.Int)
	result.Div(amt, etherInMWei)
	return result.Int64()
}

func weiToAmount(amt *big.Int) int64 {
	etherInWei := new(big.Int)
	etherInWei.SetString("1000000000000000000", 10)

	result := new(big.Int)
	result.Div(amt, etherInWei)
	return result.Int64()
}

// CreatePayment creates a userop for an initiated payment,
// generates the userop hash, sets the payment status to a pending state
// returns a message to be signed.
func (ws *WalletService) CreatePayment(payment *models.Payment) (map[string]any, common.Hash, error) {

	// payment, _ := ws.database.FindPaymentByReference(reference)

	// bundler, err := erc4337.NewAlchemyService(payment.Chain)
	// if err != nil {
	// 	err = errors.Wrap(err, "initialising alchemy service failed")
	// 	log.Err(err).Send()
	// 	return nil, common.Hash{}, err
	// }
	erc20Token := erc20.GetTokenAddress(payment.Token, payment.Chain)
	tokenAddress := common.HexToAddress(erc20Token)

	// wallet, err := ws.database.FetchAccountByAddress(payment.Source)
	// if err != nil {
	// 	err = errors.Wrapf(err, "smart account lookup for address [%v] failed", payment.Source)
	// 	log.Err(err).Caller().Send()
	// 	return "", fmt.Errorf("sender address [%v] not found", payment.Source)
	// }

	actualAmount, err := amountToMwei(payment.Amount)
	if err != nil {
		err = errors.Wrapf(err, "converting amount [%v] to wei value failed", payment.Amount)
		log.Err(err).Caller().Send()
		return nil, common.Hash{}, fmt.Errorf("internal server error")
	}
	data, err := erc4337.TransferErc20Action(tokenAddress, common.HexToAddress(payment.Destination), actualAmount)
	if err != nil {
		err = errors.Wrap(err, "creating TransferErc20Action call data failed")
		log.Err(err).Caller().Send()
		return nil, common.Hash{}, fmt.Errorf("internal server error")
	}

	nonce, err := ws.bundlerService.GetAccountNonce(common.HexToAddress(payment.Source))
	if err != nil {
		err = errors.Wrapf(err, "error occured fetching nonce for account [%v]", payment.Source)
		log.Err(err).Caller().Send()
		return nil, common.Hash{}, fmt.Errorf("internal server error")
	}

	op, err := ws.bundlerService.CreateUnsignedUserOperation(payment.Source, nil, data, nonce, payment.Sponsored, payment.Chain)
	if err != nil {
		err = errors.Wrapf(err, "error occured creating user operation for account [%v]", payment.Source)
		log.Err(err).Caller().Send()
		return nil, common.Hash{}, fmt.Errorf("internal server error")
	}

	operation, err := userop.New(op)
	if err != nil {
		err = errors.Wrapf(err, "error occured creating user operation for account [%v]", payment.Source)
		log.Err(err).Caller().Send()
		return nil, common.Hash{}, fmt.Errorf("internal server error")
	}

	entrypoint := erc4337.GetEntryPointAddress()

	chainId := big.NewInt(payment.Chain)
	userOpHash := operation.GetUserOpHash(entrypoint, chainId)
	// hash := userOpHash.Bytes()

	payment.Status = string(PaymentStatusPending)
	payment.UserOpHash = userOpHash.Hex()
	err = ws.database.CreatePayment(payment)
	if err != nil {
		log.Err(err).Send()
	}

	return op, userOpHash, nil
}

// SignPaymentOperation takes in the userop alongside it's computed userop hash and returns a signature.
// It signs the hash with the key created for the subscription that is initiating the payment
// note this method is intended only for recurring automated payments.
func (ws *WalletService) SignPaymentOperation(op map[string]any, hash common.Hash) (string, error) {
	payment, err := ws.database.FindPaymentByUseropHash(hash.Hex())
	if err != nil {
		log.Err(err).Msgf("failed to fetch payment with user op hash %v", hash)
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("payment not found")
		}
		return "", err
	}
	// find the subscription the payment is for
	subscriptionKeyId, err := ws.database.GetSubscriptionKey(payment.SubscriptionPublicKey)
	if err != nil {
		err = errors.Wrap(err, "invalid subscription key")
		log.Err(err).Send()
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("this payment was not authorised, key not found")
		}
		return "", err
	}
	message := hexutil.Encode(ecrecover.ToEthSignedMessageHash(hash.Bytes()))

	wallet, err := ws.database.FetchAccountByAddress(payment.Source)
	if err != nil {
		err = errors.Wrapf(err, "smart account lookup for address [%v] failed", payment.Source)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("sender address [%v] not found", payment.Source)
	}
	org := wallet.TurnkeySubOrgID

	fmt.Println("Signing user op with key id - ", subscriptionKeyId)
	turnkeyActivityId, err := ws.turnkey.SignMessage(org, subscriptionKeyId, message)
	if err != nil {
		err = errors.Wrapf(err, "turnkey failed to sign user operation for account [%v], keyId: [%v]", payment.Source, subscriptionKeyId)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	result, err := ws.turnkey.GetActivity(org, turnkeyActivityId)
	if err != nil {
		err = errors.Wrapf(err, "fetching turnkey activity failed, activityId: [%v]", turnkeyActivityId)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	sig, err := turnkey.ExctractTurnkeySignatureFromResult(result)
	if err != nil {
		err = errors.Wrap(err, "failed to extract signature from turnkey result")
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}
	// op["signature"] = sig.ParseSignature(erc4337.VALIDATOR_MODE)
	signature := sig.ParseSignature(erc4337.VALIDATOR_MODE)
	return signature, nil
}

// executePaymentOperation sends the userop for the payment to the bundler,
// waits to get the transaction hash, updates the payment status
// and returns the transaction hash.
func (ws *WalletService) ExecutePaymentOperation(signedOp map[string]any, chain int64) (string, error) {
	opHash, err := ws.bundlerService.SendUserOperation(signedOp)
	if err != nil {
		err = errors.Wrap(err, "sending user op failed")
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	// TODO: use the userop hash to create a reciept for the transsaction
	fmt.Println("user operation hash -", opHash) // 0x28b45cf378c23fbdbbcb4f4c4d085791eb6d660214ff4a2402e40fd1c73751c6 0xfcd3b481cc3ba345fcf24c777463baf60dbb1f7475ca297b9259d020044565be

	// Fetch transaction hash
	transactionHash, err := ws.getTransactionHash(opHash)
	if err != nil {
		err = errors.Wrapf(err, "fetching the transction hash failed. userop hash - [%v]", opHash)
		log.Err(err).Caller().Send()
	}

	payment, err := ws.database.FindPaymentByUseropHash(opHash)
	if err != nil {
		err = errors.Wrapf(err, "couldn't find payment with user_op_hash [%v] - weird. Transaction hash on chain [%v] - [%v]", opHash, chain, transactionHash)
		log.Err(err).Send()
		return transactionHash, err
	}

	explorerUrl, err := erc20.GetChainExplorer(chain)
	if err != nil {
		log.Err(err).Msg("failed to get chain explorer url")
	}
	blockExplorerTx := fmt.Sprintf("%v/tx/%v", explorerUrl, transactionHash)
	update := map[string]any{
		"status":            PaymentStatusSuccess,
		"transaction_hash":  transactionHash,
		"block_explorer_tx": blockExplorerTx,
	}
	err = ws.database.UpdatePayment(payment.ID, update)
	if err != nil {
		err = errors.Wrapf(err, "updating payment status failed but transaction was successful on chain [%v]: useropHash - [%v]; transactionHash - [%v]", chain, opHash, transactionHash)
		log.Err(err).Caller().Send()
		return transactionHash, err
	}

	// should probably trigger a webhook event

	return blockExplorerTx, nil
}

// TODO: delete
// Execute a charge on an AA wallet, currently limited to USDC
func (ws *WalletService) ExecuteCharge(sender, target, token, key string, amount, chain int64, sponsored bool) (string, error) {
	// bundler, err := erc4337.NewAlchemyService(chain)
	// if err != nil {
	// 	err = errors.Wrap(err, "initialising alchemy service failed")
	// 	log.Err(err).Send()
	// 	return "", err
	// }
	erc20Token := erc20.GetTokenAddress(token, chain)
	tokenAddress := common.HexToAddress(erc20Token)

	wallet, err := ws.database.FetchAccountByAddress(sender)
	if err != nil {
		err = errors.Wrapf(err, "smart account lookup for address [%v] failed", sender)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("sender address [%v] not found", sender)
	}
	org := wallet.TurnkeySubOrgID

	actualAmount, err := amountToMwei(amount)
	if err != nil {
		err = errors.Wrapf(err, "converting amount [%v] to wei value failed", amount)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}
	data, err := erc4337.TransferErc20Action(tokenAddress, common.HexToAddress(target), actualAmount)
	if err != nil {
		err = errors.Wrap(err, "creating TransferErc20Action call data failed")
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	nonce, err := ws.bundlerService.GetAccountNonce(common.HexToAddress(sender))
	if err != nil {
		err = errors.Wrapf(err, "error occured fetching nonce for account [%v]", sender)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	op, err := ws.bundlerService.CreateUnsignedUserOperation(sender, nil, data, nonce, sponsored, chain)
	if err != nil {
		err = errors.Wrapf(err, "error occured creating user operation for account [%v]", sender)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	operation, err := userop.New(op)
	if err != nil {
		err = errors.Wrapf(err, "error occured creating user operation for account [%v]", sender)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	entrypoint := erc4337.GetEntryPointAddress()

	chainId := big.NewInt(chain)
	userOpHash := operation.GetUserOpHash(entrypoint, chainId)
	hash := userOpHash.Bytes()

	message := hexutil.Encode(ecrecover.ToEthSignedMessageHash(hash))

	fmt.Println("Signing user op with key - ", key)
	turnkeyActivityId, err := ws.turnkey.SignMessage(org, key, message)
	if err != nil {
		err = errors.Wrapf(err, "turnkey failed to sign user operation for account [%v], keyId: [%v]", sender, key)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	result, err := ws.turnkey.GetActivity(org, turnkeyActivityId)
	if err != nil {
		err = errors.Wrapf(err, "fetching turnkey activity failed, activityId: [%v]", turnkeyActivityId)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	sig, err := turnkey.ExctractTurnkeySignatureFromResult(result)
	if err != nil {
		err = errors.Wrap(err, "failed to extract signature from turnkey result")
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}
	op["signature"] = sig.ParseSignature(erc4337.VALIDATOR_MODE)

	opHash, err := ws.bundlerService.SendUserOperation(op)
	if err != nil {
		err = errors.Wrap(err, "sending user op failed")
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	// TODO: use the userop hash to create a reciept for the transsaction
	fmt.Println("user operation hash -", opHash) // 0x28b45cf378c23fbdbbcb4f4c4d085791eb6d660214ff4a2402e40fd1c73751c6 0xfcd3b481cc3ba345fcf24c777463baf60dbb1f7475ca297b9259d020044565be

	// Fetch transaction hash
	transactionHash, err := ws.getTransactionHash(opHash)
	if err != nil {
		err = errors.Wrap(err, "fetching the transction hash failed")
		log.Err(err).Caller().Send()
	}

	return transactionHash, err
}

// Transfer tokens from a smart wallet,
// authorised by the user's EOA. This is necessary to provide users an interface to move assets
// while this method also works for erc20 tokens, note that USDc is the primary supported token
// hence, using this method to transfer othet tokens with higher decimals will result in unexpected behavior
func (ws *WalletService) InitiateTransfer(sender, target, token string, amount float64, chain int64, sponsored bool) (*model.ValidationData, map[string]any, error) {
	var callData []byte

	bundler, err := erc4337.NewAlchemyService(chain)
	if err != nil {
		log.Err(err).Msg("failed to initialise bundler")
		return nil, nil, err
	}

	transferAmount := parseTransferAmount(token, chain, amount)
	callData, err = erc4337.CreateTransferCallData(target, token, chain, transferAmount)
	if err != nil {
		err = errors.Wrapf(err, "creating transfer call data failed")
		return nil, nil, err
	}

	nonce, err := bundler.GetAccountNonce(common.HexToAddress(sender))
	if err != nil {
		err = errors.Wrapf(err, "fetching account nonce failed for acoount - %v", sender)
		return nil, nil, err
	}
	op, err := bundler.CreateUnsignedUserOperation(sender, nil, callData, nonce, sponsored, chain)
	if err != nil {
		err = errors.Wrapf(err, "creating userop failed")
		return nil, nil, err
	}
	// fmt.Println(userop)

	entrypoint := erc4337.GetEntryPointAddress()
	operation, err := userop.New(op)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, nil, err
	}
	opHash := operation.GetUserOpHash(entrypoint, big.NewInt(chain))

	vd := &model.ValidationData{
		UserOpHash: opHash.Hex(),
	}
	// userophash has to be returned for the user to sign
	return vd, op, nil
}
func (ws *WalletService) ValidateTransfer(userop map[string]any, chain int64) (*model.TransactionData, error) {
	bundler, err := erc4337.NewAlchemyService(chain)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}

	opHash, err := bundler.SendUserOperation(userop)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}

	// Fetch transaction hash
	useropResult, err := bundler.GetUserOperationByHash(opHash)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}

	transactionHash := useropResult["transactionHash"].(string)

	explorer, err := erc20.GetChainExplorer(chain)
	if err != nil {
		log.Err(err).Send()
	}
	blockExplorerTx := fmt.Sprintf("%v/tx/%v", explorer, transactionHash)

	transactionDetails := &model.TransactionData{
		Chain:               int(chain),
		TransactionHash:     &transactionHash,
		TransactionExplorer: &blockExplorerTx,
	}

	return transactionDetails, nil
}



func daysToNanoSeconds(days int64) time.Duration {
	nanoSsecondsInt := days * 24 * 60 * 60 * 1e9
	return time.Duration(nanoSsecondsInt)
}

func nanoSecondsToDay(ns int64) int64 {
	interval := time.Duration(ns)
	hours := interval.Hours()

	days := hours / 24
	return int64(days)
}

// creats the calldata that scopes a kernel executor to a validator
func setValidatorExecutor(sessionKey, validatorAddress, ownerAddress string, chain int64) ([]byte, error) {
	mode := erc4337.ENABLE_MODE
	validator, err := erc4337.InitialiseValidator(validatorAddress, sessionKey, mode, chain)
	if err != nil {
		return nil, err
	}

	enableData, err := validator.GetEnableData()
	if err != nil {
		return nil, err
	}

	callData, err := validator.SetExecution(enableData, ownerAddress)
	if err != nil {
		err = errors.Wrap(err, "validator.SetExecution():")
		return nil, err
	}
	return callData, nil
}

func GetContractInitCode(owner common.Address, index *big.Int) ([]byte, error) {
	initCode := []byte{}
	factoryAddress := os.Getenv("KERNEL_FACTORY_ADDRESS")
	implementation := os.Getenv("KERNEL_IMPLEMENTATION_ADDRESS")
	defaultValidator := os.Getenv("DEFAULT_VALIDATOR")

	kernelImplementation := common.HexToAddress(implementation)
	// fmt.Println("accountAddress ", accountAddress)

	/** inputs to createAddress
		- account implementation
		- calldata:abi.encodeWithSelector(
	 		KernelStorage.initialize.selector, defaultValidator, abi.encodePacked(owner)),
		- index
	*/

	callData, err := erc4337.EncodeKernelStorageWithSelector("initialize", common.HexToAddress(defaultValidator), owner.Bytes())

	fmt.Println("callData", hexutil.Encode(callData))
	if err != nil {
		return nil, err
	}

	data := owner.Bytes()
	fmt.Println("enable data ", hexutil.Encode(data))
	code, err := erc4337.GetCreateAccountFnData(kernelImplementation, callData, index)
	if err != nil {
		return nil, err
	}
	factoryAddressToBytes := common.FromHex(factoryAddress)
	initCode = append(initCode, factoryAddressToBytes...)
	initCode = append(initCode, code...)

	return initCode, nil
}

func (ws *WalletService) isAccountDeployed(address string, chain int64) bool {
	bundler, err := erc4337.NewAlchemyService(chain)
	if err != nil {
		err = errors.Wrap(err, "failed to initialise bundler")
		log.Panic().Err(err).Send()
		return false
	}

	code, err := bundler.GetAccountCode(common.HexToAddress(address))
	if err != nil {
		log.Err(err).Caller().Send()
		return false
	}
	fmt.Println("Code ", code)
	if len(code) == 0 {
		log.Info().Msg("account not deployed, should be deployed first!")
		return false
	}
	return true
}

// TODO: not finished
// CancelSubscription will remove the subscription key from the wallet
func (ws *WalletService) CancelSubscription(subscriptionId string) (string, error) {
	id, err := uuid.Parse(subscriptionId)
	if err != nil {
		err = errors.Wrapf(err, "parsing subscription id %v failed", subscriptionId)
		log.Err(err).Send()
		return "", err
	}
	_, err = ws.database.FindSubscriptionById(id)
	if err != nil {
		err = errors.Wrapf(err, "failed to fetch subscription %v", subscriptionId)
		log.Err(err).Send()
	}

	// ws.database.DeactivateSubscription()
	return "", errors.New("unimplemented")
}

// DisableSubscription only toggles a subscription status to inactive
func (ws *WalletService) DisableSubscription(subscriptionId string) (string, error) {
	id, err := uuid.Parse(subscriptionId)
	if err != nil {
		err = errors.Wrapf(err, "parsing subscription id %v failed", subscriptionId)
		log.Err(err).Send()
		return "", fmt.Errorf("could not disable subscription with id %v", subscriptionId)
	}

	update := map[string]any{
		"active":     false,
		"updated_at": time.Now(),
	}

	err = ws.database.UpdateSubscription(id, update)
	if err != nil {
		err = errors.Wrapf(err, "modifying subscription status failed for sub %v ", subscriptionId)
		log.Err(err).Send()
		return "", fmt.Errorf("could not disable subscription with id %v", subscriptionId)
	}

	return subscriptionId, nil
}

// Enable subscriptions reactivates a disabled subscription, won't work for cancelled subscriptions
func (ws *WalletService) EnableSubscription(subscriptionId string) (string, error) {
	id, err := uuid.Parse(subscriptionId)
	if err != nil {
		err = errors.Wrapf(err, "parsing subscription id %v failed", subscriptionId)
		log.Err(err).Send()
		return "", fmt.Errorf("could not enable subscription with id %v", subscriptionId)
	}

	update := map[string]any{
		"active":     true,
		"updated_at": time.Now(),
	}

	err = ws.database.UpdateSubscription(id, update)
	if err != nil {
		err = errors.Wrapf(err, "modifying subscription status failed for sub %v ", subscriptionId)
		log.Err(err).Send()
		return "", fmt.Errorf("could not enable subscription with id %v", subscriptionId)
	}

	return subscriptionId, nil
}

func (ws *WalletService) getTransactionHash(useropHash string) (string, error) {
	useropResult, err := ws.bundlerService.GetUserOperationByHash(useropHash)
	if err != nil {
		err = errors.Wrap(err, "fetching the transction hash failed")
		log.Err(err).Caller().Send()
		return "", err
	}

	transactionHash := useropResult["transactionHash"].(string)
	return transactionHash, nil
}

func randKey(length int) string {
	key := make([]byte, length)

	_, err := rand.Read(key)
	if err != nil {
		// handle error here
	}
	// fmt.Println(key)
	return hexutil.Encode(key)
}

func parseTransferAmount(token string, chain int64, amount float64) *big.Int {
	var divisor int
	if token == "USDC" || token == "USDT" {
		divisor = 6
	} else {
		divisor = 18
	}
	minorFactor := math.Pow10(divisor)
	parsedAmount := int64(amount * minorFactor)

	return big.NewInt(parsedAmount)
}

func isPaymentDue(dueDate time.Time) bool {
	return dueDate.Before(time.Now())
}
