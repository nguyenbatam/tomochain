package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"runtime"
	"time"
)

var (
	from       = flag.String("from", "/data/tomo/chaindata", "directory to TomoChain chaindata")
	to         = flag.String("to", "/data/tomo/chaindata_copy", "directory to clean chaindata")
	length     = flag.Uint64("length", 100, "minimum backup block data")
	sercureKey = []byte("secure-key-") // preimagePrefix + hash -> preimage
	nWorker    = runtime.NumCPU() / 2
	finish     = int32(0)
	running    = true
	stateRoots = make(chan TrieRoot)
	emptyRoot  = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	emptyState = crypto.Keccak256Hash(nil)
	batch      ethdb.Batch
	count      = 0
)

type TrieRoot struct {
	trie   *trie.SecureTrie
	number uint64
}
type StateNode struct {
	node trie.Node
	path []byte
}

func main() {
	flag.Parse()
	fromDB, _ := ethdb.NewLDBDatabase(*from, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	tridb := trie.NewDatabase(fromDB)
	head := core.GetHeadBlockHash(fromDB)
	header := core.GetHeader(fromDB, head, core.GetBlockNumber(fromDB, head))
	number := header.Number.Uint64()
	lastestRoot := common.Hash{}
	lastestRootNumber := uint64(0)
	backupRoot := common.Hash{}
	backupNumber := number
	for number >= 1 {
		number = number - 1
		hash := core.GetCanonicalHash(fromDB, number)
		root := core.GetHeader(fromDB, hash, number).Root
		_, err := trie.NewSecure(root, tridb, 0)
		if err != nil {
			continue
		}
		if common.EmptyHash(lastestRoot) {
			lastestRoot = root
			lastestRootNumber = number
		} else if root != lastestRoot && number < lastestRootNumber-*length {
			backupRoot = root
			backupNumber = number
			break
		}
	}
	fmt.Println("lastestRoot", lastestRoot.Hex(), "lastestRootNumber", lastestRootNumber, "backupRoot", backupRoot.Hex(), "backupNumber", backupNumber, "currentNumber", header.Number.Uint64())
	fromDB.Close()
	err := copyHeadData(*from, *to)
	if err != nil {
		fmt.Println("copyHeadData", err)
		return
	}
	err = copyBlockData(*from, *to, backupNumber)
	if err != nil {
		fmt.Println("copyBlockData", err)
		return
	}
	err = copyStateData(*from, *to, lastestRoot)
	if err != nil {
		fmt.Println("copyStateData lastestRoot", lastestRoot.Hex(), "err", err)
		return
	}
	//err = copyStateData(*from, *to, backupRoot)
	//if err != nil {
	//	fmt.Println("copyStateData backupRoot", backupRoot.Hex(),"err",err)
	//	return
	//}
}
func copyHeadData(from string, to string) error {
	fmt.Println(time.Now(), "copyHeadData")
	fromDB, err := ethdb.NewLDBDatabase(from, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	defer fromDB.Close()
	if err != nil {
		return err
	}
	toDB, err := ethdb.NewLDBDatabase(to, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	defer toDB.Close()
	if err != nil {
		return err
	}
	//headHeaderKey = []byte("LastHeader")
	hash := core.GetHeadHeaderHash(fromDB)
	core.WriteHeadHeaderHash(toDB, hash)
	//headBlockKey  = []byte("LastBlock")
	hash = core.GetHeadBlockHash(fromDB)
	core.WriteHeadBlockHash(toDB, hash)
	//headFastKey   = []byte("LastFast")
	hash = core.GetHeadFastBlockHash(fromDB)
	core.WriteHeadFastBlockHash(toDB, hash)
	//trieSyncKey   = []byte("TrieSync")
	trie := core.GetTrieSyncProgress(fromDB)
	core.WriteTrieSyncProgress(toDB, trie)
	fmt.Println(time.Now(), "compact")
	toDB.LDB().CompactRange(util.Range{})
	fmt.Println(time.Now(), "end")
	return nil
}
func copyBlockData(from string, to string, backupNumber uint64) error {
	fmt.Println(time.Now(), "copyBlockData", "backupNumber", backupNumber)
	fromDB, err := ethdb.NewLDBDatabase(from, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	defer fromDB.Close()
	if err != nil {
		return err
	}
	toDB, err := ethdb.NewLDBDatabase(to, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	defer toDB.Close()
	if err != nil {
		return err
	}
	head := core.GetHeadBlockHash(fromDB)
	header := core.GetHeader(fromDB, head, core.GetBlockNumber(fromDB, head))
	number := header.Number.Uint64()
	for number >= backupNumber {
		hash := header.Hash()
		//bodyPrefix          = []byte("b") // bodyPrefix + num (uint64 big endian) + hash -> block body
		//blockHashPrefix     = []byte("H") // blockHashPrefix + hash -> num (uint64 big endian)
		//headerPrefix        = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
		block := core.GetBlock(fromDB, hash, number)
		core.WriteBlock(toDB, block)
		//tdSuffix            = []byte("t") // headerPrefix + num (uint64 big endian) + hash + tdSuffix -> td
		td := core.GetTd(fromDB, hash, number)
		core.WriteTd(toDB, hash, number, td)
		//numSuffix           = []byte("n") // headerPrefix + num (uint64 big endian) + numSuffix -> hash
		hash = core.GetCanonicalHash(fromDB, number)
		core.WriteCanonicalHash(toDB, hash, number)
		if number == 0 {
			break
		}
		header = core.GetHeader(fromDB, block.ParentHash(), number-1)
		number = header.Number.Uint64()
	}
	fmt.Println(time.Now(), "compact")
	toDB.LDB().CompactRange(util.Range{})
	fmt.Println(time.Now(), "end")
	return nil
}

func copyStateData(from, to string, root common.Hash) error {
	fmt.Println(time.Now(), "run copy state data ", "root", root.Hex())
	fromDB, err := ethdb.NewLDBDatabase(from, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	defer fromDB.Close()
	if err != nil {
		return err
	}
	toDB, err := ethdb.NewLDBDatabase(to, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	defer toDB.Close()
	if err != nil {
		return err
	}
	batch = toDB.NewBatch()
	rootNode, valueDB, err := resolveHash(root[:], fromDB.LDB())
	if err != nil {
		return err
	}

	err = processNode(rootNode, nil, fromDB.LDB(), true)
	if err != nil {
		return err
	}
	err = toDB.LDB().Put(root[:], valueDB, nil)
	if err != nil {
		return err
	}
	err = batch.Write()
	if err != nil {
		return err
	}
	return nil
}
func putToData(key []byte, value []byte) {
	batch.Put(key, value)
	count++
	if count%1000 == 0 {
		err := batch.Write()
		count = 0
		if err != nil {
			fmt.Println("Error when put data to copy db")
			panic(err)
		}
	}
}
func processNode(n trie.Node, path []byte, fromDB *leveldb.DB, checkAddr bool) error {
	switch node := n.(type) {
	case *trie.FullNode:
		// Full Node, move to the first non-nil child.
		for i := 0; i < len(node.Children); i++ {
			child := node.Children[i]
			if child != nil {
				childNode := child
				var err error = nil
				var valueDB []byte
				if _, ok := child.(trie.HashNode); ok {
					childNode, valueDB, err = resolveHash(child.(trie.HashNode), fromDB)
				}
				if err == nil {
					err = processNode(childNode, append(path, byte(i)), fromDB, checkAddr)
					if err != nil {
						return err
					}
					putToData(child.(trie.HashNode), valueDB)
				} else if err != nil {
					_, ok := err.(*trie.MissingNodeError)
					if !ok {
						return err
					}
				}
			}
		}
	case *trie.ShortNode:
		// Short Node, return the pointer singleton child
		childNode := node.Val
		var err error = nil
		var valueDB []byte
		if _, ok := node.Val.(trie.HashNode); ok {
			childNode, valueDB, err = resolveHash(node.Val.(trie.HashNode), fromDB)
		}
		if err == nil {
			err = processNode(childNode, append(path, node.Key...), fromDB, checkAddr)
			if err != nil {
				return err
			}
			if _, ok := node.Val.(trie.HashNode); ok {
				putToData(node.Val.(trie.HashNode), valueDB)
			}
		} else if err != nil {
			_, ok := err.(*trie.MissingNodeError)
			if !ok {
				return err
			}
		}
	case trie.ValueNode:
		if checkAddr {
			var data state.Account
			if err := rlp.DecodeBytes(node, &data); err != nil {
				fmt.Println("Failed to decode state object", "path", common.Bytes2Hex(path), "value", common.Bytes2Hex(node))
				return err
			}
			if common.EmptyHash(data.Root) && data.Root != emptyRoot && data.Root != emptyState {
				fmt.Println("Try copy data in a smart contract ")
				newNode, valueDB, err := resolveHash(data.Root[:], fromDB)
				if err != nil {
					return err
				}
				err = processNode(newNode, nil, fromDB, false)
				if err != nil {
					return err
				}
				putToData(data.Root[:], valueDB)
			}
		}
	}
	return nil
}

func resolveHash(n trie.HashNode, db *leveldb.DB) (trie.Node, []byte, error) {
	enc, err := db.Get(n, nil)
	if err != nil || enc == nil {
		return nil, nil, &trie.MissingNodeError{}
	}
	return trie.MustDecodeNode(n, enc, 0), enc, nil
}

func hexToKeybytes(hex []byte) []byte {
	if hasTerm(hex) {
		hex = hex[:len(hex)-1]
	}
	if len(hex)&1 != 0 {
		panic("can't convert hex key of odd length")
	}
	key := make([]byte, (len(hex)+1)/2)
	decodeNibbles(hex, key)
	return key
}

// hasTerm returns whether a hex key has the terminator flag.
func hasTerm(s []byte) bool {
	return len(s) > 0 && s[len(s)-1] == 16
}

func decodeNibbles(nibbles []byte, bytes []byte) {
	for bi, ni := 0, 0; ni < len(nibbles); bi, ni = bi+1, ni+2 {
		bytes[bi] = nibbles[ni]<<4 | nibbles[ni+1]
	}
}
