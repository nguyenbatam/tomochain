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
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

var (
	// emptyState is the known orderId of an empty state trie entry.
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
	db   Database
	trie Trie

	// This map holds 'live' objects, which will get modified while processing a state transition.
	stateExhangeObjects      map[common.Hash]*stateExchanges
	stateExhangeObjectsDirty map[common.Hash]struct{}

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	lock sync.Mutex
}

// Create a new state from a given trie.
func New(root common.Hash, db Database) (*StateDB, error) {
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
	return nil
}

// Exist reports whether the given orderId address exists in the state.
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
func (self *StateDB) Database() Database {
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

func (self *StateDB) InsertOrderItem(orderBook common.Hash, orderId common.Hash, price *big.Int, amount *big.Int, orderType string) {
	priceHash := common.BigToHash(price)
	stateObject := self.getStateExchangeObject(orderBook)
	if stateObject == nil {
		stateObject = self.createExchangeObject(orderBook)
	}
	var stateOrderList *stateOrderList
	switch orderType {
	case Ask:
		stateOrderList = stateObject.getStateOrderListAskObject(self.db, priceHash)
		if stateOrderList == nil {
			stateOrderList, _ = stateObject.createStateOrderListAskObject(self.db, priceHash)
		}
	case Bid:
		stateOrderList = stateObject.getStateBidOrderListObject(self.db, priceHash)
		if stateOrderList == nil {
			stateOrderList, _ = stateObject.createStateBidOrderListObject(self.db, priceHash)
		}

	default:
		return
	}
	stateOrderList.insertOrderItem(self.db, orderId, common.BigToHash(amount))
	stateOrderList.AddVolume(amount)
}

func (self *StateDB) GetOrderAmount(orderBook common.Hash, orderId common.Hash, price *big.Int, orderType string) (*big.Int, error) {
	priceHash := common.BigToHash(price)
	stateObject := self.GetOrNewStateExchangeObject(orderBook)
	if stateObject == nil {
		return Zero, fmt.Errorf("Order book not found : %s ", orderBook.Hex())
	}
	var stateOrderList *stateOrderList
	switch orderType {
	case Ask:
		stateOrderList = stateObject.getStateOrderListAskObject(self.db, priceHash)
	case Bid:
		stateOrderList = stateObject.getStateBidOrderListObject(self.db, priceHash)
	default:
		return Zero, fmt.Errorf("Order type not found : %s ", orderType)
	}
	if stateOrderList.empty() {
		return Zero, fmt.Errorf("Order list empty  order book : %s , order id  : %s , price  : %s ", orderBook, orderId.Hex(), priceHash.Hex())
	}
	amountHash := stateOrderList.GetOrderAmount(self.db, orderId)
	return new(big.Int).SetBytes(amountHash[:]), nil
}
func (self *StateDB) SubAmountOrderItem(orderBook common.Hash, orderId common.Hash, price *big.Int, amount *big.Int, orderType string) error {
	priceHash := common.BigToHash(price)
	stateObject := self.GetOrNewStateExchangeObject(orderBook)
	if stateObject == nil {
		return fmt.Errorf("Order book not found : %s ", orderBook.Hex())
	}
	var stateOrderList *stateOrderList
	switch orderType {
	case Ask:
		stateOrderList = stateObject.getStateOrderListAskObject(self.db, priceHash)
	case Bid:
		stateOrderList = stateObject.getStateBidOrderListObject(self.db, priceHash)
	default:
		return fmt.Errorf("Order type not found : %s ", orderType)
	}
	if stateOrderList.empty() {
		return fmt.Errorf("Order list empty  order book : %s , order id  : %s , price  : %s ", orderBook, orderId.Hex(), priceHash.Hex())
	}
	currentAmount := new(big.Int).SetBytes(stateOrderList.GetOrderAmount(self.db, orderId).Bytes()[:])
	if currentAmount.Cmp(amount) < 0 {
		return fmt.Errorf("Order amount not enough : %s , have : %d , want : %d ", orderId.Hex(), currentAmount, amount)
	}
	newAmount := new(big.Int).Sub(currentAmount, amount)
	stateOrderList.subVolume(amount)
	if newAmount.Sign() == 0 {
		stateOrderList.removeOrderItem(self.db, orderId)
	} else {
		stateOrderList.setOrderItem(orderId, common.BigToHash(newAmount))
	}
	if stateOrderList.empty() {
		switch orderType {
		case Ask:
			stateObject.removeStateOrderListAskObject(self.db, stateOrderList)
		case Bid:
			stateObject.removeStateOrderListBidObject(self.db, stateOrderList)
		default:
		}
	}
	return nil
}

func (self *StateDB) GetVolume(orderBook common.Hash, price *big.Int, orderType string) *big.Int {
	stateObject := self.GetOrNewStateExchangeObject(orderBook)
	var volume *big.Int = nil
	if stateObject != nil {
		switch orderType {
		case Ask:
			volume = stateObject.getStateOrderListAskObject(self.db, common.BigToHash(price)).Volume()
		case Bid:
			volume = stateObject.getStateBidOrderListObject(self.db, common.BigToHash(price)).Volume()
		default:
		}
	}
	return volume
}
func (self *StateDB) GetBestAskPrice(orderBook common.Hash) (*big.Int, error) {
	stateObject := self.GetOrNewStateExchangeObject(orderBook)
	if stateObject != nil {
		price := stateObject.getBestAsksTrie(self.db)
		return price, nil
	}
	return big.NewInt(0), fmt.Errorf("can not get best ask orderId %s ", orderBook.Hex())
}

func (self *StateDB) GetBestBidPrice(orderBook common.Hash) (*big.Int, error) {
	stateObject := self.GetOrNewStateExchangeObject(orderBook)
	if stateObject != nil {
		price := stateObject.getBestBidsTrie(self.db)
		return price, nil
	}
	return big.NewInt(0), fmt.Errorf("can not get best ask orderId %s ", orderBook.Hex())
}

func (self *StateDB) GetBestOrderIdAndAmount(orderBook common.Hash, price *big.Int) (common.Hash, *big.Int, error) {
	stateObject := self.GetOrNewStateExchangeObject(orderBook)
	if stateObject != nil {
		orderList := stateObject.getStateOrderListAskObject(self.db, common.BigToHash(price))
		if orderList != nil {
			key, value, err := orderList.getTrie(self.db).TryGetBestLeftKeyAndValue()
			if err != nil {
				return EmptyHash, Zero, err
			}
			return common.BytesToHash(key), new(big.Int).SetBytes(value[:]), nil
		}
	}
	return EmptyHash, Zero, fmt.Errorf("can not get best order : %s & orderId : %d", orderBook.Hex(), price)
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
	var data exchangeObject
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state object", "addr", addr, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newStateExchanges(self, addr, data, self.MarkStateExchangeObjectDirty)
	self.stateExhangeObjects[addr] = obj
	return obj
}

func (self *StateDB) setStateExchangeObject(object *stateExchanges) {
	self.stateExhangeObjects[object.Hash()] = object
	self.stateExhangeObjectsDirty[object.Hash()] = struct{}{}
}

// Retrieve a state object or create a new state object if nil.
func (self *StateDB) GetOrNewStateExchangeObject(addr common.Hash) *stateExchanges {
	stateExchangeObject := self.getStateExchangeObject(addr)
	if stateExchangeObject == nil {
		stateExchangeObject = self.createExchangeObject(addr)
	}
	return stateExchangeObject
}

// MarkStateAskObjectDirty adds the specified object to the dirty map to avoid costly
// state object cache iteration to find a handful of modified ones.
func (self *StateDB) MarkStateExchangeObjectDirty(addr common.Hash) {
	self.stateExhangeObjectsDirty[addr] = struct{}{}
}

// createStateOrderListObject creates a new state object. If there is an existing orderId with
// the given address, it is overwritten and returned as the second return value.
func (self *StateDB) createExchangeObject(addr common.Hash) (newobj *stateExchanges) {
	newobj = newStateExchanges(self, addr, exchangeObject{}, self.MarkStateExchangeObjectDirty)
	newobj.setNonce(0) // sets the object to dirty
	self.setStateExchangeObject(newobj)
	return newobj
}

func (db *StateDB) ForEachStorage(addr common.Hash, cb func(key, value interface{}) bool) {
	//so := db.getStateExchangeObject(addr)
	//if so == nil {
	//	return
	//}
	//// When iterating over the storage check the cache first
	//for h, value := range so.cachedAsksStorage {
	//	cb(h, value)
	//}
	//for h, value := range so.cachedBidsStorage {
	//	cb(h, value)
	//}
	//it := trie.NewIterator(so.getAsksTrie(db.db).NodeIterator(nil))
	//for it.Next() {
	//	// ignore cached values
	//	key := common.BytesToHash(db.trie.GetKey(it.Key))
	//	if _, ok := so.cachedAsksStorage[key]; !ok {
	//		cb(key, common.BytesToHash(it.Quantity))
	//	}
	//	if _, ok := so.cachedAsksStorage[key]; !ok {
	//		cb(key, common.BytesToHash(it.Quantity))
	//	}
	//}
	//it = trie.NewIterator(so.getBidsTrie(db.db).NodeIterator(nil))
	//for it.Next() {
	//	// ignore cached values
	//	key := common.BytesToHash(db.trie.GetKey(it.Key))
	//	if _, ok := so.cachedBidsStorage[key]; !ok {
	//		cb(key, common.BytesToHash(it.Quantity))
	//	}
	//	if _, ok := so.cachedBidsStorage[key]; !ok {
	//		cb(key, common.BytesToHash(it.Quantity))
	//	}
	//}
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
		state.stateExhangeObjectsDirty[addr] = struct{}{}
	}
	for addr, exchangeObject := range self.stateExhangeObjects {
		state.stateExhangeObjects[addr] = exchangeObject.deepCopy(state, state.MarkStateExchangeObjectDirty)
	}

	return state
}

