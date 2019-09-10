// Copyright 2014 The go-ethereum Authors
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
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// stateObject represents an Ethereum orderId which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// ExchangeObject values can be accessed and modified through the object.
// Finally, call CommitAskTrie to write the modified storage trie into a database.
type stateOrderItem struct {
	orderId   common.Hash
	orderBook common.Hash
	data      OrderItem
	db        *StateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	deleted bool
	onDirty func(price common.Hash) // Callback method to mark a state object newly dirty
}

// empty returns whether the orderId is considered empty.
func (s *stateOrderItem) empty() bool {
	return s.data.Amount.Cmp(tomox.Zero()) == 0
}

// ExchangeObject is the Ethereum consensus representation of exchanges.
// These objects are stored in the main orderId trie.
type OrderItem struct {
	OrderType string
	Amount    big.Int
}

// newObject creates a state object.
func newStateOrderItem(db *StateDB, orderBook common.Hash, orderId common.Hash, data OrderItem, onDirty func(price common.Hash)) *stateOrderItem {
	return &stateOrderItem{
		db:        db,
		orderBook: orderBook,
		orderId:   orderId,
		data:      data,
		onDirty:   onDirty,
	}
}

// EncodeRLP implements rlp.Encoder.
func (c *stateOrderItem) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, c.data)
}

// setError remembers the first non-nil error it is called with.
func (self *stateOrderItem) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (self *stateOrderItem) deepCopy(db *StateDB, onDirty func(price common.Hash)) *stateOrderItem {
	stateOrderItem := newStateOrderItem(db, self.orderBook, self.orderId, self.data, onDirty)
	stateOrderItem.deleted = self.deleted
	return stateOrderItem
}

// AddVolume removes amount from c's balance.
// It is used to add funds to the destination exchanges of a transfer.
func (c *stateOrderItem) SubAmount(amount *big.Int) {
	c.SetAmount(*new(big.Int).Sub(&c.data.Amount, amount))
}

//
// Attribute accessors
//
func (self *stateOrderItem) SetAmount(amount big.Int) {
	self.db.journal = append(self.db.journal, AmountOrderItemChange{
		orderBook: &self.orderBook,
		orderId:   &self.orderId,
		amount:    &amount,
		prev:      self.data.Amount,
	})
	self.setAmount(amount)
}

func (self *stateOrderItem) setAmount(quantity big.Int) {
	self.data.Amount = quantity
	if self.onDirty != nil {
		self.onDirty(self.orderId)
		self.onDirty = nil
	}
}

func (self *stateOrderItem) Amount() big.Int {
	return self.data.Amount
}

func (self *stateOrderItem) OrderId() common.Hash {
	return self.orderId
}
