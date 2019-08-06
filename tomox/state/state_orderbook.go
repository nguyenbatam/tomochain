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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/tomox"
	"io"
)

// stateObject represents an Ethereum price which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// ExchangeObject values can be accessed and modified through the object.
// Finally, call CommitAskTrie to write the modified storage trie into a database.
type stateExchanges struct {
	hash common.Hash // orderbookHashprice of ethereum address of the price
	data ExchangeObject
	db   *StateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Write caches.
	asksTrie   state.Trie // storage trie, which becomes non-nil on first access
	bidsTrie   state.Trie // storage trie, which becomes non-nil on first access
	ordersTrie state.Trie // storage trie, which becomes non-nil on first access

	cachedAsksStorage map[common.Hash]common.Hash // Storage entry cache to avoid duplicate reads
	dirtyAsksStorage  map[common.Hash]common.Hash // Storage entries that need to be flushed to disk

	cachedBidsStorage map[common.Hash]common.Hash // Storage entry cache to avoid duplicate reads
	dirtyBidsStorage  map[common.Hash]common.Hash // Storage entries that need to be flushed to d

	cachedOrdersStorage map[common.Hash]common.Hash // Storage entry cache to avoid duplicate reads
	dirtyOrdersStorage  map[common.Hash]common.Hash // Storage entries that need to be flushed to d

	stateAskObjects      map[common.Hash]*stateOrderList
	stateAskObjectsDirty map[common.Hash]struct{}

	stateBidObjects      map[common.Hash]*stateOrderList
	stateBidObjectsDirty map[common.Hash]struct{}

	onDirty func(hash common.Hash) // Callback method to mark a state object newly dirty
}

// empty returns whether the price is considered empty.
func (s *stateExchanges) empty() bool {
	return s.data.Nonce == 0 && common.EmptyHash(s.data.AskRoot) && common.EmptyHash(s.data.BidRoot) && common.EmptyHash(s.data.OrderRoot)
}

// ExchangeObject is the Ethereum consensus representation of exchanges.
// These objects are stored in the main price trie.
type ExchangeObject struct {
	Nonce     uint64
	AskRoot   common.Hash // merkle root of the storage trie
	BidRoot   common.Hash // merkle root of the storage trie
	OrderRoot common.Hash
}

// newObject creates a state object.
func newStateExchanges(db *StateDB, hash common.Hash, data ExchangeObject, onDirty func(addr common.Hash)) *stateExchanges {
	return &stateExchanges{
		db:                  db,
		hash:                hash,
		data:                data,
		cachedAsksStorage:   make(map[common.Hash]common.Hash),
		dirtyAsksStorage:    make(map[common.Hash]common.Hash),
		cachedBidsStorage:   make(map[common.Hash]common.Hash),
		dirtyBidsStorage:    make(map[common.Hash]common.Hash),
		cachedOrdersStorage: make(map[common.Hash]common.Hash),
		dirtyOrdersStorage:  make(map[common.Hash]common.Hash),
		onDirty:             onDirty,
	}
}

// EncodeRLP implements rlp.Encoder.
func (c *stateExchanges) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, c.data)
}

// setError remembers the first non-nil error it is called with.
func (self *stateExchanges) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (c *stateExchanges) getAsksTrie(db state.Database) state.Trie {
	if c.asksTrie == nil {
		var err error
		c.asksTrie, err = db.OpenStorageTrie(c.hash, c.data.AskRoot)
		if err != nil {
			c.asksTrie, _ = db.OpenStorageTrie(c.hash, common.Hash{})
			c.setError(fmt.Errorf("can't create asks trie: %v", err))
		}
	}
	return c.asksTrie
}

// GetState returns a value in price storage.
func (self *stateExchanges) GetAskOrderList(db state.Database, price common.Hash) common.Hash {
	value, exists := self.cachedAsksStorage[price]
	if exists {
		return value
	}
	// Load from DB in case it is missing.
	enc, err := self.getAsksTrie(db).TryGet(price[:])
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
		self.cachedAsksStorage[price] = value
	}
	return value
}

// SetState updates a value in price storage.
func (self *stateExchanges) SetAskPrice(db state.Database, price common.Hash, value common.Hash) {
	self.db.journal = append(self.db.journal, asksChange{
		orderBook: &self.hash,
		price:     &price,
		prevalue:  self.GetAskOrderList(db, price),
	})
	self.setAskPrice(price, value)
}

