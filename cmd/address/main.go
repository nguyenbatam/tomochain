package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	lru "github.com/hashicorp/golang-lru"
	"math/big"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	dir          = flag.String("dir", "/data/tomo/chaindata", "directory to TomoChain chaindata")
	address      = flag.String("address", "/data/tomo/address.txt", "output list address in block")
	from         = flag.Uint64("from", 0, "from block number")
	cache, _     = lru.NewARC(1000000)
	batch        *bytes.Buffer
	lengthBuffer = 1000000
	addrChan     chan string
	nWorker      = runtime.NumCPU() / 2
)

func main() {
	flag.Parse()
	db, err := ethdb.NewLDBDatabase(*dir, eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.OpenFile(*address, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
	}
	addrChan = make(chan string)
	batch = bytes.NewBuffer(make([]byte, 0, lengthBuffer))
	signer := types.NewEIP155Signer(big.NewInt(88))
	head := core.GetHeadBlockHash(db)
	header := core.GetHeader(db, head, core.GetBlockNumber(db, head))
	mapNonces := map[common.Address]uint64{}
	for number := *from; number <= header.Number.Uint64(); number++ {
		if number%1000 == 0 {
			fmt.Println(time.Now(), number)
		}
		hash := core.GetCanonicalHash(db, number)
		body := core.GetBody(db, hash, number)
		length := len(body.Transactions)
		froms := make([]common.Address, length)
		wg := sync.WaitGroup{}
		wg.Add(length)
		for i := 0; i < length; i++ {
			go func(index int, tx *types.Transaction) {
				from, _ := signer.Sender(tx)
				froms[index] = from
				wg.Done()
			}(i, body.Transactions[i])
		}
		wg.Wait()
		for i, tx := range body.Transactions {
			from := froms[i]
			oldNonce := mapNonces[from]
			mapNonces[from] = oldNonce + 1
			if tx.To() == nil {
				smc := crypto.CreateAddress(from, tx.Nonce())
				go func() {
					addrChan <- smc.Hex()
				}()
			} else {
				if tx.To().Hex() != common.BlockSigners {
					go func() {
						addrChan <- tx.To().Hex()
					}()
				}
			}
			go func() {
				addrChan <- from.Hex()
			}()
		}
	}
	go func() {
		for addr := range addrChan {
			if !cache.Contains(addr) {
				cache.Add(addr, true)
				f.WriteString(addr + "\n")
			}
		}
	}()
	f.Close()
	db.Close()
}
