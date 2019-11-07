package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"runtime"
	"time"
)

var (
	from = flag.String("from", "/data/tomo/chaindata", "directory to TomoChain chaindata")
	to   = flag.String("to", "/data/tomo/chaindata_copy", "directory to clean chaindata")

	sercureKey = []byte("secure-key-") // preimagePrefix + hash -> preimage
	nWorker    = runtime.NumCPU() / 2
	finish     = int32(0)
	running    = true
	stateRoots = make(chan TrieRoot)
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
	backupRoot := common.Hash{}
	backupNumber := uint64(0)
	for number >= 0 {
		hash := core.GetCanonicalHash(fromDB, number)
		root := core.GetHeader(fromDB, hash, number).Root
		_, err := trie.NewSecure(root, tridb, 0)
		if err != nil {
			continue
		}
		if common.EmptyHash(lastestRoot) {
			lastestRoot = root
		} else {
			backupRoot = root
			backupNumber = header.Number.Uint64()
			break
		}
	}
	copyHeadData(*from, *to)
	copyBlockData(*from, *to, backupNumber)
	copyStateData(*from, *to, lastestRoot)
	copyStateData(*from, *to, backupRoot)
}
func copyHeadData(from string, to string) {
	fmt.Println(time.Now(), "copyHeadData")
	fromDB, _ := ethdb.NewLDBDatabase(from, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	toDB, _ := ethdb.NewLDBDatabase(to, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
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
	fromDB.Close()
	fmt.Println(time.Now(), "compact")
	toDB.LDB().CompactRange(util.Range{})
	toDB.Close()
	fmt.Println(time.Now(), "end")
}
func copyBlockData(from string, to string, backupNumber uint64) {
	fmt.Println(time.Now(), "copyBlockData")
	fromDB, _ := ethdb.NewLDBDatabase(from, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	toDB, _ := ethdb.NewLDBDatabase(to, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
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
	fromDB.Close()
	fmt.Println(time.Now(), "compact")
	toDB.LDB().CompactRange(util.Range{})
	toDB.Close()
	fmt.Println(time.Now(), "end")
}

func copyStateData(from string, to string, root common.Hash) {
	fmt.Println(time.Now(), "copyBlockData")
	fromDB, _ := ethdb.NewLDBDatabase(from, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	toDB, _ := ethdb.NewLDBDatabase(to, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	head := core.GetHeadBlockHash(fromDB)
	header := core.GetHeader(fromDB, head, core.GetBlockNumber(fromDB, head))
	number := header.Number.Uint64()
	for i := length; i >= 0; i-- {
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
	fromDB.Close()
	fmt.Println(time.Now(), "compact")
	toDB.LDB().CompactRange(util.Range{})
	toDB.Close()
	fmt.Println(time.Now(), "end")
}

func copyTrieData(from *trie.SecureTrie, to *trie.SecureTrie) {
	nodeIterator := from.NodeIterator(nil)
	for nodeIterator.Next(true) {
		nodeIterator.
	}
}

func copyNodeData(from *ethdb.LDBDatabase, to *ethdb.LDBDatabase, node common.Hash) {
	if to.Has(node)
		nodeIterator := from.NodeIterator(nil)
	for nodeIterator.Next(true) {
		nodeIterator.
	}
}

func processNodeAddress(n trie.Node, path []byte, fromDB *leveldb.DB, toDB *leveldb.DB) error {
	keyDB := []byte{}
	if _, ok := n.(trie.ValueNode); ok {
		buf := append(sercureKey, path...)
		keyDB = buf
	} else {
		hash, _ := n.Cache()
		keyDB = hash
	}
	exist, err := toDB.Has(keyDB, nil)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	value, err := fromDB.Get(keyDB, nil)
	if err != nil {
		return err
	}
	err = toDB.Put(keyDB, value, nil)
	if err != nil {
		return err
	}
	switch node := n.(type) {
	case *trie.FullNode:
		// Full Node, move to the first non-nil child.
		for i := 0; i < len(node.Children); i++ {
			child := node.Children[i]
			if child != nil {
				childNode := child
				var err error = nil
				if _, ok := child.(trie.HashNode); ok {
					childNode, err = resolveHash(child.(trie.HashNode), fromDB)
				}
				if err == nil {
					err = processNode(childNode, append(path, byte(i)), fromDB, toDB)
					if err != nil {
						return err
					}
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
		if _, ok := node.Val.(trie.HashNode); ok {
			childNode, err = resolveHash(node.Val.(trie.HashNode), fromDB)
		}
		if err == nil {
			err = processNode(childNode, append(path, node.Key...), fromDB, toDB)
			if err != nil {
				return err
			}
		} else if err != nil {
			_, ok := err.(*trie.MissingNodeError)
			if !ok {
				return err
			}
		}
	case *trie.ValueNode :

	}
	return nil
}

func processNode(n trie.Node, path []byte, fromDB *leveldb.DB, toDB *leveldb.DB) error {
	keyDB := []byte{}
	if _, ok := n.(trie.ValueNode); ok {
		buf := append(sercureKey, path...)
		keyDB = buf
	} else {
		hash, _ := n.Cache()
		keyDB = hash
	}
	exist, err := toDB.Has(keyDB, nil)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	value, err := fromDB.Get(keyDB, nil)
	if err != nil {
		return err
	}
	err = toDB.Put(keyDB, value, nil)
	if err != nil {
		return err
	}
	switch node := n.(type) {
	case *trie.FullNode:
		// Full Node, move to the first non-nil child.
		for i := 0; i < len(node.Children); i++ {
			child := node.Children[i]
			if child != nil {
				childNode := child
				var err error = nil
				if _, ok := child.(trie.HashNode); ok {
					childNode, err = resolveHash(child.(trie.HashNode), fromDB)
				}
				if err == nil {
					err = processNode(childNode, append(path, byte(i)), fromDB, toDB)
					if err != nil {
						return err
					}
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
		if _, ok := node.Val.(trie.HashNode); ok {
			childNode, err = resolveHash(node.Val.(trie.HashNode), fromDB)
		}
		if err == nil {
			err = processNode(childNode, append(path, node.Key...), fromDB, toDB)
			if err != nil {
				return err
			}
		} else if err != nil {
			_, ok := err.(*trie.MissingNodeError)
			if !ok {
				return err
			}
		}
	}
	return nil
}

func resolveHash(n trie.HashNode, db *leveldb.DB) (trie.Node, error) {
	enc, err := db.Get(n, nil)
	if err != nil || enc == nil {
		return nil, &trie.MissingNodeError{}
	}
	return trie.MustDecodeNode(n, enc, 0), nil
}

//func removeNodesNil(list [][17]*StateNode, length int) []*StateNode {
//	results := make([]*StateNode, length)
//	index := 0
//	for _, nodes := range list {
//		for _, node := range nodes {
//			if node != nil {
//				results[index] = node
//				index++
//			}
//		}
//	}
//	return results
//}
//func catchEventInterupt(db *leveldb.DB) {
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, os.Interrupt)
//	go func() {
//		for sig := range c {
//			fmt.Println("catch event interrupt ", sig, running, finish)
//			running = false
//			if atomic.LoadInt32(&finish) == 0 {
//				close(stateRoots)
//				db.Close()
//				os.Exit(1)
//			}
//		}
//	}()
//}
//func resolveHash(n trie.HashNode, db *leveldb.DB) (trie.Node, error) {
//	if cache.Contains(common.BytesToHash(n)) {
//		return nil, &trie.MissingNodeError{}
//	}
//	enc, err := db.Get(n, nil)
//	if err != nil || enc == nil {
//		return nil, &trie.MissingNodeError{}
//	}
//	return trie.MustDecodeNode(n, enc, 0), nil
//}
//
//func getAllChilds(n StateNode, db *leveldb.DB) ([17]*StateNode, error) {
//	childs := [17]*StateNode{}
//	switch node := n.node.(type) {
//	case *trie.FullNode:
//		// Full Node, move to the first non-nil child.
//		for i := 0; i < len(node.Children); i++ {
//			child := node.Children[i]
//			if child != nil {
//				childNode := child
//				var err error = nil
//				if _, ok := child.(trie.HashNode); ok {
//					childNode, err = resolveHash(child.(trie.HashNode), db)
//				}
//				if err == nil {
//					childs[i] = &StateNode{node: childNode, path: append(n.path, byte(i))}
//				} else if err != nil {
//					_, ok := err.(*trie.MissingNodeError)
//					if !ok {
//						return childs, err
//					}
//				}
//			}
//		}
//	case *trie.ShortNode:
//		// Short Node, return the pointer singleton child
//		childNode := node.Val
//		var err error = nil
//		if _, ok := node.Val.(trie.HashNode); ok {
//			childNode, err = resolveHash(node.Val.(trie.HashNode), db)
//		}
//		if err == nil {
//			childs[0] = &StateNode{node: childNode, path: append(n.path, node.Key...)}
//		} else if err != nil {
//			_, ok := err.(*trie.MissingNodeError)
//			if !ok {
//				return childs, err
//			}
//		}
//	}
//	return childs, nil
//}
//func processNodes(node StateNode, db *leveldb.DB) ([17]*StateNode, [17]*[]byte, int) {
//	hash, _ := node.node.Cache()
//	commonHash := common.BytesToHash(hash)
//	newNodes := [17]*StateNode{}
//	keys := [17]*[]byte{}
//	number := 0
//	if !cache.Contains(commonHash) {
//		childNodes, err := getAllChilds(node, db)
//		if err != nil {
//			fmt.Println("Error when get all childs node : ", common.Bytes2Hex(node.path), err)
//			os.Exit(1)
//		}
//		for i, child := range childNodes {
//			if child != nil {
//				if _, ok := child.node.(trie.ValueNode); ok {
//					buf := append(sercureKey, child.path...)
//					keys[i] = &buf
//				} else {
//					hash, _ := child.node.Cache()
//					var bytes []byte = hash
//					keys[i] = &bytes
//					newNodes[i] = child
//					number++
//				}
//			}
//		}
//		cache.Add(commonHash, true)
//	}
//	return newNodes, keys, number
//}
//
//func findNewNodes(nodes []*StateNode, db *leveldb.DB, batchlvdb *leveldb.Batch) ([][17]*StateNode, int) {
//	length := len(nodes)
//	chunkSize := length / nWorker
//	if len(nodes)%nWorker != 0 {
//		chunkSize++
//	}
//	childNodes := make([][17]*StateNode, length)
//	results := make(chan ResultProcessNode)
//	wg := sync.WaitGroup{}
//	wg.Add(length)
//	for i := 0; i < nWorker; i++ {
//		from := i * chunkSize
//		to := from + chunkSize
//		if to > length {
//			to = length
//		}
//		go func(from int, to int) {
//			for j := from; j < to; j++ {
//				childs, keys, number := processNodes(*nodes[j], db)
//				go func(result ResultProcessNode) {
//					results <- result
//				}(ResultProcessNode{j, number, childs, keys})
//			}
//		}(from, to)
//	}
//	total := 0
//	go func() {
//		for result := range results {
//			childNodes[result.index] = result.newNodes
//			total = total + result.number
//			for _, key := range result.keys {
//				if key != nil {
//					batchlvdb.Delete(*key)
//				}
//			}
//			wg.Done()
//		}
//	}()
//	wg.Wait()
//	close(results)
//	return childNodes, total
//}
