// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sessionKeyOwnedValidator

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// UserOperation is an auto generated low-level Go binding around an user-defined struct.
type UserOperation struct {
	Sender               common.Address
	Nonce                *big.Int
	InitCode             []byte
	CallData             []byte
	CallGasLimit         *big.Int
	VerificationGasLimit *big.Int
	PreVerificationGas   *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	PaymasterAndData     []byte
	Signature            []byte
}

// SessionKeyOwnedValidatorMetaData contains all meta data concerning the SessionKeyOwnedValidator contract.
var SessionKeyOwnedValidatorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"kernel\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnerChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"disable\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"enable\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sessionKey\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kernel\",\"type\":\"address\"}],\"name\":\"sessionKeyStorage\",\"outputs\":[{\"internalType\":\"ValidUntil\",\"name\":\"validUntil\",\"type\":\"uint48\"},{\"internalType\":\"ValidAfter\",\"name\":\"validAfter\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_caller\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"validCaller\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"validateSignature\",\"outputs\":[{\"internalType\":\"ValidationData\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"_userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"_userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validateUserOp\",\"outputs\":[{\"internalType\":\"ValidationData\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x6080806040523461001657610681908161001c8239f35b600080fdfe608060408181526004918236101561001657600080fd5b600092833560e01c9182630c9595561461020157508163333daf92146101c15781633a871cdd1461017a578382638fc925aa14610124575081639ea9bd59146100cf575063fbab6a7a1461006957600080fd5b346100cb57806003193601126100cb57610081610386565b916024356001600160a01b03818116918290036100c757839416825281602052828220908252602052205481519065ffffffffffff90818116835260301c166020820152f35b8280fd5b5080fd5b9050346100c757816003193601126100c7576100e9610386565b926024359067ffffffffffffffff821161012157509261011161011892602095369101610353565b5050610635565b90519015158152f35b80fd5b92915060203660031901126101715780359067ffffffffffffffff82116101755761015191369101610353565b601411610171573560601c82528160205280822033835260205281205580f35b5050fd5b505050fd5b82600319856060368301126101215783359167ffffffffffffffff83116100cb5761016090833603011261012157506020926101ba916024359101610438565b9051908152f35b8284346101215781600319360112610121576024359067ffffffffffffffff821161012157506020926101fa6101ba9236908301610353565b91356105b7565b8385916020918260031936011261034f5767ffffffffffffffff92853584811161034b576102329036908801610353565b93846014116103475784601a1161034757601482013560d01c94831161034757601a82013560d01c90858211156102e05750835195868501908111878210176102cd5765ffffffffffff9596979850845286528186019485523560601c8652858152818620903387525284209251166bffffffffffff0000000000008354925160301b16916bffffffffffffffffffffffff19161717905580f35b634e487b7160e01b885260418952602488fd5b62461bcd60e51b8152888101849052603760248201527f53657373696f6e4b65794f776e656456616c696461746f723a20696e76616c6960448201527f642076616c6964556e74696c2f76616c696441667465720000000000000000006064820152608490fd5b8680fd5b8580fd5b8380fd5b9181601f840112156103815782359167ffffffffffffffff8311610381576020838186019501011161038157565b600080fd5b600435906001600160a01b038216820361038157565b903590601e1981360301821215610381570180359067ffffffffffffffff82116103815760200191813603831361038157565b92919267ffffffffffffffff918281116104225760405192601f8201601f19908116603f011684019081118482101761042257604052829481845281830111610381578281602093846000960137010152565b634e487b7160e01b600052604160045260246000fd5b906101408201610448818461039c565b6001600160a01b03918291610468916104629136916103cf565b85610541565b16916000928352826020526040832033845260205260408320549365ffffffffffff9485811661051b5750946104d66104cf6104dc93604097986020527b19457468657265756d205369676e6564204d6573736167653a0a33328752603c6004209361039c565b36916103cf565b90610541565b168152806020528181203382526020522054908116156105155760a01b65ffffffffffff60a01b81169065ffffffffffff60d01b161790565b50600190565b9550505050505060a01b65ffffffffffff60a01b81169065ffffffffffff60d01b161790565b60207f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a19392608060405193838301516040526040830151928360605260416000988995869485526060840151851a885210915114165afa508251923d156105aa57606052604052565b638baa579f90526004601cfd5b906001600160a01b03806105cf6104623687866103cf565b16916000928352826020526040832033845260205260408320549365ffffffffffff9485811661051b5750946104d66104dc92604096976020527b19457468657265756d205369676e6564204d6573736167653a0a33328652603c6004209236916103cf565b6001600160a01b031660009081526020818152604080832033845290915281205465ffffffffffff90603081901c821642111561067c571642116106795750600190565b90565b50509056",
}

