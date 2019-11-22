// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"math/big"
	"strings"
)

// ConfigRewardABI is the input ABI used to generate the binding from.
const ConfigRewardABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"owners\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"revokeConfirmation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"confirmations\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pending\",\"type\":\"bool\"},{\"name\":\"executed\",\"type\":\"bool\"}],\"name\":\"getTransactionCount\",\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"isConfirmed\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"getConfirmationCount\",\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transactions\",\"outputs\":[{\"name\":\"rate\",\"type\":\"uint256\"},{\"name\":\"executed\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwners\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"from\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"uint256\"},{\"name\":\"pending\",\"type\":\"bool\"},{\"name\":\"executed\",\"type\":\"bool\"}],\"name\":\"getTransactionIds\",\"outputs\":[{\"name\":\"_transactionIds\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"getConfirmations\",\"outputs\":[{\"name\":\"_confirmations\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"transactionCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_required\",\"type\":\"uint256\"}],\"name\":\"changeRequirement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"confirmTransaction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"submitTransaction\",\"outputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_OWNER_COUNT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"required\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"replaceOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"executeTransaction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_owners\",\"type\":\"address[]\"},{\"name\":\"_required\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"Confirmation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"Revocation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"Submission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"Execution\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"ExecutionFailure\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnerAddition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnerRemoval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"required\",\"type\":\"uint256\"}],\"name\":\"RequirementChange\",\"type\":\"event\"}]"

// ConfigRewardBin is the compiled bytecode used for deploying new contracts.
const ConfigRewardBin = `0x60806040523480156200001157600080fd5b50604051620013373803806200133783398101604052805160208201519101805190919060009082603282118015906200004b5750818111155b80156200005757508015155b80156200006357508115155b15156200006f57600080fd5b600092505b845183101562000147576003600086858151811015156200009157fe5b6020908102909101810151600160a060020a031682528101919091526040016000205460ff16158015620000e757508483815181101515620000cf57fe5b90602001906020020151600160a060020a0316600014155b1515620000f357600080fd5b60016003600087868151811015156200010857fe5b602090810291909101810151600160a060020a03168252810191909152604001600020805460ff19169115159190911790556001929092019162000074565b84516200015c90600490602088019062000172565b5050506005919091555050600a60005562000206565b828054828255906000526020600020908101928215620001ca579160200282015b82811115620001ca5782518254600160a060020a031916600160a060020a0390911617825560209092019160019091019062000193565b50620001d8929150620001dc565b5090565b6200020391905b80821115620001d8578054600160a060020a0319168155600101620001e3565b90565b61112180620002166000396000f3006080604052600436106101275763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663025e7c278114610169578063173825d91461019d57806320ea8d86146101be5780632f54bf6e146101d65780633411c81c1461020b578063547415251461022f578063679aefce146102605780637065cb4814610275578063784547a7146102965780638b51d13f146102ae5780639ace38c2146102c6578063a0e67e2b146102f7578063a8abe69a1461035c578063b5dc40c314610381578063b77bf60014610399578063ba51a6df146103ae578063c01a8c84146103c6578063c34a0ead146103de578063d74f8edd146103f6578063dc8452cd1461040b578063e20056e614610420578063ee22610b14610447575b60003411156101675760408051348152905133917fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c919081900360200190a25b005b34801561017557600080fd5b5061018160043561045f565b60408051600160a060020a039092168252519081900360200190f35b3480156101a957600080fd5b50610167600160a060020a0360043516610487565b3480156101ca57600080fd5b506101676004356105fe565b3480156101e257600080fd5b506101f7600160a060020a03600435166106b9565b604080519115158252519081900360200190f35b34801561021757600080fd5b506101f7600435600160a060020a03602435166106ce565b34801561023b57600080fd5b5061024e600435151560243515156106ee565b60408051918252519081900360200190f35b34801561026c57600080fd5b5061024e61075c565b34801561028157600080fd5b50610167600160a060020a0360043516610763565b3480156102a257600080fd5b506101f7600435610888565b3480156102ba57600080fd5b5061024e60043561090c565b3480156102d257600080fd5b506102de60043561097b565b6040805192835290151560208301528051918290030190f35b34801561030357600080fd5b5061030c610997565b60408051602080825283518183015283519192839290830191858101910280838360005b83811015610348578181015183820152602001610330565b505050509050019250505060405180910390f35b34801561036857600080fd5b5061030c600435602435604435151560643515156109f9565b34801561038d57600080fd5b5061030c600435610b34565b3480156103a557600080fd5b5061024e610cad565b3480156103ba57600080fd5b50610167600435610cb3565b3480156103d257600080fd5b50610167600435610d32565b3480156103ea57600080fd5b5061024e600435610dd5565b34801561040257600080fd5b5061024e610df9565b34801561041757600080fd5b5061024e610dfe565b34801561042c57600080fd5b50610167600160a060020a0360043581169060243516610e04565b34801561045357600080fd5b50610167600435610f8e565b600480548290811061046d57fe5b600091825260209091200154600160a060020a0316905081565b600033301461049557600080fd5b600160a060020a038216600090815260036020526040902054829060ff1615156104be57600080fd5b600160a060020a0383166000908152600360205260408120805460ff1916905591505b600454600019018210156105995782600160a060020a031660048381548110151561050857fe5b600091825260209091200154600160a060020a0316141561058e5760048054600019810190811061053557fe5b60009182526020909120015460048054600160a060020a03909216918490811061055b57fe5b9060005260206000200160006101000a815481600160a060020a030219169083600160a060020a03160217905550610599565b6001909101906104e1565b6004805460001901906105ac90826110ae565b5060045460055411156105c5576004546105c590610cb3565b604051600160a060020a038416907f8001553a916ef2f495d26a907cc54d96ed840d7bda71e73194bf5a9df7a76b9090600090a2505050565b3360008181526003602052604090205460ff16151561061c57600080fd5b60008281526002602090815260408083203380855292529091205483919060ff16151561064857600080fd5b60008481526001602081905260409091200154849060ff161561066a57600080fd5b6000858152600260209081526040808320338085529252808320805460ff191690555187927ff6a317157440607f36269043eb55f1287a5a19ba2216afeab88cd46cbcfb88e991a35050505050565b60036020526000908152604090205460ff1681565b600260209081526000928352604080842090915290825290205460ff1681565b6000805b6006548110156107555783801561071c57506000818152600160208190526040909120015460ff16155b80610741575082801561074157506000818152600160208190526040909120015460ff165b1561074d576001820191505b6001016106f2565b5092915050565b6000545b90565b33301461076f57600080fd5b600160a060020a038116600090815260036020526040902054819060ff161561079757600080fd5b81600160a060020a03811615156107ad57600080fd5b600480549050600101600554603282111580156107ca5750818111155b80156107d557508015155b80156107e057508115155b15156107eb57600080fd5b600160a060020a038516600081815260036020526040808220805460ff1916600190811790915560048054918201815583527f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b01805473ffffffffffffffffffffffffffffffffffffffff191684179055517ff39e6e1eb0edcf53c221607b54b00cd28f3196fed0a24994dc308b8f611b682d9190a25050505050565b600080805b60045481101561090557600084815260026020526040812060048054919291849081106108b657fe5b6000918252602080832090910154600160a060020a0316835282019290925260400190205460ff16156108ea576001820191505b6005548214156108fd5760019250610905565b60010161088d565b5050919050565b6000805b600454811015610975576000838152600260205260408120600480549192918490811061093957fe5b6000918252602080832090910154600160a060020a0316835282019290925260400190205460ff161561096d576001820191505b600101610910565b50919050565b6001602081905260009182526040909120805491015460ff1682565b606060048054806020026020016040519081016040528092919081815260200182805480156109ef57602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116109d1575b5050505050905090565b606080600080600654604051908082528060200260200182016040528015610a2b578160200160208202803883390190505b50925060009150600090505b600654811015610ab457858015610a6157506000818152600160208190526040909120015460ff16155b80610a865750848015610a8657506000818152600160208190526040909120015460ff165b15610aac57808383815181101515610a9a57fe5b60209081029091010152600191909101905b600101610a37565b878703604051908082528060200260200182016040528015610ae0578160200160208202803883390190505b5093508790505b86811015610b29578281815181101515610afd57fe5b9060200190602002015184898303815181101515610b1757fe5b60209081029091010152600101610ae7565b505050949350505050565b606080600080600480549050604051908082528060200260200182016040528015610b69578160200160208202803883390190505b50925060009150600090505b600454811015610c265760008581526002602052604081206004805491929184908110610b9e57fe5b6000918252602080832090910154600160a060020a0316835282019290925260400190205460ff1615610c1e576004805482908110610bd957fe5b6000918252602090912001548351600160a060020a0390911690849084908110610bff57fe5b600160a060020a03909216602092830290910190910152600191909101905b600101610b75565b81604051908082528060200260200182016040528015610c50578160200160208202803883390190505b509350600090505b81811015610ca5578281815181101515610c6e57fe5b906020019060200201518482815181101515610c8657fe5b600160a060020a03909216602092830290910190910152600101610c58565b505050919050565b60065481565b333014610cbf57600080fd5b6004548160328211801590610cd45750818111155b8015610cdf57508015155b8015610cea57508115155b1515610cf557600080fd5b60058390556040805184815290517fa3f1ee9126a074d9326c682f561767f710e927faa811f7a99829d49dc421797a9181900360200190a1505050565b3360008181526003602052604090205460ff161515610d5057600080fd5b60008281526002602090815260408083203380855292529091205483919060ff1615610d7b57600080fd5b6000848152600260209081526040808320338085529252808320805460ff191660011790555186927f4a504a94899432a9846e1aa406dceb1bcfd538bb839071d49d1e5e23f5be30ef91a3610dcf84610f8e565b50505050565b60008160648110610de557600080fd5b610dee8361103a565b915061097582610d32565b603281565b60055481565b6000333014610e1257600080fd5b600160a060020a038316600090815260036020526040902054839060ff161515610e3b57600080fd5b600160a060020a038316600090815260036020526040902054839060ff1615610e6357600080fd5b600092505b600454831015610ef45784600160a060020a0316600484815481101515610e8b57fe5b600091825260209091200154600160a060020a03161415610ee95783600484815481101515610eb657fe5b9060005260206000200160006101000a815481600160a060020a030219169083600160a060020a03160217905550610ef4565b600190920191610e68565b600160a060020a03808616600081815260036020526040808220805460ff1990811690915593881682528082208054909416600117909355915190917f8001553a916ef2f495d26a907cc54d96ed840d7bda71e73194bf5a9df7a76b9091a2604051600160a060020a038516907ff39e6e1eb0edcf53c221607b54b00cd28f3196fed0a24994dc308b8f611b682d90600090a25050505050565b3360008181526003602052604081205490919060ff161515610faf57600080fd5b60008381526002602090815260408083203380855292529091205484919060ff161515610fdb57600080fd5b60008581526001602081905260409091200154859060ff1615610ffd57600080fd5b61100686610888565b156110325760008681526001602081905260408220808201805460ff1916909217909155805490915594505b505050505050565b60068054604080518082018252848152600060208083018281528583526001918290528483209351845551928101805460ff191693151593909317909255845490910190935551909182917fc0ba8fe4b176c1714197d43b9cc6bcf797a4a7461c5fe8d0ef6e184ae7601e519190a2919050565b8154818355818111156110d2576000838152602090206110d29181019083016110d7565b505050565b61076091905b808211156110f157600081556001016110dd565b50905600a165627a7a72305820c2c97e76052676997bf519226d91bb1b954fca253ce3881ee237a7c8f29f8e420029`

