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
	number := *from
	before := uint64(0)
	go func() {
		for addr := range addrChan {
			if !cache.Contains(addr) {
				cache.Add(addr, true)
				f.WriteString(addr + "\n")
			}
		}
	}()
	for number <= header.Number.Uint64() {
		if number > before+1000 {
			fmt.Println(time.Now(), number)
			before = number
		}
		txs := types.Transactions{}
		for i := number; i <= number+20; i++ {
			hash := core.GetCanonicalHash(db, i)
			if common.EmptyHash(hash) {
				break
			}
			body := core.GetBody(db, hash, i)
			if len(body.Transactions) > 0 {
				txs = append(txs, body.Transactions...)
			}
		}
		number = number + 21
		if len(txs) == 0 {
			continue
		}
		length := len(txs)
		froms := make([]common.Address, length)
		wg := sync.WaitGroup{}
		wg.Add(length)
		for i := 0; i < length; i++ {
			go func(index int, tx *types.Transaction) {
				from, _ := signer.Sender(tx)
				froms[index] = from
				wg.Done()
			}(i, txs[i])
		}
		wg.Wait()
		for i, tx := range txs {
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
						fmt.Println(addrChan)
						addrChan <- tx.To().Hex()
					}()
				}
			}
			go func() {
				addrChan <- from.Hex()
			}()
		}
	}
	time.Sleep(10 * time.Second)
	fmt.Println("close addrChan")
	close(addrChan)
	f.Close()
	db.Close()
}
