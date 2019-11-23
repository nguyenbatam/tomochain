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

// Contract2ABI is the input ABI used to generate the binding from.
const Contract2ABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"subB\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"subA\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getB\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getA\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"addB\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"addA\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Contract2Bin is the compiled bytecode used for deploying new contracts.
const Contract2Bin = `0x608060405234801561001057600080fd5b50610173806100206000396000f3006080604052600436106100775763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166319a32cf9811461007c57806347860cd314610093578063a1c51915146100a8578063d46300fd146100cf578063f06106d9146100e4578063f33b8597146100f9575b600080fd5b34801561008857600080fd5b5061009161010e565b005b34801561009f57600080fd5b5061009161011a565b3480156100b457600080fd5b506100bd610126565b60408051918252519081900360200190f35b3480156100db57600080fd5b506100bd61012c565b3480156100f057600080fd5b50610091610132565b34801561010557600080fd5b5061009161013c565b60018054600019019055565b60008054600019019055565b60015490565b60005490565b6001805481019055565b6000805460010190555600a165627a7a72305820a3d48103bb82d45efa9723c5d08d43ced6eb08b7884cd98bcb841e5994ea4b870029`

// DeployContract2 deploys a new Ethereum contract, binding an instance of Contract2 to it.60806040526004361060525763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166347860cd381146057578063d46300fd14606b578063f33b859714608f575b600080fd5b348015606257600080fd5b50606960a1565b005b348015607657600080fd5b50607d60ad565b60408051918252519081900360200190f35b348015609a57600080fd5b50606960b3565b60008054600019019055565b60005490565b6000805460010190555600a165627a7a72305820582b2f70e6306b2c4829ef14f20900d071de2d8281edb51b53fb39783caa88b70029
func DeployContract2(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract2, error) {
	parsed, err := abi.JSON(strings.NewReader(Contract2ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Contract2Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract2{Contract2Caller: Contract2Caller{contract: contract}, Contract2Transactor: Contract2Transactor{contract: contract}, Contract2Filterer: Contract2Filterer{contract: contract}}, nil
}

// Contract2 is an auto generated Go binding around an Ethereum contract.
type Contract2 struct {
	Contract2Caller     // Read-only binding to the contract
	Contract2Transactor // Write-only binding to the contract
	Contract2Filterer   // Log filterer for contract events
}

// Contract2Caller is an auto generated read-only Go binding around an Ethereum contract.
type Contract2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Contract2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Contract2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Contract2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Contract2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Contract2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Contract2Session struct {
	Contract     *Contract2        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Contract2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Contract2CallerSession struct {
	Contract *Contract2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// Contract2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Contract2TransactorSession struct {
	Contract     *Contract2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// Contract2Raw is an auto generated low-level Go binding around an Ethereum contract.
type Contract2Raw struct {
	Contract *Contract2 // Generic contract binding to access the raw methods on
}

// Contract2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Contract2CallerRaw struct {
	Contract *Contract2Caller // Generic read-only contract binding to access the raw methods on
}

// Contract2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Contract2TransactorRaw struct {
	Contract *Contract2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewContract2 creates a new instance of Contract2, bound to a specific deployed contract.
func NewContract2(address common.Address, backend bind.ContractBackend) (*Contract2, error) {
	contract, err := bindContract2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract2{Contract2Caller: Contract2Caller{contract: contract}, Contract2Transactor: Contract2Transactor{contract: contract}, Contract2Filterer: Contract2Filterer{contract: contract}}, nil
}

// NewContract2Caller creates a new read-only instance of Contract2, bound to a specific deployed contract.
func NewContract2Caller(address common.Address, caller bind.ContractCaller) (*Contract2Caller, error) {
	contract, err := bindContract2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Contract2Caller{contract: contract}, nil
}

// NewContract2Transactor creates a new write-only instance of Contract2, bound to a specific deployed contract.
func NewContract2Transactor(address common.Address, transactor bind.ContractTransactor) (*Contract2Transactor, error) {
	contract, err := bindContract2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Contract2Transactor{contract: contract}, nil
}

// NewContract2Filterer creates a new log filterer instance of Contract2, bound to a specific deployed contract.
func NewContract2Filterer(address common.Address, filterer bind.ContractFilterer) (*Contract2Filterer, error) {
	contract, err := bindContract2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Contract2Filterer{contract: contract}, nil
}

// bindContract2 binds a generic wrapper to an already deployed contract.
func bindContract2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Contract2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract2 *Contract2Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract2.Contract.Contract2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract2 *Contract2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract2.Contract.Contract2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract2 *Contract2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract2.Contract.Contract2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract2 *Contract2CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract2 *Contract2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract2 *Contract2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract2.Contract.contract.Transact(opts, method, params...)
}

// GetA is a free data retrieval call binding the contract method 0xd46300fd.
//
// Solidity: function getA() constant returns(uint256)
func (_Contract2 *Contract2Caller) GetA(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract2.contract.Call(opts, out, "getA")
	return *ret0, err
}

// GetA is a free data retrieval call binding the contract method 0xd46300fd.
//
// Solidity: function getA() constant returns(uint256)
func (_Contract2 *Contract2Session) GetA() (*big.Int, error) {
	return _Contract2.Contract.GetA(&_Contract2.CallOpts)
}

// GetA is a free data retrieval call binding the contract method 0xd46300fd.
//
// Solidity: function getA() constant returns(uint256)
func (_Contract2 *Contract2CallerSession) GetA() (*big.Int, error) {
	return _Contract2.Contract.GetA(&_Contract2.CallOpts)
}

// GetB is a free data retrieval call binding the contract method 0xa1c51915.
//
// Solidity: function getB() constant returns(uint256)
func (_Contract2 *Contract2Caller) GetB(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract2.contract.Call(opts, out, "getB")
	return *ret0, err
}

// GetB is a free data retrieval call binding the contract method 0xa1c51915.
//
// Solidity: function getB() constant returns(uint256)
func (_Contract2 *Contract2Session) GetB() (*big.Int, error) {
	return _Contract2.Contract.GetB(&_Contract2.CallOpts)
}

// GetB is a free data retrieval call binding the contract method 0xa1c51915.
//
// Solidity: function getB() constant returns(uint256)
func (_Contract2 *Contract2CallerSession) GetB() (*big.Int, error) {
	return _Contract2.Contract.GetB(&_Contract2.CallOpts)
}

// AddA is a paid mutator transaction binding the contract method 0xf33b8597.
//
// Solidity: function addA() returns()
func (_Contract2 *Contract2Transactor) AddA(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract2.contract.Transact(opts, "addA")
}

// AddA is a paid mutator transaction binding the contract method 0xf33b8597.
//
// Solidity: function addA() returns()
func (_Contract2 *Contract2Session) AddA() (*types.Transaction, error) {
	return _Contract2.Contract.AddA(&_Contract2.TransactOpts)
}

// AddA is a paid mutator transaction binding the contract method 0xf33b8597.
//
// Solidity: function addA() returns()
func (_Contract2 *Contract2TransactorSession) AddA() (*types.Transaction, error) {
	return _Contract2.Contract.AddA(&_Contract2.TransactOpts)
}

// AddB is a paid mutator transaction binding the contract method 0xf06106d9.
//
// Solidity: function addB() returns()
func (_Contract2 *Contract2Transactor) AddB(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract2.contract.Transact(opts, "addB")
}

// AddB is a paid mutator transaction binding the contract method 0xf06106d9.
//
// Solidity: function addB() returns()
func (_Contract2 *Contract2Session) AddB() (*types.Transaction, error) {
	return _Contract2.Contract.AddB(&_Contract2.TransactOpts)
}

// AddB is a paid mutator transaction binding the contract method 0xf06106d9.
//
// Solidity: function addB() returns()
func (_Contract2 *Contract2TransactorSession) AddB() (*types.Transaction, error) {
	return _Contract2.Contract.AddB(&_Contract2.TransactOpts)
}

// SubA is a paid mutator transaction binding the contract method 0x47860cd3.
//
// Solidity: function subA() returns()
func (_Contract2 *Contract2Transactor) SubA(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract2.contract.Transact(opts, "subA")
}

// SubA is a paid mutator transaction binding the contract method 0x47860cd3.
//
// Solidity: function subA() returns()
func (_Contract2 *Contract2Session) SubA() (*types.Transaction, error) {
	return _Contract2.Contract.SubA(&_Contract2.TransactOpts)
}

// SubA is a paid mutator transaction binding the contract method 0x47860cd3.
//
// Solidity: function subA() returns()
func (_Contract2 *Contract2TransactorSession) SubA() (*types.Transaction, error) {
	return _Contract2.Contract.SubA(&_Contract2.TransactOpts)
}

// SubB is a paid mutator transaction binding the contract method 0x19a32cf9.
//
// Solidity: function subB() returns()
func (_Contract2 *Contract2Transactor) SubB(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract2.contract.Transact(opts, "subB")
}

// SubB is a paid mutator transaction binding the contract method 0x19a32cf9.
//
// Solidity: function subB() returns()
func (_Contract2 *Contract2Session) SubB() (*types.Transaction, error) {
	return _Contract2.Contract.SubB(&_Contract2.TransactOpts)
}

// SubB is a paid mutator transaction binding the contract method 0x19a32cf9.
//
// Solidity: function subB() returns()
func (_Contract2 *Contract2TransactorSession) SubB() (*types.Transaction, error) {
	return _Contract2.Contract.SubB(&_Contract2.TransactOpts)
}
