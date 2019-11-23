package test

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	configReward "github.com/ethereum/go-ethereum/contracts/configreward"
	"github.com/ethereum/go-ethereum/contracts/test/contract"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
	"time"
)

var (
	mainKey, _  = crypto.HexToECDSA("65ec4d4dfbcac594a14c36baa462d6f73cd86134840f6cf7b80a1e1cd33473e2")
	mainAddr    = crypto.PubkeyToAddress(mainKey.PublicKey)
	acc1Key, _  = crypto.HexToECDSA("be6573530c6d04ba50c1294b2c0e73a2275238d6eaa78990dccb996589c87d43")
	acc2Key, _  = crypto.HexToECDSA("f7753824da8ed6b8fc82a26fd7a638a12cba2bdf063c566ecbe84d07f040a0d6")
	acc1Addr    = crypto.PubkeyToAddress(acc1Key.PublicKey)
	acc2Addr    = crypto.PubkeyToAddress(acc2Key.PublicKey)
	rcpEndPoint = "http://localhost:8501"
)

func TestLocalhost(t *testing.T) {
	client, err := ethclient.Dial(rcpEndPoint)
	if err != nil {
		fmt.Println(err, client)
	}
	nonce, _ := client.NonceAt(context.Background(), mainAddr, nil)
	fmt.Println(nonce, mainAddr.Hex())
	mainAccount := bind.NewKeyedTransactor(mainKey)
	mainAccount.Nonce = big.NewInt(int64(nonce))
	mainAccount.Value = big.NewInt(0)      // in wei
	mainAccount.GasLimit = uint64(4000000) // in units
	mainAccount.GasPrice = big.NewInt(0)
	configRewardInstance, _ := configReward.NewConfigReward(mainAccount, common.HexToAddress(common.ConfigRewardAddr), client)
	fmt.Println(configRewardInstance.GetRate())

	//foudationInstance, _ := multisigwallet.NewMultiSigWallet(mainAccount, common.HexToAddress(common.FoudationAddr), client)
	//foudationInstance.SubmitTransaction()
	smcTestAddr, smcTestInsance, err := DeployContract1(mainAccount, client)
	if err != nil {
		fmt.Println("DeployContract1", err)
	}
	fmt.Println("smcTestAddr", smcTestAddr.Hex())
	time.Sleep(10 * time.Second) // wait process transaction : deploy smart contract 1
	a, err := smcTestInsance.GetA()
	fmt.Println("GetA", a, err)
	smcTestInsance.TransactOpts.Nonce = big.NewInt(int64(nonce + 1))
	_, err = smcTestInsance.AddA()
	if err != nil {
		fmt.Println("AddA", err)
	}
	time.Sleep(10 * time.Second) // wait process transaction : func addA() at smart contract 1
	fmt.Println(smcTestInsance.GetA())

	nonce, _ = client.NonceAt(context.Background(), mainAddr, nil)
	dataUpgrade := append([]byte{}, smcTestAddr.Bytes()...)
	dataUpgrade = append(dataUpgrade, common.FromHex(contract.Contract2Bin)[32:]...)
	upgradeSmcTx := types.NewTransaction(nonce, common.HexToAddress(common.SMCUpgradeAddr), big.NewInt(0), 200000, big.NewInt(100), dataUpgrade)

	signedTx, _ := types.SignTx(upgradeSmcTx, types.NewEIP155Signer(big.NewInt(66)), mainKey)
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(10 * time.Second) // wait process transaction : upgrade smart contract 1 -> smart contract 2

	smcTest2Insance, _ := NewContract2(mainAccount, smcTestAddr, client)
	fmt.Println(smcTest2Insance.GetA())
	smcTest2Insance.TransactOpts.Nonce = big.NewInt(int64(nonce + 1))
	_, err = smcTest2Insance.AddA()
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(10 * time.Second) // wait process transaction : func addA() at smart contract 2
	fmt.Println(smcTest2Insance.GetA())
	smcTest2Insance.TransactOpts.Nonce = big.NewInt(int64(nonce + 2))
	_, err = smcTest2Insance.AddB()
	if err != nil {
		fmt.Println("AddB", err)
	}
	time.Sleep(10 * time.Second) // wait process transaction : func addB() at smart contract 2
	fmt.Println(smcTest2Insance.GetB())

}