// DeployConfigReward deploys a new Ethereum contract, binding an instance of ConfigReward to it.
func DeployConfigReward(auth *bind.TransactOpts, backend bind.ContractBackend, _owners []common.Address, _required *big.Int) (common.Address, *types.Transaction, *ConfigReward, error) {
	parsed, err := abi.JSON(strings.NewReader(ConfigRewardABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ConfigRewardBin), backend, _owners, _required)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConfigReward{ConfigRewardCaller: ConfigRewardCaller{contract: contract}, ConfigRewardTransactor: ConfigRewardTransactor{contract: contract}, ConfigRewardFilterer: ConfigRewardFilterer{contract: contract}}, nil
}

// ConfigReward is an auto generated Go binding around an Ethereum contract.
type ConfigReward struct {
	ConfigRewardCaller     // Read-only binding to the contract
	ConfigRewardTransactor // Write-only binding to the contract
	ConfigRewardFilterer   // Log filterer for contract events
}

// ConfigRewardCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConfigRewardCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigRewardTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfigRewardTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigRewardFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfigRewardFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigRewardSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfigRewardSession struct {
	Contract     *ConfigReward     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConfigRewardCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfigRewardCallerSession struct {
	Contract *ConfigRewardCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ConfigRewardTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfigRewardTransactorSession struct {
	Contract     *ConfigRewardTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ConfigRewardRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConfigRewardRaw struct {
	Contract *ConfigReward // Generic contract binding to access the raw methods on
}

// ConfigRewardCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfigRewardCallerRaw struct {
	Contract *ConfigRewardCaller // Generic read-only contract binding to access the raw methods on
}

// ConfigRewardTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfigRewardTransactorRaw struct {
	Contract *ConfigRewardTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConfigReward creates a new instance of ConfigReward, bound to a specific deployed contract.
func NewConfigReward(address common.Address, backend bind.ContractBackend) (*ConfigReward, error) {
	contract, err := bindConfigReward(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConfigReward{ConfigRewardCaller: ConfigRewardCaller{contract: contract}, ConfigRewardTransactor: ConfigRewardTransactor{contract: contract}, ConfigRewardFilterer: ConfigRewardFilterer{contract: contract}}, nil
}

// NewConfigRewardCaller creates a new read-only instance of ConfigReward, bound to a specific deployed contract.
func NewConfigRewardCaller(address common.Address, caller bind.ContractCaller) (*ConfigRewardCaller, error) {
	contract, err := bindConfigReward(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfigRewardCaller{contract: contract}, nil
}

// NewConfigRewardTransactor creates a new write-only instance of ConfigReward, bound to a specific deployed contract.
func NewConfigRewardTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfigRewardTransactor, error) {
	contract, err := bindConfigReward(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfigRewardTransactor{contract: contract}, nil
}

// NewConfigRewardFilterer creates a new log filterer instance of ConfigReward, bound to a specific deployed contract.
func NewConfigRewardFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfigRewardFilterer, error) {
	contract, err := bindConfigReward(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfigRewardFilterer{contract: contract}, nil
}

// bindConfigReward binds a generic wrapper to an already deployed contract.
func bindConfigReward(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConfigRewardABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfigReward *ConfigRewardRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ConfigReward.Contract.ConfigRewardCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfigReward *ConfigRewardRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfigReward.Contract.ConfigRewardTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfigReward *ConfigRewardRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfigReward.Contract.ConfigRewardTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfigReward *ConfigRewardCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ConfigReward.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfigReward *ConfigRewardTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfigReward.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfigReward *ConfigRewardTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfigReward.Contract.contract.Transact(opts, method, params...)
}

// MAXOWNERCOUNT is a free data retrieval call binding the contract method 0xd74f8edd.
//
// Solidity: function MAX_OWNER_COUNT() constant returns(uint256)
func (_ConfigReward *ConfigRewardCaller) MAXOWNERCOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "MAX_OWNER_COUNT")
	return *ret0, err
}

// MAXOWNERCOUNT is a free data retrieval call binding the contract method 0xd74f8edd.
//
// Solidity: function MAX_OWNER_COUNT() constant returns(uint256)
func (_ConfigReward *ConfigRewardSession) MAXOWNERCOUNT() (*big.Int, error) {
	return _ConfigReward.Contract.MAXOWNERCOUNT(&_ConfigReward.CallOpts)
}

// MAXOWNERCOUNT is a free data retrieval call binding the contract method 0xd74f8edd.
//
// Solidity: function MAX_OWNER_COUNT() constant returns(uint256)
func (_ConfigReward *ConfigRewardCallerSession) MAXOWNERCOUNT() (*big.Int, error) {
	return _ConfigReward.Contract.MAXOWNERCOUNT(&_ConfigReward.CallOpts)
}

// Confirmations is a free data retrieval call binding the contract method 0x3411c81c.
//
// Solidity: function confirmations( uint256,  address) constant returns(bool)
func (_ConfigReward *ConfigRewardCaller) Confirmations(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "confirmations", arg0, arg1)
	return *ret0, err
}

