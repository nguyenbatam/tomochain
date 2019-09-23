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
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tomochain/tomox-sdk/types"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"testing"
)

// Tests that updating a state trie does not leak any database writes prior to
// actually committing the state.
func TestUpdateLeaks(t *testing.T) {
	// Create an empty statedb database
	db, _ := ethdb.NewMemDatabase()
	statedb, _ := New(common.Hash{}, NewDatabase(db))

	// Update it with some exchanges
	for i := byte(0); i < 255; i++ {
		addr := common.BytesToHash([]byte{i})
		statedb.SetNonce(addr, uint64(42*i))
		statedb.IntermediateRoot()
	}
	// Ensure that no data was leaked into the database
	for _, key := range db.Keys() {
		value, _ := db.Get(key)
		t.Errorf("State leaked into database: %x -> %x", key, value)
	}
}

// Tests that no intermediate state of an object is stored into the database,
// only the one right before the commit.
func TestIntermediateLeaks(t *testing.T) {
	// Create two state databases, one transitioning to the final state, the other final from the beginning
	transDb, _ := ethdb.NewMemDatabase()
	finalDb, _ := ethdb.NewMemDatabase()
	transState, _ := New(common.Hash{}, NewDatabase(transDb))
	finalState, _ := New(common.Hash{}, NewDatabase(finalDb))

	modify := func(state *StateDB, addr common.Hash, i, tweak byte) {
		state.SetNonce(addr, uint64(42*i+tweak))
	}

	// Modify the transient state.
	for i := byte(0); i < 255; i++ {
		modify(transState, common.Hash{byte(i)}, i, 0)
	}
	// Write modifications to trie.
	transState.IntermediateRoot()

	// Overwrite all the data with new values in the transient database.
	for i := byte(0); i < 255; i++ {
		modify(transState, common.Hash{byte(i)}, i, 99)
		modify(finalState, common.Hash{byte(i)}, i, 99)
	}

	// Commit and cross check the databases.
	if _, err := transState.Commit(false); err != nil {
		t.Fatalf("failed to commit transition state: %v", err)
	}
	if _, err := finalState.Commit(false); err != nil {
		t.Fatalf("failed to commit final state: %v", err)
	}
	for _, key := range finalDb.Keys() {
		if _, err := transDb.Get(key); err != nil {
			val, _ := finalDb.Get(key)
			t.Errorf("entry missing from the transition database: %x -> %x", key, val)
		}
	}
	for _, key := range transDb.Keys() {
		if _, err := finalDb.Get(key); err != nil {
			val, _ := transDb.Get(key)
			t.Errorf("extra entry in the transition database: %x -> %x", key, val)
		}
	}
}

// TestCopy tests that copying a statedb object indeed makes the original and
// the copy independent of each other. This test is a regression test against
// https://github.com/ethereum/go-ethereum/pull/15549.
func TestCopy(t *testing.T) {
	// Create a random state test to copy and modify "independently"
	db, _ := ethdb.NewMemDatabase()
	orig, _ := New(common.Hash{}, NewDatabase(db))

	for i := byte(0); i < 255; i++ {
		obj := orig.GetOrNewStateExchangeObject(common.BytesToHash([]byte{i}))
		obj.SetNonce(uint64(i))
		orig.updateStateExchangeObject(obj)
	}
	orig.Finalise()

	// Copy the state, modify both in-memory
	copy := orig.Copy()

	for i := byte(0); i < 255; i++ {
		origObj := orig.GetOrNewStateExchangeObject(common.BytesToHash([]byte{i}))
		copyObj := copy.GetOrNewStateExchangeObject(common.BytesToHash([]byte{i}))

		origObj.SetNonce(2 * uint64(i))
		copyObj.SetNonce(3 * uint64(i))

		orig.updateStateExchangeObject(origObj)
		copy.updateStateExchangeObject(copyObj)
	}
	// Finalise the changes on both concurrently
	done := make(chan struct{})
	go func() {
		orig.Finalise()
		close(done)
	}()
	copy.Finalise()
	<-done

	// Verify that the two states have been updated independently
	for i := byte(0); i < 255; i++ {
		origObj := orig.GetOrNewStateExchangeObject(common.BytesToHash([]byte{i}))
		copyObj := copy.GetOrNewStateExchangeObject(common.BytesToHash([]byte{i}))

		if want := 2 * uint64(i); origObj.Nonce() != want {
			t.Errorf("orig obj %d: balance mismatch: have %v, want %v", i, origObj.Nonce(), want)
		}
		if want := 3 * uint64(i); copyObj.Nonce() != want {
			t.Errorf("copy obj %d: balance mismatch: have %v, want %v", i, copyObj.Nonce(), want)
		}
	}
}

