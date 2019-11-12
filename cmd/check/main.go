package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"os"
	"runtime"
)

var (
	from    = flag.String("from", "/data/tomo/chaindata", "directory to TomoChain chaindata")
	to      = flag.String("to", "/data/tomo/copy", "directory to clean chaindata")
	root    = flag.String("root", "0xc96e205f8e0d7dccea94c04fde6e6f7e508a3bbb91d2630335b37fbe23ec3a87", "state root compare")
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
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stdout, log.TerminalFormat(true))))
	fmt.Println("flag")
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
	fromState, err := state.New(common.HexToHash(*root), state.NewDatabase(fromDB))
	if err != nil {
		fmt.Println("fromState", *root, err)
		return
	}
	toStateCache := state.NewDatabase(toDB)
	toState, err := state.NewEmpty(common.HexToHash(*root), toStateCache)
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
			fmt.Println("Fail when compare 2 address ", addr.Hex(), common.Bytes2Hex(byteFrom), common.Bytes2Hex(byteTo))
			break
		}
		fmt.Println("addr", addr.Hex())
		check := fromState.ForEachStorageAndCheck(addr, func(key, value common.Hash) bool {
			value = fromState.GetStateNotCache(addr, key)
			toObject := toState.GetStateObjectNotCache(addr)
			if toObject == nil {
				fmt.Println("Fail when get state in address ", addr.Hex(), toState.Error())
				return false
			}
			toValue := toObject.GetStateNotCache(toState.Database(), key)
			if bytes.Compare(toValue.Bytes(), value.Bytes()) != 0 {
				fmt.Println("Fail when compare 2 state in address ", addr.Hex(), "key", key.Hex(), "decode", value.Hex(), "toValue", toValue.Hex(), "err", toState.Error())
				return false
			}
			return true
		})
		if !check {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Crit("scan", "err", err)
	}
}
