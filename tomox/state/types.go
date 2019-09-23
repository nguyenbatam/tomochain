package state

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var EmptyHash = common.Hash{}
var Zero = big.NewInt(0)
var EmptyOrderList = orderList{
	Volume: nil,
	Root:   EmptyHash,
}
var EmptyExchangeOnject = exchangeObject{
	Nonce:   0,
	AskRoot: EmptyHash,
	BidRoot: EmptyHash,
}

// exchangeObject is the Ethereum consensus representation of exchanges.
// These objects are stored in the main orderId trie.
type orderList struct {
	Volume *big.Int
	Root   common.Hash // merkle root of the storage trie
}

// exchangeObject is the Ethereum consensus representation of exchanges.
// These objects are stored in the main orderId trie.
type orderItem struct {
	OrderType string
	Amount    *big.Int
}

// exchangeObject is the Ethereum consensus representation of exchanges.
// These objects are stored in the main orderId trie.
type exchangeObject struct {
	Nonce   uint64
	AskRoot common.Hash // merkle root of the storage trie
	BidRoot common.Hash // merkle root of the storage trie
}
