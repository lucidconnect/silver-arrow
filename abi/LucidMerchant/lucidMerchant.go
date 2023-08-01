// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package LucidMerchant

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

// LucidMerchantMerchant is an auto generated low-level Go binding around an user-defined struct.
type LucidMerchantMerchant struct {
	Name             [32]byte
	Owner            common.Address
	ReceivingAddress common.Address
}

// LucidMerchantMutableMerchantData is an auto generated low-level Go binding around an user-defined struct.
type LucidMerchantMutableMerchantData struct {
	Name             [32]byte
	ReceivingAddress common.Address
}

// LucidMerchantMetaData contains all meta data concerning the LucidMerchant contract.
var LucidMerchantMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"receivingAddress\",\"type\":\"address\"}],\"internalType\":\"structLucidMerchant.MutableMerchantData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"createMerchant\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"}],\"name\":\"deleteMerchant\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"}],\"name\":\"getMerchant\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receivingAddress\",\"type\":\"address\"}],\"internalType\":\"structLucidMerchant.Merchant\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"merchants\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receivingAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"receivingAddress\",\"type\":\"address\"}],\"internalType\":\"structLucidMerchant.MutableMerchantData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"updateMerchant\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b50610f218061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610055575f3560e01c80632000302214610059578063359b582d1461008b57806365faba86146100bb578063c0f1d997146100d7578063f7bee303146100f3575b5f80fd5b610073600480360381019061006e91906109e2565b61010f565b60405161008293929190610a5b565b60405180910390f35b6100a560048036038101906100a091906109e2565b610172565b6040516100b29190610aee565b60405180910390f35b6100d560048036038101906100d09190610b29565b6102c6565b005b6100f160048036038101906100ec91906109e2565b6105f4565b005b61010d60048036038101906101089190610c6a565b6106fe565b005b5f602052805f5260405f205f91509050805f015490806001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806002015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905083565b61017a610955565b5f805f8481526020019081526020015f206040518060600160405290815f8201548152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600282015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090505f73ffffffffffffffffffffffffffffffffffffffff16816020015173ffffffffffffffffffffffffffffffffffffffff16036102bd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102b490610d02565b60405180910390fd5b80915050919050565b815f808281526020019081526020015f206001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610367576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161035e90610d6a565b60405180910390fd5b5f805f8581526020019081526020015f206040518060600160405290815f8201548152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600282015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090505f73ffffffffffffffffffffffffffffffffffffffff16816020015173ffffffffffffffffffffffffffffffffffffffff16036104aa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104a190610d02565b60405180910390fd5b5f801b835f0135036104f1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104e890610dd2565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff1683602001602081019061051b9190610df0565b73ffffffffffffffffffffffffffffffffffffffff1603610571576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161056890610e65565b60405180910390fd5b825f01355f808681526020019081526020015f205f018190555082602001602081019061059e9190610df0565b5f808681526020019081526020015f206002015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505050565b805f808281526020019081526020015f206001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610695576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161068c90610d6a565b60405180910390fd5b5f808381526020019081526020015f205f8082015f9055600182015f6101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600282015f6101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905550505050565b5f73ffffffffffffffffffffffffffffffffffffffff165f808481526020019081526020015f206001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461079e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161079590610ecd565b60405180910390fd5b5f801b815f0151036107e5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107dc90610dd2565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff16816020015173ffffffffffffffffffffffffffffffffffffffff1603610857576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161084e90610e65565b60405180910390fd5b5f6040518060600160405280835f015181526020013373ffffffffffffffffffffffffffffffffffffffff168152602001836020015173ffffffffffffffffffffffffffffffffffffffff168152509050805f808581526020019081526020015f205f820151815f01556020820151816001015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550905050505050565b60405180606001604052805f80191681526020015f73ffffffffffffffffffffffffffffffffffffffff1681526020015f73ffffffffffffffffffffffffffffffffffffffff1681525090565b5f604051905090565b5f80fd5b5f819050919050565b6109c1816109af565b81146109cb575f80fd5b50565b5f813590506109dc816109b8565b92915050565b5f602082840312156109f7576109f66109ab565b5b5f610a04848285016109ce565b91505092915050565b610a16816109af565b82525050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610a4582610a1c565b9050919050565b610a5581610a3b565b82525050565b5f606082019050610a6e5f830186610a0d565b610a7b6020830185610a4c565b610a886040830184610a4c565b949350505050565b610a99816109af565b82525050565b610aa881610a3b565b82525050565b606082015f820151610ac25f850182610a90565b506020820151610ad56020850182610a9f565b506040820151610ae86040850182610a9f565b50505050565b5f606082019050610b015f830184610aae565b92915050565b5f80fd5b5f60408284031215610b2057610b1f610b07565b5b81905092915050565b5f8060608385031215610b3f57610b3e6109ab565b5b5f610b4c858286016109ce565b9250506020610b5d85828601610b0b565b9150509250929050565b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b610bb182610b6b565b810181811067ffffffffffffffff82111715610bd057610bcf610b7b565b5b80604052505050565b5f610be26109a2565b9050610bee8282610ba8565b919050565b610bfc81610a3b565b8114610c06575f80fd5b50565b5f81359050610c1781610bf3565b92915050565b5f60408284031215610c3257610c31610b67565b5b610c3c6040610bd9565b90505f610c4b848285016109ce565b5f830152506020610c5e84828501610c09565b60208301525092915050565b5f8060608385031215610c8057610c7f6109ab565b5b5f610c8d858286016109ce565b9250506020610c9e85828601610c1d565b9150509250929050565b5f82825260208201905092915050565b7f4d65726368616e7420646f6573206e6f742065786973740000000000000000005f82015250565b5f610cec601783610ca8565b9150610cf782610cb8565b602082019050919050565b5f6020820190508181035f830152610d1981610ce0565b9050919050565b7f556e617574686f72697a656400000000000000000000000000000000000000005f82015250565b5f610d54600c83610ca8565b9150610d5f82610d20565b602082019050919050565b5f6020820190508181035f830152610d8181610d48565b9050919050565b7f496e76616c6964206e616d6500000000000000000000000000000000000000005f82015250565b5f610dbc600c83610ca8565b9150610dc782610d88565b602082019050919050565b5f6020820190508181035f830152610de981610db0565b9050919050565b5f60208284031215610e0557610e046109ab565b5b5f610e1284828501610c09565b91505092915050565b7f496e76616c696420726563656976696e672061646472657373000000000000005f82015250565b5f610e4f601983610ca8565b9150610e5a82610e1b565b602082019050919050565b5f6020820190508181035f830152610e7c81610e43565b9050919050565b7f4964656e74696669657220616c726561647920657869737473000000000000005f82015250565b5f610eb7601983610ca8565b9150610ec282610e83565b602082019050919050565b5f6020820190508181035f830152610ee481610eab565b905091905056fea2646970667358221220b6ef96b15694a3921e6ca559a916ca3207b1a60c02fe459f6d676e722ade882064736f6c63430008140033",
}