// Confirmations is a free data retrieval call binding the contract method 0x3411c81c.
//
// Solidity: function confirmations( uint256,  address) constant returns(bool)
func (_ConfigReward *ConfigRewardSession) Confirmations(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _ConfigReward.Contract.Confirmations(&_ConfigReward.CallOpts, arg0, arg1)
}

// Confirmations is a free data retrieval call binding the contract method 0x3411c81c.
//
// Solidity: function confirmations( uint256,  address) constant returns(bool)
func (_ConfigReward *ConfigRewardCallerSession) Confirmations(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _ConfigReward.Contract.Confirmations(&_ConfigReward.CallOpts, arg0, arg1)
}

// GetConfirmationCount is a free data retrieval call binding the contract method 0x8b51d13f.
//
// Solidity: function getConfirmationCount(transactionId uint256) constant returns(count uint256)
func (_ConfigReward *ConfigRewardCaller) GetConfirmationCount(opts *bind.CallOpts, transactionId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "getConfirmationCount", transactionId)
	return *ret0, err
}

// GetConfirmationCount is a free data retrieval call binding the contract method 0x8b51d13f.
//
// Solidity: function getConfirmationCount(transactionId uint256) constant returns(count uint256)
func (_ConfigReward *ConfigRewardSession) GetConfirmationCount(transactionId *big.Int) (*big.Int, error) {
	return _ConfigReward.Contract.GetConfirmationCount(&_ConfigReward.CallOpts, transactionId)
}

// GetConfirmationCount is a free data retrieval call binding the contract method 0x8b51d13f.
//
// Solidity: function getConfirmationCount(transactionId uint256) constant returns(count uint256)
func (_ConfigReward *ConfigRewardCallerSession) GetConfirmationCount(transactionId *big.Int) (*big.Int, error) {
	return _ConfigReward.Contract.GetConfirmationCount(&_ConfigReward.CallOpts, transactionId)
}

// GetConfirmations is a free data retrieval call binding the contract method 0xb5dc40c3.
//
// Solidity: function getConfirmations(transactionId uint256) constant returns(_confirmations address[])
func (_ConfigReward *ConfigRewardCaller) GetConfirmations(opts *bind.CallOpts, transactionId *big.Int) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "getConfirmations", transactionId)
	return *ret0, err
}

// GetConfirmations is a free data retrieval call binding the contract method 0xb5dc40c3.
//
// Solidity: function getConfirmations(transactionId uint256) constant returns(_confirmations address[])
func (_ConfigReward *ConfigRewardSession) GetConfirmations(transactionId *big.Int) ([]common.Address, error) {
	return _ConfigReward.Contract.GetConfirmations(&_ConfigReward.CallOpts, transactionId)
}

// GetConfirmations is a free data retrieval call binding the contract method 0xb5dc40c3.
//
// Solidity: function getConfirmations(transactionId uint256) constant returns(_confirmations address[])
func (_ConfigReward *ConfigRewardCallerSession) GetConfirmations(transactionId *big.Int) ([]common.Address, error) {
	return _ConfigReward.Contract.GetConfirmations(&_ConfigReward.CallOpts, transactionId)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() constant returns(address[])
func (_ConfigReward *ConfigRewardCaller) GetOwners(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "getOwners")
	return *ret0, err
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() constant returns(address[])
func (_ConfigReward *ConfigRewardSession) GetOwners() ([]common.Address, error) {
	return _ConfigReward.Contract.GetOwners(&_ConfigReward.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() constant returns(address[])
func (_ConfigReward *ConfigRewardCallerSession) GetOwners() ([]common.Address, error) {
	return _ConfigReward.Contract.GetOwners(&_ConfigReward.CallOpts)
}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() constant returns(uint256)
func (_ConfigReward *ConfigRewardCaller) GetRate(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "getRate")
	return *ret0, err
}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() constant returns(uint256)
func (_ConfigReward *ConfigRewardSession) GetRate() (*big.Int, error) {
	return _ConfigReward.Contract.GetRate(&_ConfigReward.CallOpts)
}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() constant returns(uint256)
func (_ConfigReward *ConfigRewardCallerSession) GetRate() (*big.Int, error) {
	return _ConfigReward.Contract.GetRate(&_ConfigReward.CallOpts)
}

// GetTransactionCount is a free data retrieval call binding the contract method 0x54741525.
//
// Solidity: function getTransactionCount(pending bool, executed bool) constant returns(count uint256)
func (_ConfigReward *ConfigRewardCaller) GetTransactionCount(opts *bind.CallOpts, pending bool, executed bool) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "getTransactionCount", pending, executed)
	return *ret0, err
}

