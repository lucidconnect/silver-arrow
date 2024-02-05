package gateway

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/conversions"
	"github.com/lucidconnect/silver-arrow/core"
	"github.com/lucidconnect/silver-arrow/core/service/erc4337"
	"github.com/lucidconnect/silver-arrow/core/service/turnkey"
	"github.com/lucidconnect/silver-arrow/core/wallet"
	"github.com/lucidconnect/silver-arrow/erc20"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/server/graphql/wallet/graph/model"
	"github.com/pkg/errors"
	"github.com/rmanzoku/ethutils/ecrecover"
	"github.com/rs/zerolog/log"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"gorm.io/gorm"
)

const (
	rate    int64 = 5
	divisor int64 = 1000
)

// PaymentGateway is the entrypoint for all things payment
// initialises a checkout session, manages the lifecycle of the session
// This gateway also manages the lifecycle of a payment, from when the charge is initiated to completion
type PaymentGateway struct {
	// ProductId        uuid.UUID
	MerchantId uuid.UUID

	executorAddress  string
	validatorAddress string
	cache            repository.CacheWizard
	database         repository.Database
	turnkey          *turnkey.TurnkeyService
	bundlerService   *erc4337.AlchemyService
}

func NewPaymentGateway(r repository.Database, chain int64) *PaymentGateway {
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

	return &PaymentGateway{
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

func (p *PaymentGateway) InitialiseCheckoutSession(product, customer, callbackUrl string) (*core.CheckoutSession, error) {
	var url string
	sessionId := uuid.New()
	productId, err := uuid.Parse(product)
	if err != nil {
		log.Err(err).Caller().Msg("invalid product id")
		return nil, err
	}

	paymentLink, err := p.database.FetchPaymentLinkByProduct(productId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// create a payment link
			environment := os.Getenv("APP_ENV")
			switch environment {
			case "staging":
				url = fmt.Sprintf("https://pay.staging.lucidconnect.xyz/c/%v", sessionId.String())
			case "production":
				url = fmt.Sprintf("https://pay.lucidconnect.xyz/c/%v", sessionId.String())
			}
			paymentLink = &models.PaymentLink{
				ID:          uuid.New(),
				MerchantID:  p.MerchantId,
				CallbackURL: callbackUrl,
				ProductID:   productId,
				Url:         url,
				CreatedAt:   time.Now(),
			}
			err = p.database.CreatePaymentLink(paymentLink)
			if err != nil {
				log.Err(err).Msgf("could not create payment link for product [%v]", productId)
				return nil, err
			}
		} else {
			log.Err(err).Caller().Send()
			return nil, err
		}
	}

	newSession := &models.CheckoutSession{
		ID:            sessionId,
		Customer:      customer,
		ProductID:     productId,
		MerchantID:    p.MerchantId,
		PaymentLinkID: paymentLink.ID,
		CallbackURL:   callbackUrl,
		State:         string(core.SessionActive),
	}

	if err = p.database.CreateCheckoutSession(newSession); err != nil {
		log.Err(err).Send()
		return nil, err
	}

	session := &core.CheckoutSession{
		Id:          sessionId,
		ProductId:   productId,
		MerchantId:  p.MerchantId,
		CustomerId:  customer,
		CallbackURL: callbackUrl,
		PaymentLink: url,
	}
	return session, nil
}

func (p *PaymentGateway) CreatePaymentIntent(intent core.PaymentIntent) (map[string]any, string, error) {
	var usePaymaster bool
	var userOp map[string]any
	var useropHash string

	switch os.Getenv("USE_PAYMASTER") {
	case "TRUE":
		usePaymaster = true
	default:
		usePaymaster = false
	}

	// fetch product
	productId, err := uuid.Parse(intent.ProductId)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, "", err
	}

	product, err := p.database.FetchProduct(productId)
	if err != nil {
		log.Err(err).Send()
		return nil, "", err
	}

	priceId, err := uuid.Parse(intent.PriceId)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, "", err
	}

	price, err := p.database.FetchPrice(priceId)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, "", err
	}

	sessionId := intent.CheckoutSessionId
	switch intent.Type {
	case core.PaymentTypeRecurring:
		var nextCharge time.Time

		if intent.FirstChargeNow {
			nextCharge = time.Now()
		} else {
			interval := price.Interval
			intervalUnit := price.IntervalUnit
			nextCharge = CalculateNextChargeDate(intervalUnit, interval)
		}

		var email string
		if intent.Email != "" {
			email = intent.Email
		}

		newSubscription := core.NewSubscription{
			Chain:             price.Chain,
			Token:             price.Token,
			Email:             email,
			Amount:            price.Amount,
			Interval:          price.Interval,
			IntervalUnit:      price.IntervalUnit,
			ProductID:         productId,
			CheckoutSessionID: sessionId,
			ProductName:       product.Name,
			OwnerAddress:      intent.OwnerAddress,
			WalletAddress:     intent.WalletAddress,
			DepositAddress:    product.DepositAddress,
			NextChargeDate:    nextCharge,
		}

		validationData, op, err := p.AddSubscription(p.MerchantId, newSubscription, usePaymaster, common.Big0, price.Chain)
		if err != nil {
			return nil, "", err
		}
		userOp = op
		fmt.Println("Userop hash", validationData.UserOpHash)
		useropHash = validationData.UserOpHash
	case core.PaymentTypeSingle:
		// no need to create a subscription
		// requires user validation
		session, err := p.database.FetchCheckoutSession(sessionId)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, "", err
		}
		var sponsored bool
		switch os.Getenv("USE_PAYMASTER") {
		case "TRUE":
			sponsored = true
		default:
			sponsored = false
		}
		reference := uuid.New()
		token, err := p.database.FetchOneToken(price.Token, price.Chain)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, "", err
		}

		tokenAddress := token.Address

		// tokenAddress := erc20.GetTokenAddress(price.Token, price.Chain)

		_, _, walletID, err := p.database.GetWalletMetadata(intent.WalletAddress)
		if err != nil {
			log.Err(err).Msgf("failed to fetch private key tag for wallet - %v", intent.WalletAddress)
			return nil, "", errors.New("failed to fetch wallet metadata")
		}

		// fetch the deposit address for the product
		destinationAddress := product.DepositWallets
		// get deposit address
		payment := &models.Payment{
			Type:      models.PaymentTypeSingle,
			Chain:     price.Chain,
			Token:     price.Token,
			Amount:    price.Amount,
			Source:    intent.WalletAddress,
			WalletID:  walletID,
			ProductID: productId,
			Sponsored: sponsored,
			Reference: reference,
			// Destination:       product.DepositAddress,
			DestinationAddress: destinationAddress,
			TokenAddress:       tokenAddress,
			Customer:           session.Customer,
			CheckoutSessionID:  session.ID,
			MerchantID:         session.MerchantID,
		}

		userop, hash, err := p.CreatePayment(payment)
		useropHash = hash.Hex()
		userOp = userop
		if err != nil {
			err = errors.Wrap(err, "creating payment operation failed")
			log.Err(err).Caller().Send()
			return nil, "", err
		}
	default:
		return nil, "", errors.New("unsupported payment type")
	}

	return userOp, useropHash, nil
}

