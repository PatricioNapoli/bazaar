// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chain

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
)

// ChainMetaData contains all meta data concerning the Chain contract.
var ChainMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"name\":\"viewPair\",\"outputs\":[{\"internalType\":\"uint112[]\",\"name\":\"\",\"type\":\"uint112[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ChainABI is the input ABI used to generate the binding from.
// Deprecated: Use ChainMetaData.ABI instead.
var ChainABI = ChainMetaData.ABI

// Chain is an auto generated Go binding around an Ethereum contract.
type Chain struct {
	ChainCaller     // Read-only binding to the contract
	ChainTransactor // Write-only binding to the contract
	ChainFilterer   // Log filterer for contract events
}

// ChainCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainSession struct {
	Contract     *Chain            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainCallerSession struct {
	Contract *ChainCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ChainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainTransactorSession struct {
	Contract     *ChainTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChainRaw struct {
	Contract *Chain // Generic contract binding to access the raw methods on
}

// ChainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainCallerRaw struct {
	Contract *ChainCaller // Generic read-only contract binding to access the raw methods on
}

// ChainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainTransactorRaw struct {
	Contract *ChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChain creates a new instance of Chain, bound to a specific deployed contract.
func NewChain(address common.Address, backend bind.ContractBackend) (*Chain, error) {
	contract, err := bindChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Chain{ChainCaller: ChainCaller{contract: contract}, ChainTransactor: ChainTransactor{contract: contract}, ChainFilterer: ChainFilterer{contract: contract}}, nil
}

// NewChainCaller creates a new read-only instance of Chain, bound to a specific deployed contract.
func NewChainCaller(address common.Address, caller bind.ContractCaller) (*ChainCaller, error) {
	contract, err := bindChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainCaller{contract: contract}, nil
}

// NewChainTransactor creates a new write-only instance of Chain, bound to a specific deployed contract.
func NewChainTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainTransactor, error) {
	contract, err := bindChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainTransactor{contract: contract}, nil
}

// NewChainFilterer creates a new log filterer instance of Chain, bound to a specific deployed contract.
func NewChainFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainFilterer, error) {
	contract, err := bindChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainFilterer{contract: contract}, nil
}

// bindChain binds a generic wrapper to an already deployed contract.
func bindChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chain *ChainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chain.Contract.ChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chain *ChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chain.Contract.ChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chain *ChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chain.Contract.ChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chain *ChainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chain *ChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chain *ChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chain.Contract.contract.Transact(opts, method, params...)
}

// ViewPair is a free data retrieval call binding the contract method 0x2245f986.
//
// Solidity: function viewPair(address[] ) view returns(uint112[])
func (_Chain *ChainCaller) ViewPair(opts *bind.CallOpts, arg0 []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _Chain.contract.Call(opts, &out, "viewPair", arg0)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// ViewPair is a free data retrieval call binding the contract method 0x2245f986.
//
// Solidity: function viewPair(address[] ) view returns(uint112[])
func (_Chain *ChainSession) ViewPair(arg0 []common.Address) ([]*big.Int, error) {
	return _Chain.Contract.ViewPair(&_Chain.CallOpts, arg0)
}

// ViewPair is a free data retrieval call binding the contract method 0x2245f986.
//
// Solidity: function viewPair(address[] ) view returns(uint112[])
func (_Chain *ChainCallerSession) ViewPair(arg0 []common.Address) ([]*big.Int, error) {
	return _Chain.Contract.ViewPair(&_Chain.CallOpts, arg0)
}
