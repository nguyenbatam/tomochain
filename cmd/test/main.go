package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
	"os"
	"sync/atomic"
	"time"
)

// Block represents an entire block in the Ethereum blockchain.
type OldBlock struct {
	header       *types.Header
	uncles       []*types.Header
	transactions types.Transactions
	// caches
	hash atomic.Value
	size atomic.Value

	// Td is used by package core to store the total difficulty
	// of the chain up to and including the block.
	td *big.Int

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

// "external" block encoding. used for eth protocol, etc.
type oldextblock struct {
	Header *types.Header
	Txs    []*types.Transaction
	Uncles []*types.Header
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (b *OldBlock) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, oldextblock{
		Header: b.header,
		Txs:    b.transactions,
		Uncles: b.uncles,
	})
}

// DecodeRLP decodes the Ethereum
func (b *OldBlock) DecodeRLP(s *rlp.Stream) error {
	var eb oldextblock
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.header, b.uncles, b.transactions = eb.Header, eb.Uncles, eb.Txs
	return nil
}

// newBlockData is the network packet for the block propagation message.
type oldBlockData struct {
	Block *OldBlock
	TD    *big.Int
}

// newBlockData is the network packet for the block propagation message.
type newBlockData struct {
	Block *types.Block
	TD    *big.Int
}

func main() {
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.LvlTrace)
	log.Root().SetHandler(glogger)
	lddb, _ := ethdb.NewLDBDatabase("/data/tomo/chaindata", eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	head := core.GetHeadBlockHash(lddb)
	block := core.GetBlock(lddb, head, core.GetBlockNumber(lddb, head))
	newBlockData := newBlockData{Block: block, TD: big.NewInt(100000)}
	data, err := rlp.EncodeToBytes(newBlockData)
	fmt.Println("encode", len(data), err)
	var oldBlockData oldBlockData
	err = rlp.DecodeBytes(data, &oldBlockData)
	fmt.Println("decode", len(data), err)
	fmt.Println(oldBlockData.Block.header.Hash())
}