// SessionKeyOwnedValidatorABI is the input ABI used to generate the binding from.
// Deprecated: Use SessionKeyOwnedValidatorMetaData.ABI instead.
var SessionKeyOwnedValidatorABI = SessionKeyOwnedValidatorMetaData.ABI

// SessionKeyOwnedValidatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SessionKeyOwnedValidatorMetaData.Bin instead.
var SessionKeyOwnedValidatorBin = SessionKeyOwnedValidatorMetaData.Bin

// DeploySessionKeyOwnedValidator deploys a new Ethereum contract, binding an instance of SessionKeyOwnedValidator to it.
func DeploySessionKeyOwnedValidator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SessionKeyOwnedValidator, error) {
	parsed, err := SessionKeyOwnedValidatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SessionKeyOwnedValidatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SessionKeyOwnedValidator{SessionKeyOwnedValidatorCaller: SessionKeyOwnedValidatorCaller{contract: contract}, SessionKeyOwnedValidatorTransactor: SessionKeyOwnedValidatorTransactor{contract: contract}, SessionKeyOwnedValidatorFilterer: SessionKeyOwnedValidatorFilterer{contract: contract}}, nil
}

// SessionKeyOwnedValidator is an auto generated Go binding around an Ethereum contract.
type SessionKeyOwnedValidator struct {
	SessionKeyOwnedValidatorCaller     // Read-only binding to the contract
	SessionKeyOwnedValidatorTransactor // Write-only binding to the contract
	SessionKeyOwnedValidatorFilterer   // Log filterer for contract events
}

// SessionKeyOwnedValidatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type SessionKeyOwnedValidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SessionKeyOwnedValidatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SessionKeyOwnedValidatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SessionKeyOwnedValidatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SessionKeyOwnedValidatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SessionKeyOwnedValidatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SessionKeyOwnedValidatorSession struct {
	Contract     *SessionKeyOwnedValidator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// SessionKeyOwnedValidatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SessionKeyOwnedValidatorCallerSession struct {
	Contract *SessionKeyOwnedValidatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// SessionKeyOwnedValidatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SessionKeyOwnedValidatorTransactorSession struct {
	Contract     *SessionKeyOwnedValidatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// SessionKeyOwnedValidatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type SessionKeyOwnedValidatorRaw struct {
	Contract *SessionKeyOwnedValidator // Generic contract binding to access the raw methods on
}

// SessionKeyOwnedValidatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SessionKeyOwnedValidatorCallerRaw struct {
	Contract *SessionKeyOwnedValidatorCaller // Generic read-only contract binding to access the raw methods on
}

// SessionKeyOwnedValidatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SessionKeyOwnedValidatorTransactorRaw struct {
	Contract *SessionKeyOwnedValidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSessionKeyOwnedValidator creates a new instance of SessionKeyOwnedValidator, bound to a specific deployed contract.
func NewSessionKeyOwnedValidator(address common.Address, backend bind.ContractBackend) (*SessionKeyOwnedValidator, error) {
	contract, err := bindSessionKeyOwnedValidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SessionKeyOwnedValidator{SessionKeyOwnedValidatorCaller: SessionKeyOwnedValidatorCaller{contract: contract}, SessionKeyOwnedValidatorTransactor: SessionKeyOwnedValidatorTransactor{contract: contract}, SessionKeyOwnedValidatorFilterer: SessionKeyOwnedValidatorFilterer{contract: contract}}, nil
}

