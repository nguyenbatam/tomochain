// Copyright (c) 2018 Tomochain
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package configreward

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/configreward/contract"

	"math/big"
)

type ConfigReward struct {
	*contract.ConfigRewardSession
	contractBackend bind.ContractBackend
}

func NewConfigReward(transactOpts *bind.TransactOpts, contractAddr common.Address, contractBackend bind.ContractBackend) (*ConfigReward, error) {
	smcInstance, err := contract.NewConfigReward(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &ConfigReward{
		&contract.ConfigRewardSession{
			Contract:     smcInstance,
			TransactOpts: *transactOpts,
		},
		contractBackend,
	}, nil
}

func DeployConfigReward(transactOpts *bind.TransactOpts, contractBackend bind.ContractBackend, _owners []common.Address, _required *big.Int) (common.Address, *ConfigReward, error) {
	smcAddr, _, _, err := contract.DeployConfigReward(transactOpts, contractBackend, _owners, _required)
	if err != nil {
		return smcAddr, nil, err
	}

	smcInstance, err := NewConfigReward(transactOpts, smcAddr, contractBackend)
	if err != nil {
		return smcAddr, nil, err
	}

	return smcAddr, smcInstance, nil
}