// Finalise finalises the state by removing the self destructed objects
// and clears the journal as well as the refunds.
func (s *StateDB) Finalise() {
	for addr, stateObject := range s.stateExhangeObjects {
		if _, exist := s.stateExhangeObjectsDirty[addr]; exist {
			delete(s.stateExhangeObjectsDirty, addr)
		}
		stateObject.updateAsksRoot(s.db)
		stateObject.updateBidsRoot(s.db)
		s.updateStateExchangeObject(stateObject)
	}
}

// IntermediateRoot computes the current root orderId of the state trie.
// It is called in between transactions to get the root orderId that
// goes into transaction receipts.
func (s *StateDB) IntermediateRoot() common.Hash {
	s.Finalise()
	return s.trie.Hash()
}

// Commit writes the state to the underlying in-memory trie database.
func (s *StateDB) Commit(deleteEmptyObjects bool) (root common.Hash, err error) {
	// Commit objects to the trie.
	for addr, stateObject := range s.stateExhangeObjects {
		if _, isDirty := s.stateExhangeObjectsDirty[addr]; isDirty {
			// Write any storage changes in the state object to its storage trie.
			if err := stateObject.CommitAsksTrie(s.db); err != nil {
				return EmptyHash, err
			}
			if err := stateObject.CommitBidsTrie(s.db); err != nil {
				return EmptyHash, err
			}
			if s.dbErr != nil {
				fmt.Println("dbError", s.dbErr)
			}
			// Update the object in the main orderId trie.
			s.updateStateExchangeObject(stateObject)
			delete(s.stateExhangeObjectsDirty, addr)
		}
	}
	// Write trie changes.
	root, err = s.trie.Commit(func(leaf []byte, parent common.Hash) error {
		var account exchangeObject
		if err := rlp.DecodeBytes(leaf, &account); err != nil {
			return nil
		}
		if account.AskRoot != emptyState {
			s.db.TrieDB().Reference(account.AskRoot, parent)
		}
		if account.BidRoot != emptyState {
			s.db.TrieDB().Reference(account.BidRoot, parent)
		}
		return nil
	})
	log.Debug("Trie cache stats after commit", "misses", trie.CacheMisses(), "unloads", trie.CacheUnloads())
	return root, err
}