// NewSessionKeyOwnedValidatorCaller creates a new read-only instance of SessionKeyOwnedValidator, bound to a specific deployed contract.
func NewSessionKeyOwnedValidatorCaller(address common.Address, caller bind.ContractCaller) (*SessionKeyOwnedValidatorCaller, error) {
	contract, err := bindSessionKeyOwnedValidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SessionKeyOwnedValidatorCaller{contract: contract}, nil
}

// NewSessionKeyOwnedValidatorTransactor creates a new write-only instance of SessionKeyOwnedValidator, bound to a specific deployed contract.
func NewSessionKeyOwnedValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*SessionKeyOwnedValidatorTransactor, error) {
	contract, err := bindSessionKeyOwnedValidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SessionKeyOwnedValidatorTransactor{contract: contract}, nil
}

// NewSessionKeyOwnedValidatorFilterer creates a new log filterer instance of SessionKeyOwnedValidator, bound to a specific deployed contract.
func NewSessionKeyOwnedValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*SessionKeyOwnedValidatorFilterer, error) {
	contract, err := bindSessionKeyOwnedValidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SessionKeyOwnedValidatorFilterer{contract: contract}, nil
}

// bindSessionKeyOwnedValidator binds a generic wrapper to an already deployed contract.
func bindSessionKeyOwnedValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SessionKeyOwnedValidatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SessionKeyOwnedValidator.Contract.SessionKeyOwnedValidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.Contract.SessionKeyOwnedValidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.Contract.SessionKeyOwnedValidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SessionKeyOwnedValidator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.Contract.contract.Transact(opts, method, params...)
}

// SessionKeyStorage is a free data retrieval call binding the contract method 0xfbab6a7a.
//
// Solidity: function sessionKeyStorage(address sessionKey, address kernel) view returns(uint48 validUntil, uint48 validAfter)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorCaller) SessionKeyStorage(opts *bind.CallOpts, sessionKey common.Address, kernel common.Address) (struct {
	ValidUntil *big.Int
	ValidAfter *big.Int
}, error) {
	var out []interface{}
	err := _SessionKeyOwnedValidator.contract.Call(opts, &out, "sessionKeyStorage", sessionKey, kernel)

	outstruct := new(struct {
		ValidUntil *big.Int
		ValidAfter *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ValidUntil = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ValidAfter = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// SessionKeyStorage is a free data retrieval call binding the contract method 0xfbab6a7a.
//
// Solidity: function sessionKeyStorage(address sessionKey, address kernel) view returns(uint48 validUntil, uint48 validAfter)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorSession) SessionKeyStorage(sessionKey common.Address, kernel common.Address) (struct {
	ValidUntil *big.Int
	ValidAfter *big.Int
}, error) {
	return _SessionKeyOwnedValidator.Contract.SessionKeyStorage(&_SessionKeyOwnedValidator.CallOpts, sessionKey, kernel)
}

// SessionKeyStorage is a free data retrieval call binding the contract method 0xfbab6a7a.
//
// Solidity: function sessionKeyStorage(address sessionKey, address kernel) view returns(uint48 validUntil, uint48 validAfter)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorCallerSession) SessionKeyStorage(sessionKey common.Address, kernel common.Address) (struct {
	ValidUntil *big.Int
	ValidAfter *big.Int
}, error) {
	return _SessionKeyOwnedValidator.Contract.SessionKeyStorage(&_SessionKeyOwnedValidator.CallOpts, sessionKey, kernel)
}

// ValidCaller is a free data retrieval call binding the contract method 0x9ea9bd59.
//
// Solidity: function validCaller(address _caller, bytes ) view returns(bool)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorCaller) ValidCaller(opts *bind.CallOpts, _caller common.Address, arg1 []byte) (bool, error) {
	var out []interface{}
	err := _SessionKeyOwnedValidator.contract.Call(opts, &out, "validCaller", _caller, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidCaller is a free data retrieval call binding the contract method 0x9ea9bd59.
//
// Solidity: function validCaller(address _caller, bytes ) view returns(bool)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorSession) ValidCaller(_caller common.Address, arg1 []byte) (bool, error) {
	return _SessionKeyOwnedValidator.Contract.ValidCaller(&_SessionKeyOwnedValidator.CallOpts, _caller, arg1)
}

