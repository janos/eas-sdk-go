// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// SchemaRecord is an auto generated low-level Go binding around an user-defined struct.
type SchemaRecord struct {
	Uid       [32]byte
	Resolver  common.Address
	Revocable bool
	Schema    string
}

// SchemaRegistryMetaData contains all meta data concerning the SchemaRegistry contract.
var SchemaRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AlreadyExists\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"registerer\",\"type\":\"address\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"}],\"name\":\"getSchema\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"uid\",\"type\":\"bytes32\"},{\"internalType\":\"contractISchemaResolver\",\"name\":\"resolver\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"revocable\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"schema\",\"type\":\"string\"}],\"internalType\":\"structSchemaRecord\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"schema\",\"type\":\"string\"},{\"internalType\":\"contractISchemaResolver\",\"name\":\"resolver\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"revocable\",\"type\":\"bool\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e060405234801561001057600080fd5b5060016080819052600060a081905260c0819052806108e56100478239600060fe0152600060d50152600060ac01526108e56000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806354fd4d501461004657806360d7a27814610064578063a2ea7c6e14610085575b600080fd5b61004e6100a5565b60405161005b9190610598565b60405180910390f35b6100776100723660046105b2565b610148565b60405190815260200161005b565b610098610093366004610656565b610295565b60405161005b919061066f565b60606100d07f000000000000000000000000000000000000000000000000000000000000000061039f565b6100f97f000000000000000000000000000000000000000000000000000000000000000061039f565b6101227f000000000000000000000000000000000000000000000000000000000000000061039f565b604051602001610134939291906106ba565b604051602081830303815290604052905090565b60008060405180608001604052806000801b8152602001856001600160a01b03168152602001841515815260200187878080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250939094525092935091506101bd905082610431565b600081815260208190526040902054909150156101ed5760405163119b4fd360e11b815260040160405180910390fd5b808252600081815260208181526040918290208451815590840151600182018054938601511515600160a01b026001600160a81b03199094166001600160a01b0390921691909117929092179091556060830151839190600282019061025390826107b3565b50506040513381528291507f7d917fcbc9a29a9705ff9936ffa599500e4fd902e4486bae317414fe967b307c9060200160405180910390a29695505050505050565b60408051608081018252600080825260208201819052918101919091526060808201526000828152602081815260409182902082516080810184528154815260018201546001600160a01b03811693820193909352600160a01b90920460ff161515928201929092526002820180549192916060840191906103169061072a565b80601f01602080910402602001604051908101604052809291908181526020018280546103429061072a565b801561038f5780601f106103645761010080835404028352916020019161038f565b820191906000526020600020905b81548152906001019060200180831161037257829003601f168201915b5050505050815250509050919050565b606060006103ac83610471565b60010190506000816001600160401b038111156103cb576103cb610714565b6040519080825280601f01601f1916602001820160405280156103f5576020820181803683370190505b5090508181016020015b600019016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a85049450846103ff57509392505050565b600081606001518260200151836040015160405160200161045493929190610872565b604051602081830303815290604052805190602001209050919050565b60008072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b83106104b05772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6904ee2d6d415b85acef8160201b83106104da576904ee2d6d415b85acef8160201b830492506020015b662386f26fc1000083106104f857662386f26fc10000830492506010015b6305f5e1008310610510576305f5e100830492506008015b612710831061052457612710830492506004015b60648310610536576064830492506002015b600a8310610542576001015b92915050565b60005b8381101561056357818101518382015260200161054b565b50506000910152565b60008151808452610584816020860160208601610548565b601f01601f19169290920160200192915050565b6020815260006105ab602083018461056c565b9392505050565b600080600080606085870312156105c857600080fd5b84356001600160401b03808211156105df57600080fd5b818701915087601f8301126105f357600080fd5b81358181111561060257600080fd5b88602082850101111561061457600080fd5b602092830196509450508501356001600160a01b038116811461063657600080fd5b91506040850135801515811461064b57600080fd5b939692955090935050565b60006020828403121561066857600080fd5b5035919050565b602081528151602082015260018060a01b036020830151166040820152604082015115156060820152600060608301516080808401526106b260a084018261056c565b949350505050565b600084516106cc818460208901610548565b8083019050601760f91b80825285516106ec816001850160208a01610548565b60019201918201528351610707816002840160208801610548565b0160020195945050505050565b634e487b7160e01b600052604160045260246000fd5b600181811c9082168061073e57607f821691505b60208210810361075e57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156107ae57600081815260208120601f850160051c8101602086101561078b5750805b601f850160051c820191505b818110156107aa57828155600101610797565b5050505b505050565b81516001600160401b038111156107cc576107cc610714565b6107e0816107da845461072a565b84610764565b602080601f83116001811461081557600084156107fd5750858301515b600019600386901b1c1916600185901b1785556107aa565b600085815260208120601f198616915b8281101561084457888601518255948401946001909101908401610825565b50858210156108625787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60008451610884818460208901610548565b60609490941b6001600160601b0319169190930190815290151560f81b60148201526015019291505056fea2646970667358221220ee8383cd197230693296542c3fed7c3e499a9f09b40c2eedf7215d426152c4bd64736f6c63430008130033",
}

// SchemaRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use SchemaRegistryMetaData.ABI instead.
var SchemaRegistryABI = SchemaRegistryMetaData.ABI

// SchemaRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SchemaRegistryMetaData.Bin instead.
var SchemaRegistryBin = SchemaRegistryMetaData.Bin

// DeploySchemaRegistry deploys a new Ethereum contract, binding an instance of SchemaRegistry to it.
func DeploySchemaRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SchemaRegistry, error) {
	parsed, err := SchemaRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SchemaRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SchemaRegistry{SchemaRegistryCaller: SchemaRegistryCaller{contract: contract}, SchemaRegistryTransactor: SchemaRegistryTransactor{contract: contract}, SchemaRegistryFilterer: SchemaRegistryFilterer{contract: contract}}, nil
}

// SchemaRegistry is an auto generated Go binding around an Ethereum contract.
type SchemaRegistry struct {
	SchemaRegistryCaller     // Read-only binding to the contract
	SchemaRegistryTransactor // Write-only binding to the contract
	SchemaRegistryFilterer   // Log filterer for contract events
}

// SchemaRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SchemaRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SchemaRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SchemaRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SchemaRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SchemaRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SchemaRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SchemaRegistrySession struct {
	Contract     *SchemaRegistry   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SchemaRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SchemaRegistryCallerSession struct {
	Contract *SchemaRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// SchemaRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SchemaRegistryTransactorSession struct {
	Contract     *SchemaRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// SchemaRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SchemaRegistryRaw struct {
	Contract *SchemaRegistry // Generic contract binding to access the raw methods on
}

// SchemaRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SchemaRegistryCallerRaw struct {
	Contract *SchemaRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// SchemaRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SchemaRegistryTransactorRaw struct {
	Contract *SchemaRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSchemaRegistry creates a new instance of SchemaRegistry, bound to a specific deployed contract.
func NewSchemaRegistry(address common.Address, backend bind.ContractBackend) (*SchemaRegistry, error) {
	contract, err := bindSchemaRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SchemaRegistry{SchemaRegistryCaller: SchemaRegistryCaller{contract: contract}, SchemaRegistryTransactor: SchemaRegistryTransactor{contract: contract}, SchemaRegistryFilterer: SchemaRegistryFilterer{contract: contract}}, nil
}

// NewSchemaRegistryCaller creates a new read-only instance of SchemaRegistry, bound to a specific deployed contract.
func NewSchemaRegistryCaller(address common.Address, caller bind.ContractCaller) (*SchemaRegistryCaller, error) {
	contract, err := bindSchemaRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SchemaRegistryCaller{contract: contract}, nil
}

// NewSchemaRegistryTransactor creates a new write-only instance of SchemaRegistry, bound to a specific deployed contract.
func NewSchemaRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*SchemaRegistryTransactor, error) {
	contract, err := bindSchemaRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SchemaRegistryTransactor{contract: contract}, nil
}

// NewSchemaRegistryFilterer creates a new log filterer instance of SchemaRegistry, bound to a specific deployed contract.
func NewSchemaRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*SchemaRegistryFilterer, error) {
	contract, err := bindSchemaRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SchemaRegistryFilterer{contract: contract}, nil
}

// bindSchemaRegistry binds a generic wrapper to an already deployed contract.
func bindSchemaRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SchemaRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SchemaRegistry *SchemaRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SchemaRegistry.Contract.SchemaRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SchemaRegistry *SchemaRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SchemaRegistry.Contract.SchemaRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SchemaRegistry *SchemaRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SchemaRegistry.Contract.SchemaRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SchemaRegistry *SchemaRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SchemaRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SchemaRegistry *SchemaRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SchemaRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SchemaRegistry *SchemaRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SchemaRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetSchema is a free data retrieval call binding the contract method 0xa2ea7c6e.
//
// Solidity: function getSchema(bytes32 uid) view returns((bytes32,address,bool,string))
func (_SchemaRegistry *SchemaRegistryCaller) GetSchema(opts *bind.CallOpts, uid [32]byte) (SchemaRecord, error) {
	var out []interface{}
	err := _SchemaRegistry.contract.Call(opts, &out, "getSchema", uid)

	if err != nil {
		return *new(SchemaRecord), err
	}

	out0 := *abi.ConvertType(out[0], new(SchemaRecord)).(*SchemaRecord)

	return out0, err

}

