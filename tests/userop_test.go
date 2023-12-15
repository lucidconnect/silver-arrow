package tests

import (
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/joho/godotenv"
	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/lucidconnect/silver-arrow/erc20"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/service/erc4337"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	nodeUrl, entrypointAddress, paymasterUrl string
	defaultChain                             int64
	db                                       *gorm.DB
	// nodeClient                               *erc4337.Client
)

func TestMain(m *testing.M) {
	var err error
	if err = godotenv.Load("../.env.test"); err != nil {
		log.Fatal().Err(err)
	}
	db, err = repository.SetupDatabase(nil)
	if err != nil {
		log.Fatal().Err(err)
	}

	chain := os.Getenv("DEFAULT_CHAIN")
	defaultChain, err = strconv.ParseInt(chain, 10, 64)
	if err != nil {
		panic(err)
	}

	network, err := erc4337.GetNetwork(defaultChain)
	if err != nil {
		panic(err)
	}
	// seedWalletsTable(db)
	nodeUrl = os.Getenv(fmt.Sprintf("%s_NODE_URL", network))

	entrypointAddress = os.Getenv("ENTRY_POINT")
	paymasterUrl = os.Getenv("PAYMASTER_URL")

	// nodeClient, err = erc4337.Dial(nodeUrl, paymasterUrl)
	// if err != nil {
	// 	log.Fatal().Err(err)
	// }

	exitVal := m.Run()
	// clearTables(db)
	os.Exit(exitVal)
}

func getType(strukt interface{}) string {
	if t := reflect.TypeOf(strukt); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func seedWalletsTable(db *gorm.DB) {
	q := "INSERT INTO wallets (id, email, signer_address, wallet_address, turnkey_sub_org_id, turnkey_sub_org_name) VALUES (0, 'gb@backdrop.photo', '0x85fc2E4425d0DAba7426F50091a384ee05D37Cd2', '0x6a6F07c5c32F5fb20393a2110B2Bf0925e59571b','123','random-123')"

	if err := db.Exec(q).Error; err != nil {
		log.Fatal().Err(err)
	}
}

func clearTables(db *gorm.DB) {
	for _, table := range []interface{}{&models.Subscription{}, &models.Key{}, &models.Wallet{}} {
		log.Info().Msgf("Clearing %v table", getType(table))
		if err := db.Where("TRUE").Delete(table).Error; err != nil {
			log.Fatal().Err(err)
		}
	}
}

func TestGetUserOperationByHash(t *testing.T) {
	// nodeClient, err := erc4337.Dial(nodeUrl, "")
	// if !assert.NoError(t, err) {
	// 	t.FailNow()
	// }
	ercBundler, err := erc4337.NewAlchemyService(defaultChain)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	useropshash := "0x28b45cf378c23fbdbbcb4f4c4d085791eb6d660214ff4a2402e40fd1c73751c6"
	_, err = ercBundler.GetUserOperationByHash(useropshash)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	t.Fail()
}

/*
curl --location 'https://api.stackup.sh/v1/node/fc4b8aee3102327ddd59941bfa616d631f0d458032ef71b8a9a28b005c1c2f06' \
--header 'Content-Type: application/json' \
--data '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "eth_getUserOperationByHash",
  "params": ["0x12dbc9b0412c10728dd08dc89ec1b79c89675c3080c009ff590456e2b7cda5a7"]
}
'
*/

func TestValidator(t *testing.T) {
	sender := "0x3D073632A7a29b2AdcbF12D2712fA3E72fABc3dD"
	key := "" // sensitive private key, need better test flow
	// sessionKey := "0x6574f281AAaA788cf89e5269E9c842E50c5713fe"
	// privKey := "0xc1fce60cfb4b32bf4584e577904d806f8c5af28104d34e9923466eb8ca6faeff"
	validatorAddress := "0xc621C75900dbF0234dD7c9be8f8A3bCA7433E3E8"
	executorAddress := "0x701d90246bbd60ED47134CD16Fdd14c570C1Ca51"
	sessionKey, privKey, _ := auth.CreateAccessKey()
	mode := erc4337.ENABLE_MODE
	chainId := 80001

	ercBundler, _ := erc4337.NewAlchemyService(defaultChain)

	nonce, err := ercBundler.GetAccountNonce(common.HexToAddress(sender))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	validator, err := erc4337.InitialiseValidator(validatorAddress, executorAddress, sessionKey, mode, int64(chainId))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	enableData, err := validator.GetEnableData()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	callData, err := validator.SetExecution(enableData, sender)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	op, err := ercBundler.CreateUnsignedUserOperation(sender, nil, callData, nonce, false, int64(chainId))
	assert.NoError(t, err)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	sig, _, err := erc4337.SignUserOp(op, key, erc4337.SUDO_MODE, nil, int64(chainId))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	op["signature"] = hexutil.Encode(sig)

	opHash, err := ercBundler.SendUserOperation(op)
	fmt.Println("user operation hash -", opHash)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	time.Sleep(8 * time.Second)
	fmt.Println("----------------------")

	nonce2, err := ercBundler.GetAccountNonce(common.HexToAddress(sender))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	target := "0xB77ce6ec08B85DcC468B94Cea7Cc539a3BbF9510"

	amount := big.NewInt(0)
	data, err := erc4337.TransferErc20Action(common.HexToAddress("0x0FA8781a83E46826621b3BC094Ea2A0212e71B23"), common.HexToAddress(target), amount)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	op, err = ercBundler.CreateUnsignedUserOperation(sender, nil, data, nonce2, true, int64(chainId))
	assert.NoError(t, err)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	sig, _, err = erc4337.SignUserOp(op, privKey, erc4337.VALIDATOR_MODE, nil, int64(chainId))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	fmt.Println(hexutil.Encode(sig))

	op["signature"] = hexutil.Encode(sig)

	opHash, err = ercBundler.SendUserOperation(op)
	fmt.Println("user operation hash -", opHash)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
}

func TestTokenAction(t *testing.T) {

	sender := "0x3D073632A7a29b2AdcbF12D2712fA3E72fABc3dD"
	target := common.HexToAddress("0xB77ce6ec08B85DcC468B94Cea7Cc539a3BbF9510")
	token := "USDC"

	ercBundler, _ := erc4337.NewAlchemyService(defaultChain)

	// 1000000000000000000 = 1 ether
	// 1000000000000000000 = 1 erc20Token
	// 10000000000000000 = 0.01 erc20Token
	amount := big.NewInt(1000000)
	erc20Token := erc20.GetTokenAddress(token, 80001)
	tokenAddress := common.HexToAddress(erc20Token)

	data, err := erc4337.TransferErc20Action(tokenAddress, target, amount)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	nonce, err := ercBundler.GetAccountNonce(common.HexToAddress(sender))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	key := "0xc1fce60cfb4b32bf4584e577904d806f8c5af28104d34e9923466eb8ca6faeff"
	// key := "0xc1fce60cfb4b32bf4584e577904d806f8c5af28104d34e9923466eb8ca6faeff"

	// fmt.Println(hexutil.Encode(initCode))
	chainId := 80001
	op, err := ercBundler.CreateUnsignedUserOperation(sender, nil, data, nonce, false, int64(chainId))
	assert.NoError(t, err)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	sig, _, err := erc4337.SignUserOp(op, key, erc4337.VALIDATOR_MODE, nil, int64(chainId))
	assert.NoError(t, err)
	op["signature"] = hexutil.Encode(sig)
	// send user operation
	opHash, err := ercBundler.SendUserOperation(op)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	fmt.Println("user operation hash -", opHash)
}