// A snapshotTest checks that reverting StateDB snapshots properly undoes all changes
// captured by the snapshot. Instances of this test with pseudorandom content are created
// by Generate.
//
// The test works as follows:
//
// A new state is created and all actions are applied to it. Several snapshots are taken
// in between actions. The test then reverts each snapshot. For each snapshot the actions
// leading up to it are replayed on a fresh, empty state. The behaviour of all public
// accessor methods on the reverted state must match the return value of the equivalent
// methods on the replayed state.

type testAction struct {
	name   string
	fn     func(testAction, *StateDB)
	args   []int64
	noAddr bool
}

// newTestAction creates a random action that changes state.
func newTestAction(addr common.Hash, r *rand.Rand) testAction {
	actions := []testAction{
		{
			name: "SetBalance",
			fn: func(a testAction, s *StateDB) {
				s.SetNonce(addr, uint64(a.args[0]))
			},
			args: make([]int64, 1),
		},
		{
			name: "AddBalance",
			fn: func(a testAction, s *StateDB) {
				s.SetNonce(addr, uint64(a.args[0]))
			},
			args: make([]int64, 1),
		},
		{
			name: "SetNonce",
			fn: func(a testAction, s *StateDB) {
				s.SetNonce(addr, uint64(a.args[0]))
			},
			args: make([]int64, 1),
		},
		{
			name: "CreateAccount",
			fn: func(a testAction, s *StateDB) {
				s.createExchangeObject(addr)
			},
		},
	}
	action := actions[r.Intn(len(actions))]
	var nameargs []string
	if !action.noAddr {
		nameargs = append(nameargs, addr.Hex())
	}
	for _, i := range action.args {
		action.args[i] = rand.Int63n(100)
		nameargs = append(nameargs, fmt.Sprint(action.args[i]))
	}
	action.name += strings.Join(nameargs, ", ")
	return action
}

