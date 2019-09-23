package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Signature struct
type Signature struct {
	V byte
	R common.Hash
	S common.Hash
}

// OrderItem : info that will be store in database
type OrderItem struct {
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
