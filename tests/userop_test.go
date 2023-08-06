package tests

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/helicarrierstudio/silver-arrow/erc4337"
	"github.com/helicarrierstudio/silver-arrow/repository"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	nodeUrl, entrypointAddress, paymasterUrl string
	mongoClient                              *mongo.Client
	nodeClient                               *erc4337.Client
)

func TestMain(m *testing.M) {
	var err error
	if err = godotenv.Load("../.env.test"); err != nil {
		log.Fatal(err)
	}
	mongoClient, err = repository.SetupMongoDatabase()
	if err != nil {
		log.Fatal(err)
	}
	entrypointAddress = os.Getenv("ENTRY_POINT")
	nodeUrl = os.Getenv("NODE_URL")
	paymasterUrl = os.Getenv("PAYMASTER_URL")

	nodeClient, err = erc4337.Dial(nodeUrl, paymasterUrl)
	if err != nil {
		log.Fatal(err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestSendUserOp(t *testing.T) {
	// sender := "0x6a6F07c5c32F5fb20393a2110B2Bf0925e59571b"
	// target := "0x605F2a359EFbCf1aAF708153Ec0ED402d0746ACC"
	sender := "0x14De44b6100dE479655D752ECD2230D10F8fA061"
	target := "0xB77ce6ec08B85DcC468B94Cea7Cc539a3BbF9510"
	token := "USDC"

	ercBundler := erc4337.NewERCBundler(entrypointAddress, nodeClient)

	// 1000000000000000000 = 1 ether
	// 1000000000000000000 = 1 erc20Token
	// 10000000000000000 = 0.01 erc20Token
	amount := big.NewInt(1000000)
	data, err := erc4337.CreateTransferCallData(target, token, amount)
	if !assert.NoError(t, err) {
		t.FailNow()
	}


	nonce, _ := ercBundler.AccountNonce(sender)
	key := "0xcea5314e325233134348f39363151a5fff8051a5e48f8ac96b6dd9866bc2336b"

	mId, _ := hexutil.Decode("0x829f80a98190408d9c22d06ef11ecb213c3fde8a388e9b2052bc1eeee89f2fb7")
	fmt.Println((mId))

	chainId := 80001
	op, err := ercBundler.CreateUnsignedUserOperation(sender, target, nil, data, nonce, false, int64(chainId))
	assert.NoError(t, err)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	fmt.Println("user operation", op)
	// signature := []byte{}
	sig, _ := erc4337.SignUserOp(op, key, erc4337.VALIDATOR_MODE, mId, int64(chainId))
	// signature = append(signature, mId...)
	// signature = append(signature, sig...)
	fmt.Println(sig)
	op["signature"] = hexutil.Encode(sig)
	// send user operation
	opHash, err := ercBundler.SendUserOp(op)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	fmt.Println("user operation hash -", opHash) // 0x28b45cf378c23fbdbbcb4f4c4d085791eb6d660214ff4a2402e40fd1c73751c6 0xfcd3b481cc3ba345fcf24c777463baf60dbb1f7475ca297b9259d020044565be

	t.Fail()
}

func TestGetUserOperationByHash(t *testing.T) {
	nodeClient, err := erc4337.Dial(nodeUrl, "")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	ercBundler := erc4337.NewERCBundler(entrypointAddress, nodeClient)

	useropshash := "0x28b45cf378c23fbdbbcb4f4c4d085791eb6d660214ff4a2402e40fd1c73751c6"
	_, err = ercBundler.GetUserOp(useropshash)
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