func (self *stateExchanges) setAskPrice(price common.Hash, value common.Hash) {
	self.cachedAsksStorage[price] = value
	self.dirtyAsksStorage[price] = value

	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *stateExchanges) updateAsksTrie(db state.Database) state.Trie {
	tr := self.getAsksTrie(db)
	for key, value := range self.dirtyAsksStorage {
		delete(self.dirtyAsksStorage, key)
		if (value != common.Hash{}) {
			self.setError(tr.TryDelete(key[:]))
			continue
		}
		// Encoding []byte cannot fail, ok to ignore the error.
		v, _ := rlp.EncodeToBytes(bytes.TrimLeft(value.Bytes()[:], "\x00"))
		self.setError(tr.TryUpdate(key[:], v))
	}
	return tr
}

// UpdateRoot sets the trie root to the current root price of
func (self *stateExchanges) updateAsksRoot(db state.Database) {
	self.updateAsksTrie(db)
	self.data.AskRoot = self.asksTrie.Hash()
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *stateExchanges) CommitAsksTrie(db state.Database) error {
	self.updateAsksTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.asksTrie.Commit(nil)
	if err == nil {
		self.data.AskRoot = root
	}
	return err
}

func (c *stateExchanges) getBidsTrie(db state.Database) state.Trie {
	if c.bidsTrie == nil {
		var err error
		c.bidsTrie, err = db.OpenStorageTrie(c.hash, c.data.BidRoot)
		if err != nil {
			c.bidsTrie, _ = db.OpenStorageTrie(c.hash, common.Hash{})
			c.setError(fmt.Errorf("can't create bids trie: %v", err))
		}
	}
	return c.bidsTrie
}

// GetState returns a value in price storage.
func (self *stateExchanges) GetBidPrice(db state.Database, price common.Hash) common.Hash {
	value, exists := self.cachedBidsStorage[price]
	if exists {
		return value
	}
	// Load from DB in case it is missing.
	enc, err := self.getBidsTrie(db).TryGet(price[:])
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
		self.cachedBidsStorage[price] = value
	}
	return value
}

// SetState updates a value in price storage.
func (self *stateExchanges) SetBidPrice(db state.Database, price common.Hash, value common.Hash) {
	self.db.journal = append(self.db.journal, bidsChange{
		orderHash: &self.hash,
		price:     &price,
		prevalue:  self.GetBidPrice(db, price),
	})
	self.setBidPrice(price, value)
}

func (self *stateExchanges) setBidPrice(price common.Hash, value common.Hash) {
	self.cachedBidsStorage[price] = value
	self.dirtyBidsStorage[price] = value

	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *stateExchanges) updateBidsTrie(db state.Database) state.Trie {
	tr := self.getBidsTrie(db)
	for key, value := range self.dirtyBidsStorage {
		delete(self.dirtyBidsStorage, key)
		if (value != common.Hash{}) {
			self.setError(tr.TryDelete(key[:]))
			continue
		}
		// Encoding []byte cannot fail, ok to ignore the error.
		v, _ := rlp.EncodeToBytes(bytes.TrimLeft(value.Bytes()[:], "\x00"))
		self.setError(tr.TryUpdate(key[:], v))
	}
	return tr
}

// UpdateRoot sets the trie root to the current root price of
func (self *stateExchanges) updateBidsRoot(db state.Database) {
	self.updateBidsTrie(db)
	self.data.BidRoot = self.bidsTrie.Hash()
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *stateExchanges) CommitBidsTrie(db state.Database) error {
	self.updateBidsTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.bidsTrie.Commit(nil)
	if err == nil {
		self.data.BidRoot = root
	}
	return err
}

func (c *stateExchanges) getOrdersTrie(db state.Database) state.Trie {
	if c.ordersTrie == nil {
		var err error
		c.ordersTrie, err = db.OpenStorageTrie(c.hash, c.data.OrderRoot)
		if err != nil {
			c.ordersTrie, _ = db.OpenStorageTrie(c.hash, common.Hash{})
			c.setError(fmt.Errorf("can't create orders trie: %v", err))
		}
	}
	return c.ordersTrie
}

// GetState returns a value in price storage.
func (self *stateExchanges) GetOrderHash(db state.Database, key common.Hash) common.Hash {
	value, exists := self.cachedOrdersStorage[key]
	if exists {
		return value
	}
	// Load from DB in case it is missing.
	enc, err := self.getOrdersTrie(db).TryGet(key[:])
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
		self.cachedOrdersStorage[key] = value
	}
	return value
}

