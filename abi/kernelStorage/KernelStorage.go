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
	ValidAfter *big.Int
	ValidUntil *big.Int
	Executor   common.Address
	Validator  common.Address
}

// KernelStorageMetaData contains all meta data concerning the KernelStorage contract.
var KernelStorageMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIEntryPoint\",\"name\":\"_entryPoint\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAuthorizedCaller\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldValidator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newValidator\",\"type\":\"address\"}],\"name\":\"DefaultValidatorChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"ExecutionChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_disableFlag\",\"type\":\"bytes4\"}],\"name\":\"disableMode\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"entryPoint\",\"outputs\":[{\"internalType\":\"contractIEntryPoint\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDefaultValidator\",\"outputs\":[{\"internalType\":\"contractIKernelValidator\",\"name\":\"validator\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDisabledMode\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"disabled\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_selector\",\"type\":\"bytes4\"}],\"name\":\"getExecution\",\"outputs\":[{\"components\":[{\"internalType\":\"uint48\",\"name\":\"validAfter\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"validUntil\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"internalType\":\"contractIKernelValidator\",\"name\":\"validator\",\"type\":\"address\"}],\"internalType\":\"structExecutionDetail\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastDisabledTime\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKernelValidator\",\"name\":\"_defaultValidator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKernelValidator\",\"name\":\"_defaultValidator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"setDefaultValidator\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_selector\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"_executor\",\"type\":\"address\"},{\"internalType\":\"contractIKernelValidator\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"_validUntil\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"_validAfter\",\"type\":\"uint48\"},{\"internalType\":\"bytes\",\"name\":\"_enableData\",\"type\":\"bytes\"}],\"name\":\"setExecution\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a0346100d257601f610b2138819003918201601f19168301916001600160401b038311848410176100d7578084926020946040528339810103126100d257516001600160a01b03811681036100d2576080527f439ffe7df606b78489639bc0b827913bd09e1246fa6802968a5b3694c53e0dd98054600160501b600160f01b0319166a0100000000000000000000179055604051610a3390816100ee823960805181818160bf0152818161022b015281816102c90152818161038001528181610556015281816105f901526107000152f35b600080fd5b634e487b7160e01b600052604160045260246000fdfe6080604081815260048036101561001557600080fd5b600092833560e01c9081630b3dc354146108775750806329f8b1741461067b5780633659cfe6146105d25780633e1b08121461050c57806351166ba01461043a57806355b14f501461036557806357b750471461032d57806388e7fd06146102f8578063b0d691fe146102b4578063d087d288146101f5578063d1f578941461014b5763d5416221146100a757600080fd5b6020366003190112610147576100bb6108ac565b91337f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614158061013d575b610130575050600080516020610a1383398151915290815469ffffffffffff000000004260201b169160e01c9069ffffffffffffffffffff19161717905580f35b51637046c88d60e01b8152fd5b50303314156100ef565b8280fd5b508290610157366108f1565b600080516020610a13833981519152549295919391926001600160a01b039060501c81166101e6578661018a87986109c5565b16803b156101e2576101b39486809486519788958694859363064acaab60e11b8552840161099d565b03925af19081156101d957506101c65750f35b6101cf90610967565b6101d65780f35b80fd5b513d84823e3d90fd5b8580fd5b835162dc149f60e41b81528390fd5b509190346102b057816003193601126102b0578051631aab3f0d60e11b81523093810193909352602483018290526020836044817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9182156102a5579161026c575b6020925051908152f35b90506020823d821161029d575b816102866020938361097b565b81010312610298576020915190610262565b600080fd5b3d9150610279565b9051903d90823e3d90fd5b5080fd5b5050346102b057816003193601126102b057517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b5050346102b057816003193601126102b057600080516020610a13833981519152549051602091821c65ffffffffffff168152f35b5050346102b057816003193601126102b057602090600080516020610a138339815191525460e01b90519063ffffffff60e01b168152f35b509190610371366108f1565b919291906001600160a01b03337f00000000000000000000000000000000000000000000000000000000000000008216141580610430575b6104205795868697600080516020610a138339815191525460501c16956103cf816109c5565b1690818551967fa35f5cdc5fbabb614b4cd5064ce5543f43dc8fab0e4da41255230eb8aba2531c8980a3813b1561041c5786866101b382968296839563064acaab60e11b8552840161099d565b8680fd5b8351637046c88d60e01b81528790fd5b50303314156103a9565b5050346102b05760203660031901126102b05760018160809361045b6108ac565b816060845161046981610935565b8281528260208201528286820152015263ffffffff60e01b1681527f439ffe7df606b78489639bc0b827913bd09e1246fa6802968a5b3694c53e0dda60205220918051906104b682610935565b83549365ffffffffffff948581169586855260208501818360301c1681528486019260601c83526060878060a01b0380988196015416960195865284519788525116602087015251169084015251166060820152f35b509190346102b057602092836003193601126101475780356001600160c01b038116908190036105ce578251631aab3f0d60e11b81523092810192909252602482015283816044817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9283156105c35792610594575b5051908152f35b9091508281813d83116105bc575b6105ac818361097b565b810103126102985751903861058d565b503d6105a2565b8251903d90823e3d90fd5b8380fd5b506020366003190112610147576001600160a01b03813581811693909291848403610298577f00000000000000000000000000000000000000000000000000000000000000001633141580610671575b6101305750507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b8280a280f35b5030331415610622565b5060c0366003190112610147576106906108ac565b916001600160a01b03602435818116939192908490036101e2576044359483861680960361041c576064359265ffffffffffff948585168095036108735760843586811680910361086f5760a43567ffffffffffffffff811161086b576106fa90369085016108c3565b969094837f00000000000000000000000000000000000000000000000000000000000000001633141580610861575b610851578a928a8d979695936bffffffffffff0000000000006001948b519461075186610935565b8552602085019283528b85019384526060850197885263ffffffff60e01b169c8d8b527f439ffe7df606b78489639bc0b827913bd09e1246fa6802968a5b3694c53e0dda6020528b8b20945116915160301b16916bffffffffffffffffffffffff19905160601b169117178155019151166bffffffffffffffffffffffff60a01b825416179055873b15610147576107fa8451958693849363064acaab60e11b8552840161099d565b038183895af19081156108485750610835575b507fed03d2572564284398470d3f266a693e29ddfff3eba45fc06c5e91013d3213538480a480f35b61084190949194610967565b923861080d565b513d87823e3d90fd5b8651637046c88d60e01b81528590fd5b5030331415610729565b8a80fd5b8980fd5b8880fd5b8490346102b057816003193601126102b057600080516020610a138339815191525460501c6001600160a01b03168152602090f35b600435906001600160e01b03198216820361029857565b9181601f840112156102985782359167ffffffffffffffff8311610298576020838186019501011161029857565b906040600319830112610298576004356001600160a01b038116810361029857916024359067ffffffffffffffff821161029857610931916004016108c3565b9091565b6080810190811067ffffffffffffffff82111761095157604052565b634e487b7160e01b600052604160045260246000fd5b67ffffffffffffffff811161095157604052565b90601f8019910116810190811067ffffffffffffffff82111761095157604052565b90918060409360208452816020850152848401376000828201840152601f01601f1916010190565b600080516020610a1383398151915280547fffff0000000000000000000000000000000000000000ffffffffffffffffffff1660509290921b600160501b600160f01b031691909117905556fe439ffe7df606b78489639bc0b827913bd09e1246fa6802968a5b3694c53e0dd9",
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
// Solidity: function getDefaultValidator() view returns(address validator)
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
// Solidity: function getDefaultValidator() view returns(address validator)
func (_KernelStorage *KernelStorageSession) GetDefaultValidator() (common.Address, error) {
	return _KernelStorage.Contract.GetDefaultValidator(&_KernelStorage.CallOpts)
}

