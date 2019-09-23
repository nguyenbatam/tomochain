package types

import (
	"container/heap"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/tomox"
	"io"
	"math/big"
	"sync/atomic"
)

// Signature struct
type Signature struct {
	V byte
	R common.Hash
	S common.Hash
}

type OrderItem struct {
	data orderdata
	// caches
	hash     atomic.Value
	size     atomic.Value
	from     atomic.Value
	Accepted int64
}

type orderdata struct {
	Quantity        *big.Int       `json:"quantity,omitempty"`
	Price           *big.Int       `json:"price,omitempty"`
	ExchangeAddress common.Address `json:"exchangeAddress,omitempty"`
	UserAddress     common.Address `json:"userAddress,omitempty"`
	BaseToken       common.Address `json:"baseToken,omitempty"`
	QuoteToken      common.Address `json:"quoteToken,omitempty"`
	Status          string         `json:"status,omitempty"`
	Side            string         `json:"side,omitempty"`
	Type            string         `json:"type,omitempty"`
	Hash            common.Hash    `json:"hash,omitempty"`
	Signature       *Signature     `json:"signature,omitempty"`
	FilledAmount    *big.Int       `json:"filledAmount,omitempty"`
	Nonce           *big.Int       `json:"nonce,omitempty"`
	MakeFee         *big.Int       `json:"makeFee,omitempty"`
	TakeFee         *big.Int       `json:"takeFee,omitempty"`
	PairName        string         `json:"pairName,omitempty"`
	CreatedAt       uint64         `json:"createdAt,omitempty"`
	UpdatedAt       uint64         `json:"updatedAt,omitempty"`
	OrderID         uint64         `json:"orderID,omitempty"`
	// *OrderMeta
	NextOrder []byte `json:"-"`
	PrevOrder []byte `json:"-"`
	OrderList []byte `json:"-"`
	Key       string `json:"key"`
}

// EncodeRLP implements rlp.Encoder
func (order *OrderItem) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &order.data)
}

// DecodeRLP implements rlp.Decoder
func (order *OrderItem) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	err := s.Decode(&order.data)
	if err == nil {
		order.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}
func (order *OrderItem) From() common.Address {
	return order.data.UserAddress
}

func (order *OrderItem) ExchangeAddress() common.Address {
	return order.data.ExchangeAddress
}

// Hash hashes the RLP encoding of order.
// It uniquely identifies the transaction.
func (order *OrderItem) Hash() common.Hash {
	if hash := order.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(order)
	order.hash.Store(v)
	return v
}

// Hash hashes the RLP encoding of order.
// It uniquely identifies the transaction.
func (order *OrderItem) Received() int64 {
	return order.Accepted
}

// MarshalJSON encodes the web3 RPC transaction format.
func (order *OrderItem) MarshalJSON() ([]byte, error) {
	return order.data.MarshalJSON()
}

// UnmarshalJSON decodes the web3 RPC transaction format.
func (order *OrderItem) UnmarshalJSON(input []byte) error {
	var dec orderdata
	if err := dec.UnmarshalJSON(input); err != nil {
		return err
	}
	*order = OrderItem{data: dec}
	return nil
}
func (order *OrderItem) Price() *big.Int    { return new(big.Int).Set(order.data.Price) }
func (order *OrderItem) Quantity() *big.Int { return new(big.Int).Set(order.data.Quantity) }
func (order *OrderItem) Nonce() uint64      { return order.data.Nonce.Uint64() }

func (order *OrderItem) CacheHash() {
	v := rlpHash(order)
	order.hash.Store(v)
}

