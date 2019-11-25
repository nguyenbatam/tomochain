package test

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	configReward "github.com/ethereum/go-ethereum/contracts/configreward"
	"github.com/ethereum/go-ethereum/contracts/multisigwallet"
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

func testUpgradeSmc() {
	client, err := ethclient.Dial(rcpEndPoint)
	if err != nil {
		fmt.Println(err, client)
	}
	nonce, _ := client.NonceAt(context.Background(), mainAddr, nil)
	mainAccount := bind.NewKeyedTransactor(mainKey)
	mainAccount.Nonce = big.NewInt(int64(nonce))
	mainAccount.Value = big.NewInt(0)      // in wei
	mainAccount.GasLimit = uint64(4000000) // in units
	mainAccount.GasPrice = big.NewInt(0)
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
func testChangConfigReward() {
	client, err := ethclient.Dial(rcpEndPoint)
	if err != nil {
		fmt.Println(err, client)
	}
	newRate := big.NewInt(15)
	nonce, _ := client.NonceAt(context.Background(), acc1Addr, nil)
	acc1Account := bind.NewKeyedTransactor(acc1Key)
	acc1Account.Nonce = big.NewInt(int64(nonce))
	acc1Account.Value = big.NewInt(0)      // in wei
	acc1Account.GasLimit = uint64(4000000) // in units
	acc1Account.GasPrice = big.NewInt(0)
	acc1Instance, _ := configReward.NewConfigReward(acc1Account, common.HexToAddress(common.ConfigRewardAddr), client)
	acc1Instance.SubmitTransaction(newRate)
	time.Sleep(10 * time.Second) // wait process transaction : func  at smart contract
	transactionId, err := acc1Instance.GetTransactionCount(true, true)
	if err != nil {
		fmt.Println("get transactionId", err)
	}
	nonce, _ = client.NonceAt(context.Background(), acc2Addr, nil)
	acc2Account := bind.NewKeyedTransactor(acc2Key)
	acc2Account.Nonce = big.NewInt(int64(nonce))
	acc2Account.Value = big.NewInt(0)      // in wei
	acc2Account.GasLimit = uint64(4000000) // in units
	acc2Account.GasPrice = big.NewInt(0)
	acc2Instance, _ := configReward.NewConfigReward(acc2Account, common.HexToAddress(common.ConfigRewardAddr), client)
	acc2Instance.ConfirmTransaction(new(big.Int).Sub(transactionId, big.NewInt(1)))
	time.Sleep(10 * time.Second) // wait process transaction : func  at smart contract

	fmt.Println("new rate")
	fmt.Println(acc2Instance.GetRate())
}

func testFoundationWallet() {
	client, err := ethclient.Dial(rcpEndPoint)
	if err != nil {
		fmt.Println(err, client)
	}
	amount := big.NewInt(10000000)
	to := common.HexToAddress("0x2da72c9c4792e2ba85e2b980af4c6aa9afa9f3df")
	toBalance, _ := client.BalanceAt(context.Background(), to, nil)
	foundationBalance, _ := client.BalanceAt(context.Background(), common.HexToAddress(common.FoudationAddr), nil)
	if foundationBalance.Cmp(toBalance) < 0 {
		fmt.Println("foundation address not enough balance")
		return
	}
	fmt.Println("before sent", toBalance)
	nonce, _ := client.NonceAt(context.Background(), acc1Addr, nil)
	acc1Account := bind.NewKeyedTransactor(acc1Key)
	acc1Account.Nonce = big.NewInt(int64(nonce))
	acc1Account.Value = big.NewInt(0)      // in wei
	acc1Account.GasLimit = uint64(4000000) // in units
	acc1Account.GasPrice = big.NewInt(0)
	acc1Instance, _ := multisigwallet.NewMultiSigWallet(acc1Account, common.HexToAddress(common.FoudationAddr), client)
	acc1Instance.SubmitTransaction(to, amount, nil)
	time.Sleep(10 * time.Second) // wait process transaction : func  at smart contract
	transactionId, err := acc1Instance.GetTransactionCount(true, true)
	if err != nil {
		fmt.Println("get transactionId", err)
	}
	nonce, _ = client.NonceAt(context.Background(), acc2Addr, nil)
	acc2Account := bind.NewKeyedTransactor(acc2Key)
	acc2Account.Nonce = big.NewInt(int64(nonce))
	acc2Account.Value = big.NewInt(0)      // in wei
	acc2Account.GasLimit = uint64(4000000) // in units
	acc2Account.GasPrice = big.NewInt(0)
	acc2Instance, _ := multisigwallet.NewMultiSigWallet(acc2Account, common.HexToAddress(common.FoudationAddr), client)
	acc2Instance.ConfirmTransaction(new(big.Int).Sub(transactionId, big.NewInt(1)))
	time.Sleep(10 * time.Second) // wait process transaction : func  at smart contract

	fmt.Println("new balance")
	newBalance, _ := client.BalanceAt(context.Background(), to, nil)
	fmt.Println(newBalance)
}
func TestLocalhost(t *testing.T) {
	testUpgradeSmc()
	testChangConfigReward()
	testFoundationWallet()
}