func (p *PaymentGateway) AddSubscription(merchantId uuid.UUID, newSub core.NewSubscription, usePaymaster bool, index *big.Int, chain int64) (*model.ValidationData, map[string]any, error) {
	var nextChargeAt time.Time
	var initCode []byte
	var nonce, amount *big.Int

	_, err := p.database.FetchCheckoutSession(newSub.CheckoutSessionID)
	if err != nil {
		log.Err(err).Caller().Msg("invalid session id")
		return nil, nil, errors.New("invalid sesison id")
	}

	// check if a subscription already exists for this product
	pid := newSub.ProductID
	existingSub, _ := p.database.FindSubscriptionByProductId(pid, newSub.WalletAddress)
	if existingSub != nil {
		log.Info().Msg("an active subscription exists for this product")
		return nil, nil, errors.New("an active subscription exists for this product cancel subscription before creating a new one")
	}

	// NB: figure out a way to check if the subscription exist without having to do the above operation
	// fetch the product by id, use the details in the product to create a subscription
	product, err := p.database.FetchProduct(pid)
	if err != nil {
		log.Err(err).Msgf("failed to fetch product with id [%v]", pid)
		return nil, nil, errors.New("product not found")
	}

	tagId, orgId, walletID, err := p.database.GetWalletMetadata(newSub.WalletAddress)
	if err != nil {
		log.Err(err).Msgf("failed to fetch private key tag for wallet - %v", newSub.WalletAddress)
		return nil, nil, errors.New("failed to fetch wallet metadata")
	}

	randomSalt := randKey(4)
	keyName := fmt.Sprintf("sub-%v-%v", randomSalt, newSub.ProductID)
	activityId, err := p.turnkey.CreatePrivateKey(orgId, keyName, tagId)
	if err != nil {
		log.Err(err).Msg("failed to create subscription private key")
		return nil, nil, err
	}

	result, err := p.turnkey.GetActivity(orgId, activityId)
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

	// interval := daysToNanoSeconds(int64(newSub.Interval))
	// intervalCount := newSub.Interval
	// newSub.Amount was gotten from the price object which is already stored in the tokens denominational value
	amount = conversions.ParseAmountToMwei(newSub.Amount)

	// if !newSub.NextChargeDate.IsZero() {
	nextChargeAt = newSub.NextChargeDate
	// } else {
	// 	nextChargeAt = time.Now().Add(interval)
	// }

	isAccountDeployed := p.bundlerService.IsAccountDeployed(newSub.WalletAddress, chain)
	if !isAccountDeployed {
		initCode, err = wallet.GetContractInitCode(common.HexToAddress(newSub.OwnerAddress), index)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, nil, err
		}
		nonce = common.Big0
	} else {
		nonce, err = p.bundlerService.GetAccountNonce(common.HexToAddress(newSub.WalletAddress))
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, nil, err
		}
	}

	callData, err := SetValidatorExecutor(sessionKey, p.validatorAddress, p.executorAddress, newSub.WalletAddress, int64(newSub.Chain))
	if err != nil {
		log.Err(err).Msg("failed to set a validator")
		return nil, nil, err
	}

	op, err := p.bundlerService.CreateUnsignedUserOperation(newSub.WalletAddress, initCode, callData, nonce, usePaymaster, int64(newSub.Chain))
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
	opHash := operation.GetUserOpHash(entrypoint, big.NewInt(int64(newSub.Chain)))

	key := &models.Key{
		WalletID:     walletID,
		PublicKey:    sessionKey,
		PrivateKeyId: privateKeyID,
	}

	tokenAddress := erc20.GetTokenAddress(newSub.Token, chain)
	sub := &models.Subscription{
		Token:                  newSub.Token,
		Amount:                 amount.Int64(),
		Active:                 false,
		Interval:               newSub.Interval,
		IntervalUnit:           newSub.IntervalUnit,
		UserOpHash:             opHash.Hex(),
		MerchantId:             merchantId.String(),
		ProductID:              newSub.ProductID,
		ProductName:            product.Name,
		CheckoutSessionID:      newSub.CheckoutSessionID,
		MerchantDepositAddress: product.DepositAddress,
		NextChargeAt:           nextChargeAt,
		ExpiresAt:              nextChargeAt,
		WalletID:               walletID,
		WalletAddress:          newSub.WalletAddress,
		Chain:                  chain,
		Key:                    *key,
		TokenAddress:           tokenAddress,
	}

	err = p.database.AddSubscription(sub, key)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, nil, err
	}

	// fmt.Println("New Subscription added", opHash.Hex())
	return &model.ValidationData{
		UserOpHash: opHash.Hex(),
	}, op, nil
}

