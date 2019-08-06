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
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/tomox"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// stateObject represents an Ethereum price which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// ExchangeObject values can be accessed and modified through the object.
// Finally, call CommitAskTrie to write the modified storage trie into a database.
type stateOrderList struct {
	price     common.Hash
	orderBook common.Hash
	orderType string
	data      OrderList
	db        *StateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Write caches.
	trie state.Trie // storage trie, which becomes non-nil on first access

	cachedStorage map[uint64]common.Hash // Storage entry cache to avoid duplicate reads
	dirtyStorage  map[uint64]common.Hash // Storage entries that need to be flushed to disk

	deleted bool
	onDirty func(price common.Hash) // Callback method to mark a state object newly dirty
}

// empty returns whether the price is considered empty.
func (s *stateOrderList) empty() bool {
	return s.data.Volume.Cmp(tomox.Zero()) == 0
}

// ExchangeObject is the Ethereum consensus representation of exchanges.
// These objects are stored in the main price trie.
type OrderList struct {
	Volume big.Int
	Root   common.Hash // merkle root of the storage trie
}

// newObject creates a state object.
func newStateOrderList(db *StateDB, orderType string, orderBook common.Hash, price common.Hash, data OrderList, onDirty func(price common.Hash)) *stateOrderList {
	return &stateOrderList{
		db:            db,
		orderType:     orderType,
		orderBook:     orderBook,
		price:         price,
		data:          data,
		cachedStorage: make(map[uint64]common.Hash),
		dirtyStorage:  make(map[uint64]common.Hash),
		onDirty:       onDirty,
	}
}

// EncodeRLP implements rlp.Encoder.
func (c *stateOrderList) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, c.data)
}

// setError remembers the first non-nil error it is called with.
func (self *stateOrderList) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (c *stateOrderList) getTrie(db state.Database) state.Trie {
	if c.trie == nil {
		var err error
		c.trie, err = db.OpenStorageTrie(c.price, c.data.Root)
		if err != nil {
			c.trie, _ = db.OpenStorageTrie(c.price, common.Hash{})
			c.setError(fmt.Errorf("can't create storage trie: %v", err))
		}
	}
	return c.trie
}

// GetState returns a value in price storage.
func (self *stateOrderList) GetOrderItem(db state.Database, orderId uint64) common.Hash {
	value, exists := self.cachedStorage[orderId]
	if exists {
		return value
	}
	// Load from DB in case it is missing.
	enc, err := self.getTrie(db).TryGet(new(big.Int).SetUint64(orderId).Bytes()[:])
	if err != nil {
		self.setError(err)
		return common.Hash{}
	}
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			self.setError(err)
		}
		value.SetBytes(content)
	}
	if (value != common.Hash{}) {
		self.cachedStorage[orderId] = value
	}
	return value
}

// SetState updates a value in price storage.
func (self *stateOrderList) SetOrderItem(db state.Database, orderId uint64, hash common.Hash) {
	self.db.journal = append(self.db.journal, storageOrderItemChange{
		orderBook: &self.orderBook,
		orderType: &self.orderType,
		price:     &self.price,
		orderId:   &orderId,
		value:     self.GetOrderItem(db, orderId),
	})
	self.setOrderItem(orderId, hash)
}

func (self *stateOrderList) setOrderItem(orderId uint64, hash common.Hash) {
	self.cachedStorage[orderId] = hash
	self.dirtyStorage[orderId] = hash

	if self.onDirty != nil {
		self.onDirty(self.Price())
		self.onDirty = nil
	}
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *stateOrderList) updateTrie(db state.Database) state.Trie {
	tr := self.getTrie(db)
	for orderId, hash := range self.dirtyStorage {
		delete(self.dirtyStorage, orderId)
		key := new(big.Int).SetUint64(orderId).Bytes()[:]
		if (hash == common.Hash{}) {
			self.setError(tr.TryDelete(key))
			continue
		}
		// Encoding []byte cannot fail, ok to ignore the error.
		v, _ := rlp.EncodeToBytes(bytes.TrimLeft(hash[:], "\x00"))
		self.setError(tr.TryUpdate(key, v))
	}
	return tr
}

// UpdateRoot sets the trie root to the current root price of
func (self *stateOrderList) updateRoot(db state.Database) {
	self.updateTrie(db)
	self.data.Root = self.trie.Hash()
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *stateOrderList) CommitTrie(db state.Database) error {
	self.updateTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.trie.Commit(nil)
	if err == nil {
		self.data.Root = root
	}
	return err
}

func (self *stateOrderList) deepCopy(db *StateDB, onDirty func(price common.Hash)) *stateOrderList {
	stateOrderList := newStateOrderList(db, self.orderType, self.orderBook, self.price, self.data, onDirty)
	if self.trie != nil {
		stateOrderList.trie = db.db.CopyTrie(self.trie)
	}
	stateOrderList.dirtyStorage = make(map[uint64]common.Hash)
	stateOrderList.cachedStorage = make(map[uint64]common.Hash)
	for key, value := range self.dirtyStorage {
		stateOrderList.dirtyStorage[key] = value
		stateOrderList.cachedStorage[key] = value
	}
	stateOrderList.deleted = self.deleted
	return stateOrderList
}

// AddVolume removes amount from c's balance.
// It is used to add funds to the destination exchanges of a transfer.
func (c *stateOrderList) AddVolume(amount *big.Int) {
	c.setVolume(*new(big.Int).Add(&c.data.Volume, amount))
}

//
// Attribute accessors
//
func (self *stateOrderList) SetVolume(volume big.Int) {
	self.db.journal = append(self.db.journal, volumeChange{
		orderBook: &self.orderBook,
		orderType: &self.orderType,
		price:     &self.price,
		prev:      self.data.Volume,
	})
	self.setVolume(volume)
}

func (self *stateOrderList) setVolume(volume big.Int) {
	self.data.Volume = volume
	if self.onDirty != nil {
		self.onDirty(self.price)
		self.onDirty = nil
	}
}

// Returns the address of the contract/price
func (c *stateOrderList) Price() common.Hash {
	return c.price
}

func (self *stateOrderList) Volume() big.Int {
	return self.data.Volume
}
