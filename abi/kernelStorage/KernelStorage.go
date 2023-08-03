// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package KernelStorage

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

// ExecutionDetail is an auto generated low-level Go binding around an user-defined struct.
type ExecutionDetail struct {
	ValidUntil *big.Int
	ValidAfter *big.Int
	Executor   common.Address
	Validator  common.Address
}

// KernelStorageMetaData contains all meta data concerning the KernelStorage contract.
var KernelStorageMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIEntryPoint\",\"name\":\"_entryPoint\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldValidator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newValidator\",\"type\":\"address\"}],\"name\":\"DefaultValidatorChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"ExecutionChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_disableFlag\",\"type\":\"bytes4\"}],\"name\":\"disableMode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"entryPoint\",\"outputs\":[{\"internalType\":\"contractIEntryPoint\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDefaultValidator\",\"outputs\":[{\"internalType\":\"contractIKernelValidator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDisabledMode\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_selector\",\"type\":\"bytes4\"}],\"name\":\"getExecution\",\"outputs\":[{\"components\":[{\"internalType\":\"uint48\",\"name\":\"validUntil\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"validAfter\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"internalType\":\"contractIKernelValidator\",\"name\":\"validator\",\"type\":\"address\"}],\"internalType\":\"structExecutionDetail\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastDisabledTime\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint192\",\"name\":\"key\",\"type\":\"uint192\"}],\"name\":\"getNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKernelValidator\",\"name\":\"_defaultValidator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKernelValidator\",\"name\":\"_defaultValidator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"setDefaultValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_selector\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"_executor\",\"type\":\"address\"},{\"internalType\":\"contractIKernelValidator\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"_validUntil\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"_validAfter\",\"type\":\"uint48\"},{\"internalType\":\"bytes\",\"name\":\"_enableData\",\"type\":\"bytes\"}],\"name\":\"setExecution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610e09380380610e0983398101604081905261002f916100a7565b6001600160a01b0381166080526001610046610073565b600101600a6101000a8154816001600160a01b0302191690836001600160a01b03160217905550506100f8565b6000806100a160017f439ffe7df606b78489639bc0b827913bd09e1246fa6802968a5b3694c53e0dd96100d7565b92915050565b6000602082840312156100b957600080fd5b81516001600160a01b03811681146100d057600080fd5b9392505050565b818103818111156100a157634e487b7160e01b600052601160045260246000fd5b608051610ccc61013d600039600081816101df0152818161025d015281816104470152818161051101528181610632015281816107b6015261092d0152610ccc6000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c806357b750471161007157806357b750471461019a57806388e7fd06146101bb578063b0d691fe146101da578063d087d28814610201578063d1f5789414610209578063d54162211461021c57600080fd5b80630b3dc354146100b957806329f8b174146100de5780633659cfe6146100f35780633e1b08121461010657806351166ba01461012757806355b14f5014610187575b600080fd5b6100c161022f565b6040516001600160a01b0390911681526020015b60405180910390f35b6100f16100ec366004610a88565b610252565b005b6100f1610101366004610b23565b61043c565b610119610114366004610b47565b6104ea565b6040519081526020016100d5565b61013a610135366004610b70565b61058a565b60408051825165ffffffffffff908116825260208085015190911690820152828201516001600160a01b0390811692820192909252606092830151909116918101919091526080016100d5565b6100f1610195366004610b8b565b610627565b6101a261075e565b6040516001600160e01b031990911681526020016100d5565b6101c3610774565b60405165ffffffffffff90911681526020016100d5565b6100c17f000000000000000000000000000000000000000000000000000000000000000081565b610119610797565b6100f1610217366004610b8b565b61082e565b6100f161022a366004610b70565b610922565b60006102396109c6565b60010154600160501b90046001600160a01b0316919050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148061028857503330145b6102ad5760405162461bcd60e51b81526004016102a490610be0565b60405180910390fd5b60405180608001604052808565ffffffffffff1681526020018465ffffffffffff168152602001876001600160a01b03168152602001866001600160a01b03168152506102f86109c6565b6001600160e01b031989166000908152600291909101602090815260409182902083518154928501518585015165ffffffffffff9283166bffffffffffffffffffffffff199095169490941766010000000000009290911691909102176bffffffffffffffffffffffff16600160601b6001600160a01b0393841602178155606090930151600190930180546001600160a01b031916938216939093179092555163064acaab60e11b815290861690630c959556906103bd9085908590600401610c2d565b600060405180830381600087803b1580156103d757600080fd5b505af11580156103eb573d6000803e3d6000fd5b50506040516001600160a01b038089169350891691506001600160e01b03198a16907fed03d2572564284398470d3f266a693e29ddfff3eba45fc06c5e91013d32135390600090a450505050505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148061047257503330145b61048e5760405162461bcd60e51b81526004016102a490610be0565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8181556040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a25050565b604051631aab3f0d60e11b81523060048201526001600160c01b03821660248201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906335567e1a90604401602060405180830381865afa158015610560573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105849190610c5c565b92915050565b6040805160808101825260008082526020820181905291810182905260608101919091526105b66109c6565b6001600160e01b0319909216600090815260029290920160209081526040928390208351608081018552815465ffffffffffff80821683526601000000000000820416938201939093526001600160a01b03600160601b909304831694810194909452600101541660608301525090565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148061065d57503330145b6106795760405162461bcd60e51b81526004016102a490610be0565b60006106836109c6565b60010154600160501b90046001600160a01b03169050836106a26109c6565b6001018054600160501b600160f01b031916600160501b6001600160a01b0393841602179055604051858216918316907fa35f5cdc5fbabb614b4cd5064ce5543f43dc8fab0e4da41255230eb8aba2531c90600090a360405163064acaab60e11b81526001600160a01b03851690630c959556906107269086908690600401610c2d565b600060405180830381600087803b15801561074057600080fd5b505af1158015610754573d6000803e3d6000fd5b5050505050505050565b60006107686109c6565b6001015460e01b919050565b600061077e6109c6565b60010154640100000000900465ffffffffffff16919050565b604051631aab3f0d60e11b8152306004820152600060248201819052907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906335567e1a90604401602060405180830381865afa158015610805573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108299190610c5c565b905090565b60006108386109c6565b6001810154909150600160501b90046001600160a01b03161561089d5760405162461bcd60e51b815260206004820152601c60248201527f6163636f756e743a20616c726561647920696e697469616c697a65640000000060448201526064016102a4565b600181018054600160501b600160f01b031916600160501b6001600160a01b038716908102919091179091556040516000907fa35f5cdc5fbabb614b4cd5064ce5543f43dc8fab0e4da41255230eb8aba2531c908290a360405163064acaab60e11b81526001600160a01b03851690630c959556906107269086908690600401610c2d565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148061095857503330145b6109745760405162461bcd60e51b81526004016102a490610be0565b8061097d6109c6565b600101805463ffffffff191660e09290921c9190911790554261099e6109c6565b60010160046101000a81548165ffffffffffff021916908365ffffffffffff16021790555050565b60008061058460017f439ffe7df606b78489639bc0b827913bd09e1246fa6802968a5b3694c53e0dd9610c75565b80356001600160e01b031981168114610a0c57600080fd5b919050565b6001600160a01b0381168114610a2657600080fd5b50565b803565ffffffffffff81168114610a0c57600080fd5b60008083601f840112610a5157600080fd5b50813567ffffffffffffffff811115610a6957600080fd5b602083019150836020828501011115610a8157600080fd5b9250929050565b600080600080600080600060c0888a031215610aa357600080fd5b610aac886109f4565b96506020880135610abc81610a11565b95506040880135610acc81610a11565b9450610ada60608901610a29565b9350610ae860808901610a29565b925060a088013567ffffffffffffffff811115610b0457600080fd5b610b108a828b01610a3f565b989b979a50959850939692959293505050565b600060208284031215610b3557600080fd5b8135610b4081610a11565b9392505050565b600060208284031215610b5957600080fd5b81356001600160c01b0381168114610b4057600080fd5b600060208284031215610b8257600080fd5b610b40826109f4565b600080600060408486031215610ba057600080fd5b8335610bab81610a11565b9250602084013567ffffffffffffffff811115610bc757600080fd5b610bd386828701610a3f565b9497909650939450505050565b6020808252602d908201527f6163636f756e743a206e6f742066726f6d20656e747279706f696e74206f722060408201526c37bbb732b91037b91039b2b63360991b606082015260800190565b60208152816020820152818360408301376000818301604090810191909152601f909201601f19160101919050565b600060208284031215610c6e57600080fd5b5051919050565b8181038181111561058457634e487b7160e01b600052601160045260246000fdfea26469706673582212206f4c0ed06fc1e447dfa12f5207dd5f6b25a7917f7d199f22ad1ab91d1f6a91d364736f6c63430008130033",
}

// KernelStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use KernelStorageMetaData.ABI instead.
var KernelStorageABI = KernelStorageMetaData.ABI

// KernelStorageBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KernelStorageMetaData.Bin instead.
var KernelStorageBin = KernelStorageMetaData.Bin

// DeployKernelStorage deploys a new Ethereum contract, binding an instance of KernelStorage to it.
func DeployKernelStorage(auth *bind.TransactOpts, backend bind.ContractBackend, _entryPoint common.Address) (common.Address, *types.Transaction, *KernelStorage, error) {
	parsed, err := KernelStorageMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KernelStorageBin), backend, _entryPoint)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KernelStorage{KernelStorageCaller: KernelStorageCaller{contract: contract}, KernelStorageTransactor: KernelStorageTransactor{contract: contract}, KernelStorageFilterer: KernelStorageFilterer{contract: contract}}, nil
}

// KernelStorage is an auto generated Go binding around an Ethereum contract.
type KernelStorage struct {
	KernelStorageCaller     // Read-only binding to the contract
	KernelStorageTransactor // Write-only binding to the contract
	KernelStorageFilterer   // Log filterer for contract events
}

// KernelStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type KernelStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KernelStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KernelStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KernelStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KernelStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KernelStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KernelStorageSession struct {
	Contract     *KernelStorage    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KernelStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KernelStorageCallerSession struct {
	Contract *KernelStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// KernelStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KernelStorageTransactorSession struct {
	Contract     *KernelStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// KernelStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type KernelStorageRaw struct {
	Contract *KernelStorage // Generic contract binding to access the raw methods on
}

// KernelStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KernelStorageCallerRaw struct {
	Contract *KernelStorageCaller // Generic read-only contract binding to access the raw methods on
}

// KernelStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KernelStorageTransactorRaw struct {
	Contract *KernelStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKernelStorage creates a new instance of KernelStorage, bound to a specific deployed contract.
func NewKernelStorage(address common.Address, backend bind.ContractBackend) (*KernelStorage, error) {
	contract, err := bindKernelStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KernelStorage{KernelStorageCaller: KernelStorageCaller{contract: contract}, KernelStorageTransactor: KernelStorageTransactor{contract: contract}, KernelStorageFilterer: KernelStorageFilterer{contract: contract}}, nil
}

// NewKernelStorageCaller creates a new read-only instance of KernelStorage, bound to a specific deployed contract.
func NewKernelStorageCaller(address common.Address, caller bind.ContractCaller) (*KernelStorageCaller, error) {
	contract, err := bindKernelStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KernelStorageCaller{contract: contract}, nil
}

// NewKernelStorageTransactor creates a new write-only instance of KernelStorage, bound to a specific deployed contract.
func NewKernelStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*KernelStorageTransactor, error) {
	contract, err := bindKernelStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KernelStorageTransactor{contract: contract}, nil
}

// NewKernelStorageFilterer creates a new log filterer instance of KernelStorage, bound to a specific deployed contract.
func NewKernelStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*KernelStorageFilterer, error) {
	contract, err := bindKernelStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KernelStorageFilterer{contract: contract}, nil
}

// bindKernelStorage binds a generic wrapper to an already deployed contract.
func bindKernelStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KernelStorageMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KernelStorage *KernelStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KernelStorage.Contract.KernelStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KernelStorage *KernelStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KernelStorage.Contract.KernelStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KernelStorage *KernelStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KernelStorage.Contract.KernelStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KernelStorage *KernelStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KernelStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KernelStorage *KernelStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KernelStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KernelStorage *KernelStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KernelStorage.Contract.contract.Transact(opts, method, params...)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_KernelStorage *KernelStorageCaller) EntryPoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KernelStorage.contract.Call(opts, &out, "entryPoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_KernelStorage *KernelStorageSession) EntryPoint() (common.Address, error) {
	return _KernelStorage.Contract.EntryPoint(&_KernelStorage.CallOpts)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_KernelStorage *KernelStorageCallerSession) EntryPoint() (common.Address, error) {
	return _KernelStorage.Contract.EntryPoint(&_KernelStorage.CallOpts)
}

