package test

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/test/contract"
)

type Contract2 struct {
	*contract.Contract2Session
	contractBackend bind.ContractBackend
}

func NewContract2(transactOpts *bind.TransactOpts, contractAddr common.Address, contractBackend bind.ContractBackend) (*Contract2, error) {
	smcInstance, err := contract.NewContract2(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &Contract2{
		&contract.Contract2Session{
			Contract:     smcInstance,
			TransactOpts: *transactOpts,
		},
		contractBackend,
	}, nil
}

func DeployContract2(transactOpts *bind.TransactOpts, contractBackend bind.ContractBackend) (common.Address, *Contract2, error) {
	smcAddr, _, _, err := contract.DeployContract2(transactOpts, contractBackend)
	if err != nil {
		return smcAddr, nil, err
	}

	smcInstance, err := NewContract2(transactOpts, smcAddr, contractBackend)
	if err != nil {
		return smcAddr, nil, err
	}

	return smcAddr, smcInstance, nil
}