// SetState updates a value in price storage.
func (self *stateExchanges) SetOrderHash(db state.Database, key, value common.Hash) {
	self.db.journal = append(self.db.journal, orderChange{
		orderBook:     &self.hash,
		key:      &key,
		prevalue: self.GetOrderHash(db, key),
	})
	self.setOrderHash(key, value)
}

func (self *stateExchanges) setOrderHash(key, value common.Hash) {
	self.cachedOrdersStorage[key] = value
	self.dirtyOrdersStorage[key] = value

	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *stateExchanges) updateOrdersTrie(db state.Database) state.Trie {
	tr := self.getOrdersTrie(db)
	for key, value := range self.dirtyOrdersStorage {
		delete(self.dirtyOrdersStorage, key)
		if (value == common.Hash{}) {
			self.setError(tr.TryDelete(key[:]))
			continue
		}
		// Encoding []byte cannot fail, ok to ignore the error.
		v, _ := rlp.EncodeToBytes(bytes.TrimLeft(value.Bytes()[:], "\x00"))
		self.setError(tr.TryUpdate(key[:], v))
	}
	return tr
}

// UpdateRoot sets the trie root to the current root price of
func (self *stateExchanges) updateOrderRoot(db state.Database) {
	self.updateOrdersTrie(db)
	self.data.OrderRoot = self.ordersTrie.Hash()
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *stateExchanges) CommitOrdersTrie(db state.Database) error {
	self.updateOrdersTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.ordersTrie.Commit(nil)
	if err == nil {
		self.data.OrderRoot = root
	}
	return err
}
func (self *stateExchanges) deepCopy(db *StateDB, onDirty func(hash common.Hash)) *stateExchanges {
	stateExchanges := newStateExchanges(db, self.hash, self.data, onDirty)
	if self.asksTrie != nil {
		stateExchanges.asksTrie = db.db.CopyTrie(self.asksTrie)
	}
	if self.bidsTrie != nil {
		stateExchanges.bidsTrie = db.db.CopyTrie(self.bidsTrie)
	}
	if self.ordersTrie != nil {
		stateExchanges.ordersTrie = db.db.CopyTrie(self.ordersTrie)
	}
	stateExchanges.dirtyAsksStorage = make(map[common.Hash]common.Hash)
	stateExchanges.cachedAsksStorage = make(map[common.Hash]common.Hash)
	for key, value := range self.dirtyAsksStorage {
		stateExchanges.dirtyAsksStorage[key] = value
		stateExchanges.cachedAsksStorage[key] = value
	}
	stateExchanges.dirtyBidsStorage = make(map[common.Hash]common.Hash)
	stateExchanges.cachedBidsStorage = make(map[common.Hash]common.Hash)
	for key, value := range self.dirtyBidsStorage {
		stateExchanges.dirtyBidsStorage[key] = value
		stateExchanges.cachedBidsStorage[key] = value
	}
	stateExchanges.dirtyOrdersStorage = make(map[common.Hash]common.Hash)
	stateExchanges.cachedOrdersStorage = make(map[common.Hash]common.Hash)
	for key, value := range self.dirtyOrdersStorage {
		stateExchanges.dirtyOrdersStorage[key] = value
		stateExchanges.cachedOrdersStorage[key] = value
	}
	return stateExchanges
}

// Returns the address of the contract/price
func (c *stateExchanges) Hash() common.Hash {
	return c.hash
}

func (self *stateExchanges) SetNonce(nonce uint64) {
	self.db.journal = append(self.db.journal, nonceChange{
		hash: &self.hash,
		prev: self.data.Nonce,
	})
	self.setNonce(nonce)
}

