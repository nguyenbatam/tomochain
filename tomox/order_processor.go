package tomox

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/tomox/state"
	"math/big"
)

func ProcessOrder(statedb *state.StateDB, orderBook common.Hash, order *OrderItem) ([]map[string]string, common.Hash, error) {
	var (
		orderInBook common.Hash
		trades      []map[string]string
		err         error
	)
	orderType := order.Type
	// if we do not use auto-increment orderid, we must set price slot to avoid conflict
	if orderType == Market {
		log.Debug("Process market order", "order", order)
		trades, orderInBook, err = processMarketOrder(statedb, orderBook, order)
		if err != nil {
			return nil, state.EmptyHash, err
		}
	} else {
		log.Debug("Process limit order", "order", order)
		trades, orderInBook, err = processLimitOrder(statedb, orderBook, order)
		if err != nil {
			return nil, state.EmptyHash, err
		}
	}

	return trades, orderInBook, nil
}

// processMarketOrder : process the market order
func processMarketOrder(statedb *state.StateDB, orderBook common.Hash, order *OrderItem) ([]map[string]string, common.Hash, error) {
	var (
		trades      []map[string]string
		newTrades   []map[string]string
		orderInBook common.Hash
	)
	quantityToTrade := order.Quantity
	side := order.Side
	// speedup the comparison, do not assign because it is pointer
	zero := Zero()
	if side == Bid {
		bestPrice, err := statedb.GetBestAskPrice(orderBook)
		if err != nil {
			return nil, state.EmptyHash, err
		}
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, orderInBook, err = processOrderList(statedb, Ask, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, orderInBook, err
			}
			trades = append(trades, newTrades...)
			bestPrice, err = statedb.GetBestAskPrice(orderBook)
			if err != nil {
				return nil, state.EmptyHash, err
			}
		}
	} else {
		bestPrice, err := statedb.GetBestBidPrice(orderBook)
		if err != nil {
			return nil, state.EmptyHash, err
		}
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, orderInBook, err = processOrderList(statedb, Bid, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, orderInBook, err
			}
			trades = append(trades, newTrades...)
			bestPrice, err = statedb.GetBestBidPrice(orderBook)
			if err != nil {
				return nil, state.EmptyHash, err
			}
		}
	}
	return trades, orderInBook, nil
}

// processLimitOrder : process the limit order, can change the quote
// If not care for performance, we should make a copy of quote to prevent further reference problem
func processLimitOrder(statedb *state.StateDB, orderBook common.Hash, order *OrderItem) ([]map[string]string, common.Hash, error) {
	var (
		trades      []map[string]string
		newTrades   []map[string]string
		orderInBook common.Hash
	)
	quantityToTrade := order.Quantity
	side := order.Side
	price := order.Price

	// speedup the comparison, do not assign because it is pointer
	zero := Zero()

	if side == Bid {
		minPrice, err := statedb.GetBestAskPrice(orderBook)
		if err != nil {
			return nil, state.EmptyHash, err
		}
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(minPrice) >= 0 {
			log.Debug("Min price in asks tree", "price", minPrice.String())
			quantityToTrade, newTrades, orderInBook, err = processOrderList(statedb, Ask, orderBook, minPrice, quantityToTrade, order)
			if err != nil {
				return nil, state.EmptyHash, err
			}
			trades = append(trades, newTrades...)
			log.Debug("New trade found", "newTrades", newTrades, "orderInBook", orderInBook, "quantityToTrade", quantityToTrade)
			minPrice, err = statedb.GetBestAskPrice(orderBook)
			if err != nil {
				return nil, state.EmptyHash, err
			}
		}
	} else {
		maxPrice, err := statedb.GetBestBidPrice(orderBook)
		if err != nil {
			return nil, state.EmptyHash, err
		}
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(maxPrice) >= 0 {
			log.Debug("Max price in bids tree", "price", maxPrice.String())
			quantityToTrade, newTrades, orderInBook, err = processOrderList(statedb, Bid, orderBook, maxPrice, quantityToTrade, order)
			if err != nil {
				return nil, state.EmptyHash, err
			}
			trades = append(trades, newTrades...)
			log.Debug("New trade found", "newTrades", newTrades, "orderInBook", orderInBook, "quantityToTrade", quantityToTrade)
			maxPrice, err = statedb.GetBestBidPrice(orderBook)
			if err != nil {
				return nil, state.EmptyHash, err
			}
		}
	}
	if quantityToTrade.Cmp(zero) > 0 {
		orderId := statedb.GetNonce(orderBook)
		order.OrderID = orderId + 1
		order.Quantity = quantityToTrade
		statedb.SetNonce(orderBook, orderId+1)
		orderIdHash := common.BigToHash(new(big.Int).SetUint64(orderId))
		statedb.InsertOrderItem(orderBook, orderIdHash, price, quantityToTrade, side)
		log.Debug("After matching, order (unmatched part) is now added to bids tree", "order", order)
		orderInBook = orderIdHash
	}
	return trades, orderInBook, nil
}

// processOrderList : process the order list
func processOrderList(statedb *state.StateDB, side string, orderBook common.Hash, price *big.Int, quantityStillToTrade *big.Int, order *OrderItem) (*big.Int, []map[string]string, common.Hash, error) {
	log.Debug("Process matching between order and orderlist")
	quantityToTrade := CloneBigInt(quantityStillToTrade)
	var (
		trades      []map[string]string
		orderInBook common.Hash
	)
	// speedup the comparison, do not assign because it is pointer
	zero := Zero()
	orderId, amount, err := statedb.GetBestOrderIdAndAmount(orderBook, price)
	if err != nil {
		return nil, nil, state.EmptyHash, err
	}
	for amount.Cmp(zero) > 0 && quantityToTrade.Cmp(zero) > 0 {
		var (
			tradedQuantity *big.Int
		)
		if quantityToTrade.Cmp(amount) <= 0 {
			tradedQuantity = CloneBigInt(quantityToTrade)
			quantityToTrade = Zero()
			orderInBook = orderId
		} else {
			tradedQuantity = CloneBigInt(amount)
			quantityToTrade = Sub(quantityToTrade, tradedQuantity)
		}
		statedb.SubAmountOrderItem(orderBook, orderId, price, amount, side)
		log.Debug("Update quantity for orderId", "orderId", orderId.Hex())
		log.Debug("TRADE", "orderBook", orderBook, "Price", price, "Amount", tradedQuantity, "orderId", orderId, "side", side)

		transactionRecord := make(map[string]string)
		transactionRecord[TradedTakerOrderHash] = hex.EncodeToString(order.Hash.Bytes())
		//transactionRecord[TradedMakerOrderHash] = hex.EncodeToString(headOrder.Item.Hash.Bytes())
		//transactionRecord[TradedTimestamp] = strconv.FormatUint(orderBook.Timestamp, 10)
		transactionRecord[TradedQuantity] = tradedQuantity.String()
		//transactionRecord[TradedMakerExchangeAddress] = headOrder.Item.ExchangeAddress.String()
		//transactionRecord[TradedMaker] = headOrder.Item.UserAddress.String()
		//transactionRecord[TradedBaseToken] = headOrder.Item.BaseToken.String()
		//transactionRecord[TradedQuoteToken] = headOrder.Item.QuoteToken.String()
		//// maker price is actual price
		//// taker price is offer price
		//// tradedPrice is always actual price
		//transactionRecord[TradedPrice] = headOrder.Item.Price.String()
		trades = append(trades, transactionRecord)
	}
	return quantityToTrade, trades, orderInBook, nil
}