// GetDefaultValidator is a free data retrieval call binding the contract method 0x0b3dc354.
//
// Solidity: function getDefaultValidator() view returns(address validator)
func (_KernelStorage *KernelStorageCallerSession) GetDefaultValidator() (common.Address, error) {
	return _KernelStorage.Contract.GetDefaultValidator(&_KernelStorage.CallOpts)
}

// GetDisabledMode is a free data retrieval call binding the contract method 0x57b75047.
//
// Solidity: function getDisabledMode() view returns(bytes4 disabled)
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
// Solidity: function getDisabledMode() view returns(bytes4 disabled)
func (_KernelStorage *KernelStorageSession) GetDisabledMode() ([4]byte, error) {
	return _KernelStorage.Contract.GetDisabledMode(&_KernelStorage.CallOpts)
}

// GetDisabledMode is a free data retrieval call binding the contract method 0x57b75047.
//
// Solidity: function getDisabledMode() view returns(bytes4 disabled)
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

// GetNonce is a free data retrieval call binding the contract method 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (_KernelStorage *KernelStorageCaller) GetNonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KernelStorage.contract.Call(opts, &out, "getNonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNonce is a free data retrieval call binding the contract method 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (_KernelStorage *KernelStorageSession) GetNonce() (*big.Int, error) {
	return _KernelStorage.Contract.GetNonce(&_KernelStorage.CallOpts)
}

// GetNonce is a free data retrieval call binding the contract method 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (_KernelStorage *KernelStorageCallerSession) GetNonce() (*big.Int, error) {
	return _KernelStorage.Contract.GetNonce(&_KernelStorage.CallOpts)
}

// DisableMode is a paid mutator transaction binding the contract method 0xd5416221.
//
// Solidity: function disableMode(bytes4 _disableFlag) payable returns()
func (_KernelStorage *KernelStorageTransactor) DisableMode(opts *bind.TransactOpts, _disableFlag [4]byte) (*types.Transaction, error) {
	return _KernelStorage.contract.Transact(opts, "disableMode", _disableFlag)
}

// DisableMode is a paid mutator transaction binding the contract method 0xd5416221.
//
// Solidity: function disableMode(bytes4 _disableFlag) payable returns()
func (_KernelStorage *KernelStorageSession) DisableMode(_disableFlag [4]byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.DisableMode(&_KernelStorage.TransactOpts, _disableFlag)
}

