package erc4337

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	KernelStorage "github.com/lucidconnect/silver-arrow/abi/kernelStorage"
	"github.com/lucidconnect/silver-arrow/abi/sessionKeyOwnedValidator"
	"github.com/pkg/errors"
)

/*
*
- The enableData will be appended to the userop signature
- signature[56:88] contains the enableData
*/
type SessionKeyOwnedValidator struct {
	// privatekey       string
	Mode             []byte
	Chain            *big.Int
	ValidUntil       []byte
	ValidAfter       []byte
	SessionKey       common.Address
	ExecutorAddress  common.Address
	ValidatorAddress common.Address
}

func InitialiseValidator(validatorAddress, executorAddress, sessionKey, mode string, chainId int64) (*SessionKeyOwnedValidator, error) {
	validator := common.HexToAddress(validatorAddress)
	executor := common.HexToAddress(executorAddress)
	session := common.HexToAddress(sessionKey)
	md, err := hexutil.Decode(mode)
	if err != nil {
		return nil, err
	}

	return newSessionKeyOwnedValidator(validator, executor, session, md, big.NewInt(chainId)), nil
}

func newSessionKeyOwnedValidator(validator, executor, sessionKey common.Address, mode []byte, chain *big.Int) *SessionKeyOwnedValidator {
	validUntil, _ := parseUint48(uint64(math.Pow(2, 48)) - 1)
	validAfter, _ := parseUint48(0)
	return &SessionKeyOwnedValidator{
		Mode:             mode,
		Chain:            chain,
		ValidUntil:       validUntil,
		ValidAfter:       validAfter,
		SessionKey:       sessionKey,
		ExecutorAddress:  executor,
		ValidatorAddress: validator,
		// privatekey:       privateKey[2:],
	}
}

/**
op.signature = abi.encodePacked(
    bytes4(0x00000002), 4 bytes
    uint48(0), 6 bytes
    uint48(0), 6 bytes
    address(sessionKeyValidator), 20 bytes
    address(0), 20 bytes
    uint256(enableData.length), 32 bytes
    enableData, 32 bytes
    uint256(65), 32 bytes
    r,
    s,
    v
);
*/

func (v *SessionKeyOwnedValidator) GetEnableData() ([]byte, error) {
	// var validAfter, validUntil []bytes{}
	var data []byte
	fmt.Println("valid until ", v.ValidUntil)
	fmt.Println("valid after ", v.ValidAfter)
	data = append(data, v.SessionKey.Bytes()...)
	data = append(data, v.ValidAfter...)
	data = append(data, v.ValidUntil...)
	return data, nil
}

func (v *SessionKeyOwnedValidator) SetExecution(enableData []byte, ownerAccount string) ([]byte, error) {
	kernelAbi, err := KernelStorage.KernelStorageMetaData.GetAbi()
	if err != nil {
		err = errors.Wrap(err, "KernelStorage.GetAbi():")
		return nil, err
	}

	selector := [4]byte{}
	sel, err := hexutil.Decode("0x84189294")
	if err != nil {
		err = errors.Wrap(err, "hexutil.Decode():")
		return nil, err
	}
	copy(selector[:], sel)
	// executorAddress := common.HexToAddress(os.Getenv("EXECUTOR_ADDRESS"))
	// validatorAddress := v.ValidatorAddress
	validUntil := big.NewInt(99999999999)
	validAfter := big.NewInt(0)

	callData, err := kernelAbi.Pack("setExecution", selector, v.ExecutorAddress, v.ValidatorAddress, validUntil, validAfter, enableData)
	if err != nil {
		return nil, err
	}
	// log.Fatalln("exit")

	return callData, nil
}

// DisableValidator creates the call data which is a call to the disable method on the validator
func DisableValidator(sessionKey common.Address) ([]byte, error) {
	sessionKeyValidatorAbi, err := sessionKeyOwnedValidator.SessionKeyOwnedValidatorMetaData.GetAbi()
	if err != nil {
		err = errors.Wrap(err, "SessionKeyOwnedValidator abi error")
		return nil, err
	}

	callData, err := sessionKeyValidatorAbi.Pack("disable", sessionKey.Bytes())
	if err != nil {
		return nil, err
	}
	return callData, nil
}

func parseUint48(value uint64) ([]byte, error) {
	const max48BitValue = 0xFFFFFFFFFFFF
	if value > max48BitValue {
		return nil, fmt.Errorf("value is larger than 48 bits")
	}

	byteArray := make([]byte, 6)
	byteArray[0] = byte(value >> 40)
	byteArray[1] = byte(value >> 32)
	byteArray[2] = byte(value >> 24)
	byteArray[3] = byte(value >> 16)
	byteArray[4] = byte(value >> 8)
	byteArray[5] = byte(value)

	return common.LeftPadBytes(byteArray, 6), nil
}