// following: https://github.com/tomochain/tomox-sdk/blob/master/types/order.go#L125
func (o *OrderItem) computeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(o.data.ExchangeAddress.Bytes())
	sha.Write(o.data.UserAddress.Bytes())
	sha.Write(o.data.BaseToken.Bytes())
	sha.Write(o.data.QuoteToken.Bytes())
	sha.Write(common.BigToHash(o.data.Quantity).Bytes())
	sha.Write(common.BigToHash(o.data.Price).Bytes())
	sha.Write(common.BigToHash(o.encodedSide()).Bytes())
	sha.Write(common.BigToHash(o.data.Nonce).Bytes())
	sha.Write(common.BigToHash(o.data.MakeFee).Bytes())
	sha.Write(common.BigToHash(o.data.TakeFee).Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

//verify signatures
func (o *OrderItem) VerifySignature() error {
	var (
		hash           common.Hash
		err            error
		signatureBytes []byte
	)
	hash = o.computeHash()
	if hash != o.Hash() {
		return errWrongHash
	}
	signatureBytes = append(signatureBytes, o.data.Signature.R.Bytes()...)
	signatureBytes = append(signatureBytes, o.data.Signature.S.Bytes()...)
	signatureBytes = append(signatureBytes, o.data.Signature.V-27)
	pubkey, err := crypto.Ecrecover(hash.Bytes(), signatureBytes)
	if err != nil {
		return err
	}
	var userAddress common.Address
	copy(userAddress[:], crypto.Keccak256(pubkey[1:])[12:])
	if userAddress != o.data.UserAddress {
		return errInvalidSignature
	}
	return nil
}

func (o *OrderItem) encodedSide() *big.Int {
	if o.data.Side == tomox.Bid {
		return big.NewInt(0)
	}
	return big.NewInt(1)
}

func (order *OrderItem) String() string {
	var from, to string
	if order.data.V != nil {
		// make a best guess about the signer and use that to derive
		// the sender.
		signer := deriveSigner(order.data.V)
		if f, err := Sender(signer, order); err != nil { // derive but don't cache
			from = "[invalid sender: invalid sig]"
		} else {
			from = fmt.Sprintf("%x", f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}

	if order.data.Recipient == nil {
		to = "[contract creation]"
	} else {
		to = fmt.Sprintf("%x", order.data.Recipient[:])
	}
	enc, _ := rlp.EncodeToBytes(&order.data)
	return fmt.Sprintf(`
	order(%x)
	Contract: %v
	From:     %s
	To:       %s
	Nonce:    %v
	GasPrice: %#x
	GasLimit  %#x
	Quantity:    %#x
	Data:     0x%x
	V:        %#x
	R:        %#x
	S:        %#x
	Hex:      %x
`,
		order.Hash(),
		order.data.Recipient == nil,
		from,
		to,
		order.data.AccountNonce,
		order.data.Price,
		order.data.GasLimit,
		order.data.Amount,
		order.data.Payload,
		order.data.V,
		order.data.R,
		order.data.S,
		enc,
	)
}

// Transactions is a Transaction slice type for basic sorting.
type OrderItems []*OrderItem

// Len returns the length of s.
func (s OrderItems) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s OrderItems) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (s OrderItems) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

func OrderDifference(a, b OrderItems) (keep OrderItems) {
	keep = make(OrderItems, 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, order := range b {
		remove[order.Hash()] = struct{}{}
	}

	for _, order := range a {
		if _, ok := remove[order.Hash()]; !ok {
			keep = append(keep, order)
		}
	}

	return keep
}

type OrderByNonce OrderItems

func (s OrderByNonce) Len() int           { return len(s) }
func (s OrderByNonce) Less(i, j int) bool { return s[i].data.Nonce.Cmp(s[j].data.Nonce) < 0 }
func (s OrderByNonce) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// OrderByTime implements both the sort and the heap interface, making it useful
// for all at once sorting as well as individually adding and removing elements.
type OrderByTime struct {
	orders OrderItems
}

func (s OrderByTime) Len() int { return len(s.orders) }
func (s OrderByTime) Less(i, j int) bool {
	return s.orders[i].Accepted < s.orders[j].Accepted
}
func (s OrderByTime) Swap(i, j int) { s.orders[i], s.orders[j] = s.orders[j], s.orders[i] }

func (s *OrderByTime) Push(x interface{}) {
	s.orders = append(s.orders, x.(*OrderItem))
}

func (s *OrderByTime) Pop() interface{} {
	old := s.orders
	n := len(old)
	x := old[n-1]
	s.orders = old[0 : n-1]
	return x
}

// OrderItemsByTimeAndNonce represents a set of transactions that can return
// transactions in a profit-maximizing sorted order, while supporting removing
// entire batches of transactions for non-executable accounts.
type OrderItemsByTimeAndNonce struct {
	orders map[common.Address]OrderItems // Per account nonce-sorted list of transactions
	heads  OrderByTime                   // Next transaction for each unique account (price heap)
	signer Signer                        // Signer for the set of transactions
}

// NewTransactionsByPriceAndNonce creates a transaction set that can retrieve
// price sorted transactions in a nonce-honouring way.
//
// Note, the input map is reowned so the caller should not interact any more with
// if after providing it to the constructor.

// It also classifies special orders and normal orders
func NewOrderItemsByTimeAndNonce(signer Signer, orders map[common.Address]OrderItems, signers map[common.Address]struct{}, payersSwap map[common.Address]*big.Int) *OrderItemsByTimeAndNonce {
	// Initialize a price based heap with the head transactions
	heads := OrderByTime{}
	for _, accorders := range orders {
		from := accorders[0].From()
		if len(accorders) > 0 {
			heads.orders = append(heads.orders, accorders[0])
			// Ensure the sender address is from the signer
			orders[*from] = accorders[1:]
		}
	}
	heap.Init(&heads)

	// Assemble and return the transaction set
	return &OrderItemsByTimeAndNonce{
		orders: orders,
		heads:  heads,
		signer: signer,
	}
}

// Peek returns the next transaction by price.
func (t *OrderItemsByTimeAndNonce) Peek() *OrderItem {
	if len(t.heads.orders) == 0 {
		return nil
	}
	return t.heads.orders[0]
}

// Shift replaces the current best head with the next one from the same account.
func (t *OrderItemsByTimeAndNonce) Shift() {
	acc := t.heads.orders[0].From()
	if orders, ok := t.orders[*acc]; ok && len(orders) > 0 {
		t.heads.orders[0], t.orders[*acc] = orders[0], orders[1:]
		heap.Fix(&t.heads, 0)
	} else {
		heap.Pop(&t.heads)
	}
}

// Pop removes the best transaction, *not* replacing it with the next one from
// the same account. This should be used when a transaction cannot be executed
// and hence all subsequent ones should be discarded from the same account.
func (t *OrderItemsByTimeAndNonce) Pop() {
	heap.Pop(&t.heads)
}
