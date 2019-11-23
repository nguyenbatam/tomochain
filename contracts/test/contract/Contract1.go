// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strings"
)

// Contract1ABI is the input ABI used to generate the binding from.
const Contract1ABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"subA\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getA\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"addA\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Contract1Bin is the compiled bytecode used for deploying new contracts.
const Contract1Bin = `0x608060405234801561001057600080fd5b5060ea8061001f6000396000f30060806040526004361060525763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166347860cd381146057578063d46300fd14606b578063f33b859714608f575b600080fd5b348015606257600080fd5b50606960a1565b005b348015607657600080fd5b50607d60ad565b60408051918252519081900360200190f35b348015609a57600080fd5b50606960b3565b60008054600019019055565b60005490565b6000805460010190555600a165627a7a72305820582b2f70e6306b2c4829ef14f20900d071de2d8281edb51b53fb39783caa88b70029`

// DeployContract1 deploys a new Ethereum contract, binding an instance of Contract1 to it.
func DeployContract1(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract1, error) {
	parsed, err := abi.JSON(strings.NewReader(Contract1ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Contract1Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract1{Contract1Caller: Contract1Caller{contract: contract}, Contract1Transactor: Contract1Transactor{contract: contract}, Contract1Filterer: Contract1Filterer{contract: contract}}, nil
}

// Contract1 is an auto generated Go binding around an Ethereum contract.
type Contract1 struct {
	Contract1Caller     // Read-only binding to the contract
	Contract1Transactor // Write-only binding to the contract
	Contract1Filterer   // Log filterer for contract events
}

// Contract1Caller is an auto generated read-only Go binding around an Ethereum contract.
type Contract1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Contract1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Contract1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Contract1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Contract1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Contract1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Contract1Session struct {
	Contract     *Contract1        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Contract1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Contract1CallerSession struct {
	Contract *Contract1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// Contract1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Contract1TransactorSession struct {
	Contract     *Contract1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// Contract1Raw is an auto generated low-level Go binding around an Ethereum contract.
type Contract1Raw struct {
	Contract *Contract1 // Generic contract binding to access the raw methods on
}

// Contract1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Contract1CallerRaw struct {
	Contract *Contract1Caller // Generic read-only contract binding to access the raw methods on
}

// Contract1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Contract1TransactorRaw struct {
	Contract *Contract1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewContract1 creates a new instance of Contract1, bound to a specific deployed contract.
func NewContract1(address common.Address, backend bind.ContractBackend) (*Contract1, error) {
	contract, err := bindContract1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract1{Contract1Caller: Contract1Caller{contract: contract}, Contract1Transactor: Contract1Transactor{contract: contract}, Contract1Filterer: Contract1Filterer{contract: contract}}, nil
}

// NewContract1Caller creates a new read-only instance of Contract1, bound to a specific deployed contract.
func NewContract1Caller(address common.Address, caller bind.ContractCaller) (*Contract1Caller, error) {
	contract, err := bindContract1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Contract1Caller{contract: contract}, nil
}

// NewContract1Transactor creates a new write-only instance of Contract1, bound to a specific deployed contract.
func NewContract1Transactor(address common.Address, transactor bind.ContractTransactor) (*Contract1Transactor, error) {
	contract, err := bindContract1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Contract1Transactor{contract: contract}, nil
}

// NewContract1Filterer creates a new log filterer instance of Contract1, bound to a specific deployed contract.
func NewContract1Filterer(address common.Address, filterer bind.ContractFilterer) (*Contract1Filterer, error) {
	contract, err := bindContract1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Contract1Filterer{contract: contract}, nil
}

// bindContract1 binds a generic wrapper to an already deployed contract.
func bindContract1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Contract1ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract1 *Contract1Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract1.Contract.Contract1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract1 *Contract1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract1.Contract.Contract1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract1 *Contract1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract1.Contract.Contract1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract1 *Contract1CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract1 *Contract1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract1 *Contract1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract1.Contract.contract.Transact(opts, method, params...)
}

// GetA is a free data retrieval call binding the contract method 0xd46300fd.
//
// Solidity: function getA() constant returns(uint256)
func (_Contract1 *Contract1Caller) GetA(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract1.contract.Call(opts, out, "getA")
	return *ret0, err
}

// GetA is a free data retrieval call binding the contract method 0xd46300fd.
//
// Solidity: function getA() constant returns(uint256)
func (_Contract1 *Contract1Session) GetA() (*big.Int, error) {
	return _Contract1.Contract.GetA(&_Contract1.CallOpts)
}

// GetA is a free data retrieval call binding the contract method 0xd46300fd.
//
// Solidity: function getA() constant returns(uint256)
func (_Contract1 *Contract1CallerSession) GetA() (*big.Int, error) {
	return _Contract1.Contract.GetA(&_Contract1.CallOpts)
}

// AddA is a paid mutator transaction binding the contract method 0xf33b8597.
//
// Solidity: function addA() returns()
func (_Contract1 *Contract1Transactor) AddA(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract1.contract.Transact(opts, "addA")
}

// AddA is a paid mutator transaction binding the contract method 0xf33b8597.
//
// Solidity: function addA() returns()
func (_Contract1 *Contract1Session) AddA() (*types.Transaction, error) {
	return _Contract1.Contract.AddA(&_Contract1.TransactOpts)
}

// AddA is a paid mutator transaction binding the contract method 0xf33b8597.
//
// Solidity: function addA() returns()
func (_Contract1 *Contract1TransactorSession) AddA() (*types.Transaction, error) {
	return _Contract1.Contract.AddA(&_Contract1.TransactOpts)
}

// SubA is a paid mutator transaction binding the contract method 0x47860cd3.
//
// Solidity: function subA() returns()
func (_Contract1 *Contract1Transactor) SubA(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract1.contract.Transact(opts, "subA")
}

// SubA is a paid mutator transaction binding the contract method 0x47860cd3.
//
// Solidity: function subA() returns()
func (_Contract1 *Contract1Session) SubA() (*types.Transaction, error) {
	return _Contract1.Contract.SubA(&_Contract1.TransactOpts)
}

// SubA is a paid mutator transaction binding the contract method 0x47860cd3.
//
// Solidity: function subA() returns()
func (_Contract1 *Contract1TransactorSession) SubA() (*types.Transaction, error) {
	return _Contract1.Contract.SubA(&_Contract1.TransactOpts)
}