// GetTransactionCount is a free data retrieval call binding the contract method 0x54741525.
//
// Solidity: function getTransactionCount(pending bool, executed bool) constant returns(count uint256)
func (_ConfigReward *ConfigRewardSession) GetTransactionCount(pending bool, executed bool) (*big.Int, error) {
	return _ConfigReward.Contract.GetTransactionCount(&_ConfigReward.CallOpts, pending, executed)
}

// GetTransactionCount is a free data retrieval call binding the contract method 0x54741525.
//
// Solidity: function getTransactionCount(pending bool, executed bool) constant returns(count uint256)
func (_ConfigReward *ConfigRewardCallerSession) GetTransactionCount(pending bool, executed bool) (*big.Int, error) {
	return _ConfigReward.Contract.GetTransactionCount(&_ConfigReward.CallOpts, pending, executed)
}

// GetTransactionIds is a free data retrieval call binding the contract method 0xa8abe69a.
//
// Solidity: function getTransactionIds(from uint256, to uint256, pending bool, executed bool) constant returns(_transactionIds uint256[])
func (_ConfigReward *ConfigRewardCaller) GetTransactionIds(opts *bind.CallOpts, from *big.Int, to *big.Int, pending bool, executed bool) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "getTransactionIds", from, to, pending, executed)
	return *ret0, err
}

// GetTransactionIds is a free data retrieval call binding the contract method 0xa8abe69a.
//
// Solidity: function getTransactionIds(from uint256, to uint256, pending bool, executed bool) constant returns(_transactionIds uint256[])
func (_ConfigReward *ConfigRewardSession) GetTransactionIds(from *big.Int, to *big.Int, pending bool, executed bool) ([]*big.Int, error) {
	return _ConfigReward.Contract.GetTransactionIds(&_ConfigReward.CallOpts, from, to, pending, executed)
}

// GetTransactionIds is a free data retrieval call binding the contract method 0xa8abe69a.
//
// Solidity: function getTransactionIds(from uint256, to uint256, pending bool, executed bool) constant returns(_transactionIds uint256[])
func (_ConfigReward *ConfigRewardCallerSession) GetTransactionIds(from *big.Int, to *big.Int, pending bool, executed bool) ([]*big.Int, error) {
	return _ConfigReward.Contract.GetTransactionIds(&_ConfigReward.CallOpts, from, to, pending, executed)
}

// IsConfirmed is a free data retrieval call binding the contract method 0x784547a7.
//
// Solidity: function isConfirmed(transactionId uint256) constant returns(bool)
func (_ConfigReward *ConfigRewardCaller) IsConfirmed(opts *bind.CallOpts, transactionId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "isConfirmed", transactionId)
	return *ret0, err
}

// IsConfirmed is a free data retrieval call binding the contract method 0x784547a7.
//
// Solidity: function isConfirmed(transactionId uint256) constant returns(bool)
func (_ConfigReward *ConfigRewardSession) IsConfirmed(transactionId *big.Int) (bool, error) {
	return _ConfigReward.Contract.IsConfirmed(&_ConfigReward.CallOpts, transactionId)
}

// IsConfirmed is a free data retrieval call binding the contract method 0x784547a7.
//
// Solidity: function isConfirmed(transactionId uint256) constant returns(bool)
func (_ConfigReward *ConfigRewardCallerSession) IsConfirmed(transactionId *big.Int) (bool, error) {
	return _ConfigReward.Contract.IsConfirmed(&_ConfigReward.CallOpts, transactionId)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner( address) constant returns(bool)
func (_ConfigReward *ConfigRewardCaller) IsOwner(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "isOwner", arg0)
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner( address) constant returns(bool)
func (_ConfigReward *ConfigRewardSession) IsOwner(arg0 common.Address) (bool, error) {
	return _ConfigReward.Contract.IsOwner(&_ConfigReward.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner( address) constant returns(bool)
func (_ConfigReward *ConfigRewardCallerSession) IsOwner(arg0 common.Address) (bool, error) {
	return _ConfigReward.Contract.IsOwner(&_ConfigReward.CallOpts, arg0)
}

// Owners is a free data retrieval call binding the contract method 0x025e7c27.
//
// Solidity: function owners( uint256) constant returns(address)
func (_ConfigReward *ConfigRewardCaller) Owners(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "owners", arg0)
	return *ret0, err
}

// Owners is a free data retrieval call binding the contract method 0x025e7c27.
//
// Solidity: function owners( uint256) constant returns(address)
func (_ConfigReward *ConfigRewardSession) Owners(arg0 *big.Int) (common.Address, error) {
	return _ConfigReward.Contract.Owners(&_ConfigReward.CallOpts, arg0)
}

// Owners is a free data retrieval call binding the contract method 0x025e7c27.
//
// Solidity: function owners( uint256) constant returns(address)
func (_ConfigReward *ConfigRewardCallerSession) Owners(arg0 *big.Int) (common.Address, error) {
	return _ConfigReward.Contract.Owners(&_ConfigReward.CallOpts, arg0)
}

// Required is a free data retrieval call binding the contract method 0xdc8452cd.
//
// Solidity: function required() constant returns(uint256)
func (_ConfigReward *ConfigRewardCaller) Required(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "required")
	return *ret0, err
}

// Required is a free data retrieval call binding the contract method 0xdc8452cd.
//
// Solidity: function required() constant returns(uint256)
func (_ConfigReward *ConfigRewardSession) Required() (*big.Int, error) {
	return _ConfigReward.Contract.Required(&_ConfigReward.CallOpts)
}

// Required is a free data retrieval call binding the contract method 0xdc8452cd.
//
// Solidity: function required() constant returns(uint256)
func (_ConfigReward *ConfigRewardCallerSession) Required() (*big.Int, error) {
	return _ConfigReward.Contract.Required(&_ConfigReward.CallOpts)
}

// TransactionCount is a free data retrieval call binding the contract method 0xb77bf600.
//
// Solidity: function transactionCount() constant returns(uint256)
func (_ConfigReward *ConfigRewardCaller) TransactionCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConfigReward.contract.Call(opts, out, "transactionCount")
	return *ret0, err
}

// TransactionCount is a free data retrieval call binding the contract method 0xb77bf600.
//
// Solidity: function transactionCount() constant returns(uint256)
func (_ConfigReward *ConfigRewardSession) TransactionCount() (*big.Int, error) {
	return _ConfigReward.Contract.TransactionCount(&_ConfigReward.CallOpts)
}

