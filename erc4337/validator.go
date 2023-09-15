package erc4337

import (
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	KernelStorage "github.com/helicarrierstudio/silver-arrow/abi/kernelStorage"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
)

/*
*
- The enableData will be appended to the userop signature
- signature[56:88] contains the enableData
*/
type SessionKeyOwnedValidator struct {
	privatekey       string
	Mode             []byte
	Chain            *big.Int
	ValidUntil       []byte
	ValidAfter       []byte
	SessionKey       common.Address
	ExecutorAddress  common.Address
	ValidatorAddress common.Address
}

func InitialiseValidator(validatorAddress, sessionKey, privKey, mode string, chainId int64) (*SessionKeyOwnedValidator, error) {
	validator := common.HexToAddress(validatorAddress)
	executor := common.HexToAddress("0x")
	session := common.HexToAddress(sessionKey)
	md, err := hexutil.Decode(mode)
	if err != nil {
		return nil, err
	}

	return newSessionKeyOwnedValidator(validator, executor, session, md, privKey, big.NewInt(chainId)), nil
}

func newSessionKeyOwnedValidator(validator, executor, sessionKey common.Address, mode []byte, privateKey string, chain *big.Int) *SessionKeyOwnedValidator {
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
		privatekey:       privateKey[2:],
	}
}

// func (v *SessionKeyOwnedValidator) Sign(op map[string]any) (sig, hash []byte, err error) {
// 	entrypoint := GetEntryPointAddress()

// 	// userOpSignature := []byte{}

// 	operation, err := userop.New(op)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	opHash := operation.GetUserOpHash(entrypoint, v.Chain)
// 	fmt.Println("opHash", opHash)
// 	hash = opHash.Bytes()

// 	pk, err := crypto.HexToECDSA(v.privatekey)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	sig, err = crypto.Sign(ecrecover.ToEthSignedMessageHash(hash), pk)
// 	if err != nil {
// 		err = errors.Wrap(err, "generating signature failed.")
// 		return nil, nil, err
// 	}
// 	sig[64] += 27

// 	// userOpSignature = append(userOpSignature, v.Mode...)
// 	// userOpSignature = append(userOpSignature, sig...)

// 	return sig, ecrecover.ToEthSignedMessageHash(hash), nil
// }

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

func (v *SessionKeyOwnedValidator) CreateEnableDigest(sender common.Address, sig, enableData []byte) ([]byte, error) {
	digest := getTypedDataHash(sender, v.ValidatorAddress, v.ExecutorAddress, sig, v.ValidUntil, v.ValidAfter, enableData, v.Chain)

	return digest, nil
}

// func (v *SessionKeyOwnedValidator) GetEnableSignature(signature, enbaleData []byte) []byte {
// 	userOpSignature := []byte{}

// 	enableDataLength := uint256.MustFromBig(big.NewInt(int64(len(enbaleData))))
// 	signatureLength := uint256.MustFromBig(big.NewInt(65))

// 	userOpSignature = append(userOpSignature, v.Mode...)
// 	userOpSignature = append(userOpSignature, v.getValidatorData()...)
// 	userOpSignature = append(userOpSignature, enableDataLength.PaddedBytes(32)...)
// 	userOpSignature = append(userOpSignature, enbaleData...)
// 	userOpSignature = append(userOpSignature, signatureLength.PaddedBytes(32)...)
// 	userOpSignature = append(userOpSignature, signature...)

// 	// fmt.Println("signature ", userOpSignature)
// 	return userOpSignature
// }

func (v *SessionKeyOwnedValidator) PackEnableModeSignature(enableSignature, validatorSignature, enableData []byte) ([]byte, error) {
	var userOpSignature []byte

	// validate signature length
	enableSigLength := len(enableSignature)
	if enableSigLength < 65 {
		return nil, fmt.Errorf("invalid enableSignature length of %v, should be %v", enableSigLength, 65)
	}
	enableDataLength := uint256.MustFromBig(big.NewInt(int64(len(enableData))))
	enableSignatureLength := uint256.MustFromBig(big.NewInt(int64(enableSigLength)))

	userOpSignature = append(userOpSignature, v.Mode...)
	userOpSignature = append(userOpSignature, v.getValidatorData()...)
	userOpSignature = append(userOpSignature, enableDataLength.PaddedBytes(32)...)
	userOpSignature = append(userOpSignature, enableData...)
	userOpSignature = append(userOpSignature, enableSignatureLength.PaddedBytes(32)...)
	userOpSignature = append(userOpSignature, enableSignature...)
	userOpSignature = append(userOpSignature, validatorSignature...)

	return userOpSignature, nil
}

