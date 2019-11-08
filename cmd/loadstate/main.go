package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
)

func main()  {
	db, err := ethdb.NewLDBDatabase("/data/tomo/copy", eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	fmt.Println("NewLDBDatabase",err)
	fmt.Println(loadLastBlock(db))
}


func loadLastBlock(db *ethdb.LDBDatabase) (uint64, uint64, uint64) {
	// Restore the last known head block
	head := core.GetHeadBlockHash(db)
	if head == (common.Hash{}) {
		// Corrupt or empty database, init from scratch
		fmt.Println("Empty database, resetting chain")
		return 0, 0, 0
	}
	// Make sure the entire head block is available
	currentBlock := GetBlockByHash(db, head)
	if currentBlock == nil {
		// Corrupt or empty database, init from scratch
		fmt.Println("Head block missing, resetting chain", "hash", head)
		return 0, 0, 0
	}

	// Restore the last known head header
	currentHeader := currentBlock.Header()
	if head := core.GetHeadHeaderHash(db); head != (common.Hash{}) {
		if header := GetHeaderByHash(db, head); header != nil {
			currentHeader = header
		}
	}
	var currentFastBlock *types.Block
	if head := core.GetHeadFastBlockHash(db); head != (common.Hash{}) {
		if block := GetBlockByHash(db, head); block != nil {
			currentFastBlock = block
		}
	}
	headerTd := core.GetTd(db, currentHeader.Hash(), currentHeader.Number.Uint64())
	blockTd := core.GetTd(db, currentBlock.Hash(), currentBlock.NumberU64())
	fastTd := core.GetTd(db, currentFastBlock.Hash(), currentFastBlock.NumberU64())

	fmt.Println("Loaded most recent local header", "number", currentHeader.Number, "hash", currentHeader.Hash(), "td", headerTd)
	fmt.Println("Loaded most recent local full block", "number", currentBlock.Number(), "hash", currentBlock.Hash(), "td", blockTd)
	fmt.Println("Loaded most recent local fast block", "number", currentFastBlock.Number(), "hash", currentFastBlock.Hash(), "td", fastTd)

	return currentHeader.Number.Uint64(), currentBlock.Number().Uint64(), currentFastBlock.Number().Uint64()
}

// GetBlockByHash retrieves a block from the database by hash, caching it if found.
func GetBlockByHash(db *ethdb.LDBDatabase, hash common.Hash) *types.Block {
	return GetBlock(db, hash, core.GetBlockNumber(db, hash))
}

// GetBlock retrieves a block from the database by hash and number,
// caching it if found.
func GetBlock(db *ethdb.LDBDatabase, hash common.Hash, number uint64) *types.Block {
	block := core.GetBlock(db, hash, number)
	if block == nil {
		return nil
	}
	return block
}

func GetHeaderByHash(db *ethdb.LDBDatabase, hash common.Hash) *types.Header {
	return core.GetHeader(db, hash, core.GetBlockNumber(db, hash))
}