func TestEchangeStates(t *testing.T) {
	dir := "/home/tamnb/_projects/tomochain/src/github.com/ethereum/go-ethereum/devnet/test"
	os.RemoveAll(dir)
	orderBook := common.StringToHash("BTC/TOMO")
	numberOrder := 2500;
	orderItems := []orderItem{}
	relayers := []common.Hash{}
	for i := 0; i < numberOrder; i++ {
		relayers = append(relayers, common.BigToHash(big.NewInt(int64(i))))
		orderItems = append(orderItems, orderItem{Amount: big.NewInt(int64(2*i + 1)), OrderType: types.SELL})
		orderItems = append(orderItems, orderItem{Amount: big.NewInt(int64(2*i + 1)), OrderType: types.BUY})
	}
	// Create an empty statedb database
	db, _ := ethdb.NewLDBDatabase(dir, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	stateCache := NewDatabase(db)
	statedb, _ := New(common.Hash{}, stateCache)

	// Update it with some exchanges
	for i := 0; i < numberOrder; i++ {
		statedb.SetNonce(relayers[i], uint64(1))
	}
	mapPriceSell := map[uint64]uint64{}
	mapPriceBuy := map[uint64]uint64{}
	orderBookId := uint64(1)
	for i := 0; i < len(orderItems); i++ {
		id := new(big.Int).SetUint64(orderBookId)
		orderId := common.BigToHash(id)
		amount := orderItems[i].Amount.Uint64()
		statedb.InsertOrderItem(orderBook, orderId, orderId, orderItems[i].Amount, orderItems[i].Amount, orderItems[i].OrderType)
		statedb.SetNonce(orderBook, orderBookId)
		orderBookId = orderBookId + 1

		switch orderItems[i].OrderType {
		case types.SELL:
			old := mapPriceSell[amount]
			mapPriceSell[amount] = old + amount
		case types.BUY:
			old := mapPriceBuy[amount]
			mapPriceBuy[amount] = old + amount
		default:
		}
	}

	copyStateDb := statedb.Copy()
	copyStateDb.SubAmountOrderItem(orderBook, common.BigToHash(new(big.Int).SetUint64(1)), orderItems[0].Amount, orderItems[0].Amount, orderItems[0].OrderType)
	orderId := new(big.Int).SetUint64(1)
	fmt.Println(copyStateDb.GetOrderAmount(orderBook, common.BigToHash(orderId), orderId, orderItems[0].OrderType))
	fmt.Println(statedb.GetOrderAmount(orderBook, common.BigToHash(orderId), orderId, orderItems[0].OrderType))
	root, err := statedb.Commit(false)
	if err != nil {
		t.Fatalf("Error when commit into database: %v", err)
	}
	fmt.Println("root", root.Hex())
	err = stateCache.TrieDB().Commit(root, false)
	if err != nil {
		t.Errorf("Error when commit into database: %v", err)
	}
	db.Close()

	db, _ = ethdb.NewLDBDatabase("/home/tamnb/_projects/tomochain/src/github.com/ethereum/go-ethereum/devnet/test", eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	stateCache = NewDatabase(db)
	statedb, err = New(root, stateCache)
	if err != nil {
		t.Fatalf("Error when get trie in database: %s , err: %v", root.Hex(), err)
	}
	for i := 0; i < numberOrder; i++ {
		nonce := statedb.GetNonce(relayers[i])
		if nonce != uint64(1) {
			t.Fatalf("Error when get nonce save in database: got : %d , wanted : %d ", nonce, i)
		}
	}
	statedb.GetNonce(orderBook)
	for i := 0; i < len(orderItems); i++ {
		id := new(big.Int).SetUint64(uint64(i + 1))
		orderId := common.BigToHash(id)
		orderItem := statedb.GetOrderItem(orderBook, orderId)
		if orderItem == nil {
			t.Fatalf("Error == nil when get Order Item save in database: orderId %s ", orderId.Hex())
		}
		if orderItem.Amount.Cmp(orderItems[i].Amount) != 0 {
			t.Fatalf("Error when get nonce save in database: orderId %s ,got : %d , wanted : %d ", orderId.Hex(), orderItem.Amount.Uint64(), orderItems[i].Amount.Uint64())
		}
	}

	minSell := uint64(math.MaxUint64)
	for price, amount := range mapPriceSell {
		data := statedb.GetVolume(orderBook, new(big.Int).SetUint64(price), Ask)
		if data.Uint64() != amount {
			t.Fatalf("Error when get volume save in database: price  %d ,got : %d , wanted : %d ", price, data.Uint64(), amount)
		}
		if price < minSell {
			minSell = price
		}
	}
	maxBuy := uint64(0)
	for price, amount := range mapPriceBuy {
		data := statedb.GetVolume(orderBook, new(big.Int).SetUint64(price), Bid)
		if data.Uint64() != amount {
			t.Fatalf("Error when get volume save in database: price  %d ,got : %d , wanted : %d ", price, data.Uint64(), amount)
		}
		if price > maxBuy {
			maxBuy = price
		}
	}
	fmt.Println("===============>")
	bestAsk, err := statedb.GetBestAskPrice(orderBook)
	if err != nil {
		t.Fatalf("Error when get best ask trie in orderBook: %s , err: %v", orderBook.Hex(), err)
	}
	fmt.Println("===============>")
	bestBid, err := statedb.GetBestBidPrice(orderBook)
	if err != nil {
		t.Fatalf("Error when get best bid trie in orderBook: %s , err: %v", orderBook.Hex(), err)
	}
	fmt.Println("best price ", bestBid, bestAsk, minSell, maxBuy)
	db.Close()
}

func TestMemmory(t *testing.T) {
	dir := "/home/tamnb/_projects/tomochain/src/github.com/ethereum/go-ethereum/devnet/test"
	err := os.RemoveAll(dir)
	orderBook := common.StringToHash("BTC/TOMO")
	numberOrder := 500000;
	orderItems := []orderItem{}
	relayers := []common.Hash{}
	for i := 0; i < numberOrder; i++ {
		relayers = append(relayers, common.BigToHash(big.NewInt(int64(i))))
		rand := rand.Intn(numberOrder/10 + i)
		orderItems = append(orderItems, orderItem{Amount: big.NewInt(int64(rand)), OrderType: types.SELL})
		orderItems = append(orderItems, orderItem{Amount: big.NewInt(int64(rand)), OrderType: types.BUY})
	}
	// Create an empty statedb database
	db, _ := ethdb.NewLDBDatabase(dir, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	stateCache := NewDatabase(db)
	statedb, _ := New(common.Hash{}, stateCache)

	// Update it with some exchanges
	for i := 0; i < numberOrder; i++ {
		statedb.SetNonce(relayers[i], uint64(1))
	}
	mapPriceSell := map[uint64]uint64{}
	mapPriceBuy := map[uint64]uint64{}
	orderBookId := uint64(1)
	for i := 0; i < len(orderItems); i++ {
		id := new(big.Int).SetUint64(orderBookId)
		orderId := common.BigToHash(id)
		amount := orderItems[i].Amount.Uint64()
		statedb.InsertOrderItem(orderBook, orderId, orderId, orderItems[i].Amount, orderItems[i].Amount, orderItems[i].OrderType)
		statedb.SetNonce(orderBook, orderBookId)
		orderBookId = orderBookId + 1

		switch orderItems[i].OrderType {
		case types.SELL:
			old := mapPriceSell[amount]
			mapPriceSell[amount] = old + amount
		case types.BUY:
			old := mapPriceBuy[amount]
			mapPriceBuy[amount] = old + amount
		default:
		}
	}
	fmt.Println("finish insert")
	//copyStateDb := statedb.Copy()
	//copyStateDb.SubAmountOrderItem(orderBook,common.BigToHash(new(big.Int).SetUint64(1)),orderItems[0].Amount,orderItems[0].Amount,orderItems[0].OrderType)
	//fmt.Println(copyStateDb.GetOrderAmount(orderBook,common.BigToHash(new(big.Int).SetUint64(1))))
	//fmt.Println(statedb.GetOrderAmount(orderBook,common.BigToHash(new(big.Int).SetUint64(1))))
	root, err := statedb.Commit(false)
	if err != nil {
		t.Fatalf("Error when commit into database: %v", err)
	}
	fmt.Println("root", root.Hex())
	err = stateCache.TrieDB().Commit(root, false)
	if err != nil {
		t.Errorf("Error when commit into database: %v", err)
	}
	db.Close()
	//
	//db, _ = ethdb.NewLDBDatabase("/home/tamnb/_projects/tomochain/src/github.com/ethereum/go-ethereum/devnet/test", eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	//stateCache = NewDatabase(db)
	//statedb, err = New(root, stateCache)
	//if err != nil {
	//	t.Fatalf("Error when get trie in database: %s , err: %v", root.Hex(), err)
	//}
	//for i := 0; i < numberOrder; i++ {
	//	nonce := statedb.GetNonce(relayers[i])
	//	if nonce != uint64(1) {
	//		t.Fatalf("Error when get nonce save in database: got : %d , wanted : %d ", nonce, i)
	//	}
	//}
	//statedb.GetNonce(orderBook)
	//for i := 0; i < len(orderItems); i++ {
	//	id := new(big.Int).SetUint64(uint64(i + 1))
	//	orderId := common.BigToHash(id)
	//	orderItem := statedb.GetOrderAmount(orderBook, orderId)
	//	if orderItem == nil {
	//		t.Fatalf("Error == nil when get Order Item save in database: orderId %s ", orderId.Hex())
	//	}
	//	if orderItem.Amount.Cmp(orderItems[i].Amount) != 0 {
	//		t.Fatalf("Error when get nonce save in database: orderId %s ,got : %d , wanted : %d ", orderId.Hex(), orderItem.Amount.Uint64(), orderItems[i].Amount.Uint64())
	//	}
	//}
	//
	//minSell := uint64(math.MaxUint64)
	//for price, amount := range mapPriceSell {
	//	data := statedb.GetVolume(orderBook, new(big.Int).SetUint64(price), Ask)
	//	if data.Uint64() != amount {
	//		t.Fatalf("Error when get volume save in database: price  %d ,got : %d , wanted : %d ", price, data.Uint64(), amount)
	//	}
	//	if price < minSell {
	//		minSell = price
	//	}
	//}
	//maxBuy := uint64(0)
	//for price, amount := range mapPriceBuy {
	//	data := statedb.GetVolume(orderBook, new(big.Int).SetUint64(price), Bid)
	//	if data.Uint64() != amount {
	//		t.Fatalf("Error when get volume save in database: price  %d ,got : %d , wanted : %d ", price, data.Uint64(), amount)
	//	}
	//	if price > maxBuy {
	//		maxBuy = price
	//	}
	//}
	//fmt.Println("===============>")
	//bestAsk, err := statedb.GetBestAskPrice(orderBook)
	//if err != nil {
	//	t.Fatalf("Error when get best ask trie in orderBook: %s , err: %v", orderBook.Hex(), err)
	//}
	//fmt.Println("===============>")
	//bestBid, err := statedb.GetBestBidPrice(orderBook)
	//if err != nil {
	//	t.Fatalf("Error when get best bid trie in orderBook: %s , err: %v", orderBook.Hex(), err)
	//}
	//fmt.Println("best price ", bestBid, bestAsk, minSell, maxBuy)
	//db.Close()
}
