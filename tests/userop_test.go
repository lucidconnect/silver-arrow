package tests

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/helicarrierstudio/silver-arrow/bundlerclient"
	"github.com/helicarrierstudio/silver-arrow/useroperation"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var nodeUrl, entrypointAddress, paymasterUrl string

func TestMain(m *testing.M) {
	if err := godotenv.Load(".env.test"); err != nil {
		log.Fatal(err)
	}
	entrypointAddress = os.Getenv("ENTRY_POINT")
	nodeUrl = os.Getenv("NODE_URL")
	paymasterUrl = os.Getenv("PAYMASTER_URL")

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestSendUserOp(t *testing.T) {
	sender := "0x6a6F07c5c32F5fb20393a2110B2Bf0925e59571b"
	target := "0x605F2a359EFbCf1aAF708153Ec0ED402d0746ACC"
	token := "ETH"

	nodeClient, err := bundlerclient.Dial(nodeUrl, paymasterUrl)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	ercBundler := useroperation.NewERCBundler(entrypointAddress, nodeClient)

	// 1000000000000000000 = 1 ether
	// 1000000000000000000 = 1 erc20Token
	// 10000000000000000 = 0.01 erc20Token
	amount := big.NewInt(6000000000000000)
	data, err := useroperation.CreateTransferCallData(target, token, amount)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	nonce := big.NewInt(8)
	op, err := ercBundler.CreateUserOperation(sender, target, token, data, nonce, amount, true)
	assert.NoError(t, err)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	fmt.Println("user operation", op)

	// send user operation
	opHash, err := ercBundler.SendUserOp(op)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	fmt.Println("user operation hash -", opHash) // 0x28b45cf378c23fbdbbcb4f4c4d085791eb6d660214ff4a2402e40fd1c73751c6 0xfcd3b481cc3ba345fcf24c777463baf60dbb1f7475ca297b9259d020044565be

	t.Fail()
}

func TestGetUserOperationByHash(t *testing.T) {
	nodeClient, err := bundlerclient.Dial(nodeUrl, "")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	ercBundler := useroperation.NewERCBundler(entrypointAddress, nodeClient)

	useropshash := "0x28b45cf378c23fbdbbcb4f4c4d085791eb6d660214ff4a2402e40fd1c73751c6"
	err = ercBundler.GetUserOp(useropshash)
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