func (v *SessionKeyOwnedValidator) getValidatorData() []byte {
	var data []byte
	// v.ValidAfter // 4:10
	// v.ValidUntil //10:16
	// v.ValidatorAddress //16:36
	data = append(data, v.ValidUntil...)
	data = append(data, v.ValidAfter...)
	data = append(data, v.ValidatorAddress.Bytes()...)
	data = append(data, v.ExecutorAddress.Bytes()...)
	return data
}

func (v *SessionKeyOwnedValidator) GetEnableData() ([]byte, error) {
	// var validAfter, validUntil []bytes{}
	var data []byte
	fmt.Println("valid until ", v.ValidUntil)
	fmt.Println("valid after ", v.ValidAfter)
	data = append(data, v.SessionKey.Bytes()...)
	data = append(data, v.ValidUntil...)
	data = append(data, v.ValidAfter...)
	return data, nil
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

func getTypedDataHash(sender, validator, executor common.Address, sig, validUntil, validaAfter, enableData []byte, chainId *big.Int) []byte {
	data := []byte("\x19\x01")

	// domain := apitypes.TypedDataDomain{
	// 	Name:    "Kernel",
	// 	Version: "0.2.1",
	// 	ChainId: (*mth.HexOrDecimal256)(chainId),
	// 	VerifyingContract: sender.Hex(),
	// }
	// message := apitypes.TypedDataMessage{
	// 	"": "",

	// }

	ds := buildDomainSeparator("Kernel", "0.2.1", sender, chainId)
	structHash := getStructHash(sig, validUntil, validaAfter, enableData, validator, executor)
	// data = append(data, ds...)
	// data = append(data, getStructHash(sig, validUntil, validaAfter, enableData, validator, executor)...)
	// crypto.Keccak256()
	return crypto.Keccak256(data, ds, structHash)
}

func buildDomainSeparator(name, version string, verifyingContract common.Address, chainId *big.Int) []byte {

	hashedName := crypto.Keccak256([]byte(name))
	hashedVersion := crypto.Keccak256([]byte(version))
	hashedType := crypto.Keccak256([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))

	return crypto.Keccak256(encodeBytes(hashedType, hashedName, hashedVersion, chainId.Bytes(), verifyingContract[:]))
}

func encodeBytes(input ...[]byte) []byte {
	var data []byte
	for _, v := range input {
		data = append(data, common.LeftPadBytes(v, 32)...)
	}
	return data
}

func getStructHash(sig, validUntil, validAfter, enableData []byte, validator, executor common.Address) []byte {
	a := crypto.Keccak256([]byte("ValidatorApproved(bytes4 sig,uint256 validatorData,address executor,bytes enableData)"))
	vuntil := new(big.Int).SetBytes(validUntil)
	vafter := new(big.Int).SetBytes(validAfter)
	va := new(big.Int).Lsh(vafter, 160)
	vu := new(big.Int).Lsh(vuntil, (160 + 48))

	validatorData := new(big.Int).Or(validator.Big(), va)
	validatorData = validatorData.Or(validatorData, vu)

	return crypto.Keccak256(encodeBytes(a, sig, validatorData.Bytes(), executor[:], crypto.Keccak256(enableData)))
}

func (v *SessionKeyOwnedValidator) SetExecution(enableData []byte, ownerAccount string) ([]byte, error) {
	kernelAbi, err := KernelStorage.KernelStorageMetaData.GetAbi()
	// abiStr := getKernelStorageAbi()
	// kernelAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	selector := [4]byte{}
	sel, err := hexutil.Decode("0x84189294")
	if err != nil {
		err = errors.Wrap(err, "invalid selector hex")
		return nil, err
	}
	copy(selector[:], sel)
	executorAddress := common.HexToAddress("0x2087C7FfD0d0DAE80a00EE74325aBF3449e0eaf1")
	validatorAddress := v.ValidatorAddress
	validUntil := big.NewInt(99999999999)
	validAfter := big.NewInt(0)

	callData, err := kernelAbi.Pack("setExecution", selector, executorAddress, validatorAddress, validUntil, validAfter, enableData)
	if err != nil {
		return nil, err
	}
	fmt.Println("callData", hexutil.Encode(callData))
	// log.Fatalln("exit")

	return callData, nil
}