func (p *PaymentGateway) ValidatePaymentIntent(userop map[string]any, chain int64, paymentType string) (*model.TransactionData, error) {
	var transactionData *model.TransactionData
	opHash, err := p.bundlerService.SendUserOperation(userop)
	if err != nil {
		log.Err(err).Msg("failed to send user op")
		return nil, err
	}

	switch paymentType {
	case "single":
		result, err := p.database.FindPaymentByUseropHash(opHash)
		if err != nil {
			log.Err(err).Msgf("failed to find payment with hash %v", opHash)
			return nil, err
		}

		session, err := p.database.FetchCheckoutSession(result.CheckoutSessionID)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, err
		}
		createdAt := result.CreatedAt.Format(time.RFC3339)
		amount := int(result.Amount)
		transactionData = &model.TransactionData{
			Token:         result.Token,
			Amount:        amount,
			ProductID:     session.ProductID.String(),
			WalletAddress: result.Source,
			CreatedAt:     createdAt,
		}

		onchainTx, err := p.ExecutePaymentOperation(userop, chain)
		if err != nil {
			log.Err(err).Send()
			return transactionData, err
		}
		transactionData.TransactionExplorer = onchainTx
	case "recurring":
		fmt.Println("validating subscription with userop hash -", opHash)
		result, err := p.database.FindSubscriptionByHash(opHash)
		if err != nil {
			log.Err(err).Msgf("failed to find subscription with hash %v", opHash)
			return nil, err
		}

		session, err := p.database.FetchCheckoutSession(result.CheckoutSessionID)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, err
		}

		productId := result.ProductID.String()
		if err != nil {
			log.Err(err).Msg("encoding product id failed")
			return nil, err
		}
		token := result.Token
		createdAt := result.CreatedAt.Format(time.RFC3339)
		amount := int(result.Amount)
		// interval := int(result.Interval)
		transactionData = &model.TransactionData{
			Token:         token,
			Amount:        amount,
			Interval:      model.IntervalType(result.IntervalUnit),
			IntervalCount: int(result.Interval),
			ProductID:     productId,
			WalletAddress: result.WalletAddress,
			CreatedAt:     createdAt,
		}
		fmt.Println("subscription result - ", result)

		transactionHash, err := p.bundlerService.GetTransactionHash(opHash)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, err
		}
		explorerUrl, err := erc20.GetChainExplorer(chain)
		if err != nil {
			log.Err(err).Msg("failed to get chain explorer url")
		}
		blockExplorerTx := fmt.Sprintf("%v/tx/%v", explorerUrl, transactionHash)

		update := map[string]interface{}{
			"active":           true,
			"updated_at":       time.Now(),
			"transaction_hash": transactionHash,
			"status":           model.SubscriptionStatusActive,
		}
		err = p.database.UpdateSubscription(result.ID, update)
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
				Type:                  models.PaymentTypeRecurring,
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
				TokenAddress:          result.TokenAddress,
				Customer:              session.Customer,
				CheckoutSessionID:     session.ID,
				MerchantID:            session.MerchantID,
			}

			userop, useropHash, err := p.CreatePayment(payment)
			if err != nil {
				err = errors.Wrap(err, "creating payment operation failed")
				log.Err(err).Caller().Send()
				return nil, err
			}

			signature, err := p.SignPaymentOperation(userop, useropHash)
			if err != nil {
				err = errors.Wrap(err, "signing payment operation failed")
				log.Err(err).Caller().Send()
				return nil, err
			}
			userop["signature"] = signature

			onchainTx, err := p.ExecutePaymentOperation(userop, payment.Chain)
			if err != nil {
				log.Err(err).Send()
				return transactionData, err
			}
			intervalUnit := result.IntervalUnit
			intervalCount := result.Interval
			nextChargeAt := CalculateNextChargeDate(intervalUnit, intervalCount)

			update := map[string]interface{}{
				"expires_at":     nextChargeAt,
				"next_charge_at": nextChargeAt,
			}
			err = p.database.UpdateSubscription(result.ID, update)
			if err != nil {
				log.Err(err).Send()
			}
			transactionData.TransactionExplorer = onchainTx
		} else {
			transactionData.TransactionExplorer = blockExplorerTx
		}
	}
	// I should trigger a webhook somewhere here
	return transactionData, nil
}

