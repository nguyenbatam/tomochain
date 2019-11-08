package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	lru "github.com/hashicorp/golang-lru"
	"os"
	"time"
)

var (
	dir      = flag.String("dir", "/data/tomo/chaindata", "directory to TomoChain chaindata")
	address  = flag.String("address", "/data/tomo/address.txt", "output list address in block")
	from     = flag.Uint64("from", 0, "from block number")
	cache, _ = lru.NewARC(1000000)
)

func main() {
	flag.Parse()
	db, err := ethdb.NewLDBDatabase(*dir, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	if err != nil {
		fmt.Println(err)
	}
	err = os.Remove(*address)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create(*address)
	if err != nil {
		fmt.Println(err)
	}
	head := core.GetHeadBlockHash(db)
	header := core.GetHeader(db, head, core.GetBlockNumber(db, head))
	mapNonces := map[common.Address]uint64{}
	for number := *from; number <= header.Number.Uint64(); number++ {
		if number%1000 == 0 {
			fmt.Println(time.Now(), number)
		}
		hash := core.GetCanonicalHash(db, number)
		body := core.GetBody(db, hash, number)
		for _, tx := range body.Transactions {
			from := *tx.From();
			oldNonce := mapNonces[from]
			mapNonces[from] = oldNonce + 1
			if tx.To() == nil {
				smc := crypto.CreateAddress(from, tx.Nonce())
				write(f, smc.Hex())
			} else {
				if tx.To().Hex() != common.BlockSigners {
					write(f, tx.To().Hex())
				}
			}
			write(f, from.Hex())
		}
	}
	f.Close()
	db.Close()
}
func write(f *os.File, addr string) {
	if cache.Contains(addr) {
		return
	}
	f.WriteString(addr + "\n")
	cache.Add(addr, true)
}