// GetDefaultValidator is a free data retrieval call binding the contract method 0x0b3dc354.
//
// Solidity: function getDefaultValidator() view returns(address)
func (_KernelStorage *KernelStorageCaller) GetDefaultValidator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KernelStorage.contract.Call(opts, &out, "getDefaultValidator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDefaultValidator is a free data retrieval call binding the contract method 0x0b3dc354.
//
// Solidity: function getDefaultValidator() view returns(address)
func (_KernelStorage *KernelStorageSession) GetDefaultValidator() (common.Address, error) {
	return _KernelStorage.Contract.GetDefaultValidator(&_KernelStorage.CallOpts)
}

// GetDefaultValidator is a free data retrieval call binding the contract method 0x0b3dc354.
//
// Solidity: function getDefaultValidator() view returns(address)
func (_KernelStorage *KernelStorageCallerSession) GetDefaultValidator() (common.Address, error) {
	return _KernelStorage.Contract.GetDefaultValidator(&_KernelStorage.CallOpts)
}

// GetDisabledMode is a free data retrieval call binding the contract method 0x57b75047.
//
// Solidity: function getDisabledMode() view returns(bytes4)
func (_KernelStorage *KernelStorageCaller) GetDisabledMode(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _KernelStorage.contract.Call(opts, &out, "getDisabledMode")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// GetDisabledMode is a free data retrieval call binding the contract method 0x57b75047.
//
// Solidity: function getDisabledMode() view returns(bytes4)
func (_KernelStorage *KernelStorageSession) GetDisabledMode() ([4]byte, error) {
	return _KernelStorage.Contract.GetDisabledMode(&_KernelStorage.CallOpts)
}

// GetDisabledMode is a free data retrieval call binding the contract method 0x57b75047.
//
// Solidity: function getDisabledMode() view returns(bytes4)
func (_KernelStorage *KernelStorageCallerSession) GetDisabledMode() ([4]byte, error) {
	return _KernelStorage.Contract.GetDisabledMode(&_KernelStorage.CallOpts)
}

// GetExecution is a free data retrieval call binding the contract method 0x51166ba0.
//
// Solidity: function getExecution(bytes4 _selector) view returns((uint48,uint48,address,address))
func (_KernelStorage *KernelStorageCaller) GetExecution(opts *bind.CallOpts, _selector [4]byte) (ExecutionDetail, error) {
	var out []interface{}
	err := _KernelStorage.contract.Call(opts, &out, "getExecution", _selector)

	if err != nil {
		return *new(ExecutionDetail), err
	}

	out0 := *abi.ConvertType(out[0], new(ExecutionDetail)).(*ExecutionDetail)

	return out0, err

}

// GetExecution is a free data retrieval call binding the contract method 0x51166ba0.
//
// Solidity: function getExecution(bytes4 _selector) view returns((uint48,uint48,address,address))
func (_KernelStorage *KernelStorageSession) GetExecution(_selector [4]byte) (ExecutionDetail, error) {
	return _KernelStorage.Contract.GetExecution(&_KernelStorage.CallOpts, _selector)
}

// GetExecution is a free data retrieval call binding the contract method 0x51166ba0.
//
// Solidity: function getExecution(bytes4 _selector) view returns((uint48,uint48,address,address))
func (_KernelStorage *KernelStorageCallerSession) GetExecution(_selector [4]byte) (ExecutionDetail, error) {
	return _KernelStorage.Contract.GetExecution(&_KernelStorage.CallOpts, _selector)
}

// GetLastDisabledTime is a free data retrieval call binding the contract method 0x88e7fd06.
//
// Solidity: function getLastDisabledTime() view returns(uint48)
func (_KernelStorage *KernelStorageCaller) GetLastDisabledTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KernelStorage.contract.Call(opts, &out, "getLastDisabledTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLastDisabledTime is a free data retrieval call binding the contract method 0x88e7fd06.
//
// Solidity: function getLastDisabledTime() view returns(uint48)
func (_KernelStorage *KernelStorageSession) GetLastDisabledTime() (*big.Int, error) {
	return _KernelStorage.Contract.GetLastDisabledTime(&_KernelStorage.CallOpts)
}

// GetLastDisabledTime is a free data retrieval call binding the contract method 0x88e7fd06.
//
// Solidity: function getLastDisabledTime() view returns(uint48)
func (_KernelStorage *KernelStorageCallerSession) GetLastDisabledTime() (*big.Int, error) {
	return _KernelStorage.Contract.GetLastDisabledTime(&_KernelStorage.CallOpts)
}

// GetNonce is a free data retrieval call binding the contract method 0x3e1b0812.
//
// Solidity: function getNonce(uint192 key) view returns(uint256)
func (_KernelStorage *KernelStorageCaller) GetNonce(opts *bind.CallOpts, key *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KernelStorage.contract.Call(opts, &out, "getNonce", key)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNonce is a free data retrieval call binding the contract method 0x3e1b0812.
//
// Solidity: function getNonce(uint192 key) view returns(uint256)
func (_KernelStorage *KernelStorageSession) GetNonce(key *big.Int) (*big.Int, error) {
	return _KernelStorage.Contract.GetNonce(&_KernelStorage.CallOpts, key)
}

// GetNonce is a free data retrieval call binding the contract method 0x3e1b0812.
//
// Solidity: function getNonce(uint192 key) view returns(uint256)
func (_KernelStorage *KernelStorageCallerSession) GetNonce(key *big.Int) (*big.Int, error) {
	return _KernelStorage.Contract.GetNonce(&_KernelStorage.CallOpts, key)
}

// GetNonce0 is a free data retrieval call binding the contract method 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (_KernelStorage *KernelStorageCaller) GetNonce0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KernelStorage.contract.Call(opts, &out, "getNonce0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNonce0 is a free data retrieval call binding the contract method 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (_KernelStorage *KernelStorageSession) GetNonce0() (*big.Int, error) {
	return _KernelStorage.Contract.GetNonce0(&_KernelStorage.CallOpts)
}

// GetNonce0 is a free data retrieval call binding the contract method 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (_KernelStorage *KernelStorageCallerSession) GetNonce0() (*big.Int, error) {
	return _KernelStorage.Contract.GetNonce0(&_KernelStorage.CallOpts)
}

// DisableMode is a paid mutator transaction binding the contract method 0xd5416221.
//
// Solidity: function disableMode(bytes4 _disableFlag) returns()
func (_KernelStorage *KernelStorageTransactor) DisableMode(opts *bind.TransactOpts, _disableFlag [4]byte) (*types.Transaction, error) {
	return _KernelStorage.contract.Transact(opts, "disableMode", _disableFlag)
}

// DisableMode is a paid mutator transaction binding the contract method 0xd5416221.
//
// Solidity: function disableMode(bytes4 _disableFlag) returns()
func (_KernelStorage *KernelStorageSession) DisableMode(_disableFlag [4]byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.DisableMode(&_KernelStorage.TransactOpts, _disableFlag)
}

// DisableMode is a paid mutator transaction binding the contract method 0xd5416221.
//
// Solidity: function disableMode(bytes4 _disableFlag) returns()
func (_KernelStorage *KernelStorageTransactorSession) DisableMode(_disableFlag [4]byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.DisableMode(&_KernelStorage.TransactOpts, _disableFlag)
}

// Initialize is a paid mutator transaction binding the contract method 0xd1f57894.
//
// Solidity: function initialize(address _defaultValidator, bytes _data) returns()
func (_KernelStorage *KernelStorageTransactor) Initialize(opts *bind.TransactOpts, _defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.contract.Transact(opts, "initialize", _defaultValidator, _data)
}

// Initialize is a paid mutator transaction binding the contract method 0xd1f57894.
//
// Solidity: function initialize(address _defaultValidator, bytes _data) returns()
func (_KernelStorage *KernelStorageSession) Initialize(_defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.Initialize(&_KernelStorage.TransactOpts, _defaultValidator, _data)
}

// Initialize is a paid mutator transaction binding the contract method 0xd1f57894.
//
// Solidity: function initialize(address _defaultValidator, bytes _data) returns()
func (_KernelStorage *KernelStorageTransactorSession) Initialize(_defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.Initialize(&_KernelStorage.TransactOpts, _defaultValidator, _data)
}

// SetDefaultValidator is a paid mutator transaction binding the contract method 0x55b14f50.
//
// Solidity: function setDefaultValidator(address _defaultValidator, bytes _data) returns()
func (_KernelStorage *KernelStorageTransactor) SetDefaultValidator(opts *bind.TransactOpts, _defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.contract.Transact(opts, "setDefaultValidator", _defaultValidator, _data)
}

// SetDefaultValidator is a paid mutator transaction binding the contract method 0x55b14f50.
//
// Solidity: function setDefaultValidator(address _defaultValidator, bytes _data) returns()
func (_KernelStorage *KernelStorageSession) SetDefaultValidator(_defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.SetDefaultValidator(&_KernelStorage.TransactOpts, _defaultValidator, _data)
}

// SetDefaultValidator is a paid mutator transaction binding the contract method 0x55b14f50.
//
// Solidity: function setDefaultValidator(address _defaultValidator, bytes _data) returns()
func (_KernelStorage *KernelStorageTransactorSession) SetDefaultValidator(_defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.SetDefaultValidator(&_KernelStorage.TransactOpts, _defaultValidator, _data)
}

// SetExecution is a paid mutator transaction binding the contract method 0x29f8b174.
//
// Solidity: function setExecution(bytes4 _selector, address _executor, address _validator, uint48 _validUntil, uint48 _validAfter, bytes _enableData) returns()
func (_KernelStorage *KernelStorageTransactor) SetExecution(opts *bind.TransactOpts, _selector [4]byte, _executor common.Address, _validator common.Address, _validUntil *big.Int, _validAfter *big.Int, _enableData []byte) (*types.Transaction, error) {
	return _KernelStorage.contract.Transact(opts, "setExecution", _selector, _executor, _validator, _validUntil, _validAfter, _enableData)
}

// SetExecution is a paid mutator transaction binding the contract method 0x29f8b174.
//
// Solidity: function setExecution(bytes4 _selector, address _executor, address _validator, uint48 _validUntil, uint48 _validAfter, bytes _enableData) returns()
func (_KernelStorage *KernelStorageSession) SetExecution(_selector [4]byte, _executor common.Address, _validator common.Address, _validUntil *big.Int, _validAfter *big.Int, _enableData []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.SetExecution(&_KernelStorage.TransactOpts, _selector, _executor, _validator, _validUntil, _validAfter, _enableData)
}

// SetExecution is a paid mutator transaction binding the contract method 0x29f8b174.
//
// Solidity: function setExecution(bytes4 _selector, address _executor, address _validator, uint48 _validUntil, uint48 _validAfter, bytes _enableData) returns()
func (_KernelStorage *KernelStorageTransactorSession) SetExecution(_selector [4]byte, _executor common.Address, _validator common.Address, _validUntil *big.Int, _validAfter *big.Int, _enableData []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.SetExecution(&_KernelStorage.TransactOpts, _selector, _executor, _validator, _validUntil, _validAfter, _enableData)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address _newImplementation) returns()
func (_KernelStorage *KernelStorageTransactor) UpgradeTo(opts *bind.TransactOpts, _newImplementation common.Address) (*types.Transaction, error) {
	return _KernelStorage.contract.Transact(opts, "upgradeTo", _newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address _newImplementation) returns()
func (_KernelStorage *KernelStorageSession) UpgradeTo(_newImplementation common.Address) (*types.Transaction, error) {
	return _KernelStorage.Contract.UpgradeTo(&_KernelStorage.TransactOpts, _newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address _newImplementation) returns()
func (_KernelStorage *KernelStorageTransactorSession) UpgradeTo(_newImplementation common.Address) (*types.Transaction, error) {
	return _KernelStorage.Contract.UpgradeTo(&_KernelStorage.TransactOpts, _newImplementation)
}

// KernelStorageDefaultValidatorChangedIterator is returned from FilterDefaultValidatorChanged and is used to iterate over the raw logs and unpacked data for DefaultValidatorChanged events raised by the KernelStorage contract.
type KernelStorageDefaultValidatorChangedIterator struct {
	Event *KernelStorageDefaultValidatorChanged // Event containing the contract specifics and raw log

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
func (it *KernelStorageDefaultValidatorChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KernelStorageDefaultValidatorChanged)
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
		it.Event = new(KernelStorageDefaultValidatorChanged)
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
func (it *KernelStorageDefaultValidatorChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KernelStorageDefaultValidatorChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KernelStorageDefaultValidatorChanged represents a DefaultValidatorChanged event raised by the KernelStorage contract.
type KernelStorageDefaultValidatorChanged struct {
	OldValidator common.Address
	NewValidator common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDefaultValidatorChanged is a free log retrieval operation binding the contract event 0xa35f5cdc5fbabb614b4cd5064ce5543f43dc8fab0e4da41255230eb8aba2531c.
//
// Solidity: event DefaultValidatorChanged(address indexed oldValidator, address indexed newValidator)
func (_KernelStorage *KernelStorageFilterer) FilterDefaultValidatorChanged(opts *bind.FilterOpts, oldValidator []common.Address, newValidator []common.Address) (*KernelStorageDefaultValidatorChangedIterator, error) {

	var oldValidatorRule []interface{}
	for _, oldValidatorItem := range oldValidator {
		oldValidatorRule = append(oldValidatorRule, oldValidatorItem)
	}
	var newValidatorRule []interface{}
	for _, newValidatorItem := range newValidator {
		newValidatorRule = append(newValidatorRule, newValidatorItem)
	}

	logs, sub, err := _KernelStorage.contract.FilterLogs(opts, "DefaultValidatorChanged", oldValidatorRule, newValidatorRule)
	if err != nil {
		return nil, err
	}
	return &KernelStorageDefaultValidatorChangedIterator{contract: _KernelStorage.contract, event: "DefaultValidatorChanged", logs: logs, sub: sub}, nil
}

// WatchDefaultValidatorChanged is a free log subscription operation binding the contract event 0xa35f5cdc5fbabb614b4cd5064ce5543f43dc8fab0e4da41255230eb8aba2531c.
//
// Solidity: event DefaultValidatorChanged(address indexed oldValidator, address indexed newValidator)
func (_KernelStorage *KernelStorageFilterer) WatchDefaultValidatorChanged(opts *bind.WatchOpts, sink chan<- *KernelStorageDefaultValidatorChanged, oldValidator []common.Address, newValidator []common.Address) (event.Subscription, error) {

	var oldValidatorRule []interface{}
	for _, oldValidatorItem := range oldValidator {
		oldValidatorRule = append(oldValidatorRule, oldValidatorItem)
	}
	var newValidatorRule []interface{}
	for _, newValidatorItem := range newValidator {
		newValidatorRule = append(newValidatorRule, newValidatorItem)
	}

	logs, sub, err := _KernelStorage.contract.WatchLogs(opts, "DefaultValidatorChanged", oldValidatorRule, newValidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KernelStorageDefaultValidatorChanged)
				if err := _KernelStorage.contract.UnpackLog(event, "DefaultValidatorChanged", log); err != nil {
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

// ParseDefaultValidatorChanged is a log parse operation binding the contract event 0xa35f5cdc5fbabb614b4cd5064ce5543f43dc8fab0e4da41255230eb8aba2531c.
//
// Solidity: event DefaultValidatorChanged(address indexed oldValidator, address indexed newValidator)
func (_KernelStorage *KernelStorageFilterer) ParseDefaultValidatorChanged(log types.Log) (*KernelStorageDefaultValidatorChanged, error) {
	event := new(KernelStorageDefaultValidatorChanged)
	if err := _KernelStorage.contract.UnpackLog(event, "DefaultValidatorChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KernelStorageExecutionChangedIterator is returned from FilterExecutionChanged and is used to iterate over the raw logs and unpacked data for ExecutionChanged events raised by the KernelStorage contract.
type KernelStorageExecutionChangedIterator struct {
	Event *KernelStorageExecutionChanged // Event containing the contract specifics and raw log

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
func (it *KernelStorageExecutionChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KernelStorageExecutionChanged)
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
		it.Event = new(KernelStorageExecutionChanged)
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
func (it *KernelStorageExecutionChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KernelStorageExecutionChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KernelStorageExecutionChanged represents a ExecutionChanged event raised by the KernelStorage contract.
type KernelStorageExecutionChanged struct {
	Selector  [4]byte
	Executor  common.Address
	Validator common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterExecutionChanged is a free log retrieval operation binding the contract event 0xed03d2572564284398470d3f266a693e29ddfff3eba45fc06c5e91013d321353.
//
// Solidity: event ExecutionChanged(bytes4 indexed selector, address indexed executor, address indexed validator)
func (_KernelStorage *KernelStorageFilterer) FilterExecutionChanged(opts *bind.FilterOpts, selector [][4]byte, executor []common.Address, validator []common.Address) (*KernelStorageExecutionChangedIterator, error) {

	var selectorRule []interface{}
	for _, selectorItem := range selector {
		selectorRule = append(selectorRule, selectorItem)
	}
	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _KernelStorage.contract.FilterLogs(opts, "ExecutionChanged", selectorRule, executorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &KernelStorageExecutionChangedIterator{contract: _KernelStorage.contract, event: "ExecutionChanged", logs: logs, sub: sub}, nil
}

// WatchExecutionChanged is a free log subscription operation binding the contract event 0xed03d2572564284398470d3f266a693e29ddfff3eba45fc06c5e91013d321353.
//
// Solidity: event ExecutionChanged(bytes4 indexed selector, address indexed executor, address indexed validator)
func (_KernelStorage *KernelStorageFilterer) WatchExecutionChanged(opts *bind.WatchOpts, sink chan<- *KernelStorageExecutionChanged, selector [][4]byte, executor []common.Address, validator []common.Address) (event.Subscription, error) {

	var selectorRule []interface{}
	for _, selectorItem := range selector {
		selectorRule = append(selectorRule, selectorItem)
	}
	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _KernelStorage.contract.WatchLogs(opts, "ExecutionChanged", selectorRule, executorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KernelStorageExecutionChanged)
				if err := _KernelStorage.contract.UnpackLog(event, "ExecutionChanged", log); err != nil {
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

// ParseExecutionChanged is a log parse operation binding the contract event 0xed03d2572564284398470d3f266a693e29ddfff3eba45fc06c5e91013d321353.
//
// Solidity: event ExecutionChanged(bytes4 indexed selector, address indexed executor, address indexed validator)
func (_KernelStorage *KernelStorageFilterer) ParseExecutionChanged(log types.Log) (*KernelStorageExecutionChanged, error) {
	event := new(KernelStorageExecutionChanged)
	if err := _KernelStorage.contract.UnpackLog(event, "ExecutionChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KernelStorageUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the KernelStorage contract.
type KernelStorageUpgradedIterator struct {
	Event *KernelStorageUpgraded // Event containing the contract specifics and raw log

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
func (it *KernelStorageUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KernelStorageUpgraded)
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
		it.Event = new(KernelStorageUpgraded)
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
func (it *KernelStorageUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KernelStorageUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KernelStorageUpgraded represents a Upgraded event raised by the KernelStorage contract.
type KernelStorageUpgraded struct {
	NewImplementation common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed newImplementation)
func (_KernelStorage *KernelStorageFilterer) FilterUpgraded(opts *bind.FilterOpts, newImplementation []common.Address) (*KernelStorageUpgradedIterator, error) {

	var newImplementationRule []interface{}
	for _, newImplementationItem := range newImplementation {
		newImplementationRule = append(newImplementationRule, newImplementationItem)
	}

	logs, sub, err := _KernelStorage.contract.FilterLogs(opts, "Upgraded", newImplementationRule)
	if err != nil {
		return nil, err
	}
	return &KernelStorageUpgradedIterator{contract: _KernelStorage.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed newImplementation)
func (_KernelStorage *KernelStorageFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *KernelStorageUpgraded, newImplementation []common.Address) (event.Subscription, error) {

	var newImplementationRule []interface{}
	for _, newImplementationItem := range newImplementation {
		newImplementationRule = append(newImplementationRule, newImplementationItem)
	}

	logs, sub, err := _KernelStorage.contract.WatchLogs(opts, "Upgraded", newImplementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KernelStorageUpgraded)
				if err := _KernelStorage.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed newImplementation)
func (_KernelStorage *KernelStorageFilterer) ParseUpgraded(log types.Log) (*KernelStorageUpgraded, error) {
	event := new(KernelStorageUpgraded)
	if err := _KernelStorage.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