// CreatePayment creates a userop for an initiated payment, amount is already in the minor factor form
// generates the userop hash, sets the payment status to a pending state
// returns a message to be signed.
func (p *PaymentGateway) CreatePayment(payment *models.Payment) (map[string]any, common.Hash, error) {
	tokenAddress := common.HexToAddress(payment.TokenAddress)

	actualAmount := conversions.ParseAmountToMwei(payment.Amount)
	debitInstructions := getDebitInstructions(tokenAddress, payment.DestinationAddress, actualAmount)
	data, err := erc4337.BatchTransferErc20Action(debitInstructions)
	if err != nil {
		err = errors.Wrap(err, "creating TransferErc20Action call data failed")
		log.Err(err).Caller().Send()
		return nil, common.Hash{}, fmt.Errorf("internal server error")
	}

	nonce, err := p.bundlerService.GetAccountNonce(common.HexToAddress(payment.Source))
	if err != nil {
		err = errors.Wrapf(err, "error occured fetching nonce for account [%v]", payment.Source)
		log.Err(err).Caller().Send()
		return nil, common.Hash{}, fmt.Errorf("internal server error")
	}

	op, err := p.bundlerService.CreateUnsignedUserOperation(payment.Source, nil, data, nonce, payment.Sponsored, payment.Chain)
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
	err = p.database.CreatePayment(payment)
	if err != nil {
		log.Err(err).Send()
	}

	return op, userOpHash, nil
}

