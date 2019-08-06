// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package state

import (
	"github.com/ethereum/go-ethereum/tomox"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type journalEntry interface {
	undo(*StateDB)
}

type journal []journalEntry

type (
	// Changes to the price trie.
	createExchangeObjectChange struct {
		hash *common.Hash
	}
	resetExchangeObjectChange struct {
		prev *stateExchanges
	}
	createOrderListAskChange struct {
		orderBook *common.Hash
		price     *common.Hash
	}
	resetOrderListAskChange struct {
		orderBook *common.Hash
		prev      *stateOrderList
	}

	createBidOrderListChange struct {
		orderBook *common.Hash
		price     *common.Hash
	}
	resetBidOrderListChange struct {
		orderBook *common.Hash
		prev      *stateOrderList
	}

	nonceChange struct {
		hash *common.Hash
		prev uint64
	}

	volumeChange struct {
		orderBook *common.Hash
		orderType *string
		price     *common.Hash
		prev      big.Int
	}
	asksChange struct {
		orderBook *common.Hash
		price     *common.Hash
		prevalue  common.Hash
	}
	bidsChange struct {
		orderHash *common.Hash
		price     *common.Hash
		prevalue  common.Hash
	}

	orderChange struct {
		orderBook, key *common.Hash
		prevalue       common.Hash
	}

	storageOrderItemChange struct {
		orderBook, price *common.Hash
		orderType        *string
		orderId          *uint64
		value            common.Hash
	}
)

func (ch createExchangeObjectChange) undo(s *StateDB) {
	delete(s.stateExhangeObjects, *ch.hash)
	delete(s.stateExhangeObjectsDirty, *ch.hash)
}

func (ch resetExchangeObjectChange) undo(s *StateDB) {
	s.setStateExchangeObject(ch.prev)
}

func (ch createOrderListAskChange) undo(s *StateDB) {
	delete(s.getStateExchangeObject(*ch.orderBook).stateAskObjects, *ch.price)
	delete(s.getStateExchangeObject(*ch.orderBook).stateAskObjectsDirty, *ch.price)
}

func (ch resetOrderListAskChange) undo(s *StateDB) {
	s.getStateExchangeObject(*ch.orderBook).setStateOrderListAskObject(ch.prev)
}

func (ch createBidOrderListChange) undo(s *StateDB) {
	delete(s.getStateExchangeObject(*ch.orderBook).stateBidObjects, *ch.price)
	delete(s.getStateExchangeObject(*ch.orderBook).stateBidObjectsDirty, *ch.price)
}

func (ch resetBidOrderListChange) undo(s *StateDB) {
	s.getStateExchangeObject(*ch.orderBook).setStateBidOrderListObject(ch.prev)
}

func (ch nonceChange) undo(s *StateDB) {
	s.getStateExchangeObject(*ch.hash).setNonce(ch.prev)
}

func (ch volumeChange) undo(s *StateDB) {
	switch *ch.orderType {
	case tomox.Bid:
		s.getStateExchangeObject(*ch.orderBook).stateBidObjects[*ch.price].SetVolume(ch.prev)
	case tomox.Ask:
		s.getStateExchangeObject(*ch.orderBook).stateAskObjects[*ch.price].SetVolume(ch.prev)
	}
}

func (ch asksChange) undo(s *StateDB) {
	s.getStateExchangeObject(*ch.orderBook).setAskPrice(*ch.price, ch.prevalue)
}

func (ch bidsChange) undo(s *StateDB) {
	s.getStateExchangeObject(*ch.orderHash).setBidPrice(*ch.price, ch.prevalue)
}

func (ch orderChange) undo(s *StateDB) {
	s.getStateExchangeObject(*ch.orderBook).setOrderHash(*ch.key, ch.prevalue)
}

func (ch storageOrderItemChange) undo(s *StateDB) {
	switch *ch.orderType {
	case tomox.Bid:
		s.getStateExchangeObject(*ch.orderBook).getStateBidOrderListObject(*ch.price).setOrderItem(*ch.orderId, ch.value)
	case tomox.Ask:
		s.getStateExchangeObject(*ch.orderBook).getStateOrderListAskObject(*ch.price).setOrderItem(*ch.orderId, ch.value)
	}
}
