package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/trc21issuer/simulation"
	"github.com/ethereum/go-ethereum/contracts/validator"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"sort"
)

var (
	MainKey, _           = crypto.HexToECDSA("65ec4d4dfbcac594a14c36baa462d6f73cd86134840f6cf7b80a1e1cd33473e2")
	MainAddr             = crypto.PubkeyToAddress(MainKey.PublicKey) //0x17F2beD710ba50Ed27aEa52fc4bD7Bda5ED4a037
	slotValidatorMapping = map[string]uint64{
		"withdrawsState":         0,
		"validatorsState":        1,
		"voters":                 2,
		"candidates":             3,
		"candidateCount":         4,
		"minCandidateCap":        5,
		"minVoterCap":            6,
		"maxValidatorNumber":     7,
		"candidateWithdrawDelay": 8,
		"voterWithdrawDelay":     9,
	}
)

type Withdraw struct {
	number uint64
	cap    uint64
}

func main() {
	_1Ether := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(100000000000000000)) // 1 TOMO
	client, err := ethclient.Dial("http://127.0.0.1:8501")
	if err != nil {
		fmt.Println(err, client)
	}
	currentBlock, _ := client.BlockByNumber(context.Background(), nil)
	fmt.Println(currentBlock.Number())
	nonce, _ := client.NonceAt(context.Background(), simulation.MainAddr, nil)
	mainAccount := bind.NewKeyedTransactor(simulation.MainKey)
	mainAccount.Nonce = big.NewInt(int64(nonce))
	mainAccount.Value = big.NewInt(0)      // in wei
	mainAccount.GasLimit = uint64(4000000) // in units
	mainAccount.GasPrice = big.NewInt(21000)
	zero := big.NewInt(0)
	validatorInstance, _ := validator.NewValidator(mainAccount, common.HexToAddress(common.MasternodeVotingSMC), client)
	candidates, _ := validatorInstance.GetCandidates()

	unvotes := map[common.Address]bool{}
	for _, candidate := range candidates {
		voters, _ := validatorInstance.GetVoters(candidate)
		for _, voter := range voters {
			cap, _ := validatorInstance.GetVoterCap(candidate, voter)
			if cap.Cmp(zero) == 0 {
				unvotes[voter] = true
			}
		}
	}
	withdraws := map[uint64]uint64{}
	for unvote, _ := range unvotes {
		slot := slotValidatorMapping["withdrawsState"]
		locState := state.GetLocMappingAtKey(unvote.Hash(), slot)
		locState = locState.Add(locState, new(big.Int).SetUint64(uint64(1)))
		data, err := client.StorageAt(context.Background(), common.HexToAddress(common.MasternodeVotingSMC), common.BigToHash(locState), nil)
		if err != nil {
			fmt.Println(err)
		}
		arrLength := common.BytesToHash(data)
		for i := uint64(0); i < arrLength.Big().Uint64(); i++ {
			key := state.GetLocDynamicArrAtElement(common.BigToHash(locState), i, 1)
			data, _ = client.StorageAt(context.Background(), common.HexToAddress(common.MasternodeVotingSMC), key, nil)
			number := common.BytesToHash(data).Big()
			if number.Uint64() > 0 {
				slot := slotValidatorMapping["withdrawsState"]
				locState := state.GetLocMappingAtKey(unvote.Hash(), slot)
				locCandidateVoters := locState.Add(locState, new(big.Int).SetUint64(uint64(0)))
				retByte := crypto.Keccak256(data, common.BigToHash(locCandidateVoters).Bytes())
				data, _ = client.StorageAt(context.Background(), common.HexToAddress(common.MasternodeVotingSMC), common.BytesToHash(retByte), nil)
				cap := common.BytesToHash(data).Big()
				cap = cap.Div(cap, _1Ether)
				if cap.Cmp(zero) > 0 {
					old := withdraws[number.Uint64()]
					withdraws[number.Uint64()] = old + cap.Uint64()
				}
			}
		}
	}
	timeWithdraw := map[uint64]uint64{}
	for key, value := range withdraws {
		index := uint64(0)
		if key > currentBlock.Number().Uint64() {
			index = (key - currentBlock.Number().Uint64()) / 1800
		}
		old := timeWithdraw[index]
		timeWithdraw[index] = old + value
	}

	list := []Withdraw{}
	for key, value := range timeWithdraw {
		list = append(list, Withdraw{key, value})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].number < list[j].number
	})
	total := -list[0].cap
	for _, v := range list {
		total = total + v.cap
		fmt.Println(v.number, v.cap, total)
	}
}