// SignPaymentOperation takes in the userop alongside it's computed userop hash and returns a signature.
// It signs the hash with the key created for the subscription that is initiating the payment
// note this method is intended only for recurring automated payments.
func (p *PaymentGateway) SignPaymentOperation(op map[string]any, hash common.Hash) (string, error) {
	payment, err := p.database.FindPaymentByUseropHash(hash.Hex())
	if err != nil {
		log.Err(err).Msgf("failed to fetch payment with user op hash %v", hash)
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("payment not found")
		}
		return "", err
	}
	// find the subscription the payment is for
	subscriptionKeyId, err := p.database.GetSubscriptionKey(payment.SubscriptionPublicKey)
	if err != nil {
		err = errors.Wrap(err, "invalid subscription key")
		log.Err(err).Send()
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("this payment was not authorised, key not found")
		}
		return "", err
	}
	message := hexutil.Encode(ecrecover.ToEthSignedMessageHash(hash.Bytes()))

	wallet, err := p.database.FetchAccountByAddress(payment.Source)
	if err != nil {
		err = errors.Wrapf(err, "smart account lookup for address [%v] failed", payment.Source)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("sender address [%v] not found", payment.Source)
	}
	org := wallet.TurnkeySubOrgID

	fmt.Println("Signing user op with key id - ", subscriptionKeyId)
	turnkeyActivityId, err := p.turnkey.SignMessage(org, subscriptionKeyId, message)
	if err != nil {
		err = errors.Wrapf(err, "turnkey failed to sign user operation for account [%v], keyId: [%v]", payment.Source, subscriptionKeyId)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	result, err := p.turnkey.GetActivity(org, turnkeyActivityId)
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
func (p *PaymentGateway) ExecutePaymentOperation(signedOp map[string]any, chain int64) (string, error) {
	opHash, err := p.bundlerService.SendUserOperation(signedOp)
	if err != nil {
		err = errors.Wrap(err, "sending user op failed")
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	// TODO: use the userop hash to create a reciept for the transsaction
	fmt.Println("user operation hash -", opHash) // 0x28b45cf378c23fbdbbcb4f4c4d085791eb6d660214ff4a2402e40fd1c73751c6 0xfcd3b481cc3ba345fcf24c777463baf60dbb1f7475ca297b9259d020044565be

	// Fetch transaction hash
	transactionHash, err := p.bundlerService.GetTransactionHash(opHash)
	if err != nil {
		err = errors.Wrapf(err, "fetching the transction hash failed. userop hash - [%v]", opHash)
		log.Err(err).Caller().Send()
		return "", fmt.Errorf("internal server error")
	}

	payment, err := p.database.FindPaymentByUseropHash(opHash)
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
	err = p.database.UpdatePayment(payment.ID, update)
	if err != nil {
		err = errors.Wrapf(err, "updating payment status failed but transaction was successful on chain [%v]: useropHash - [%v]; transactionHash - [%v]", chain, opHash, transactionHash)
		log.Err(err).Caller().Send()
		return transactionHash, err
	}

	// should probably trigger a webhook event

	return blockExplorerTx, nil
}

func CalculateNextChargeDate(interval string, count int64) time.Time {
	var intervalNs time.Duration
	switch core.RecuringInterval(interval) {
	case core.RecuringIntervalDay:
		intervalNs = conversions.ParseDaysToNanoSeconds(count)
	case core.RecuringIntervalWeek:
		// 7 days make a week
		days := count * 7
		intervalNs = conversions.ParseDaysToNanoSeconds(days)
	case core.RecuringIntervalMonth:
		// I'm using 1 month = 30 days so a customer only get's debited 12 times in a calendar year
		days := count * 30
		intervalNs = conversions.ParseDaysToNanoSeconds(days)
	}

	nextChargeAt := time.Now().Add(intervalNs)

	return nextChargeAt
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

func isPaymentDue(dueDate time.Time) bool {
	return dueDate.Before(time.Now())
}

func calculateFees(transferAmount *big.Int) *big.Int {
	x := transferAmount.Int64() * rate

	fee := x / divisor
	return big.NewInt(fee)
}

func calculateWalletTransferAmount(amount *big.Int, percentage float64) *big.Int {
	// percent = 2.5 (2.5%) == 2.5/100 = 0.025
	// ((percentage * 10) * amount) / 1000

	normalisedPercent := percentage * 10
	walletShare := (amount.Int64() * int64(normalisedPercent)) / divisor
	return big.NewInt(walletShare)
}

func getDebitInstructions(token common.Address, depositAddress []*models.DepositWallet, amount *big.Int) []erc4337.DebitInstruction {
	var debitInstructions []erc4337.DebitInstruction
	fees := calculateFees(amount)
	lucidFeeAddress := os.Getenv("LUCID_FEE_WALLET")

	if lucidFeeAddress != "" {
		feeInstruction := erc4337.DebitInstruction{
			Token:       token,
			Destination: common.HexToAddress(lucidFeeAddress),
			Amount:      fees,
		}

		debitInstructions = append(debitInstructions, feeInstruction)
	}

	for _, address := range depositAddress {
		amount := calculateWalletTransferAmount(amount, address.Percentage)
		i := erc4337.DebitInstruction{
			Token:       token,
			Destination: common.HexToAddress(address.WalletAddress),
			Amount:      amount,
		}
		debitInstructions = append(debitInstructions, i)
	}

	return debitInstructions
}

func (p *PaymentGateway) AddNewToken(name, address string, chain int64) {
	token := &models.Token{
		Chain:   chain,
		Name:    name,
		Address: address,
	}

	if err := p.database.AddToken(token); err != nil {
		log.Err(err).Caller().Send()
		return
	}
}

func (p *PaymentGateway) FetchSupportedTokens(chain int64) ([]models.Token, error) {
	return p.database.FetchAllTokens(chain)
}

func (p *PaymentGateway) FindToken(chain int64, name string) (*models.Token, error) {
	token, err := p.database.FetchOneToken(name, chain)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}

	return token, nil
}