// GetSchema is a free data retrieval call binding the contract method 0xa2ea7c6e.
//
// Solidity: function getSchema(bytes32 uid) view returns((bytes32,address,bool,string))
func (_SchemaRegistry *SchemaRegistrySession) GetSchema(uid [32]byte) (SchemaRecord, error) {
	return _SchemaRegistry.Contract.GetSchema(&_SchemaRegistry.CallOpts, uid)
}

// GetSchema is a free data retrieval call binding the contract method 0xa2ea7c6e.
//
// Solidity: function getSchema(bytes32 uid) view returns((bytes32,address,bool,string))
func (_SchemaRegistry *SchemaRegistryCallerSession) GetSchema(uid [32]byte) (SchemaRecord, error) {
	return _SchemaRegistry.Contract.GetSchema(&_SchemaRegistry.CallOpts, uid)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_SchemaRegistry *SchemaRegistryCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SchemaRegistry.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_SchemaRegistry *SchemaRegistrySession) Version() (string, error) {
	return _SchemaRegistry.Contract.Version(&_SchemaRegistry.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_SchemaRegistry *SchemaRegistryCallerSession) Version() (string, error) {
	return _SchemaRegistry.Contract.Version(&_SchemaRegistry.CallOpts)
}

// Register is a paid mutator transaction binding the contract method 0x60d7a278.
//
// Solidity: function register(string schema, address resolver, bool revocable) returns(bytes32)
func (_SchemaRegistry *SchemaRegistryTransactor) Register(opts *bind.TransactOpts, schema string, resolver common.Address, revocable bool) (*types.Transaction, error) {
	return _SchemaRegistry.contract.Transact(opts, "register", schema, resolver, revocable)
}

// Register is a paid mutator transaction binding the contract method 0x60d7a278.
//
// Solidity: function register(string schema, address resolver, bool revocable) returns(bytes32)
func (_SchemaRegistry *SchemaRegistrySession) Register(schema string, resolver common.Address, revocable bool) (*types.Transaction, error) {
	return _SchemaRegistry.Contract.Register(&_SchemaRegistry.TransactOpts, schema, resolver, revocable)
}

// Register is a paid mutator transaction binding the contract method 0x60d7a278.
//
// Solidity: function register(string schema, address resolver, bool revocable) returns(bytes32)
func (_SchemaRegistry *SchemaRegistryTransactorSession) Register(schema string, resolver common.Address, revocable bool) (*types.Transaction, error) {
	return _SchemaRegistry.Contract.Register(&_SchemaRegistry.TransactOpts, schema, resolver, revocable)
}

// SchemaRegistryRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the SchemaRegistry contract.
type SchemaRegistryRegisteredIterator struct {
	Event *SchemaRegistryRegistered // Event containing the contract specifics and raw log

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
func (it *SchemaRegistryRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SchemaRegistryRegistered)
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
		it.Event = new(SchemaRegistryRegistered)
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
func (it *SchemaRegistryRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SchemaRegistryRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SchemaRegistryRegistered represents a Registered event raised by the SchemaRegistry contract.
type SchemaRegistryRegistered struct {
	Uid        [32]byte
	Registerer common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x7d917fcbc9a29a9705ff9936ffa599500e4fd902e4486bae317414fe967b307c.
//
// Solidity: event Registered(bytes32 indexed uid, address registerer)
func (_SchemaRegistry *SchemaRegistryFilterer) FilterRegistered(opts *bind.FilterOpts, uid [][32]byte) (*SchemaRegistryRegisteredIterator, error) {

	var uidRule []interface{}
	for _, uidItem := range uid {
		uidRule = append(uidRule, uidItem)
	}

	logs, sub, err := _SchemaRegistry.contract.FilterLogs(opts, "Registered", uidRule)
	if err != nil {
		return nil, err
	}
	return &SchemaRegistryRegisteredIterator{contract: _SchemaRegistry.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x7d917fcbc9a29a9705ff9936ffa599500e4fd902e4486bae317414fe967b307c.
//
// Solidity: event Registered(bytes32 indexed uid, address registerer)
func (_SchemaRegistry *SchemaRegistryFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *SchemaRegistryRegistered, uid [][32]byte) (event.Subscription, error) {

	var uidRule []interface{}
	for _, uidItem := range uid {
		uidRule = append(uidRule, uidItem)
	}

	logs, sub, err := _SchemaRegistry.contract.WatchLogs(opts, "Registered", uidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SchemaRegistryRegistered)
				if err := _SchemaRegistry.contract.UnpackLog(event, "Registered", log); err != nil {
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

// ParseRegistered is a log parse operation binding the contract event 0x7d917fcbc9a29a9705ff9936ffa599500e4fd902e4486bae317414fe967b307c.
//
// Solidity: event Registered(bytes32 indexed uid, address registerer)
func (_SchemaRegistry *SchemaRegistryFilterer) ParseRegistered(log types.Log) (*SchemaRegistryRegistered, error) {
	event := new(SchemaRegistryRegistered)
	if err := _SchemaRegistry.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
