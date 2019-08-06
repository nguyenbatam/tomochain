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

// Package state provides a caching layer atop the Ethereum state trie.
package state

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/tomox"
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

type revision struct {
	id           int
	journalIndex int
}

var (
	// emptyState is the known price of an empty state trie entry.
	emptyState = crypto.Keccak256Hash(nil)
	Ask        = "SELL"
	Bid        = "BUY"
)

// StateDBs within the ethereum protocol are used to store anything
// within the merkle trie. StateDBs take care of caching and storing
// nested states. It's the general query interface to retrieve:
// * Contracts
// * Accounts
type StateDB struct {
	db   state.Database
	trie state.Trie

	// This map holds 'live' objects, which will get modified while processing a state transition.
	stateExhangeObjects      map[common.Hash]*stateExchanges
	stateExhangeObjectsDirty map[common.Hash]struct{}

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        journal
	validRevisions []revision
	nextRevisionId int

	lock sync.Mutex
}

// Create a new state from a given trie.
func New(root common.Hash, db state.Database) (*StateDB, error) {
	tr, err := db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	return &StateDB{
		db:                       db,
		trie:                     tr,
		stateExhangeObjects:      make(map[common.Hash]*stateExchanges),
		stateExhangeObjectsDirty: make(map[common.Hash]struct{}),
	}, nil
}

// setError remembers the first non-nil error it is called with.
func (self *StateDB) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (self *StateDB) Error() error {
	return self.dbErr
}

// Reset clears out all ephemeral state objects from the state db, but keeps
// the underlying state trie to avoid reloading data for the next operations.
func (self *StateDB) Reset(root common.Hash) error {
	tr, err := self.db.OpenTrie(root)
	if err != nil {
		return err
	}
	self.trie = tr
	self.stateExhangeObjects = make(map[common.Hash]*stateExchanges)
	self.stateExhangeObjectsDirty = make(map[common.Hash]struct{})
	self.clearJournalAndRefund()
	return nil
}

// Exist reports whether the given price address exists in the state.
// Notably this also returns true for suicided exchanges.
func (self *StateDB) Exist(addr common.Hash) bool {
	return self.getStateExchangeObject(addr) != nil
}

// Empty returns whether the state object is either non-existent
// or empty according to the EIP161 specification (balance = nonce = code = 0)
func (self *StateDB) Empty(addr common.Hash) bool {
	so := self.getStateExchangeObject(addr)
	return so == nil || so.empty()
}

func (self *StateDB) GetNonce(addr common.Hash) uint64 {
	stateObject := self.getStateExchangeObject(addr)
	if stateObject != nil {
		return stateObject.Nonce()
	}

	return 0
}

// Database retrieves the low level database supporting the lower level trie ops.
func (self *StateDB) Database() state.Database {
	return self.db
}

/*
 * SETTERS
 */
func (self *StateDB) SetNonce(addr common.Hash, nonce uint64) {
	stateObject := self.GetOrNewStateExchangeObject(addr)
	if stateObject != nil {
		stateObject.SetNonce(nonce)
	}
}

func (self *StateDB) SetOrderItem(orderBook common.Hash, orderHash common.Hash, orderId uint64, price *big.Int, amount *big.Int, side string) {
	stateObject := self.GetOrNewStateExchangeObject(orderBook)
	if stateObject != nil {
		var stateOrderList *stateOrderList
		switch side {
		case Ask:
			stateOrderList = stateObject.GetOrNewStateOrderListAskObject(common.BigToHash(price))
		case Bid:
			stateOrderList = stateObject.GetOrNewStateOrderListBidObject(common.BigToHash(price))
		default:
			return
		}
		stateOrderList.SetOrderItem(self.db, orderId, orderHash)
		stateOrderList.AddVolume(amount)
	}
}

func (self *StateDB) GetBestAskPrice(orderBook common.Hash) (*big.Int, error) {
	stateObject := self.GetOrNewStateExchangeObject(orderBook)
	if stateObject != nil {
		price, err := stateObject.getAsksTrie(self.db).TryGetBestLeft()
		if err != nil {
			return big.NewInt(0), err
		}
		return new(big.Int).SetBytes(price), err
	}
	return big.NewInt(0), fmt.Errorf("can not get best ask price %s ", orderBook.Hex())
}

func (self *StateDB) GetBestBidPrice(orderBook common.Hash) (*big.Int, error) {
	stateObject := self.GetOrNewStateExchangeObject(orderBook)
	if stateObject != nil {
		price, err := stateObject.getBidsTrie(self.db).TryGetBestRight()
		if err != nil {
			return big.NewInt(0), err
		}
		return new(big.Int).SetBytes(price), err
	}
	return big.NewInt(0), fmt.Errorf("can not get best ask price %s ", orderBook.Hex())
}

