package test

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/test/contract"
)

type Contract1 struct {
	*contract.Contract1Session
	contractBackend bind.ContractBackend
}

func NewContract1(transactOpts *bind.TransactOpts, contractAddr common.Address, contractBackend bind.ContractBackend) (*Contract1, error) {
	smcInstance, err := contract.NewContract1(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &Contract1{
		&contract.Contract1Session{
			Contract:     smcInstance,
			TransactOpts: *transactOpts,
		},
		contractBackend,
	}, nil
}

func DeployContract1(transactOpts *bind.TransactOpts, contractBackend bind.ContractBackend) (common.Address, *Contract1, error) {
	smcAddr, _, _, err := contract.DeployContract1(transactOpts, contractBackend)
	if err != nil {
		return smcAddr, nil, err
	}

	smcInstance, err := NewContract1(transactOpts, smcAddr, contractBackend)
	if err != nil {
		return smcAddr, nil, err
	}

	return smcAddr, smcInstance, nil
}