// TransactionCount is a free data retrieval call binding the contract method 0xb77bf600.
//
// Solidity: function transactionCount() constant returns(uint256)
func (_ConfigReward *ConfigRewardCallerSession) TransactionCount() (*big.Int, error) {
	return _ConfigReward.Contract.TransactionCount(&_ConfigReward.CallOpts)
}

// Transactions is a free data retrieval call binding the contract method 0x9ace38c2.
//
// Solidity: function transactions( uint256) constant returns(rate uint256, executed bool)
func (_ConfigReward *ConfigRewardCaller) Transactions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Rate     *big.Int
	Executed bool
}, error) {
	ret := new(struct {
		Rate     *big.Int
		Executed bool
	})
	out := ret
	err := _ConfigReward.contract.Call(opts, out, "transactions", arg0)
	return *ret, err
}

// Transactions is a free data retrieval call binding the contract method 0x9ace38c2.
//
// Solidity: function transactions( uint256) constant returns(rate uint256, executed bool)
func (_ConfigReward *ConfigRewardSession) Transactions(arg0 *big.Int) (struct {
	Rate     *big.Int
	Executed bool
}, error) {
	return _ConfigReward.Contract.Transactions(&_ConfigReward.CallOpts, arg0)
}

// Transactions is a free data retrieval call binding the contract method 0x9ace38c2.
//
// Solidity: function transactions( uint256) constant returns(rate uint256, executed bool)
func (_ConfigReward *ConfigRewardCallerSession) Transactions(arg0 *big.Int) (struct {
	Rate     *big.Int
	Executed bool
}, error) {
	return _ConfigReward.Contract.Transactions(&_ConfigReward.CallOpts, arg0)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(owner address) returns()
func (_ConfigReward *ConfigRewardTransactor) AddOwner(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _ConfigReward.contract.Transact(opts, "addOwner", owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(owner address) returns()
func (_ConfigReward *ConfigRewardSession) AddOwner(owner common.Address) (*types.Transaction, error) {
	return _ConfigReward.Contract.AddOwner(&_ConfigReward.TransactOpts, owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(owner address) returns()
func (_ConfigReward *ConfigRewardTransactorSession) AddOwner(owner common.Address) (*types.Transaction, error) {
	return _ConfigReward.Contract.AddOwner(&_ConfigReward.TransactOpts, owner)
}

// ChangeRequirement is a paid mutator transaction binding the contract method 0xba51a6df.
//
// Solidity: function changeRequirement(_required uint256) returns()
func (_ConfigReward *ConfigRewardTransactor) ChangeRequirement(opts *bind.TransactOpts, _required *big.Int) (*types.Transaction, error) {
	return _ConfigReward.contract.Transact(opts, "changeRequirement", _required)
}

// ChangeRequirement is a paid mutator transaction binding the contract method 0xba51a6df.
//
// Solidity: function changeRequirement(_required uint256) returns()
func (_ConfigReward *ConfigRewardSession) ChangeRequirement(_required *big.Int) (*types.Transaction, error) {
	return _ConfigReward.Contract.ChangeRequirement(&_ConfigReward.TransactOpts, _required)
}

// ChangeRequirement is a paid mutator transaction binding the contract method 0xba51a6df.
//
// Solidity: function changeRequirement(_required uint256) returns()
func (_ConfigReward *ConfigRewardTransactorSession) ChangeRequirement(_required *big.Int) (*types.Transaction, error) {
	return _ConfigReward.Contract.ChangeRequirement(&_ConfigReward.TransactOpts, _required)
}

// ConfirmTransaction is a paid mutator transaction binding the contract method 0xc01a8c84.
//
// Solidity: function confirmTransaction(transactionId uint256) returns()
func (_ConfigReward *ConfigRewardTransactor) ConfirmTransaction(opts *bind.TransactOpts, transactionId *big.Int) (*types.Transaction, error) {
	return _ConfigReward.contract.Transact(opts, "confirmTransaction", transactionId)
}

// ConfirmTransaction is a paid mutator transaction binding the contract method 0xc01a8c84.
//
// Solidity: function confirmTransaction(transactionId uint256) returns()
func (_ConfigReward *ConfigRewardSession) ConfirmTransaction(transactionId *big.Int) (*types.Transaction, error) {
	return _ConfigReward.Contract.ConfirmTransaction(&_ConfigReward.TransactOpts, transactionId)
}

// ConfirmTransaction is a paid mutator transaction binding the contract method 0xc01a8c84.
//
// Solidity: function confirmTransaction(transactionId uint256) returns()
func (_ConfigReward *ConfigRewardTransactorSession) ConfirmTransaction(transactionId *big.Int) (*types.Transaction, error) {
	return _ConfigReward.Contract.ConfirmTransaction(&_ConfigReward.TransactOpts, transactionId)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0xee22610b.
//
// Solidity: function executeTransaction(transactionId uint256) returns()
func (_ConfigReward *ConfigRewardTransactor) ExecuteTransaction(opts *bind.TransactOpts, transactionId *big.Int) (*types.Transaction, error) {
	return _ConfigReward.contract.Transact(opts, "executeTransaction", transactionId)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0xee22610b.
//
// Solidity: function executeTransaction(transactionId uint256) returns()
func (_ConfigReward *ConfigRewardSession) ExecuteTransaction(transactionId *big.Int) (*types.Transaction, error) {
	return _ConfigReward.Contract.ExecuteTransaction(&_ConfigReward.TransactOpts, transactionId)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0xee22610b.
//
// Solidity: function executeTransaction(transactionId uint256) returns()
func (_ConfigReward *ConfigRewardTransactorSession) ExecuteTransaction(transactionId *big.Int) (*types.Transaction, error) {
	return _ConfigReward.Contract.ExecuteTransaction(&_ConfigReward.TransactOpts, transactionId)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(owner address) returns()
func (_ConfigReward *ConfigRewardTransactor) RemoveOwner(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _ConfigReward.contract.Transact(opts, "removeOwner", owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(owner address) returns()
func (_ConfigReward *ConfigRewardSession) RemoveOwner(owner common.Address) (*types.Transaction, error) {
	return _ConfigReward.Contract.RemoveOwner(&_ConfigReward.TransactOpts, owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(owner address) returns()
func (_ConfigReward *ConfigRewardTransactorSession) RemoveOwner(owner common.Address) (*types.Transaction, error) {
	return _ConfigReward.Contract.RemoveOwner(&_ConfigReward.TransactOpts, owner)
}

// ReplaceOwner is a paid mutator transaction binding the contract method 0xe20056e6.
//
// Solidity: function replaceOwner(owner address, newOwner address) returns()
func (_ConfigReward *ConfigRewardTransactor) ReplaceOwner(opts *bind.TransactOpts, owner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _ConfigReward.contract.Transact(opts, "replaceOwner", owner, newOwner)
}

// ReplaceOwner is a paid mutator transaction binding the contract method 0xe20056e6.
//
// Solidity: function replaceOwner(owner address, newOwner address) returns()
func (_ConfigReward *ConfigRewardSession) ReplaceOwner(owner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _ConfigReward.Contract.ReplaceOwner(&_ConfigReward.TransactOpts, owner, newOwner)
}

// ReplaceOwner is a paid mutator transaction binding the contract method 0xe20056e6.
//
// Solidity: function replaceOwner(owner address, newOwner address) returns()
func (_ConfigReward *ConfigRewardTransactorSession) ReplaceOwner(owner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _ConfigReward.Contract.ReplaceOwner(&_ConfigReward.TransactOpts, owner, newOwner)
}

// RevokeConfirmation is a paid mutator transaction binding the contract method 0x20ea8d86.
//
// Solidity: function revokeConfirmation(transactionId uint256) returns()
func (_ConfigReward *ConfigRewardTransactor) RevokeConfirmation(opts *bind.TransactOpts, transactionId *big.Int) (*types.Transaction, error) {
	return _ConfigReward.contract.Transact(opts, "revokeConfirmation", transactionId)
}

// RevokeConfirmation is a paid mutator transaction binding the contract method 0x20ea8d86.
//
// Solidity: function revokeConfirmation(transactionId uint256) returns()
func (_ConfigReward *ConfigRewardSession) RevokeConfirmation(transactionId *big.Int) (*types.Transaction, error) {
	return _ConfigReward.Contract.RevokeConfirmation(&_ConfigReward.TransactOpts, transactionId)
}

// RevokeConfirmation is a paid mutator transaction binding the contract method 0x20ea8d86.
//
// Solidity: function revokeConfirmation(transactionId uint256) returns()
func (_ConfigReward *ConfigRewardTransactorSession) RevokeConfirmation(transactionId *big.Int) (*types.Transaction, error) {
	return _ConfigReward.Contract.RevokeConfirmation(&_ConfigReward.TransactOpts, transactionId)
}

// SubmitTransaction is a paid mutator transaction binding the contract method 0xc34a0ead.
//
// Solidity: function submitTransaction(rate uint256) returns(transactionId uint256)
func (_ConfigReward *ConfigRewardTransactor) SubmitTransaction(opts *bind.TransactOpts, rate *big.Int) (*types.Transaction, error) {
	return _ConfigReward.contract.Transact(opts, "submitTransaction", rate)
}

// SubmitTransaction is a paid mutator transaction binding the contract method 0xc34a0ead.
//
// Solidity: function submitTransaction(rate uint256) returns(transactionId uint256)
func (_ConfigReward *ConfigRewardSession) SubmitTransaction(rate *big.Int) (*types.Transaction, error) {
	return _ConfigReward.Contract.SubmitTransaction(&_ConfigReward.TransactOpts, rate)
}

// SubmitTransaction is a paid mutator transaction binding the contract method 0xc34a0ead.
//
// Solidity: function submitTransaction(rate uint256) returns(transactionId uint256)
func (_ConfigReward *ConfigRewardTransactorSession) SubmitTransaction(rate *big.Int) (*types.Transaction, error) {
	return _ConfigReward.Contract.SubmitTransaction(&_ConfigReward.TransactOpts, rate)
}

// ConfigRewardConfirmationIterator is returned from FilterConfirmation and is used to iterate over the raw logs and unpacked data for Confirmation events raised by the ConfigReward contract.
type ConfigRewardConfirmationIterator struct {
	Event *ConfigRewardConfirmation // Event containing the contract specifics and raw log

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
func (it *ConfigRewardConfirmationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfigRewardConfirmation)
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
		it.Event = new(ConfigRewardConfirmation)
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
func (it *ConfigRewardConfirmationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfigRewardConfirmationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfigRewardConfirmation represents a Confirmation event raised by the ConfigReward contract.
type ConfigRewardConfirmation struct {
	Sender        common.Address
	TransactionId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterConfirmation is a free log retrieval operation binding the contract event 0x4a504a94899432a9846e1aa406dceb1bcfd538bb839071d49d1e5e23f5be30ef.
//
// Solidity: event Confirmation(sender indexed address, transactionId indexed uint256)
func (_ConfigReward *ConfigRewardFilterer) FilterConfirmation(opts *bind.FilterOpts, sender []common.Address, transactionId []*big.Int) (*ConfigRewardConfirmationIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _ConfigReward.contract.FilterLogs(opts, "Confirmation", senderRule, transactionIdRule)
	if err != nil {
		return nil, err
	}
	return &ConfigRewardConfirmationIterator{contract: _ConfigReward.contract, event: "Confirmation", logs: logs, sub: sub}, nil
}

// WatchConfirmation is a free log subscription operation binding the contract event 0x4a504a94899432a9846e1aa406dceb1bcfd538bb839071d49d1e5e23f5be30ef.
//
// Solidity: event Confirmation(sender indexed address, transactionId indexed uint256)
func (_ConfigReward *ConfigRewardFilterer) WatchConfirmation(opts *bind.WatchOpts, sink chan<- *ConfigRewardConfirmation, sender []common.Address, transactionId []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _ConfigReward.contract.WatchLogs(opts, "Confirmation", senderRule, transactionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfigRewardConfirmation)
				if err := _ConfigReward.contract.UnpackLog(event, "Confirmation", log); err != nil {
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

// ConfigRewardDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the ConfigReward contract.
type ConfigRewardDepositIterator struct {
	Event *ConfigRewardDeposit // Event containing the contract specifics and raw log

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
func (it *ConfigRewardDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfigRewardDeposit)
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
		it.Event = new(ConfigRewardDeposit)
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
func (it *ConfigRewardDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfigRewardDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfigRewardDeposit represents a Deposit event raised by the ConfigReward contract.
type ConfigRewardDeposit struct {
	Sender common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(sender indexed address, value uint256)
func (_ConfigReward *ConfigRewardFilterer) FilterDeposit(opts *bind.FilterOpts, sender []common.Address) (*ConfigRewardDepositIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ConfigReward.contract.FilterLogs(opts, "Deposit", senderRule)
	if err != nil {
		return nil, err
	}
	return &ConfigRewardDepositIterator{contract: _ConfigReward.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(sender indexed address, value uint256)
func (_ConfigReward *ConfigRewardFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *ConfigRewardDeposit, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ConfigReward.contract.WatchLogs(opts, "Deposit", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfigRewardDeposit)
				if err := _ConfigReward.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ConfigRewardExecutionIterator is returned from FilterExecution and is used to iterate over the raw logs and unpacked data for Execution events raised by the ConfigReward contract.
type ConfigRewardExecutionIterator struct {
	Event *ConfigRewardExecution // Event containing the contract specifics and raw log

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
func (it *ConfigRewardExecutionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfigRewardExecution)
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
		it.Event = new(ConfigRewardExecution)
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
func (it *ConfigRewardExecutionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfigRewardExecutionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfigRewardExecution represents a Execution event raised by the ConfigReward contract.
type ConfigRewardExecution struct {
	TransactionId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterExecution is a free log retrieval operation binding the contract event 0x33e13ecb54c3076d8e8bb8c2881800a4d972b792045ffae98fdf46df365fed75.
//
// Solidity: event Execution(transactionId indexed uint256)
func (_ConfigReward *ConfigRewardFilterer) FilterExecution(opts *bind.FilterOpts, transactionId []*big.Int) (*ConfigRewardExecutionIterator, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _ConfigReward.contract.FilterLogs(opts, "Execution", transactionIdRule)
	if err != nil {
		return nil, err
	}
	return &ConfigRewardExecutionIterator{contract: _ConfigReward.contract, event: "Execution", logs: logs, sub: sub}, nil
}

// WatchExecution is a free log subscription operation binding the contract event 0x33e13ecb54c3076d8e8bb8c2881800a4d972b792045ffae98fdf46df365fed75.
//
// Solidity: event Execution(transactionId indexed uint256)
func (_ConfigReward *ConfigRewardFilterer) WatchExecution(opts *bind.WatchOpts, sink chan<- *ConfigRewardExecution, transactionId []*big.Int) (event.Subscription, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _ConfigReward.contract.WatchLogs(opts, "Execution", transactionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfigRewardExecution)
				if err := _ConfigReward.contract.UnpackLog(event, "Execution", log); err != nil {
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

// ConfigRewardExecutionFailureIterator is returned from FilterExecutionFailure and is used to iterate over the raw logs and unpacked data for ExecutionFailure events raised by the ConfigReward contract.
type ConfigRewardExecutionFailureIterator struct {
	Event *ConfigRewardExecutionFailure // Event containing the contract specifics and raw log

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
func (it *ConfigRewardExecutionFailureIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfigRewardExecutionFailure)
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
		it.Event = new(ConfigRewardExecutionFailure)
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
func (it *ConfigRewardExecutionFailureIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfigRewardExecutionFailureIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfigRewardExecutionFailure represents a ExecutionFailure event raised by the ConfigReward contract.
type ConfigRewardExecutionFailure struct {
	TransactionId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterExecutionFailure is a free log retrieval operation binding the contract event 0x526441bb6c1aba3c9a4a6ca1d6545da9c2333c8c48343ef398eb858d72b79236.
//
// Solidity: event ExecutionFailure(transactionId indexed uint256)
func (_ConfigReward *ConfigRewardFilterer) FilterExecutionFailure(opts *bind.FilterOpts, transactionId []*big.Int) (*ConfigRewardExecutionFailureIterator, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _ConfigReward.contract.FilterLogs(opts, "ExecutionFailure", transactionIdRule)
	if err != nil {
		return nil, err
	}
	return &ConfigRewardExecutionFailureIterator{contract: _ConfigReward.contract, event: "ExecutionFailure", logs: logs, sub: sub}, nil
}

// WatchExecutionFailure is a free log subscription operation binding the contract event 0x526441bb6c1aba3c9a4a6ca1d6545da9c2333c8c48343ef398eb858d72b79236.
//
// Solidity: event ExecutionFailure(transactionId indexed uint256)
func (_ConfigReward *ConfigRewardFilterer) WatchExecutionFailure(opts *bind.WatchOpts, sink chan<- *ConfigRewardExecutionFailure, transactionId []*big.Int) (event.Subscription, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _ConfigReward.contract.WatchLogs(opts, "ExecutionFailure", transactionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfigRewardExecutionFailure)
				if err := _ConfigReward.contract.UnpackLog(event, "ExecutionFailure", log); err != nil {
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

// ConfigRewardOwnerAdditionIterator is returned from FilterOwnerAddition and is used to iterate over the raw logs and unpacked data for OwnerAddition events raised by the ConfigReward contract.
type ConfigRewardOwnerAdditionIterator struct {
	Event *ConfigRewardOwnerAddition // Event containing the contract specifics and raw log

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
func (it *ConfigRewardOwnerAdditionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfigRewardOwnerAddition)
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
		it.Event = new(ConfigRewardOwnerAddition)
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
func (it *ConfigRewardOwnerAdditionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfigRewardOwnerAdditionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfigRewardOwnerAddition represents a OwnerAddition event raised by the ConfigReward contract.
type ConfigRewardOwnerAddition struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOwnerAddition is a free log retrieval operation binding the contract event 0xf39e6e1eb0edcf53c221607b54b00cd28f3196fed0a24994dc308b8f611b682d.
//
// Solidity: event OwnerAddition(owner indexed address)
func (_ConfigReward *ConfigRewardFilterer) FilterOwnerAddition(opts *bind.FilterOpts, owner []common.Address) (*ConfigRewardOwnerAdditionIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ConfigReward.contract.FilterLogs(opts, "OwnerAddition", ownerRule)
	if err != nil {
		return nil, err
	}
	return &ConfigRewardOwnerAdditionIterator{contract: _ConfigReward.contract, event: "OwnerAddition", logs: logs, sub: sub}, nil
}

// WatchOwnerAddition is a free log subscription operation binding the contract event 0xf39e6e1eb0edcf53c221607b54b00cd28f3196fed0a24994dc308b8f611b682d.
//
// Solidity: event OwnerAddition(owner indexed address)
func (_ConfigReward *ConfigRewardFilterer) WatchOwnerAddition(opts *bind.WatchOpts, sink chan<- *ConfigRewardOwnerAddition, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ConfigReward.contract.WatchLogs(opts, "OwnerAddition", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfigRewardOwnerAddition)
				if err := _ConfigReward.contract.UnpackLog(event, "OwnerAddition", log); err != nil {
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

// ConfigRewardOwnerRemovalIterator is returned from FilterOwnerRemoval and is used to iterate over the raw logs and unpacked data for OwnerRemoval events raised by the ConfigReward contract.
type ConfigRewardOwnerRemovalIterator struct {
	Event *ConfigRewardOwnerRemoval // Event containing the contract specifics and raw log

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
func (it *ConfigRewardOwnerRemovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfigRewardOwnerRemoval)
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
		it.Event = new(ConfigRewardOwnerRemoval)
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
func (it *ConfigRewardOwnerRemovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfigRewardOwnerRemovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfigRewardOwnerRemoval represents a OwnerRemoval event raised by the ConfigReward contract.
type ConfigRewardOwnerRemoval struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOwnerRemoval is a free log retrieval operation binding the contract event 0x8001553a916ef2f495d26a907cc54d96ed840d7bda71e73194bf5a9df7a76b90.
//
// Solidity: event OwnerRemoval(owner indexed address)
func (_ConfigReward *ConfigRewardFilterer) FilterOwnerRemoval(opts *bind.FilterOpts, owner []common.Address) (*ConfigRewardOwnerRemovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ConfigReward.contract.FilterLogs(opts, "OwnerRemoval", ownerRule)
	if err != nil {
		return nil, err
	}
	return &ConfigRewardOwnerRemovalIterator{contract: _ConfigReward.contract, event: "OwnerRemoval", logs: logs, sub: sub}, nil
}

// WatchOwnerRemoval is a free log subscription operation binding the contract event 0x8001553a916ef2f495d26a907cc54d96ed840d7bda71e73194bf5a9df7a76b90.
//
// Solidity: event OwnerRemoval(owner indexed address)
func (_ConfigReward *ConfigRewardFilterer) WatchOwnerRemoval(opts *bind.WatchOpts, sink chan<- *ConfigRewardOwnerRemoval, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ConfigReward.contract.WatchLogs(opts, "OwnerRemoval", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfigRewardOwnerRemoval)
				if err := _ConfigReward.contract.UnpackLog(event, "OwnerRemoval", log); err != nil {
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

// ConfigRewardRequirementChangeIterator is returned from FilterRequirementChange and is used to iterate over the raw logs and unpacked data for RequirementChange events raised by the ConfigReward contract.
type ConfigRewardRequirementChangeIterator struct {
	Event *ConfigRewardRequirementChange // Event containing the contract specifics and raw log

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
func (it *ConfigRewardRequirementChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfigRewardRequirementChange)
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
		it.Event = new(ConfigRewardRequirementChange)
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
func (it *ConfigRewardRequirementChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfigRewardRequirementChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfigRewardRequirementChange represents a RequirementChange event raised by the ConfigReward contract.
type ConfigRewardRequirementChange struct {
	Required *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRequirementChange is a free log retrieval operation binding the contract event 0xa3f1ee9126a074d9326c682f561767f710e927faa811f7a99829d49dc421797a.
//
// Solidity: event RequirementChange(required uint256)
func (_ConfigReward *ConfigRewardFilterer) FilterRequirementChange(opts *bind.FilterOpts) (*ConfigRewardRequirementChangeIterator, error) {

	logs, sub, err := _ConfigReward.contract.FilterLogs(opts, "RequirementChange")
	if err != nil {
		return nil, err
	}
	return &ConfigRewardRequirementChangeIterator{contract: _ConfigReward.contract, event: "RequirementChange", logs: logs, sub: sub}, nil
}

// WatchRequirementChange is a free log subscription operation binding the contract event 0xa3f1ee9126a074d9326c682f561767f710e927faa811f7a99829d49dc421797a.
//
// Solidity: event RequirementChange(required uint256)
func (_ConfigReward *ConfigRewardFilterer) WatchRequirementChange(opts *bind.WatchOpts, sink chan<- *ConfigRewardRequirementChange) (event.Subscription, error) {

	logs, sub, err := _ConfigReward.contract.WatchLogs(opts, "RequirementChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfigRewardRequirementChange)
				if err := _ConfigReward.contract.UnpackLog(event, "RequirementChange", log); err != nil {
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

// ConfigRewardRevocationIterator is returned from FilterRevocation and is used to iterate over the raw logs and unpacked data for Revocation events raised by the ConfigReward contract.
type ConfigRewardRevocationIterator struct {
	Event *ConfigRewardRevocation // Event containing the contract specifics and raw log

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
func (it *ConfigRewardRevocationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfigRewardRevocation)
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
		it.Event = new(ConfigRewardRevocation)
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
func (it *ConfigRewardRevocationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfigRewardRevocationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfigRewardRevocation represents a Revocation event raised by the ConfigReward contract.
type ConfigRewardRevocation struct {
	Sender        common.Address
	TransactionId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRevocation is a free log retrieval operation binding the contract event 0xf6a317157440607f36269043eb55f1287a5a19ba2216afeab88cd46cbcfb88e9.
//
// Solidity: event Revocation(sender indexed address, transactionId indexed uint256)
func (_ConfigReward *ConfigRewardFilterer) FilterRevocation(opts *bind.FilterOpts, sender []common.Address, transactionId []*big.Int) (*ConfigRewardRevocationIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _ConfigReward.contract.FilterLogs(opts, "Revocation", senderRule, transactionIdRule)
	if err != nil {
		return nil, err
	}
	return &ConfigRewardRevocationIterator{contract: _ConfigReward.contract, event: "Revocation", logs: logs, sub: sub}, nil
}

// WatchRevocation is a free log subscription operation binding the contract event 0xf6a317157440607f36269043eb55f1287a5a19ba2216afeab88cd46cbcfb88e9.
//
// Solidity: event Revocation(sender indexed address, transactionId indexed uint256)
func (_ConfigReward *ConfigRewardFilterer) WatchRevocation(opts *bind.WatchOpts, sink chan<- *ConfigRewardRevocation, sender []common.Address, transactionId []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _ConfigReward.contract.WatchLogs(opts, "Revocation", senderRule, transactionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfigRewardRevocation)
				if err := _ConfigReward.contract.UnpackLog(event, "Revocation", log); err != nil {
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

// ConfigRewardSubmissionIterator is returned from FilterSubmission and is used to iterate over the raw logs and unpacked data for Submission events raised by the ConfigReward contract.
type ConfigRewardSubmissionIterator struct {
	Event *ConfigRewardSubmission // Event containing the contract specifics and raw log

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
func (it *ConfigRewardSubmissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfigRewardSubmission)
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
		it.Event = new(ConfigRewardSubmission)
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
func (it *ConfigRewardSubmissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfigRewardSubmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfigRewardSubmission represents a Submission event raised by the ConfigReward contract.
type ConfigRewardSubmission struct {
	TransactionId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSubmission is a free log retrieval operation binding the contract event 0xc0ba8fe4b176c1714197d43b9cc6bcf797a4a7461c5fe8d0ef6e184ae7601e51.
//
// Solidity: event Submission(transactionId indexed uint256)
func (_ConfigReward *ConfigRewardFilterer) FilterSubmission(opts *bind.FilterOpts, transactionId []*big.Int) (*ConfigRewardSubmissionIterator, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _ConfigReward.contract.FilterLogs(opts, "Submission", transactionIdRule)
	if err != nil {
		return nil, err
	}
	return &ConfigRewardSubmissionIterator{contract: _ConfigReward.contract, event: "Submission", logs: logs, sub: sub}, nil
}

// WatchSubmission is a free log subscription operation binding the contract event 0xc0ba8fe4b176c1714197d43b9cc6bcf797a4a7461c5fe8d0ef6e184ae7601e51.
//
// Solidity: event Submission(transactionId indexed uint256)
func (_ConfigReward *ConfigRewardFilterer) WatchSubmission(opts *bind.WatchOpts, sink chan<- *ConfigRewardSubmission, transactionId []*big.Int) (event.Subscription, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _ConfigReward.contract.WatchLogs(opts, "Submission", transactionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfigRewardSubmission)
				if err := _ConfigReward.contract.UnpackLog(event, "Submission", log); err != nil {
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