func (self *StateDB) GetBestOrder(orderBook common.Hash, price *big.Int) (common.Hash, error) {
	stateObject := self.GetOrNewStateExchangeObject(orderBook)
	if stateObject != nil {
		orderList := stateObject.getStateOrderListAskObject(common.BigToHash(price))
		volume:=orderList.Volume()
		if volume.Cmp(tomox.Zero()) > 0 {
			if orderList != nil {
				hash, err := orderList.getTrie(self.db).TryGetBestLeft()
				if err != nil {
					return common.Hash{}, err
				}
				return common.BytesToHash(hash), nil
			}
		}
	}
	return common.Hash{}, fmt.Errorf("can not get best order : %s & price : %d", orderBook.Hex(), price)
}

// updateStateExchangeObject writes the given object to the trie.
func (self *StateDB) updateStateExchangeObject(stateObject *stateExchanges) {
	addr := stateObject.Hash()
	data, err := rlp.EncodeToBytes(stateObject)
	if err != nil {
		panic(fmt.Errorf("can't encode object at %x: %v", addr[:], err))
	}
	self.setError(self.trie.TryUpdate(addr[:], data))
}

// Retrieve a state object given my the address. Returns nil if not found.
func (self *StateDB) getStateExchangeObject(addr common.Hash) (stateObject *stateExchanges) {
	// Prefer 'live' objects.
	if obj := self.stateExhangeObjects[addr]; obj != nil {
		return obj
	}

	// Load the object from the database.
	enc, err := self.trie.TryGet(addr[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data ExchangeObject
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state object", "addr", addr, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newStateExchanges(self, addr, data, self.MarkStateExchangeObjectDirty)
	self.setStateExchangeObject(obj)
	return obj
}

func (self *StateDB) setStateExchangeObject(object *stateExchanges) {
	self.stateExhangeObjects[object.Hash()] = object
}

// Retrieve a state object or create a new state object if nil.
func (self *StateDB) GetOrNewStateExchangeObject(addr common.Hash) *stateExchanges {
	stateExchangeObject := self.getStateExchangeObject(addr)
	if stateExchangeObject == nil {
		stateExchangeObject, _ = self.createExchangeObject(addr)
	}
	return stateExchangeObject
}

// MarkStateAskObjectDirty adds the specified object to the dirty map to avoid costly
// state object cache iteration to find a handful of modified ones.
func (self *StateDB) MarkStateExchangeObjectDirty(addr common.Hash) {
	self.stateExhangeObjectsDirty[addr] = struct{}{}
}

// createStateOrderListObject creates a new state object. If there is an existing price with
// the given address, it is overwritten and returned as the second return value.
func (self *StateDB) createExchangeObject(addr common.Hash) (newobj, prev *stateExchanges) {
	prev = self.getStateExchangeObject(addr)
	newobj = newStateExchanges(self, addr, ExchangeObject{}, self.MarkStateExchangeObjectDirty)
	newobj.setNonce(0) // sets the object to dirty
	if prev == nil {
		self.journal = append(self.journal, createExchangeObjectChange{hash: &addr})
	} else {
		self.journal = append(self.journal, resetExchangeObjectChange{prev: prev})
	}
	self.setStateExchangeObject(newobj)
	return newobj, prev
}

func (db *StateDB) ForEachStorage(addr common.Hash, cb func(key, value common.Hash) bool) {
	so := db.getStateExchangeObject(addr)
	if so == nil {
		return
	}
	// When iterating over the storage check the cache first
	for h, value := range so.cachedAsksStorage {
		cb(h, value)
	}
	for h, value := range so.cachedBidsStorage {
		cb(h, value)
	}
	for h, value := range so.cachedOrdersStorage {
		cb(h, value)
	}
	it := trie.NewIterator(so.getAsksTrie(db.db).NodeIterator(nil))
	for it.Next() {
		// ignore cached values
		key := common.BytesToHash(db.trie.GetKey(it.Key))
		if _, ok := so.cachedAsksStorage[key]; !ok {
			cb(key, common.BytesToHash(it.Value))
		}
		if _, ok := so.cachedAsksStorage[key]; !ok {
			cb(key, common.BytesToHash(it.Value))
		}
	}
	it = trie.NewIterator(so.getBidsTrie(db.db).NodeIterator(nil))
	for it.Next() {
		// ignore cached values
		key := common.BytesToHash(db.trie.GetKey(it.Key))
		if _, ok := so.cachedBidsStorage[key]; !ok {
			cb(key, common.BytesToHash(it.Value))
		}
		if _, ok := so.cachedBidsStorage[key]; !ok {
			cb(key, common.BytesToHash(it.Value))
		}
	}
	it = trie.NewIterator(so.getOrdersTrie(db.db).NodeIterator(nil))
	for it.Next() {
		// ignore cached values
		key := common.BytesToHash(db.trie.GetKey(it.Key))
		if _, ok := so.cachedOrdersStorage[key]; !ok {
			cb(key, common.BytesToHash(it.Value))
		}
		if _, ok := so.cachedOrdersStorage[key]; !ok {
			cb(key, common.BytesToHash(it.Value))
		}
	}
}

// Copy creates a deep, independent copy of the state.
// Snapshots of the copied state cannot be applied to the copy.
func (self *StateDB) Copy() *StateDB {
	self.lock.Lock()
	defer self.lock.Unlock()

	// Copy all the basic fields, initialize the memory ones
	state := &StateDB{
		db:                       self.db,
		trie:                     self.db.CopyTrie(self.trie),
		stateExhangeObjects:      make(map[common.Hash]*stateExchanges, len(self.stateExhangeObjectsDirty)),
		stateExhangeObjectsDirty: make(map[common.Hash]struct{}, len(self.stateExhangeObjectsDirty)),
	}
	// Copy the dirty states, logs, and preimages
	for addr := range self.stateExhangeObjectsDirty {
		state.stateExhangeObjects[addr] = self.stateExhangeObjects[addr].deepCopy(state, state.MarkStateExchangeObjectDirty)
		state.stateExhangeObjectsDirty[addr] = struct{}{}
	}
	return state
}

// Snapshot returns an identifier for the current revision of the state.
func (self *StateDB) Snapshot() int {
	id := self.nextRevisionId
	self.nextRevisionId++
	self.validRevisions = append(self.validRevisions, revision{id, len(self.journal)})
	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (self *StateDB) RevertToSnapshot(revid int) {
	// Find the snapshot in the stack of valid snapshots.
	idx := sort.Search(len(self.validRevisions), func(i int) bool {
		return self.validRevisions[i].id >= revid
	})
	if idx == len(self.validRevisions) || self.validRevisions[idx].id != revid {
		panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	}
	snapshot := self.validRevisions[idx].journalIndex

	// Replay the journal to undo changes.
	for i := len(self.journal) - 1; i >= snapshot; i-- {
		self.journal[i].undo(self)
	}
	self.journal = self.journal[:snapshot]

	// Remove invalidated snapshots from the stack.
	self.validRevisions = self.validRevisions[:idx]
}

// Finalise finalises the state by removing the self destructed objects
// and clears the journal as well as the refunds.
func (s *StateDB) Finalise(deleteEmptyObjects bool) {
	for addr := range s.stateExhangeObjectsDirty {
		stateObject := s.stateExhangeObjects[addr]
		stateObject.updateAsksRoot(s.db)
		stateObject.updateBidsTrie(s.db)
		stateObject.updateOrdersTrie(s.db)
		s.updateStateExchangeObject(stateObject)
	}
	// Invalidate journal because reverting across transactions is not allowed.
	s.clearJournalAndRefund()
}

// IntermediateRoot computes the current root price of the state trie.
// It is called in between transactions to get the root price that
// goes into transaction receipts.
func (s *StateDB) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	s.Finalise(deleteEmptyObjects)
	return s.trie.Hash()
}

func (s *StateDB) clearJournalAndRefund() {
	s.journal = nil
	s.validRevisions = s.validRevisions[:0]
}

// Commit writes the state to the underlying in-memory trie database.
func (s *StateDB) Commit(deleteEmptyObjects bool) (root common.Hash, err error) {
	defer s.clearJournalAndRefund()

	// Commit objects to the trie.
	for addr, stateObject := range s.stateExhangeObjects {
		if _, isDirty := s.stateExhangeObjectsDirty[addr]; isDirty {
			// Write any storage changes in the state object to its storage trie.
			if err := stateObject.CommitAsksTrie(s.db); err != nil {
				return common.Hash{}, err
			}
			if err := stateObject.CommitBidsTrie(s.db); err != nil {
				return common.Hash{}, err
			}
			if err := stateObject.CommitOrdersTrie(s.db); err != nil {
				return common.Hash{}, err
			}
			// Update the object in the main price trie.
			s.updateStateExchangeObject(stateObject)
			delete(s.stateExhangeObjectsDirty, addr)
		}
	}
	// Write trie changes.
	root, err = s.trie.Commit(func(leaf []byte, parent common.Hash) error {
		var account ExchangeObject
		if err := rlp.DecodeBytes(leaf, &account); err != nil {
			return nil
		}
		if account.AskRoot != emptyState {
			s.db.TrieDB().Reference(account.AskRoot, parent)
		}
		if account.BidRoot != emptyState {
			s.db.TrieDB().Reference(account.BidRoot, parent)
		}
		if account.OrderRoot != emptyState {
			s.db.TrieDB().Reference(account.OrderRoot, parent)
		}
		return nil
	})
	log.Debug("Trie cache stats after commit", "misses", trie.CacheMisses(), "unloads", trie.CacheUnloads())
	return root, err
}
