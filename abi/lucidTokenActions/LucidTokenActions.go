// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package lucidTokenActions

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

// TransferInput is an auto generated low-level Go binding around an user-defined struct.
type TransferInput struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
}

// LucidTokenActionsMetaData contains all meta data concerning the LucidTokenActions contract.
var LucidTokenActionsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structTransferInput[]\",\"name\":\"inputs\",\"type\":\"tuple[]\"}],\"name\":\"batchTransfer20Action\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transfer20Action\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferERC1155Action\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferERC721Action\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080806040523461001657610466908161001c8239f35b600080fdfe608060408181526004908136101561001657600080fd5b600092833560e01c908482631888bfd71461032a57505080631d06ac10146101ef5780637ca237e6146100f25763841892941461005257600080fd5b346100ee57602061009d928461006736610394565b865163a9059cbb60e01b81526001600160a01b039091169481019485526020850191909152958693919284929091839160400190565b03926001600160a01b03165af19081156100e557506100ba575080f35b6100da9060203d81116100de575b6100d281836103f2565b810190610414565b5080f35b503d6100c8565b513d84823e3d90fd5b8280fd5b5090346100ee5760a03660031901126100ee5780356001600160a01b0381811693918490036101ea576044359081168091036101ea5784936084359067ffffffffffffffff928383116101e657366023840112156101e657828601359384116101e65736602485850101116101e657813b156101e6578660c460249786839789519a8b9889978895637921219560e11b875230908701528286015281356044860152606435606486015260a060848601528260a486015201848401378181018301849052601f01601f191681010301925af19081156100e557506101d35750f35b6101dc906103c8565b6101e35780f35b80fd5b8680fd5b600080fd5b5090346100ee57602090816003193601126103265780359167ffffffffffffffff93848411610322573660238501121561032257838301359485116103225760249485850194863691606084020101116101e657865b818110610250578780f35b6102bf846102698161026385878c61042c565b01610452565b61027c61027785878c61042c565b610452565b8661028886888d61042c565b885163a9059cbb60e01b81526001600160a01b03909416848c0190815291013560208201529193849283918e918391604090910190565b03926001600160a01b03165af18015610318576102fb575b5060001981146102e957600101610245565b634e487b7160e01b8852601185528688fd5b61031190853d87116100de576100d281836103f2565b50386102d7565b84513d8b823e3d90fd5b8580fd5b8380fd5b915092346103905761033b36610394565b90926001600160a01b03928316803b15610322578794606494869488946323b872dd60e01b8752309087015216602485015260448401525af19081156100e55750610384575080f35b61038d906103c8565b80f35b5080fd5b60609060031901126101ea576001600160a01b0360043581811681036101ea57916024359160443590811681036101ea5790565b67ffffffffffffffff81116103dc57604052565b634e487b7160e01b600052604160045260246000fd5b90601f8019910116810190811067ffffffffffffffff8211176103dc57604052565b908160209103126101ea575180151581036101ea5790565b919081101561043c576060020190565b634e487b7160e01b600052603260045260246000fd5b356001600160a01b03811681036101ea579056",
}

// LucidTokenActionsABI is the input ABI used to generate the binding from.
// Deprecated: Use LucidTokenActionsMetaData.ABI instead.
var LucidTokenActionsABI = LucidTokenActionsMetaData.ABI

// LucidTokenActionsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LucidTokenActionsMetaData.Bin instead.
var LucidTokenActionsBin = LucidTokenActionsMetaData.Bin

// DeployLucidTokenActions deploys a new Ethereum contract, binding an instance of LucidTokenActions to it.
func DeployLucidTokenActions(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LucidTokenActions, error) {
	parsed, err := LucidTokenActionsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LucidTokenActionsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LucidTokenActions{LucidTokenActionsCaller: LucidTokenActionsCaller{contract: contract}, LucidTokenActionsTransactor: LucidTokenActionsTransactor{contract: contract}, LucidTokenActionsFilterer: LucidTokenActionsFilterer{contract: contract}}, nil
}

// LucidTokenActions is an auto generated Go binding around an Ethereum contract.
type LucidTokenActions struct {
	LucidTokenActionsCaller     // Read-only binding to the contract
	LucidTokenActionsTransactor // Write-only binding to the contract
	LucidTokenActionsFilterer   // Log filterer for contract events
}

