// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tokenActions

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

// TokenActionsMetaData contains all meta data concerning the TokenActions contract.
var TokenActionsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transfer20Action\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferERC1155Action\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferERC721Action\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60808060405234610016576102dd908161001c8239f35b600080fdfe608060408181526004908136101561001657600080fd5b600092833560e01c908482631888bfd71461021557505080637ca237e6146101125763841892941461004757600080fd5b3461010e57604491602061005a3661027f565b855163a9059cbb60e01b81526001600160a01b0391821686820152602481019290925290958692909183918991165af1801561010457610098578380f35b6020913d83116100fc575b601f8301601f191684019167ffffffffffffffff8311858410176100e75750526020908201829003126100e35751801515036100e0573880808380f35b80fd5b5080fd5b604190634e487b7160e01b6000525260246000fd5b3d92506100a3565b82513d86823e3d90fd5b8280fd5b50903461010e5760a036600319011261010e5780356001600160a01b038181169391849003610210576044359081168091036102105784936084359067ffffffffffffffff9283831161020c573660238401121561020c578286013593841161020c57366024858501011161020c57813b1561020c578660c460249786839789519a8b9889978895637921219560e11b875230908701528286015281356044860152606435606486015260a060848601528260a486015201848401378181018301849052601f01601f191681010301925af190811561020357506101f35750f35b6101fc906102b3565b6100e05780f35b513d84823e3d90fd5b8680fd5b600080fd5b915092346100e3576102263661027f565b90926001600160a01b03928316803b1561027b578794606494869488946323b872dd60e01b8752309087015216602485015260448401525af1908115610203575061026f575080f35b610278906102b3565b80f35b8580fd5b6060906003190112610210576001600160a01b03600435818116810361021057916024359160443590811681036102105790565b67ffffffffffffffff81116102c757604052565b634e487b7160e01b600052604160045260246000fd",
}

// TokenActionsABI is the input ABI used to generate the binding from.
// Deprecated: Use TokenActionsMetaData.ABI instead.
var TokenActionsABI = TokenActionsMetaData.ABI

// TokenActionsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TokenActionsMetaData.Bin instead.
var TokenActionsBin = TokenActionsMetaData.Bin

// DeployTokenActions deploys a new Ethereum contract, binding an instance of TokenActions to it.
func DeployTokenActions(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TokenActions, error) {
	parsed, err := TokenActionsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TokenActionsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenActions{TokenActionsCaller: TokenActionsCaller{contract: contract}, TokenActionsTransactor: TokenActionsTransactor{contract: contract}, TokenActionsFilterer: TokenActionsFilterer{contract: contract}}, nil
}

// TokenActions is an auto generated Go binding around an Ethereum contract.
type TokenActions struct {
	TokenActionsCaller     // Read-only binding to the contract
	TokenActionsTransactor // Write-only binding to the contract
	TokenActionsFilterer   // Log filterer for contract events
}

// TokenActionsCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenActionsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenActionsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenActionsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenActionsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenActionsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenActionsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenActionsSession struct {
	Contract     *TokenActions     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenActionsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenActionsCallerSession struct {
	Contract *TokenActionsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TokenActionsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenActionsTransactorSession struct {
	Contract     *TokenActionsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TokenActionsRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenActionsRaw struct {
	Contract *TokenActions // Generic contract binding to access the raw methods on
}

// TokenActionsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenActionsCallerRaw struct {
	Contract *TokenActionsCaller // Generic read-only contract binding to access the raw methods on
}

// TokenActionsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenActionsTransactorRaw struct {
	Contract *TokenActionsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenActions creates a new instance of TokenActions, bound to a specific deployed contract.
func NewTokenActions(address common.Address, backend bind.ContractBackend) (*TokenActions, error) {
	contract, err := bindTokenActions(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenActions{TokenActionsCaller: TokenActionsCaller{contract: contract}, TokenActionsTransactor: TokenActionsTransactor{contract: contract}, TokenActionsFilterer: TokenActionsFilterer{contract: contract}}, nil
}

// NewTokenActionsCaller creates a new read-only instance of TokenActions, bound to a specific deployed contract.
func NewTokenActionsCaller(address common.Address, caller bind.ContractCaller) (*TokenActionsCaller, error) {
	contract, err := bindTokenActions(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenActionsCaller{contract: contract}, nil
}

// NewTokenActionsTransactor creates a new write-only instance of TokenActions, bound to a specific deployed contract.
func NewTokenActionsTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenActionsTransactor, error) {
	contract, err := bindTokenActions(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenActionsTransactor{contract: contract}, nil
}

// NewTokenActionsFilterer creates a new log filterer instance of TokenActions, bound to a specific deployed contract.
func NewTokenActionsFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenActionsFilterer, error) {
	contract, err := bindTokenActions(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenActionsFilterer{contract: contract}, nil
}

// bindTokenActions binds a generic wrapper to an already deployed contract.
func bindTokenActions(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TokenActionsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenActions *TokenActionsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenActions.Contract.TokenActionsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenActions *TokenActionsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenActions.Contract.TokenActionsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenActions *TokenActionsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenActions.Contract.TokenActionsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenActions *TokenActionsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenActions.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenActions *TokenActionsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenActions.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenActions *TokenActionsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenActions.Contract.contract.Transact(opts, method, params...)
}

// Transfer20Action is a paid mutator transaction binding the contract method 0x84189294.
//
// Solidity: function transfer20Action(address _token, uint256 _amount, address _to) returns()
func (_TokenActions *TokenActionsTransactor) Transfer20Action(opts *bind.TransactOpts, _token common.Address, _amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _TokenActions.contract.Transact(opts, "transfer20Action", _token, _amount, _to)
}

// Transfer20Action is a paid mutator transaction binding the contract method 0x84189294.
//
// Solidity: function transfer20Action(address _token, uint256 _amount, address _to) returns()
func (_TokenActions *TokenActionsSession) Transfer20Action(_token common.Address, _amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _TokenActions.Contract.Transfer20Action(&_TokenActions.TransactOpts, _token, _amount, _to)
}

// Transfer20Action is a paid mutator transaction binding the contract method 0x84189294.
//
// Solidity: function transfer20Action(address _token, uint256 _amount, address _to) returns()
func (_TokenActions *TokenActionsTransactorSession) Transfer20Action(_token common.Address, _amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _TokenActions.Contract.Transfer20Action(&_TokenActions.TransactOpts, _token, _amount, _to)
}

// TransferERC1155Action is a paid mutator transaction binding the contract method 0x7ca237e6.
//
// Solidity: function transferERC1155Action(address _token, uint256 _id, address _to, uint256 amount, bytes data) returns()
func (_TokenActions *TokenActionsTransactor) TransferERC1155Action(opts *bind.TransactOpts, _token common.Address, _id *big.Int, _to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _TokenActions.contract.Transact(opts, "transferERC1155Action", _token, _id, _to, amount, data)
}

// TransferERC1155Action is a paid mutator transaction binding the contract method 0x7ca237e6.
//
// Solidity: function transferERC1155Action(address _token, uint256 _id, address _to, uint256 amount, bytes data) returns()
func (_TokenActions *TokenActionsSession) TransferERC1155Action(_token common.Address, _id *big.Int, _to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _TokenActions.Contract.TransferERC1155Action(&_TokenActions.TransactOpts, _token, _id, _to, amount, data)
}

// TransferERC1155Action is a paid mutator transaction binding the contract method 0x7ca237e6.
//
// Solidity: function transferERC1155Action(address _token, uint256 _id, address _to, uint256 amount, bytes data) returns()
func (_TokenActions *TokenActionsTransactorSession) TransferERC1155Action(_token common.Address, _id *big.Int, _to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _TokenActions.Contract.TransferERC1155Action(&_TokenActions.TransactOpts, _token, _id, _to, amount, data)
}

// TransferERC721Action is a paid mutator transaction binding the contract method 0x1888bfd7.
//
// Solidity: function transferERC721Action(address _token, uint256 _id, address _to) returns()
func (_TokenActions *TokenActionsTransactor) TransferERC721Action(opts *bind.TransactOpts, _token common.Address, _id *big.Int, _to common.Address) (*types.Transaction, error) {
	return _TokenActions.contract.Transact(opts, "transferERC721Action", _token, _id, _to)
}

// TransferERC721Action is a paid mutator transaction binding the contract method 0x1888bfd7.
//
// Solidity: function transferERC721Action(address _token, uint256 _id, address _to) returns()
func (_TokenActions *TokenActionsSession) TransferERC721Action(_token common.Address, _id *big.Int, _to common.Address) (*types.Transaction, error) {
	return _TokenActions.Contract.TransferERC721Action(&_TokenActions.TransactOpts, _token, _id, _to)
}

// TransferERC721Action is a paid mutator transaction binding the contract method 0x1888bfd7.
//
// Solidity: function transferERC721Action(address _token, uint256 _id, address _to) returns()
func (_TokenActions *TokenActionsTransactorSession) TransferERC721Action(_token common.Address, _id *big.Int, _to common.Address) (*types.Transaction, error) {
	return _TokenActions.Contract.TransferERC721Action(&_TokenActions.TransactOpts, _token, _id, _to)
}
