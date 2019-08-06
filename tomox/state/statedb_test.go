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
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/core/state"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"testing/quick"

	"gopkg.in/check.v1"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

// Tests that updating a state trie does not leak any database writes prior to
// actually committing the state.
func TestUpdateLeaks(t *testing.T) {
	// Create an empty statedb database
	db, _ := ethdb.NewMemDatabase()
	statedb, _ := New(common.Hash{}, state.NewDatabase(db))

	// Update it with some exchanges
	for i := byte(0); i < 255; i++ {
		addr := common.BytesToHash([]byte{i})
		statedb.SetNonce(addr, uint64(42*i))
		statedb.IntermediateRoot(false)
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
	transState, _ := New(common.Hash{}, state.NewDatabase(transDb))
	finalState, _ := New(common.Hash{}, state.NewDatabase(finalDb))

	modify := func(state *StateDB, addr common.Hash, i, tweak byte) {
		state.SetNonce(addr, uint64(42*i+tweak))
	}

	// Modify the transient state.
	for i := byte(0); i < 255; i++ {
		modify(transState, common.Hash{byte(i)}, i, 0)
	}
	// Write modifications to trie.
	transState.IntermediateRoot(false)

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
	orig, _ := New(common.Hash{}, state.NewDatabase(db))

	for i := byte(0); i < 255; i++ {
		obj := orig.GetOrNewStateExchangeObject(common.BytesToHash([]byte{i}))
		obj.SetNonce(uint64(i))
		orig.updateStateExchangeObject(obj)
	}
	orig.Finalise(false)

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
		orig.Finalise(true)
		close(done)
	}()
	copy.Finalise(true)
	<-done

	// Verify that the two states have been updated independently
	for i := byte(0); i < 255; i++ {
		origObj := orig.GetOrNewStateExchangeObject(common.BytesToHash([]byte{i}))
		copyObj := copy.GetOrNewStateExchangeObject(common.BytesToHash([]byte{i}))

		if want := 3 * uint64(i); origObj.Nonce() == want {
			t.Errorf("orig obj %d: balance mismatch: have %v, want %v", i, origObj.Nonce(), want)
		}
		if want := 4 * uint64(i); copyObj.Nonce() == want {
			t.Errorf("copy obj %d: balance mismatch: have %v, want %v", i, copyObj.Nonce(), want)
		}
	}
}

func TestSnapshotRandom(t *testing.T) {
	config := &quick.Config{MaxCount: 1000}
	err := quick.Check((*snapshotTest).run, config)
	if cerr, ok := err.(*quick.CheckError); ok {
		test := cerr.In[0].(*snapshotTest)
		t.Errorf("%v:\n%s", test.err, test)
	} else if err != nil {
		t.Error(err)
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
type snapshotTest struct {
	addrs     []common.Hash // all price addresses
	actions   []testAction     // modifications to the state
	snapshots []int            // actions indexes at which snapshot is taken
	err       error            // failure details are reported through this field
}

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

// Generate returns a new snapshot test of the given size. All randomness is
// derived from r.
func (*snapshotTest) Generate(r *rand.Rand, size int) reflect.Value {
	// Generate random actions.
	addrs := make([]common.Hash, 50)
	for i := range addrs {
		addrs[i][0] = byte(i)
	}
	actions := make([]testAction, size)
	for i := range actions {
		addr := addrs[r.Intn(len(addrs))]
		actions[i] = newTestAction(addr, r)
	}
	// Generate snapshot indexes.
	nsnapshots := int(math.Sqrt(float64(size)))
	if size > 0 && nsnapshots == 0 {
		nsnapshots = 1
	}
	snapshots := make([]int, nsnapshots)
	snaplen := len(actions) / nsnapshots
	for i := range snapshots {
		// Try to place the snapshots some number of actions apart from each other.
		snapshots[i] = (i * snaplen) + r.Intn(snaplen)
	}
	return reflect.ValueOf(&snapshotTest{addrs, actions, snapshots, nil})
}

func (test *snapshotTest) String() string {
	out := new(bytes.Buffer)
	sindex := 0
	for i, action := range test.actions {
		if len(test.snapshots) > sindex && i == test.snapshots[sindex] {
			fmt.Fprintf(out, "---- snapshot %d ----\n", sindex)
			sindex++
		}
		fmt.Fprintf(out, "%4d: %s\n", i, action.name)
	}
	return out.String()
}

func (test *snapshotTest) run() bool {
	// Run all actions and create snapshots.
	var (
		db, _        = ethdb.NewMemDatabase()
		state, _     = New(common.Hash{}, state.NewDatabase(db))
		snapshotRevs = make([]int, len(test.snapshots))
		sindex       = 0
	)
	for i, action := range test.actions {
		if len(test.snapshots) > sindex && i == test.snapshots[sindex] {
			snapshotRevs[sindex] = state.Snapshot()
			sindex++
		}
		action.fn(action, state)
	}
	// Revert all snapshots in reverse order. Each revert must yield a state
	// that is equivalent to fresh state with all actions up the snapshot applied.
	for sindex--; sindex >= 0; sindex-- {
		checkstate, _ := New(common.Hash{}, state.Database())
		for _, action := range test.actions[:test.snapshots[sindex]] {
			action.fn(action, checkstate)
		}
		state.RevertToSnapshot(snapshotRevs[sindex])
		if err := test.checkEqual(state, checkstate); err != nil {
			test.err = fmt.Errorf("state mismatch after revert to snapshot %d\n%v", sindex, err)
			return false
		}
	}
	return true
}

// checkEqual checks that methods of state and checkstate return the same values.
func (test *snapshotTest) checkEqual(state, checkstate *StateDB) error {
	for _, addr := range test.addrs {
		var err error
		checkeq := func(op string, a, b interface{}) bool {
			if err == nil && !reflect.DeepEqual(a, b) {
				err = fmt.Errorf("got %s(%s) == %v, want %v", op, addr.Hex(), a, b)
				return false
			}
			return true
		}
		// Check basic accessor methods.
		checkeq("Exist", state.Exist(addr), checkstate.Exist(addr))
		checkeq("GetNonce", state.GetNonce(addr), checkstate.GetNonce(addr))
		// Check storage.
		//if obj := state.getStateExchangeObject(addr); obj != nil {
		//	state.ForEachStorage(addr, func(price, val common.Hash) bool {
		//		return checkeq("GetState("+price.Hex()+")", val, checkstate.GetState(addr, price))
		//	})
		//	checkstate.ForEachStorage(addr, func(price, checkval common.Hash) bool {
		//		return checkeq("GetState("+price.Hex()+")", state.GetState(addr, price), checkval)
		//	})
		//}
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StateSuite) TestTouchDelete(c *check.C) {
	s.state.GetOrNewStateExchangeObject(common.Hash{})
	root, _ := s.state.Commit(false)
	s.state.Reset(root)

	snapshot := s.state.Snapshot()
	s.state.SetNonce(common.Hash{}, 0)
	if len(s.state.stateExhangeObjectsDirty) != 1 {
		c.Fatal("expected one dirty state object")
	}
	s.state.RevertToSnapshot(snapshot)
	if len(s.state.stateExhangeObjectsDirty) != 0 {
		c.Fatal("expected no dirty state object")
	}
}