// ValidCaller is a free data retrieval call binding the contract method 0x9ea9bd59.
//
// Solidity: function validCaller(address _caller, bytes ) view returns(bool)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorCallerSession) ValidCaller(_caller common.Address, arg1 []byte) (bool, error) {
	return _SessionKeyOwnedValidator.Contract.ValidCaller(&_SessionKeyOwnedValidator.CallOpts, _caller, arg1)
}

// ValidateSignature is a free data retrieval call binding the contract method 0x333daf92.
//
// Solidity: function validateSignature(bytes32 hash, bytes signature) view returns(uint256)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorCaller) ValidateSignature(opts *bind.CallOpts, hash [32]byte, signature []byte) (*big.Int, error) {
	var out []interface{}
	err := _SessionKeyOwnedValidator.contract.Call(opts, &out, "validateSignature", hash, signature)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidateSignature is a free data retrieval call binding the contract method 0x333daf92.
//
// Solidity: function validateSignature(bytes32 hash, bytes signature) view returns(uint256)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorSession) ValidateSignature(hash [32]byte, signature []byte) (*big.Int, error) {
	return _SessionKeyOwnedValidator.Contract.ValidateSignature(&_SessionKeyOwnedValidator.CallOpts, hash, signature)
}

// ValidateSignature is a free data retrieval call binding the contract method 0x333daf92.
//
// Solidity: function validateSignature(bytes32 hash, bytes signature) view returns(uint256)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorCallerSession) ValidateSignature(hash [32]byte, signature []byte) (*big.Int, error) {
	return _SessionKeyOwnedValidator.Contract.ValidateSignature(&_SessionKeyOwnedValidator.CallOpts, hash, signature)
}

// Disable is a paid mutator transaction binding the contract method 0x8fc925aa.
//
// Solidity: function disable(bytes _data) payable returns()
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorTransactor) Disable(opts *bind.TransactOpts, _data []byte) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.contract.Transact(opts, "disable", _data)
}

// Disable is a paid mutator transaction binding the contract method 0x8fc925aa.
//
// Solidity: function disable(bytes _data) payable returns()
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorSession) Disable(_data []byte) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.Contract.Disable(&_SessionKeyOwnedValidator.TransactOpts, _data)
}

// Disable is a paid mutator transaction binding the contract method 0x8fc925aa.
//
// Solidity: function disable(bytes _data) payable returns()
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorTransactorSession) Disable(_data []byte) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.Contract.Disable(&_SessionKeyOwnedValidator.TransactOpts, _data)
}

// Enable is a paid mutator transaction binding the contract method 0x0c959556.
//
// Solidity: function enable(bytes _data) payable returns()
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorTransactor) Enable(opts *bind.TransactOpts, _data []byte) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.contract.Transact(opts, "enable", _data)
}

// Enable is a paid mutator transaction binding the contract method 0x0c959556.
//
// Solidity: function enable(bytes _data) payable returns()
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorSession) Enable(_data []byte) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.Contract.Enable(&_SessionKeyOwnedValidator.TransactOpts, _data)
}

// Enable is a paid mutator transaction binding the contract method 0x0c959556.
//
// Solidity: function enable(bytes _data) payable returns()
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorTransactorSession) Enable(_data []byte) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.Contract.Enable(&_SessionKeyOwnedValidator.TransactOpts, _data)
}

// ValidateUserOp is a paid mutator transaction binding the contract method 0x3a871cdd.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes) _userOp, bytes32 _userOpHash, uint256 ) payable returns(uint256 validationData)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorTransactor) ValidateUserOp(opts *bind.TransactOpts, _userOp UserOperation, _userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.contract.Transact(opts, "validateUserOp", _userOp, _userOpHash, arg2)
}