// DisableMode is a paid mutator transaction binding the contract method 0xd5416221.
//
// Solidity: function disableMode(bytes4 _disableFlag) payable returns()
func (_KernelStorage *KernelStorageTransactorSession) DisableMode(_disableFlag [4]byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.DisableMode(&_KernelStorage.TransactOpts, _disableFlag)
}

// Initialize is a paid mutator transaction binding the contract method 0xd1f57894.
//
// Solidity: function initialize(address _defaultValidator, bytes _data) payable returns()
func (_KernelStorage *KernelStorageTransactor) Initialize(opts *bind.TransactOpts, _defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.contract.Transact(opts, "initialize", _defaultValidator, _data)
}

// Initialize is a paid mutator transaction binding the contract method 0xd1f57894.
//
// Solidity: function initialize(address _defaultValidator, bytes _data) payable returns()
func (_KernelStorage *KernelStorageSession) Initialize(_defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.Initialize(&_KernelStorage.TransactOpts, _defaultValidator, _data)
}

// Initialize is a paid mutator transaction binding the contract method 0xd1f57894.
//
// Solidity: function initialize(address _defaultValidator, bytes _data) payable returns()
func (_KernelStorage *KernelStorageTransactorSession) Initialize(_defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.Initialize(&_KernelStorage.TransactOpts, _defaultValidator, _data)
}

// SetDefaultValidator is a paid mutator transaction binding the contract method 0x55b14f50.
//
// Solidity: function setDefaultValidator(address _defaultValidator, bytes _data) payable returns()
func (_KernelStorage *KernelStorageTransactor) SetDefaultValidator(opts *bind.TransactOpts, _defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.contract.Transact(opts, "setDefaultValidator", _defaultValidator, _data)
}

// SetDefaultValidator is a paid mutator transaction binding the contract method 0x55b14f50.
//
// Solidity: function setDefaultValidator(address _defaultValidator, bytes _data) payable returns()
func (_KernelStorage *KernelStorageSession) SetDefaultValidator(_defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.SetDefaultValidator(&_KernelStorage.TransactOpts, _defaultValidator, _data)
}

// SetDefaultValidator is a paid mutator transaction binding the contract method 0x55b14f50.
//
// Solidity: function setDefaultValidator(address _defaultValidator, bytes _data) payable returns()
func (_KernelStorage *KernelStorageTransactorSession) SetDefaultValidator(_defaultValidator common.Address, _data []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.SetDefaultValidator(&_KernelStorage.TransactOpts, _defaultValidator, _data)
}

// SetExecution is a paid mutator transaction binding the contract method 0x29f8b174.
//
// Solidity: function setExecution(bytes4 _selector, address _executor, address _validator, uint48 _validUntil, uint48 _validAfter, bytes _enableData) payable returns()
func (_KernelStorage *KernelStorageTransactor) SetExecution(opts *bind.TransactOpts, _selector [4]byte, _executor common.Address, _validator common.Address, _validUntil *big.Int, _validAfter *big.Int, _enableData []byte) (*types.Transaction, error) {
	return _KernelStorage.contract.Transact(opts, "setExecution", _selector, _executor, _validator, _validUntil, _validAfter, _enableData)
}

// SetExecution is a paid mutator transaction binding the contract method 0x29f8b174.
//
// Solidity: function setExecution(bytes4 _selector, address _executor, address _validator, uint48 _validUntil, uint48 _validAfter, bytes _enableData) payable returns()
func (_KernelStorage *KernelStorageSession) SetExecution(_selector [4]byte, _executor common.Address, _validator common.Address, _validUntil *big.Int, _validAfter *big.Int, _enableData []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.SetExecution(&_KernelStorage.TransactOpts, _selector, _executor, _validator, _validUntil, _validAfter, _enableData)
}

// SetExecution is a paid mutator transaction binding the contract method 0x29f8b174.
//
// Solidity: function setExecution(bytes4 _selector, address _executor, address _validator, uint48 _validUntil, uint48 _validAfter, bytes _enableData) payable returns()
func (_KernelStorage *KernelStorageTransactorSession) SetExecution(_selector [4]byte, _executor common.Address, _validator common.Address, _validUntil *big.Int, _validAfter *big.Int, _enableData []byte) (*types.Transaction, error) {
	return _KernelStorage.Contract.SetExecution(&_KernelStorage.TransactOpts, _selector, _executor, _validator, _validUntil, _validAfter, _enableData)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address _newImplementation) payable returns()
func (_KernelStorage *KernelStorageTransactor) UpgradeTo(opts *bind.TransactOpts, _newImplementation common.Address) (*types.Transaction, error) {
	return _KernelStorage.contract.Transact(opts, "upgradeTo", _newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address _newImplementation) payable returns()
func (_KernelStorage *KernelStorageSession) UpgradeTo(_newImplementation common.Address) (*types.Transaction, error) {
	return _KernelStorage.Contract.UpgradeTo(&_KernelStorage.TransactOpts, _newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address _newImplementation) payable returns()
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
