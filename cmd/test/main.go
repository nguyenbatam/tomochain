package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
	"os"
	"runtime"
)

var (
	from    = flag.String("from", "/data/tomo/chaindata", "directory to TomoChain chaindata")
	to      = flag.String("to", "/data/tomo/chaindata_copy", "directory to clean chaindata")
	root    = flag.String("root", "", "state root compare")
	address = flag.String("address", "/data/tomo/adress.txt", "list address in state db")

	sercureKey = []byte("secure-key-") // preimagePrefix + hash -> preimage
	nWorker    = runtime.NumCPU() / 2
	finish     = int32(0)
	running    = true
	emptyRoot  = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	emptyState = crypto.Keccak256Hash(nil).Bytes()
	batch      ethdb.Batch
	count      = 0
	fromDB     *ethdb.LDBDatabase
	toDB       *ethdb.LDBDatabase
	err        error
)

func main() {
	flag.Parse()
	fromDB, err = ethdb.NewLDBDatabase(*from, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	defer fromDB.Close()
	if err != nil {
		fmt.Println("fromDB", err)
		return
	}
	toDB, err = ethdb.NewLDBDatabase(*to, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	defer toDB.Close()
	if err != nil {
		fmt.Println("toDB", err)
		return
	}
	fromBC, err := core.NewBlockChain(fromDB, nil, nil, nil, vm.Config{})
	if err != nil {
		fmt.Println("fromBC", err)
		return
	}
	toBC, err := core.NewBlockChain(toDB, nil, nil, nil, vm.Config{})
	if err != nil {
		fmt.Println("toBC", err)
		return
	}
	fromState, err := fromBC.StateAt(common.HexToHash(*root))
	if err != nil {
		fmt.Println("fromState", err)
		return
	}
	toState, err := toBC.StateAt(common.HexToHash(*root))
	if err != nil {
		fmt.Println("toState", err)
		return
	}

	f, err := os.Open(*address)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		addr := common.HexToAddress(scanner.Text())
		objectFrom := fromState.GetOrNewStateObject(addr)
		byteFrom, err := rlp.EncodeToBytes(objectFrom)
		if err != nil {
			fmt.Println("objectFrom", err)
		}

		objectTo := toState.GetOrNewStateObject(addr)
		byteTo, err := rlp.EncodeToBytes(objectTo)
		if err != nil {
			fmt.Println("objectTo", err)
		}

		if bytes.Compare(byteFrom, byteTo) != 0 {
			fmt.Println("Fail when compare 2 address ", addr, common.Bytes2Hex(byteFrom), common.Bytes2Hex(byteTo))
			break
		}
		if bytes.Compare(objectFrom.CodeHash(), emptyState) != 0 {
			fromState.ForEachStorage(addr, func(key, value common.Hash) bool {
				toValue := toState.GetState(addr, key)
				if value != toValue {
					fmt.Println("Fail when compare 2 state in address ", addr, "key", key.Hex(), "fromValue", value.Hex(), "toValue", toValue.Hex())
				}
				return true
			})
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