// ValidateUserOp is a paid mutator transaction binding the contract method 0x3a871cdd.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes) _userOp, bytes32 _userOpHash, uint256 ) payable returns(uint256 validationData)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorSession) ValidateUserOp(_userOp UserOperation, _userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.Contract.ValidateUserOp(&_SessionKeyOwnedValidator.TransactOpts, _userOp, _userOpHash, arg2)
}

// ValidateUserOp is a paid mutator transaction binding the contract method 0x3a871cdd.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes) _userOp, bytes32 _userOpHash, uint256 ) payable returns(uint256 validationData)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorTransactorSession) ValidateUserOp(_userOp UserOperation, _userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error) {
	return _SessionKeyOwnedValidator.Contract.ValidateUserOp(&_SessionKeyOwnedValidator.TransactOpts, _userOp, _userOpHash, arg2)
}

// SessionKeyOwnedValidatorOwnerChangedIterator is returned from FilterOwnerChanged and is used to iterate over the raw logs and unpacked data for OwnerChanged events raised by the SessionKeyOwnedValidator contract.
type SessionKeyOwnedValidatorOwnerChangedIterator struct {
	Event *SessionKeyOwnedValidatorOwnerChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SessionKeyOwnedValidatorOwnerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SessionKeyOwnedValidatorOwnerChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SessionKeyOwnedValidatorOwnerChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SessionKeyOwnedValidatorOwnerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SessionKeyOwnedValidatorOwnerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SessionKeyOwnedValidatorOwnerChanged represents a OwnerChanged event raised by the SessionKeyOwnedValidator contract.
type SessionKeyOwnedValidatorOwnerChanged struct {
	Kernel   common.Address
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnerChanged is a free log retrieval operation binding the contract event 0x381c0d11398486654573703c51ee8210ce9461764d133f9f0e53b6a539705331.
//
// Solidity: event OwnerChanged(address indexed kernel, address indexed oldOwner, address indexed newOwner)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorFilterer) FilterOwnerChanged(opts *bind.FilterOpts, kernel []common.Address, oldOwner []common.Address, newOwner []common.Address) (*SessionKeyOwnedValidatorOwnerChangedIterator, error) {

	var kernelRule []interface{}
	for _, kernelItem := range kernel {
		kernelRule = append(kernelRule, kernelItem)
	}
	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SessionKeyOwnedValidator.contract.FilterLogs(opts, "OwnerChanged", kernelRule, oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SessionKeyOwnedValidatorOwnerChangedIterator{contract: _SessionKeyOwnedValidator.contract, event: "OwnerChanged", logs: logs, sub: sub}, nil
}

// WatchOwnerChanged is a free log subscription operation binding the contract event 0x381c0d11398486654573703c51ee8210ce9461764d133f9f0e53b6a539705331.
//
// Solidity: event OwnerChanged(address indexed kernel, address indexed oldOwner, address indexed newOwner)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorFilterer) WatchOwnerChanged(opts *bind.WatchOpts, sink chan<- *SessionKeyOwnedValidatorOwnerChanged, kernel []common.Address, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var kernelRule []interface{}
	for _, kernelItem := range kernel {
		kernelRule = append(kernelRule, kernelItem)
	}
	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SessionKeyOwnedValidator.contract.WatchLogs(opts, "OwnerChanged", kernelRule, oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SessionKeyOwnedValidatorOwnerChanged)
				if err := _SessionKeyOwnedValidator.contract.UnpackLog(event, "OwnerChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnerChanged is a log parse operation binding the contract event 0x381c0d11398486654573703c51ee8210ce9461764d133f9f0e53b6a539705331.
//
// Solidity: event OwnerChanged(address indexed kernel, address indexed oldOwner, address indexed newOwner)
func (_SessionKeyOwnedValidator *SessionKeyOwnedValidatorFilterer) ParseOwnerChanged(log types.Log) (*SessionKeyOwnedValidatorOwnerChanged, error) {
	event := new(SessionKeyOwnedValidatorOwnerChanged)
	if err := _SessionKeyOwnedValidator.contract.UnpackLog(event, "OwnerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