// LucidTokenActionsCaller is an auto generated read-only Go binding around an Ethereum contract.
type LucidTokenActionsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LucidTokenActionsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LucidTokenActionsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LucidTokenActionsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LucidTokenActionsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LucidTokenActionsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LucidTokenActionsSession struct {
	Contract     *LucidTokenActions // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// LucidTokenActionsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LucidTokenActionsCallerSession struct {
	Contract *LucidTokenActionsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// LucidTokenActionsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LucidTokenActionsTransactorSession struct {
	Contract     *LucidTokenActionsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// LucidTokenActionsRaw is an auto generated low-level Go binding around an Ethereum contract.
type LucidTokenActionsRaw struct {
	Contract *LucidTokenActions // Generic contract binding to access the raw methods on
}

// LucidTokenActionsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LucidTokenActionsCallerRaw struct {
	Contract *LucidTokenActionsCaller // Generic read-only contract binding to access the raw methods on
}

// LucidTokenActionsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LucidTokenActionsTransactorRaw struct {
	Contract *LucidTokenActionsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLucidTokenActions creates a new instance of LucidTokenActions, bound to a specific deployed contract.
func NewLucidTokenActions(address common.Address, backend bind.ContractBackend) (*LucidTokenActions, error) {
	contract, err := bindLucidTokenActions(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LucidTokenActions{LucidTokenActionsCaller: LucidTokenActionsCaller{contract: contract}, LucidTokenActionsTransactor: LucidTokenActionsTransactor{contract: contract}, LucidTokenActionsFilterer: LucidTokenActionsFilterer{contract: contract}}, nil
}

// NewLucidTokenActionsCaller creates a new read-only instance of LucidTokenActions, bound to a specific deployed contract.
func NewLucidTokenActionsCaller(address common.Address, caller bind.ContractCaller) (*LucidTokenActionsCaller, error) {
	contract, err := bindLucidTokenActions(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LucidTokenActionsCaller{contract: contract}, nil
}

// NewLucidTokenActionsTransactor creates a new write-only instance of LucidTokenActions, bound to a specific deployed contract.
func NewLucidTokenActionsTransactor(address common.Address, transactor bind.ContractTransactor) (*LucidTokenActionsTransactor, error) {
	contract, err := bindLucidTokenActions(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LucidTokenActionsTransactor{contract: contract}, nil
}

// NewLucidTokenActionsFilterer creates a new log filterer instance of LucidTokenActions, bound to a specific deployed contract.
func NewLucidTokenActionsFilterer(address common.Address, filterer bind.ContractFilterer) (*LucidTokenActionsFilterer, error) {
	contract, err := bindLucidTokenActions(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LucidTokenActionsFilterer{contract: contract}, nil
}

// bindLucidTokenActions binds a generic wrapper to an already deployed contract.
func bindLucidTokenActions(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LucidTokenActionsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LucidTokenActions *LucidTokenActionsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LucidTokenActions.Contract.LucidTokenActionsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LucidTokenActions *LucidTokenActionsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.LucidTokenActionsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LucidTokenActions *LucidTokenActionsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.LucidTokenActionsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LucidTokenActions *LucidTokenActionsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LucidTokenActions.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LucidTokenActions *LucidTokenActionsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LucidTokenActions *LucidTokenActionsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.contract.Transact(opts, method, params...)
}

// BatchTransfer20Action is a paid mutator transaction binding the contract method 0x1d06ac10.
//
// Solidity: function batchTransfer20Action((address,address,uint256)[] inputs) returns()
func (_LucidTokenActions *LucidTokenActionsTransactor) BatchTransfer20Action(opts *bind.TransactOpts, inputs []TransferInput) (*types.Transaction, error) {
	return _LucidTokenActions.contract.Transact(opts, "batchTransfer20Action", inputs)
}

// BatchTransfer20Action is a paid mutator transaction binding the contract method 0x1d06ac10.
//
// Solidity: function batchTransfer20Action((address,address,uint256)[] inputs) returns()
func (_LucidTokenActions *LucidTokenActionsSession) BatchTransfer20Action(inputs []TransferInput) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.BatchTransfer20Action(&_LucidTokenActions.TransactOpts, inputs)
}

// BatchTransfer20Action is a paid mutator transaction binding the contract method 0x1d06ac10.
//
// Solidity: function batchTransfer20Action((address,address,uint256)[] inputs) returns()
func (_LucidTokenActions *LucidTokenActionsTransactorSession) BatchTransfer20Action(inputs []TransferInput) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.BatchTransfer20Action(&_LucidTokenActions.TransactOpts, inputs)
}

// Transfer20Action is a paid mutator transaction binding the contract method 0x84189294.
//
// Solidity: function transfer20Action(address _token, uint256 _amount, address _to) returns()
func (_LucidTokenActions *LucidTokenActionsTransactor) Transfer20Action(opts *bind.TransactOpts, _token common.Address, _amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _LucidTokenActions.contract.Transact(opts, "transfer20Action", _token, _amount, _to)
}

// Transfer20Action is a paid mutator transaction binding the contract method 0x84189294.
//
// Solidity: function transfer20Action(address _token, uint256 _amount, address _to) returns()
func (_LucidTokenActions *LucidTokenActionsSession) Transfer20Action(_token common.Address, _amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.Transfer20Action(&_LucidTokenActions.TransactOpts, _token, _amount, _to)
}

// Transfer20Action is a paid mutator transaction binding the contract method 0x84189294.
//
// Solidity: function transfer20Action(address _token, uint256 _amount, address _to) returns()
func (_LucidTokenActions *LucidTokenActionsTransactorSession) Transfer20Action(_token common.Address, _amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.Transfer20Action(&_LucidTokenActions.TransactOpts, _token, _amount, _to)
}

// TransferERC1155Action is a paid mutator transaction binding the contract method 0x7ca237e6.
//
// Solidity: function transferERC1155Action(address _token, uint256 _id, address _to, uint256 amount, bytes data) returns()
func (_LucidTokenActions *LucidTokenActionsTransactor) TransferERC1155Action(opts *bind.TransactOpts, _token common.Address, _id *big.Int, _to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LucidTokenActions.contract.Transact(opts, "transferERC1155Action", _token, _id, _to, amount, data)
}

// TransferERC1155Action is a paid mutator transaction binding the contract method 0x7ca237e6.
//
// Solidity: function transferERC1155Action(address _token, uint256 _id, address _to, uint256 amount, bytes data) returns()
func (_LucidTokenActions *LucidTokenActionsSession) TransferERC1155Action(_token common.Address, _id *big.Int, _to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.TransferERC1155Action(&_LucidTokenActions.TransactOpts, _token, _id, _to, amount, data)
}

// TransferERC1155Action is a paid mutator transaction binding the contract method 0x7ca237e6.
//
// Solidity: function transferERC1155Action(address _token, uint256 _id, address _to, uint256 amount, bytes data) returns()
func (_LucidTokenActions *LucidTokenActionsTransactorSession) TransferERC1155Action(_token common.Address, _id *big.Int, _to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.TransferERC1155Action(&_LucidTokenActions.TransactOpts, _token, _id, _to, amount, data)
}

// TransferERC721Action is a paid mutator transaction binding the contract method 0x1888bfd7.
//
// Solidity: function transferERC721Action(address _token, uint256 _id, address _to) returns()
func (_LucidTokenActions *LucidTokenActionsTransactor) TransferERC721Action(opts *bind.TransactOpts, _token common.Address, _id *big.Int, _to common.Address) (*types.Transaction, error) {
	return _LucidTokenActions.contract.Transact(opts, "transferERC721Action", _token, _id, _to)
}

// TransferERC721Action is a paid mutator transaction binding the contract method 0x1888bfd7.
//
// Solidity: function transferERC721Action(address _token, uint256 _id, address _to) returns()
func (_LucidTokenActions *LucidTokenActionsSession) TransferERC721Action(_token common.Address, _id *big.Int, _to common.Address) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.TransferERC721Action(&_LucidTokenActions.TransactOpts, _token, _id, _to)
}

// TransferERC721Action is a paid mutator transaction binding the contract method 0x1888bfd7.
//
// Solidity: function transferERC721Action(address _token, uint256 _id, address _to) returns()
func (_LucidTokenActions *LucidTokenActionsTransactorSession) TransferERC721Action(_token common.Address, _id *big.Int, _to common.Address) (*types.Transaction, error) {
	return _LucidTokenActions.Contract.TransferERC721Action(&_LucidTokenActions.TransactOpts, _token, _id, _to)
}
