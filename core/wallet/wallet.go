package wallet

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lucidconnect/silver-arrow/conversions"
	"github.com/lucidconnect/silver-arrow/erc20"

	// "github.com/lucidconnect/silver-arrow/erc4337"
	"github.com/lucidconnect/silver-arrow/core/service/erc4337"
	"github.com/lucidconnect/silver-arrow/core/service/turnkey"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/server/graphql/wallet/graph/model"
	"github.com/pkg/errors"
	"github.com/rmanzoku/ethutils/ecrecover"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type WalletService struct {
	turnkey          *turnkey.TurnkeyService
	database         repository.Database
	bundlerService   *erc4337.AlchemyService
	validatorAddress string
	executorAddress  string
}

func NewWalletService(r repository.Database, chain int64) *WalletService {
	// validatorAddress := os.Getenv("VALIDATOR_ADDRESS")
	var bundler *erc4337.AlchemyService
	var err error
	var validatorAddress, executorAddress string

	if chain != 0 {
		network, err := erc4337.GetNetwork(chain)
		if err != nil {
			log.Err(err).Msgf("chain %v not supported", chain)
			return nil
		}

		validatorAddress = os.Getenv(fmt.Sprintf("%s_VALIDATOR_ADDRESS", network))
		executorAddress = os.Getenv(fmt.Sprintf("%s_EXECUTOR_ADDRESS", network))
		bundler, err = initialiseBundler(chain)
		if err != nil {
			return nil
		}
	}

	tunkeyService, err := turnkey.NewTurnKeyService()
	if err != nil {
		log.Panic().Err(err).Send()
	}

	return &WalletService{
		turnkey:          tunkeyService,
		database:         r,
		bundlerService:   bundler,
		validatorAddress: validatorAddress,
		executorAddress:  executorAddress,
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

func (ws *WalletService) AddAccount(input Account) error {
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

// func (ws *WalletService) ValidateSubscription(userop map[string]any, chain int64) (*model.TransactionData, error) {
// 	opHash, err := ws.bundlerService.SendUserOperation(userop)
// 	if err != nil {
// 		log.Err(err).Msg("failed to send user op")
// 		return nil, err
// 	}
// 	fmt.Println("validating subscription with userop hash -", opHash)
// 	result, err := ws.database.FindSubscriptionByHash(opHash)
// 	if err != nil {
// 		log.Err(err).Msgf("failed to find subscription with hash %v", opHash)
// 		return nil, err
// 	}

// 	session, err := ws.database.FetchCheckoutSession(result.CheckoutSessionID)
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return nil, err
// 	}

// 	productId := result.ProductID.String()
// 	if err != nil {
// 		log.Err(err).Msg("encoding product id failed")
// 		return nil, err
// 	}
// 	token := result.Token
// 	createdAt := result.CreatedAt.Format(time.RFC3339)
// 	amount := int(result.Amount)
// 	// interval := int(result.Interval)
// 	subData := &model.TransactionData{
// 		Token:         token,
// 		Amount:        amount,
// 		Interval:      model.IntervalType(result.Interval),
// 		IntervalCount: int(result.IntervalCount),
// 		ProductID:     productId,
// 		WalletAddress: result.WalletAddress,
// 		CreatedAt:     createdAt,
// 	}
// 	fmt.Println("subscription result - ", result)

// 	transactionHash, err := ws.getTransactionHash(opHash)
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return nil, err
// 	}
// 	explorerUrl, err := erc20.GetChainExplorer(chain)
// 	if err != nil {
// 		log.Err(err).Msg("failed to get chain explorer url")
// 	}
// 	blockExplorerTx := fmt.Sprintf("%v/tx/%v", explorerUrl, transactionHash)

// 	update := map[string]interface{}{
// 		"active":           true,
// 		"updated_at":       time.Now(),
// 		"transaction_hash": transactionHash,
// 		"status":           model.SubscriptionStatusActive,
// 	}
// 	err = ws.database.UpdateSubscription(result.ID, update)
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return nil, err
// 	}

// 	// if payment is due now, create a payment
// 	if isPaymentDue(result.NextChargeAt) {
// 		// create payment
// 		var sponsored bool
// 		switch os.Getenv("USE_PAYMASTER") {
// 		case "TRUE":
// 			sponsored = true
// 		default:
// 			sponsored = false
// 		}
// 		reference := uuid.New()
// 		payment := &models.Payment{
// 			Type:                  models.PaymentTypeRecurring,
// 			Chain:                 result.Chain,
// 			Token:                 result.Token,
// 			Amount:                result.Amount,
// 			Source:                result.WalletAddress,
// 			WalletID:              result.WalletID,
// 			ProductID:             result.ProductID,
// 			Sponsored:             sponsored,
// 			Reference:             reference,
// 			Destination:           result.MerchantDepositAddress,
// 			SubscriptionID:        result.ID,
// 			SubscriptionPublicKey: result.Key.PublicKey,
// 			TokenAddress:          result.TokenAddress,
// 			Customer:              session.Customer,
// 			CheckoutSessionID:     session.ID,
// 			MerchantID:            session.MerchantID,
// 		}

// 		userop, useropHash, err := ws.CreatePayment(payment)
// 		if err != nil {
// 			err = errors.Wrap(err, "creating payment operation failed")
// 			log.Err(err).Caller().Send()
// 			return nil, err
// 		}

// 		signature, err := ws.SignPaymentOperation(userop, useropHash)
// 		if err != nil {
// 			err = errors.Wrap(err, "signing payment operation failed")
// 			log.Err(err).Caller().Send()
// 			return nil, err
// 		}
// 		userop["signature"] = signature

// 		onchainTx, err := ws.ExecutePaymentOperation(userop, payment.Chain)
// 		if err != nil {
// 			log.Err(err).Send()
// 			return subData, err
// 		}
// 		interval := result.Interval
// 		intervalCount := result.IntervalCount
// 		nextChargeAt := time.Now().Add((time.Duration(result.Interval)))

// 		update := map[string]interface{}{
// 			"expires_at":     nextChargeAt,
// 			"next_charge_at": nextChargeAt,
// 		}
// 		err = ws.database.UpdateSubscription(result.ID, update)
// 		if err != nil {
// 			log.Err(err).Send()
// 		}
// 		subData.TransactionExplorer = onchainTx
// 	} else {
// 		subData.TransactionExplorer = blockExplorerTx
// 	}
// 	// I should trigger a webhook somewhere here
// 	return subData, nil
// }

// func (ws *WalletService) AddSubscription(merchantId uuid.UUID, input NewSubscription, usePaymaster bool, index *big.Int, chain int64) (*model.ValidationData, map[string]any, error) {
// 	var nextChargeAt time.Time
// 	var initCode []byte
// 	var nonce, amount *big.Int

// 	_, err := ws.database.FetchCheckoutSession(input.CheckoutSessionID)
// 	if err != nil {
// 		log.Err(err).Caller().Msg("invalid session id")
// 		return nil, nil, errors.New("invalid sesison id")
// 	}

// 	// check if a subscription already exists for this product
// 	pid := input.ProductID
// 	existingSub, _ := ws.database.FindSubscriptionByProductId(pid, input.WalletAddress)
// 	if existingSub != nil {
// 		log.Info().Msg("an active subscription exists for this product")
// 		return nil, nil, errors.New("an active subscription exists for this product cancel subscription before creating a new one")
// 	}

// 	// NB: figure out a way to check if the subscription exist without having to do the above operation
// 	// fetch the product by id, use the details in the product to create a subscription
// 	product, err := ws.database.FetchProduct(pid)
// 	if err != nil {
// 		log.Err(err).Msgf("failed to fetch product with id [%v]", pid)
// 		return nil, nil, errors.New("product not found")
// 	}

// 	tagId, orgId, walletID, err := ws.database.GetWalletMetadata(input.WalletAddress)
// 	if err != nil {
// 		log.Err(err).Msgf("failed to fetch private key tag for wallet - %v", input.WalletAddress)
// 		return nil, nil, errors.New("failed to fetch wallet metadata")
// 	}

// 	randomSalt := randKey(4)
// 	keyName := fmt.Sprintf("sub-%v-%v", randomSalt, input.ProductID)
// 	activityId, err := ws.turnkey.CreatePrivateKey(orgId, keyName, tagId)
// 	if err != nil {
// 		log.Err(err).Msg("failed to create subscription private key")
// 		return nil, nil, err
// 	}

// 	result, err := ws.turnkey.GetActivity(orgId, activityId)
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return nil, nil, err
// 	}
// 	privateKeyID, sessionKey, err := turnkey.GetPrivateKeyIdFromResult(result)
// 	// sessionKey, signingKey, err := CreateAccessKey()
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return nil, nil, err
// 	}

// 	// interval := daysToNanoSeconds(int64(input.Interval))
// 	intervalCount := input.Interval
// 	amount = conversions.ParseTransferAmount(input.Token, input.Amount)

// 	if input.NextChargeDate != nil {
// 		nextChargeAt = *input.NextChargeDate
// 	} else {
// 		nextChargeAt = time.Now().Add(interval)
// 	}

// 	isAccountDeployed := ws.isAccountDeployed(input.WalletAddress, chain)
// 	if !isAccountDeployed {
// 		initCode, err = GetContractInitCode(common.HexToAddress(input.OwnerAddress), index)
// 		if err != nil {
// 			log.Err(err).Caller().Send()
// 			return nil, nil, err
// 		}
// 		nonce = common.Big0
// 	} else {
// 		nonce, err = ws.bundlerService.GetAccountNonce(common.HexToAddress(input.WalletAddress))
// 		if err != nil {
// 			log.Err(err).Caller().Send()
// 			return nil, nil, err
// 		}
// 	}

// 	callData, err := setValidatorExecutor(sessionKey, ws.validatorAddress, ws.executorAddress, input.WalletAddress, int64(input.Chain))
// 	if err != nil {
// 		log.Err(err).Msg("failed to set a validator")
// 		return nil, nil, err
// 	}

// 	op, err := ws.bundlerService.CreateUnsignedUserOperation(input.WalletAddress, initCode, callData, nonce, usePaymaster, int64(input.Chain))
// 	if err != nil {
// 		log.Err(err).Msg("failed to create user operation")
// 		return nil, nil, err
// 	}

// 	entrypoint := erc4337.GetEntryPointAddress()
// 	operation, err := userop.New(op)
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return nil, nil, err
// 	}
// 	opHash := operation.GetUserOpHash(entrypoint, big.NewInt(int64(input.Chain)))

// 	key := &models.Key{
// 		WalletID:     walletID,
// 		PublicKey:    sessionKey,
// 		PrivateKeyId: privateKeyID,
// 	}

// 	tokenAddress := erc20.GetTokenAddress(input.Token, chain)
// 	sub := &models.Subscription{
// 		Token:                  input.Token,
// 		Amount:                 amount.Int64(),
// 		Active:                 false,
// 		Interval:               interval.Nanoseconds(),
// 		IntervalCount:          int64(interval),
// 		UserOpHash:             opHash.Hex(),
// 		MerchantId:             merchantId.String(),
// 		ProductID:              input.ProductID,
// 		ProductName:            product.Name,
// 		CheckoutSessionID:      input.CheckoutSessionID,
// 		MerchantDepositAddress: product.DepositAddress,
// 		NextChargeAt:           nextChargeAt,
// 		ExpiresAt:              nextChargeAt,
// 		WalletID:               walletID,
// 		WalletAddress:          input.WalletAddress,
// 		Chain:                  chain,
// 		Key:                    *key,
// 		TokenAddress:           tokenAddress,
// 	}

// 	err = ws.database.AddSubscription(sub, key)
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return nil, nil, err
// 	}

// 	// fmt.Println("New Subscription added", opHash.Hex())
// 	return &model.ValidationData{
// 		UserOpHash: opHash.Hex(),
// 	}, op, nil
// }

func (w *WalletService) FetchSubscriptions(walletAddress string, status *string) ([]*model.SubscriptionData, error) {
	var subData []*model.SubscriptionData
	subs, err := w.database.FetchWalletSubscriptions(walletAddress, status)
	if err != nil {
		log.Err(err).Msgf("error while fetching subscriotions for %v", walletAddress)
		return nil, err
	}

	for _, v := range subs {
		var payments []*model.Payment
		for _, p := range v.Payments {
			payments = append(payments, &model.Payment{
				Chain:     int(p.Chain),
				Token:     p.Token,
				Status:    model.PaymentStatus(p.Status),
				Amount:    conversions.ParseTransferAmountFloat(p.Token, p.Amount),
				Reference: p.Reference.String(),
			})
		}
		// interval := conversions.ParseNanoSecondsToDay(v.Interval)
		createdAt := v.CreatedAt.Format("dd:mm:yyyy")
		sd := &model.SubscriptionData{
			ID:     v.ID.String(),
			Token:  v.Token,
			Amount: int(v.Amount),
			// Interval:       int(interval),
			MerchantID:     v.MerchantId,
			ProductID:      v.ProductID.String(),
			ProductName:    v.ProductName,
			CreatedAt:      createdAt,
			NextChargeDate: v.NextChargeAt,
			Payments:       payments,
			WalletAddress:  walletAddress,
		}
		subData = append(subData, sd)
	}

	return subData, nil
}

func (ws *WalletService) FetchPayment(reference string) (*model.Payment, error) {
	var paymentData *model.Payment
	ref, err := uuid.Parse(reference)
	if err != nil {
		log.Err(err).Msg("invalid reference")
		return nil, errors.New("invalid reference")
	}
	payment, err := ws.database.FindPaymentByReference(ref)
	if err != nil {
		log.Err(err).Msgf("payment [%v] does not exist", reference)
		return nil, errors.New("invalid payment reference")
	}

	paymentData = &model.Payment{
		Chain:     int(payment.Chain),
		Token:     payment.Token,
		Status:    model.PaymentStatus(payment.Status),
		Amount:    conversions.ParseTransferAmountFloat(payment.Token, payment.Amount),
		Source:    payment.Source,
		ProductID: payment.ProductID.String(),
		Reference: payment.Reference.String(),
	}

	return paymentData, nil
}

// func amountToWei(amount any) (*big.Int, error) {
// 	etherInWei := new(big.Int)
// 	etherInWei.SetString("1000000000000000000", 10)

// 	switch v := amount.(type) {
// 	case *big.Int:
// 		weiAmount := new(big.Int).Mul(v, etherInWei)
// 		return weiAmount, nil
// 	case *big.Float:
// 		weiAmount := new(big.Int)
// 		weiAmountFloat := new(big.Float).Mul(v, big.NewFloat(1e18))
// 		weiAmountFloat.Int(weiAmount)
// 		return weiAmount, nil
// 	default:
// 		return nil, fmt.Errorf("unsupported input type: %T", amount)
// 	}
// }

// func amountToMwei(amount int64) *big.Int {
// 	etherInMWei := new(big.Int)
// 	return etherInMWei.SetInt64(amount)
// }

// CreatePayment creates a userop for an initiated payment, amount is already in the minor factor form
// generates the userop hash, sets the payment status to a pending state
// returns a message to be signed.
func (ws *WalletService) CreatePayment(payment *models.Payment) (map[string]any, common.Hash, error) {
	tokenAddress := common.HexToAddress(payment.TokenAddress)

	actualAmount := conversions.ParseAmountToMwei(payment.Amount)
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

	payment.Status = models.PaymentStatusPending
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
		return "", fmt.Errorf("internal server error")
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
		"status":            models.PaymentStatusSuccess,
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

	actualAmount := conversions.ParseAmountToMwei(amount)
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
func isNativeToken(token string, chain int64) bool {
	nativeToken := erc20.GetNativeToken(chain)
	return token == nativeToken
}

// Transfer tokens from a smart wallet,
// authorised by the user's EOA. This is necessary to provide users an interface to move assets
// while this method also works for erc20 tokens, note that USDc is the primary supported token
// hence, using this method to transfer othet tokens with higher decimals will result in unexpected behavior
func (ws *WalletService) InitiateTransfer(sender, target, token string, amount float64, chain int64, sponsored bool) (*model.ValidationData, map[string]any, error) {
	var callData []byte
	var nativeToken bool
	bundler, err := erc4337.NewAlchemyService(chain)
	if err != nil {
		log.Err(err).Msg("failed to initialise bundler")
		return nil, nil, err
	}

	tok, err := ws.database.FetchOneToken(token, chain)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, nil, err
	}

	tokenAddress := tok.Address
	if isNativeToken(token, chain) {
		nativeToken = true
	}
	transferAmount := conversions.ParseTransferAmount(token, amount)
	callData, err = erc4337.CreateTransferCallData(target, tokenAddress, chain, transferAmount, nativeToken)
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
func (ws *WalletService) ExecuteUserop(userop map[string]any, chain int64) (*model.TransactionData, error) {
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
		TransactionHash:     transactionHash,
		TransactionExplorer: blockExplorerTx,
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
func setValidatorExecutor(sessionKey, validatorAddress, executorAddress, ownerAddress string, chain int64) ([]byte, error) {
	mode := erc4337.ENABLE_MODE
	validator, err := erc4337.InitialiseValidator(validatorAddress, executorAddress, sessionKey, mode, chain)
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
func (ws *WalletService) CancelSubscription(subscriptionId string) (string, map[string]any, error) {
	id, err := uuid.Parse(subscriptionId)
	if err != nil {
		err = errors.Wrapf(err, "parsing subscription id %v failed", subscriptionId)
		log.Err(err).Send()
		return "", nil, err
	}
	sub, err := ws.database.FindSubscriptionById(id)
	if err != nil {
		err = errors.Wrapf(err, "failed to fetch subscription %v", subscriptionId)
		log.Err(err).Send()
		return "", nil, err
	}

	address := common.HexToAddress(sub.WalletAddress)
	subKey := common.HexToAddress(sub.Key.PublicKey)

	disableValidatorCallData, err := erc4337.DisableValidator(subKey)
	if err != nil {
		err = errors.Wrap(err, "failed to fetch validator disable calldata")
		log.Err(err).Send()
		return "", nil, err
	}

	callData, err := erc4337.GetExecuteFnData(ws.validatorAddress, common.Big0, disableValidatorCallData)
	if err != nil {
		err = errors.Wrap(err, "failed to create final call data")
		return "", nil, err
	}

	nonce, err := ws.bundlerService.GetAccountNonce(address)
	if err != nil {
		return "", nil, err
	}
	op, err := ws.bundlerService.CreateUnsignedUserOperation(address.Hex(), nil, callData, nonce, true, sub.Chain)
	if err != nil {
		return "", nil, err
	}

	userop, err := userop.New(op)
	if err != nil {
		return "", nil, err
	}

	entryPoint := common.HexToAddress(ws.bundlerService.EntryPoint)
	hash := userop.GetUserOpHash(entryPoint, big.NewInt(sub.Chain))

	update := map[string]any{
		"active":       false,
		"updated_at":   time.Now(),
		"cancelled_at": time.Now(),
		"status":       model.SubscriptionStatusCancelled,
	}

	err = ws.database.UpdateSubscription(id, update)
	if err != nil {
		err = errors.Wrapf(err, "modifying subscription status failed for sub %v ", subscriptionId)
		log.Err(err).Send()
		return "", nil, fmt.Errorf("could not disable subscription with id %v", subscriptionId)
	}

	// ws.database.DeactivateSubscription()
	return hash.Hex(), op, nil
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
		"active":      false,
		"updated_at":  time.Now(),
		"disabled_at": time.Now(),
		"status":      model.SubscriptionStatusDisabled,
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
		"status":     model.SubscriptionStatusActive,
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

func (ws *WalletService) FetchUserBillingHistory(walletAddress, productID string) ([]BillingHistory, error) {
	var billingHistory []BillingHistory

	pid, err := uuid.Parse(productID)
	if err != nil {
		log.Err(err).Msg("parsing uuid failed")
		return nil, err
	}

	payments, err := ws.database.FindAllPaymentsByWallet(walletAddress)
	if err != nil {
		log.Err(err).Msgf("failed to fetch payments from wallet %v", walletAddress)
		return nil, err
	}

	for _, payment := range payments {
		if payment.ProductID == pid {
			amount := conversions.ParseTransferAmountFloat(payment.Token, payment.Amount)
			bh := BillingHistory{
				Date:        payment.CreatedAt,
				Amount:      amount,
				ExplorerURL: payment.BlockExplorerTx,
			}
			billingHistory = append(billingHistory, bh)
		}
	}

	return billingHistory, nil
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

// func parseTransferAmount(token string, amount float64) *big.Int {
// 	var divisor int
// 	if token == "USDC" || token == "USDT" {
// 		divisor = 6
// 	} else {
// 		divisor = 18
// 	}
// 	minorFactor := math.Pow10(divisor)
// 	parsedAmount := int64(amount * minorFactor)

// 	return big.NewInt(parsedAmount)
// }

// func parseTransferAmountFloat(token string, amount int64) float64 {
// 	var divisor int
// 	if token == "USDC" || token == "USDT" {
// 		divisor = 6
// 	} else {
// 		divisor = 18
// 	}
// 	minorFactor := math.Pow10(divisor)
// 	parsedAmount := float64(amount) / minorFactor

// 	return parsedAmount
// }

func isPaymentDue(dueDate time.Time) bool {
	return dueDate.Before(time.Now())
}
