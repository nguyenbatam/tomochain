package configreward

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"
)

var (
	key, _      = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	addr        = crypto.PubkeyToAddress(key.PublicKey)
	acc1Key, _  = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	acc2Key, _  = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	acc3Key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	acc4Key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee04aefe388d1e14474d32c45c72ce7b7a")
	acc1Addr    = crypto.PubkeyToAddress(acc1Key.PublicKey)
	acc2Addr    = crypto.PubkeyToAddress(acc2Key.PublicKey)
	acc3Addr    = crypto.PubkeyToAddress(acc3Key.PublicKey)
	acc4Addr    = crypto.PubkeyToAddress(acc4Key.PublicKey)
	defaultRate = big.NewInt(10)
	changeRate  = big.NewInt(15)
)

func TestValidator(t *testing.T) {
	contractBackend := backends.NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: big.NewInt(1000000000000000)},
		acc1Addr: {Balance: big.NewInt(1000000000000000)},
		acc2Addr: {Balance: big.NewInt(1000000000000000)},
	})
	owners := []common.Address{acc1Addr, acc2Addr}
	transactOpts := bind.NewKeyedTransactor(key)
	smcAddr, smc, err := DeployConfigReward(transactOpts, contractBackend, owners, big.NewInt(int64(len(owners))))
	if err != nil {
		t.Fatalf("can't deploy root registry: %v", err)
	}
	contractBackend.Commit()

	masternodeRate, err := smc.GetRate()
	if err != nil {
		t.Fatalf("can't get candidates: %v", err)
	}
	if masternodeRate.Cmp(defaultRate) != 0 {
		t.Fatalf("Fail when get master node rate , wanted : %v , got :%v ", defaultRate, masternodeRate)
	}
	acc1Smc, err := NewConfigReward(bind.NewKeyedTransactor(acc1Key), smcAddr, contractBackend)
	acc2Smc, err := NewConfigReward(bind.NewKeyedTransactor(acc2Key), smcAddr, contractBackend)
	_, err = acc1Smc.SubmitTransaction(changeRate)
	if err != nil {
		t.Fatalf("acc1 can't submit Transaction: %v", err)
	}
	contractBackend.Commit()
	_, err = acc2Smc.ConfirmTransaction(big.NewInt(0))
	if err != nil {
		t.Fatalf("acc2 can't confirm Transaction: %v", err)
	}
	contractBackend.Commit()
	masternodeRate, err = smc.GetRate()
	if masternodeRate.Cmp(changeRate) != 0 {
		t.Fatalf("Fail when get master node rate , wanted : %v , got :%v ", changeRate, masternodeRate)
	}
}
