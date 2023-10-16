package erc4337_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lucidconnect/silver-arrow/erc4337"
	"github.com/stretchr/testify/assert"
)

func TestTransferFnData(t *testing.T) {
	partialERC20TokenABI := `[{
        "constant": false,
        "inputs": [
            {
                "name": "_to",
                "type": "address"
            },
            {
                "name": "_value",
                "type": "uint256"
            }
        ],
        "name": "transfer",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }]`

	amount := big.NewInt(2)
	to := "0x0aDfe6188b916F25062b689E070Aa49bdbe8d161"
	data, err := erc4337.GetTransferFnData(partialERC20TokenABI, to, amount)
	assert.NoError(t, err)
	methodId := data[:4]
	encodedToAddress := common.TrimLeftZeroes(data[4:36])
	encodedAmount := common.TrimLeftZeroes(data[36:])

	assert.Equal(t, "0xa9059cbb", hexutil.Encode(methodId))
	assert.Equal(t, strings.ToLower(to), hexutil.Encode(encodedToAddress))
	assert.Equal(t, amount.Bytes(), encodedAmount)
}

func TestExecuteFnSignature(t *testing.T) {
	partialAccountABI := `[{
		"inputs": [
			{
				"internalType": "address",
				"name": "dest",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "value",
				"type": "uint256"
			},
			{
				"internalType": "bytes",
				"name": "func",
				"type": "bytes"
			}
		],
		"name": "execute",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}]`
	erc20TokenAddress := "0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d"
	data := common.Hex2Bytes("0xa9059cbb0000000000000000000000000adfe6188b916f25062b689e070aa49bdbe8d1610000000000000000000000000000000000000000000000000000000000000002")
	data, err := erc4337.GetExecuteFnData(partialAccountABI, erc20TokenAddress, big.NewInt(0), data)
	assert.NoError(t, err)
	methodId := data[:4]
	assert.Equal(t, "0xb61d27f6", hexutil.Encode(methodId))
}

// func TestGasFee(t *testing.T) {
// 	maxFeePerGasHex := "29260CA6A"
// 	maxFeePerGasBig, ok := new(big.Int).SetString(maxFeePerGasHex, 16)
// 	if !ok {
// 		t.Fail()
// 	}
// 	fmt.Println(maxFeePerGasBig)
// 	_maxFeePerGas := new(big.Int).Mul(maxFeePerGasBig, big.NewInt(10))
// 	maxFeePerGas := new(big.Int).Div(_maxFeePerGas, big.NewInt(7))
// 	fmt.Println("adjusted maxFeePerGas", maxFeePerGas)

// }
