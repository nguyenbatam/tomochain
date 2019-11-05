package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/tomox/simulation"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/tomox"
	"log"
	"math/big"
	"strings"
	"time"
)

func main() {
	testPair("BTC/TOMO", simulation.TokenNameAddress["BTC"], simulation.TokenNameAddress["TOMO"], simulation.RelayerCoinbaseAddr, new(big.Int).Mul(big.NewInt(1), simulation.BaseTOMO), new(big.Int).Mul(big.NewInt(5000), simulation.BaseTOMO))
}

func testPair(pair string, baseToken common.Address, quoteToken common.Address, exAdress common.Address, amount *big.Int, price *big.Int) {
	client, err := ethclient.Dial(simulation.RpcEndpoint)
	if err != nil {
		fmt.Println(err, client)
	}
	//
	// BUY
	msg := &ethapi.OrderMsg{
		Quantity:        amount,
		Price:           price,
		ExchangeAddress: exAdress,
		UserAddress:     simulation.MainAddr,
		BaseToken:       baseToken,
		QuoteToken:      quoteToken,
		Status:          tomox.OrderStatusNew,
		Side:            tomox.Bid,
		Type:            "LO",
		PairName:        pair,
	}
	oldVolume := getValue(msg.Price, getBidTree(baseToken, quoteToken))
	result, err := client.CallContext(context.Background(), "tomox_getOrderCount", msg.UserAddress)
	nonce := hexutil.Uint64(0)
	json.Unmarshal(result, &nonce)
	tx := types.NewOrderTransaction(uint64(nonce), msg.Quantity, msg.Price, msg.ExchangeAddress, msg.UserAddress, msg.BaseToken, msg.QuoteToken, msg.Status, msg.Side, msg.Type, msg.PairName, common.Hash{}, 0)
	signedTx, err := types.OrderSignTx(tx, types.OrderTxSigner{}, simulation.MainKey)
	if err != nil {
		log.Fatalln(err)
	}
	err = client.SendOrderTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalln(err)
	}
	afterSendBuy := big.NewInt(0)
	for true {
		fmt.Println("wait 2 second add order to order book")
		time.Sleep(1 * time.Second)
		bidOrders := getBidTree(baseToken, quoteToken)
		afterSendBuy = getValue(msg.Price, bidOrders)
		if afterSendBuy.Cmp(oldVolume) == 0 {
			continue
		}
		total := new(big.Int).Add(oldVolume, msg.Quantity)
		if afterSendBuy.Cmp(total) != 0 {
			log.Fatalf("Error when check total volume after sent buy order , price : %d , wanted : %d , got :%d ", msg.Price, total, afterSendBuy)
		}
		break
	}

	// SELL
	msg = &ethapi.OrderMsg{
		Quantity:        amount,
		Price:           price,
		ExchangeAddress: exAdress,
		UserAddress:     simulation.MainAddr,
		BaseToken:       baseToken,
		QuoteToken:      quoteToken,
		Status:          tomox.OrderStatusNew,
		Side:            tomox.Ask,
		Type:            "LO",
		PairName:        pair,
	}
	blockBeforeMatch, err := client.BlockByNumber(context.Background(), nil)
	numberBeforeMatch := blockBeforeMatch.Number()
	if err != nil {
		log.Fatalln("Error when get current block before match ", err)
	}
	tx = types.NewOrderTransaction(uint64(nonce+1), msg.Quantity, msg.Price, msg.ExchangeAddress, msg.UserAddress, msg.BaseToken, msg.QuoteToken, msg.Status, msg.Side, msg.Type, msg.PairName, common.Hash{}, 0)
	signedTx, err = types.OrderSignTx(tx, types.OrderTxSigner{}, simulation.MainKey)
	if err != nil {
		log.Fatalln(err)
	}
	err = client.SendOrderTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalln(err)
	}
	afterSendSell := big.NewInt(0)
	for true {
		fmt.Println("wait 2 second match order")
		time.Sleep(1 * time.Second)
		bidOrders := getBidTree(baseToken, quoteToken)
		afterSendSell = getValue(msg.Price, bidOrders)
		if afterSendSell.Cmp(afterSendBuy) == 0 {
			continue
		}
		total := new(big.Int).Sub(afterSendBuy, msg.Quantity)
		if afterSendSell.Cmp(total) != 0 {
			log.Fatalf("Error when check total volume after sent buy order , price : %d , wanted : %d , got :%d ", msg.Price, total, afterSendSell)
		}
		break
	}
	blockAfterMatch, err := client.BlockByNumber(context.Background(), nil)
	numberAfterMatch := blockAfterMatch.Number()
	if err != nil {
		log.Fatalln("Error when get current block after match ", err)
	}
	getTradeInfo(numberBeforeMatch, numberAfterMatch)

}
func getValue(key *big.Int, data map[*big.Int]*big.Int) *big.Int {
	for k, v := range data {
		if k.Cmp(key) == 0 {
			return v
		}
	}
	return big.NewInt(0)
}

func getBidTree(baseToken common.Address, quoteToken common.Address) map[*big.Int]*big.Int {
	mapVolume := map[*big.Int]*big.Int{}
	client, err := ethclient.Dial(simulation.RpcEndpoint)
	if err != nil {
		log.Fatalln(err, client)
	}
	result, err := client.CallContext(context.Background(), "tomox_getBidTree", baseToken.Hex(), quoteToken.Hex())
	fmt.Println("baseToken", baseToken.Hex(), "quoteToken", quoteToken.Hex(), "result", string(result))
	mapData := make(map[string]interface{})
	if err == nil {
		err = json.Unmarshal(result, &mapData)
	}
	if err != nil {
		fmt.Println("Fail getBidTree ", "err", err)
		return mapVolume
	}
	for k1, v1 := range mapData {
		price, _ := new(big.Int).SetString(k1, 10)
		f := v1.(map[string]interface{})["Volume"].(float64)
		text := fmt.Sprintf("%f", f)
		text = text[0:strings.Index(text, ".")]
		volume, err := new(big.Int).SetString(text, 10)
		fmt.Println("f", f, "text", text, "err", err, "v "+
			""+
			"molume", volume)
		mapVolume[price] = volume
	}
	return mapVolume
}

func getTradeInfo(from *big.Int, to *big.Int) []map[string]string {
	fmt.Println("getTradeInfo", from, to)
	client, err := ethclient.Dial(simulation.RpcEndpoint)
	if err != nil {
		log.Fatalln(err)
	}
	for number := from; number.Cmp(to) <= 0; number = new(big.Int).Add(number, big.NewInt(1)) {
		block, err := client.BlockByNumber(context.Background(), number)
		if err != nil {
			log.Fatalln(err)
		}
		trades := getTradeData(block)
		if len(trades) > 0 {
			return trades
		}
	}
	return nil
}

func getTradeData(block *types.Block) []map[string]string {
	for _, tx := range block.Transactions() {
		if tx.To() != nil && tx.To().Hex() == common.TomoXAddr {
			txMatchBatch, err := tomox.DecodeTxMatchesBatch(tx.Data())
			if err != nil {
				log.Fatalln("transaction match is corrupted. Failed to decode txMatchBatch. Error: ", err)
			}
			if len(txMatchBatch.Data) > 0 {
				for _, txMatch := range txMatchBatch.Data {
					trades := txMatch.GetTrades()
					if len(trades) > 0 {
						return trades
					}
				}
			}
		}
	}
	return nil
}