func (self *stateExchanges) setNonce(nonce uint64) {
	self.data.Nonce = nonce
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

func (self *stateExchanges) Nonce() uint64 {
	return self.data.Nonce
}

// updateStateExchangeObject writes the given object to the trie.
func (self *stateExchanges) updateStateOrderListAskObject(stateOrderList *stateOrderList) {
	price := stateOrderList.Price()
	data, err := rlp.EncodeToBytes(stateOrderList)
	if err != nil {
		panic(fmt.Errorf("can't encode order list object at %x: %v", price[:], err))
	}
	self.setError(self.asksTrie.TryUpdate(price[:], data))
}

// Retrieve a state object given my the address. Returns nil if not found.
func (self *stateExchanges) getStateOrderListAskObject(price common.Hash) (stateOrderList *stateOrderList) {
	// Prefer 'live' objects.
	if obj := self.stateAskObjects[price]; obj != nil {
		return obj
	}

	// Load the object from the database.
	enc, err := self.asksTrie.TryGet(price[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data OrderList
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state order list object", "price", price, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newStateOrderList(self.db, tomox.Bid, self.hash, price, data, self.MarkStateAskObjectDirty)
	self.setStateOrderListAskObject(obj)
	return obj
}

func (self *stateExchanges) setStateOrderListAskObject(stateOrderListObject *stateOrderList) {
	self.stateAskObjects[stateOrderListObject.Price()] = stateOrderListObject
}

// Retrieve a state object or create a new state object if nil.
func (self *stateExchanges) GetOrNewStateOrderListAskObject(price common.Hash) *stateOrderList {
	stateExchangeObject := self.getStateOrderListAskObject(price)
	if stateExchangeObject == nil {
		stateExchangeObject, _ = self.createStateOrderListAskObject(price)
	}
	return stateExchangeObject
}

// MarkStateAskObjectDirty adds the specified object to the dirty map to avoid costly
// state object cache iteration to find a handful of modified ones.
func (self *stateExchanges) MarkStateAskObjectDirty(price common.Hash) {
	self.stateAskObjectsDirty[price] = struct{}{}
}

// createStateOrderListObject creates a new state object. If there is an existing price with
// the given address, it is overwritten and returned as the second return value.
func (self *stateExchanges) createStateOrderListAskObject(price common.Hash) (newobj, prev *stateOrderList) {
	prev = self.getStateOrderListAskObject(price)
	newobj = newStateOrderList(self.db, tomox.Ask, self.hash, price, OrderList{Volume: *tomox.Zero(),}, self.MarkStateAskObjectDirty)
	if prev == nil {
		self.db.journal = append(self.db.journal, createOrderListAskChange{orderBook: &self.hash, price: &price})
	} else {
		self.db.journal = append(self.db.journal, resetOrderListAskChange{orderBook: &self.hash, prev: prev})
	}
	self.setStateOrderListAskObject(newobj)
	return newobj, prev
}

// updateStateExchangeObject writes the given object to the trie.
func (self *stateExchanges) updateStateOrderListBidObject(stateOrderList *stateOrderList) {
	price := stateOrderList.Price()
	data, err := rlp.EncodeToBytes(stateOrderList)
	if err != nil {
		panic(fmt.Errorf("can't encode order list object at %x: %v", price[:], err))
	}
	self.setError(self.bidsTrie.TryUpdate(price[:], data))
}

// Retrieve a state object given my the address. Returns nil if not found.
func (self *stateExchanges) getStateBidOrderListObject(price common.Hash) (stateOrderList *stateOrderList) {
	// Prefer 'live' objects.
	if obj := self.stateBidObjects[price]; obj != nil {
		return obj
	}

	// Load the object from the database.
	enc, err := self.bidsTrie.TryGet(price[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data OrderList
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state order list object", "price", price, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newStateOrderList(self.db, tomox.Bid, self.hash, price, data, self.MarkStateAskObjectDirty)
	self.setStateBidOrderListObject(obj)
	return obj
}

func (self *stateExchanges) setStateBidOrderListObject(stateOrderListObject *stateOrderList) {
	self.stateBidObjects[stateOrderListObject.Price()] = stateOrderListObject
}

// Retrieve a state object or create a new state object if nil.
func (self *stateExchanges) GetOrNewStateOrderListBidObject(price common.Hash) *stateOrderList {
	stateOrderListObject := self.getStateBidOrderListObject(price)
	if stateOrderListObject == nil {
		stateOrderListObject, _ = self.createStateBidOrderListObject(price)
	}
	return stateOrderListObject
}

// MarkStateAskObjectDirty adds the specified object to the dirty map to avoid costly
// state object cache iteration to find a handful of modified ones.
func (self *stateExchanges) MarkStateBidObjectDirty(price common.Hash) {
	self.stateBidObjectsDirty[price] = struct{}{}
}

// createStateOrderListObject creates a new state object. If there is an existing price with
// the given address, it is overwritten and returned as the second return value.
func (self *stateExchanges) createStateBidOrderListObject(price common.Hash) (newobj, prev *stateOrderList) {
	prev = self.getStateBidOrderListObject(price)
	newobj = newStateOrderList(self.db, tomox.Bid, self.hash, price, OrderList{Volume: *tomox.Zero()}, self.MarkStateBidObjectDirty)
	if prev == nil {
		self.db.journal = append(self.db.journal, createBidOrderListChange{orderBook: &self.hash, price: &price})
	} else {
		self.db.journal = append(self.db.journal, resetBidOrderListChange{prev: prev})
	}
	self.setStateBidOrderListObject(newobj)
	return newobj, prev
}