// LucidMerchantABI is the input ABI used to generate the binding from.
// Deprecated: Use LucidMerchantMetaData.ABI instead.
var LucidMerchantABI = LucidMerchantMetaData.ABI

// LucidMerchantBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LucidMerchantMetaData.Bin instead.
var LucidMerchantBin = LucidMerchantMetaData.Bin

// DeployLucidMerchant deploys a new Ethereum contract, binding an instance of LucidMerchant to it.
func DeployLucidMerchant(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LucidMerchant, error) {
	parsed, err := LucidMerchantMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LucidMerchantBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LucidMerchant{LucidMerchantCaller: LucidMerchantCaller{contract: contract}, LucidMerchantTransactor: LucidMerchantTransactor{contract: contract}, LucidMerchantFilterer: LucidMerchantFilterer{contract: contract}}, nil
}

// LucidMerchant is an auto generated Go binding around an Ethereum contract.
type LucidMerchant struct {
	LucidMerchantCaller     // Read-only binding to the contract
	LucidMerchantTransactor // Write-only binding to the contract
	LucidMerchantFilterer   // Log filterer for contract events
}

// LucidMerchantCaller is an auto generated read-only Go binding around an Ethereum contract.
type LucidMerchantCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LucidMerchantTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LucidMerchantTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LucidMerchantFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LucidMerchantFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LucidMerchantSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LucidMerchantSession struct {
	Contract     *LucidMerchant    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LucidMerchantCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LucidMerchantCallerSession struct {
	Contract *LucidMerchantCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// LucidMerchantTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LucidMerchantTransactorSession struct {
	Contract     *LucidMerchantTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// LucidMerchantRaw is an auto generated low-level Go binding around an Ethereum contract.
type LucidMerchantRaw struct {
	Contract *LucidMerchant // Generic contract binding to access the raw methods on
}

// LucidMerchantCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LucidMerchantCallerRaw struct {
	Contract *LucidMerchantCaller // Generic read-only contract binding to access the raw methods on
}

// LucidMerchantTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LucidMerchantTransactorRaw struct {
	Contract *LucidMerchantTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLucidMerchant creates a new instance of LucidMerchant, bound to a specific deployed contract.
func NewLucidMerchant(address common.Address, backend bind.ContractBackend) (*LucidMerchant, error) {
	contract, err := bindLucidMerchant(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LucidMerchant{LucidMerchantCaller: LucidMerchantCaller{contract: contract}, LucidMerchantTransactor: LucidMerchantTransactor{contract: contract}, LucidMerchantFilterer: LucidMerchantFilterer{contract: contract}}, nil
}

// NewLucidMerchantCaller creates a new read-only instance of LucidMerchant, bound to a specific deployed contract.
func NewLucidMerchantCaller(address common.Address, caller bind.ContractCaller) (*LucidMerchantCaller, error) {
	contract, err := bindLucidMerchant(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LucidMerchantCaller{contract: contract}, nil
}

// NewLucidMerchantTransactor creates a new write-only instance of LucidMerchant, bound to a specific deployed contract.
func NewLucidMerchantTransactor(address common.Address, transactor bind.ContractTransactor) (*LucidMerchantTransactor, error) {
	contract, err := bindLucidMerchant(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LucidMerchantTransactor{contract: contract}, nil
}

// NewLucidMerchantFilterer creates a new log filterer instance of LucidMerchant, bound to a specific deployed contract.
func NewLucidMerchantFilterer(address common.Address, filterer bind.ContractFilterer) (*LucidMerchantFilterer, error) {
	contract, err := bindLucidMerchant(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LucidMerchantFilterer{contract: contract}, nil
}

// bindLucidMerchant binds a generic wrapper to an already deployed contract.
func bindLucidMerchant(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LucidMerchantMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LucidMerchant *LucidMerchantRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LucidMerchant.Contract.LucidMerchantCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LucidMerchant *LucidMerchantRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LucidMerchant.Contract.LucidMerchantTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LucidMerchant *LucidMerchantRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LucidMerchant.Contract.LucidMerchantTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LucidMerchant *LucidMerchantCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LucidMerchant.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LucidMerchant *LucidMerchantTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LucidMerchant.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LucidMerchant *LucidMerchantTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LucidMerchant.Contract.contract.Transact(opts, method, params...)
}

// GetMerchant is a free data retrieval call binding the contract method 0x359b582d.
//
// Solidity: function getMerchant(bytes32 identifier) view returns((bytes32,address,address))
func (_LucidMerchant *LucidMerchantCaller) GetMerchant(opts *bind.CallOpts, identifier [32]byte) (LucidMerchantMerchant, error) {
	var out []interface{}
	err := _LucidMerchant.contract.Call(opts, &out, "getMerchant", identifier)

	if err != nil {
		return *new(LucidMerchantMerchant), err
	}

	out0 := *abi.ConvertType(out[0], new(LucidMerchantMerchant)).(*LucidMerchantMerchant)

	return out0, err

}

// GetMerchant is a free data retrieval call binding the contract method 0x359b582d.
//
// Solidity: function getMerchant(bytes32 identifier) view returns((bytes32,address,address))
func (_LucidMerchant *LucidMerchantSession) GetMerchant(identifier [32]byte) (LucidMerchantMerchant, error) {
	return _LucidMerchant.Contract.GetMerchant(&_LucidMerchant.CallOpts, identifier)
}

// GetMerchant is a free data retrieval call binding the contract method 0x359b582d.
//
// Solidity: function getMerchant(bytes32 identifier) view returns((bytes32,address,address))
func (_LucidMerchant *LucidMerchantCallerSession) GetMerchant(identifier [32]byte) (LucidMerchantMerchant, error) {
	return _LucidMerchant.Contract.GetMerchant(&_LucidMerchant.CallOpts, identifier)
}

// Merchants is a free data retrieval call binding the contract method 0x20003022.
//
// Solidity: function merchants(bytes32 ) view returns(bytes32 name, address owner, address receivingAddress)
func (_LucidMerchant *LucidMerchantCaller) Merchants(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Name             [32]byte
	Owner            common.Address
	ReceivingAddress common.Address
}, error) {
	var out []interface{}
	err := _LucidMerchant.contract.Call(opts, &out, "merchants", arg0)

	outstruct := new(struct {
		Name             [32]byte
		Owner            common.Address
		ReceivingAddress common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Name = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Owner = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.ReceivingAddress = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Merchants is a free data retrieval call binding the contract method 0x20003022.
//
// Solidity: function merchants(bytes32 ) view returns(bytes32 name, address owner, address receivingAddress)
func (_LucidMerchant *LucidMerchantSession) Merchants(arg0 [32]byte) (struct {
	Name             [32]byte
	Owner            common.Address
	ReceivingAddress common.Address
}, error) {
	return _LucidMerchant.Contract.Merchants(&_LucidMerchant.CallOpts, arg0)
}

// Merchants is a free data retrieval call binding the contract method 0x20003022.
//
// Solidity: function merchants(bytes32 ) view returns(bytes32 name, address owner, address receivingAddress)
func (_LucidMerchant *LucidMerchantCallerSession) Merchants(arg0 [32]byte) (struct {
	Name             [32]byte
	Owner            common.Address
	ReceivingAddress common.Address
}, error) {
	return _LucidMerchant.Contract.Merchants(&_LucidMerchant.CallOpts, arg0)
}

// CreateMerchant is a paid mutator transaction binding the contract method 0xf7bee303.
//
// Solidity: function createMerchant(bytes32 identifier, (bytes32,address) data) returns()
func (_LucidMerchant *LucidMerchantTransactor) CreateMerchant(opts *bind.TransactOpts, identifier [32]byte, data LucidMerchantMutableMerchantData) (*types.Transaction, error) {
	return _LucidMerchant.contract.Transact(opts, "createMerchant", identifier, data)
}

// CreateMerchant is a paid mutator transaction binding the contract method 0xf7bee303.
//
// Solidity: function createMerchant(bytes32 identifier, (bytes32,address) data) returns()
func (_LucidMerchant *LucidMerchantSession) CreateMerchant(identifier [32]byte, data LucidMerchantMutableMerchantData) (*types.Transaction, error) {
	return _LucidMerchant.Contract.CreateMerchant(&_LucidMerchant.TransactOpts, identifier, data)
}

// CreateMerchant is a paid mutator transaction binding the contract method 0xf7bee303.
//
// Solidity: function createMerchant(bytes32 identifier, (bytes32,address) data) returns()
func (_LucidMerchant *LucidMerchantTransactorSession) CreateMerchant(identifier [32]byte, data LucidMerchantMutableMerchantData) (*types.Transaction, error) {
	return _LucidMerchant.Contract.CreateMerchant(&_LucidMerchant.TransactOpts, identifier, data)
}

// DeleteMerchant is a paid mutator transaction binding the contract method 0xc0f1d997.
//
// Solidity: function deleteMerchant(bytes32 identifier) returns()
func (_LucidMerchant *LucidMerchantTransactor) DeleteMerchant(opts *bind.TransactOpts, identifier [32]byte) (*types.Transaction, error) {
	return _LucidMerchant.contract.Transact(opts, "deleteMerchant", identifier)
}

// DeleteMerchant is a paid mutator transaction binding the contract method 0xc0f1d997.
//
// Solidity: function deleteMerchant(bytes32 identifier) returns()
func (_LucidMerchant *LucidMerchantSession) DeleteMerchant(identifier [32]byte) (*types.Transaction, error) {
	return _LucidMerchant.Contract.DeleteMerchant(&_LucidMerchant.TransactOpts, identifier)
}

// DeleteMerchant is a paid mutator transaction binding the contract method 0xc0f1d997.
//
// Solidity: function deleteMerchant(bytes32 identifier) returns()
func (_LucidMerchant *LucidMerchantTransactorSession) DeleteMerchant(identifier [32]byte) (*types.Transaction, error) {
	return _LucidMerchant.Contract.DeleteMerchant(&_LucidMerchant.TransactOpts, identifier)
}

// UpdateMerchant is a paid mutator transaction binding the contract method 0x65faba86.
//
// Solidity: function updateMerchant(bytes32 identifier, (bytes32,address) data) returns()
func (_LucidMerchant *LucidMerchantTransactor) UpdateMerchant(opts *bind.TransactOpts, identifier [32]byte, data LucidMerchantMutableMerchantData) (*types.Transaction, error) {
	return _LucidMerchant.contract.Transact(opts, "updateMerchant", identifier, data)
}

// UpdateMerchant is a paid mutator transaction binding the contract method 0x65faba86.
//
// Solidity: function updateMerchant(bytes32 identifier, (bytes32,address) data) returns()
func (_LucidMerchant *LucidMerchantSession) UpdateMerchant(identifier [32]byte, data LucidMerchantMutableMerchantData) (*types.Transaction, error) {
	return _LucidMerchant.Contract.UpdateMerchant(&_LucidMerchant.TransactOpts, identifier, data)
}

// UpdateMerchant is a paid mutator transaction binding the contract method 0x65faba86.
//
// Solidity: function updateMerchant(bytes32 identifier, (bytes32,address) data) returns()
func (_LucidMerchant *LucidMerchantTransactorSession) UpdateMerchant(identifier [32]byte, data LucidMerchantMutableMerchantData) (*types.Transaction, error) {
	return _LucidMerchant.Contract.UpdateMerchant(&_LucidMerchant.TransactOpts, identifier, data)
}
