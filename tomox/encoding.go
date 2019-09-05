package tomox

import (
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// Node item
func EncodeBytesNodeItem(item *Item) ([]byte, error) {
	// try with order item
	start := 3 * common.HashLength

	// red-black is 1 byte
	totalLength := start + 1
	if item.Value != nil {
		totalLength += len(item.Value)
	}

	returnBytes := make([]byte, totalLength)

	if item.Keys != nil {
		copy(returnBytes[0:common.HashLength], item.Keys.Left)
		copy(returnBytes[common.HashLength:2*common.HashLength], item.Keys.Right)
		copy(returnBytes[2*common.HashLength:start], item.Keys.Parent)
	}
	returnBytes[start] = Bool2byte(item.Color)
	start++
	if start < totalLength {
		copy(returnBytes[start:], item.Value)
	}

	return returnBytes, nil
}

func DecodeBytesNodeItem(bytes []byte, item *Item) error {
	// try with OrderItem
	start := 3 * common.HashLength
	totalLength := len(bytes)
	if item.Keys == nil {
		item.Keys = &KeyMeta{
			Left:   EmptyKey(),
			Right:  EmptyKey(),
			Parent: EmptyKey(),
		}
	}
	copy(item.Keys.Left, bytes[0:common.HashLength])
	copy(item.Keys.Right, bytes[common.HashLength:2*common.HashLength])
	copy(item.Keys.Parent, bytes[2*common.HashLength:start])
	item.Color = Byte2bool(bytes[start])
	start++
	if start < totalLength {
		item.Value = make([]byte, totalLength-start)
		copy(item.Value, bytes[start:])
	}

	return nil
}

// Order item
func EncodeBytesOrderItem(item *OrderItem) ([]byte, error) {
	start := 2 * common.HashLength                //Quantity, Price
	totalLength := start + 4*common.AddressLength //ExchangeAddress,UserAddress,BaseToken,QuoteToken
	totalLength += 15                             // maxlength item.Status = 15
	totalLength += 4                              // maxlength item.Side = 4
	totalLength += 6                              // maxlength item.Type = 6
	totalLength += common.HashLength              //hash
	totalLength += 1                              //Signature.V = 1 byte
	totalLength += 2 * common.HashLength          //Signature.R, Signature.S
	totalLength += 4 * common.HashLength          //FilledAmount, Nonce, MakeFee, TakeFee
	totalLength += len(item.PairName)             //PairName
	// uint64 and int64 are 8 byte
	totalLength += 3 * 8                 // CreatedAt, UpdatedAt, OrderID
	totalLength += 3 * common.HashLength // next, prev, orderlist

	returnBytes := make([]byte, totalLength)

	if item.Quantity != nil {
		copy(returnBytes[0:common.HashLength], common.BigToHash(item.Quantity).Bytes())
	}
	if item.Price != nil {
		copy(returnBytes[common.HashLength:2*common.HashLength], common.BigToHash(item.Price).Bytes())
	}

	copy(returnBytes[start:start+common.AddressLength], item.ExchangeAddress[:])
	start += common.AddressLength
	copy(returnBytes[start:start+common.AddressLength], item.UserAddress[:])
	start += common.AddressLength
	copy(returnBytes[start:start+common.AddressLength], item.BaseToken[:])
	start += common.AddressLength
	copy(returnBytes[start:start+common.AddressLength], item.QuoteToken[:])
	start += common.AddressLength

	copy(returnBytes[start:start+15], []byte(item.Status))
	start += 15

	copy(returnBytes[start:start+4], []byte(item.Side))
	start += 4

	copy(returnBytes[start:start+6], []byte(item.Type))
	start += 6

	copy(returnBytes[start:start+common.HashLength], item.Hash.Bytes())
	start += common.HashLength

	returnBytes = append(returnBytes, item.Signature.V)
	start += 1
	copy(returnBytes[start:start+common.HashLength], item.Signature.R.Bytes())
	start += common.HashLength
	copy(returnBytes[start:start+common.HashLength], item.Signature.S.Bytes())
	start += common.HashLength

	copy(returnBytes[start:start+common.HashLength], common.BigToHash(item.FilledAmount).Bytes())
	start += common.HashLength
	copy(returnBytes[start:start+common.HashLength], common.BigToHash(item.Nonce).Bytes())
	start += common.HashLength
	copy(returnBytes[start:start+common.HashLength], common.BigToHash(item.MakeFee).Bytes())
	start += common.HashLength
	copy(returnBytes[start:start+common.HashLength], common.BigToHash(item.TakeFee).Bytes())
	start += common.HashLength

	//fixed length item.PairName = 86 bytes
	copy(returnBytes[start:start+86], []byte(item.PairName))
	start += 86

	binary.BigEndian.PutUint64(returnBytes[start:start+8], item.CreatedAt)
	start += 8
	binary.BigEndian.PutUint64(returnBytes[start:start+8], item.UpdatedAt)
	start += 8

	binary.BigEndian.PutUint64(returnBytes[start:start+8], item.OrderID)
	start += 8

	copy(returnBytes[start:start+common.HashLength], item.NextOrder)
	start += common.HashLength

	copy(returnBytes[start:start+common.HashLength], item.PrevOrder)
	start += common.HashLength

	copy(returnBytes[start:start+common.HashLength], item.OrderList)
	start += common.HashLength

	return returnBytes, nil
}

func DecodeBytesOrderItem(bytes []byte, item *OrderItem) error {
	start := 0

	if item.Quantity == nil {
		item.Quantity = new(big.Int)
	}
	item.Quantity.SetBytes(bytes[start : start+common.HashLength])
	start += common.HashLength

	if item.Price == nil {
		item.Price = new(big.Int)
	}
	item.Price.SetBytes(bytes[start : start+common.HashLength])
	start += common.HashLength

	item.ExchangeAddress = common.BytesToAddress(bytes[start : start+common.AddressLength])
	start += common.AddressLength
	item.UserAddress = common.BytesToAddress(bytes[start : start+common.AddressLength])
	start += common.AddressLength
	item.BaseToken = common.BytesToAddress(bytes[start : start+common.AddressLength])
	start += common.AddressLength
	item.QuoteToken = common.BytesToAddress(bytes[start : start+common.AddressLength])
	start += common.AddressLength

	item.Status = string(bytes[start : start+15])
	start += 15
	item.Side = string(bytes[start : start+4])
	start += 4
	item.Type = string(bytes[start : start+6])
	start += 6

	item.Hash = common.BytesToHash(bytes[start : start+common.HashLength])
	start += common.HashLength

	if item.Signature == nil {
		item.Signature = &Signature{}
	}
	item.Signature.V = bytes[start]
	start += 1
	item.Signature.R = common.BytesToHash(bytes[start : start+common.HashLength])
	start += common.HashLength
	item.Signature.S = common.BytesToHash(bytes[start : start+common.HashLength])
	start += common.HashLength

	item.FilledAmount = new(big.Int)
	item.FilledAmount.SetBytes(bytes[start : start+common.HashLength])
	start += common.HashLength

	item.Nonce = new(big.Int)
	item.Nonce.SetBytes(bytes[start : start+common.HashLength])
	start += common.HashLength

	item.MakeFee = new(big.Int)
	item.MakeFee.SetBytes(bytes[start : start+common.HashLength])
	start += common.HashLength

	item.TakeFee = new(big.Int)
	item.TakeFee.SetBytes(bytes[start : start+common.HashLength])
	start += common.HashLength

	item.PairName = string(bytes[start : start+86])
	start += 86

	item.CreatedAt = uint64(binary.BigEndian.Uint64(bytes[start : start+8]))
	start += 8

	item.UpdatedAt = uint64(binary.BigEndian.Uint64(bytes[start : start+8]))
	start += 8

	item.OrderID = uint64(binary.BigEndian.Uint64(bytes[start : start+8]))
	start += 8

	// pointers
	if item.NextOrder == nil {
		item.NextOrder = EmptyKey()
	}
	copy(item.NextOrder, bytes[start:start+common.HashLength])
	start += common.HashLength

	if item.PrevOrder == nil {
		item.PrevOrder = EmptyKey()
	}
	copy(item.PrevOrder, bytes[start:start+common.HashLength])
	start += common.HashLength

	if item.OrderList == nil {
		item.OrderList = EmptyKey()
	}
	copy(item.OrderList, bytes[start:start+common.HashLength])
	start += common.HashLength

	return nil
}

// OrderList item
func EncodeBytesOrderListItem(item *OrderListItem) ([]byte, error) {
	// try with order item from volume and price
	start := 2 * common.HashLength
	totalLength := start + 2*common.HashLength // head, tail
	// uint64 is 8 byte
	totalLength += 8 // length

	returnBytes := make([]byte, totalLength)

	if item.Volume != nil {
		copy(returnBytes[0:common.HashLength], common.BigToHash(item.Volume).Bytes())
	}
	if item.Price != nil {
		copy(returnBytes[common.HashLength:2*common.HashLength], common.BigToHash(item.Price).Bytes())
	}

	copy(returnBytes[start:start+common.HashLength], item.HeadOrder)
	start += common.HashLength
	copy(returnBytes[start:start+common.HashLength], item.TailOrder)
	start += common.HashLength
	binary.BigEndian.PutUint64(returnBytes[start:start+8], item.Length)

	return returnBytes, nil
}

func DecodeBytesOrderListItem(bytes []byte, item *OrderListItem) error {
	// try with OrderItem
	start := 0
	// make it crash it wrong format, no need to check length
	totalLength := len(bytes)

	if item.Volume == nil {
		item.Volume = new(big.Int)
	}

	item.Volume.SetBytes(bytes[start : start+common.HashLength])
	start += common.HashLength

	if item.Price == nil {
		item.Price = new(big.Int)
	}
	item.Price.SetBytes(bytes[start : start+common.HashLength])
	start += common.HashLength

	// pointers
	if item.HeadOrder == nil {
		item.HeadOrder = EmptyKey()
	}
	copy(item.HeadOrder, bytes[start:start+common.HashLength])
	start += common.HashLength

	if item.TailOrder == nil {
		item.TailOrder = EmptyKey()
	}
	copy(item.TailOrder, bytes[start:start+common.HashLength])
	start += common.HashLength

	// may have wrong format, just get next 8 bytes
	if start+8 <= totalLength {
		item.Length = binary.BigEndian.Uint64(bytes[start : start+8])
	}

	return nil
}

// order tree item
func EncodeBytesOrderTreeItem(item *OrderTreeItem) ([]byte, error) {
	// try with order item from volume
	start := 1 * common.HashLength
	totalLength := start + 1*common.HashLength // PriceTreeKey
	// uint64 is 8 byte
	totalLength += 8 * 2 // NumOrders and PriceTreeSize

	returnBytes := make([]byte, totalLength)

	if item.Volume != nil {
		copy(returnBytes[0:common.HashLength], common.BigToHash(item.Volume).Bytes())
	}

	copy(returnBytes[start:start+common.HashLength], item.PriceTreeKey)
	start += common.HashLength

	binary.BigEndian.PutUint64(returnBytes[start:start+8], item.NumOrders)
	start += 8

	binary.BigEndian.PutUint64(returnBytes[start:start+8], item.PriceTreeSize)

	return returnBytes, nil
}

func DecodeBytesOrderTreeItem(bytes []byte, item *OrderTreeItem) error {
	// try with OrderItem
	start := 0
	// make it crash it wrong format, no need to check length
	totalLength := len(bytes)

	if item.Volume == nil {
		item.Volume = new(big.Int)
	}

	item.Volume.SetBytes(bytes[start : start+common.HashLength])
	start += common.HashLength

	// pointers
	if item.PriceTreeKey == nil {
		item.PriceTreeKey = EmptyKey()
	}
	copy(item.PriceTreeKey, bytes[start:start+common.HashLength])
	start += common.HashLength

	item.NumOrders = binary.BigEndian.Uint64(bytes[start : start+8])
	start += 8

	// may have wrong format, just get next 8 bytes
	if start+8 <= totalLength {
		item.PriceTreeSize = binary.BigEndian.Uint64(bytes[start : start+8])
	}

	return nil
}

// order book item
func EncodeBytesOrderBookItem(item *OrderBookItem) ([]byte, error) {
	// try with zero
	start := 0
	totalLength := start + 2*8 // Timestamp and NextOrderID
	totalLength += len(item.Name)

	returnBytes := make([]byte, totalLength)

	binary.BigEndian.PutUint64(returnBytes[start:start+8], item.Timestamp)
	start += 8
	binary.BigEndian.PutUint64(returnBytes[start:start+8], item.NextOrderID)
	start += 8

	if start < totalLength {
		copy(returnBytes[start:], item.Name)
	}

	return returnBytes, nil
}

func DecodeBytesOrderBookItem(bytes []byte, item *OrderBookItem) error {
	// try with OrderItem
	start := 0
	totalLength := len(bytes)

	item.Timestamp = binary.BigEndian.Uint64(bytes[start : start+8])
	start += 8

	item.NextOrderID = binary.BigEndian.Uint64(bytes[start : start+8])
	start += 8

	if start < totalLength {
		item.Name = string(bytes[start:])
	}

	// fmt.Printf("Item : %d, %d\n", start+8, totalLength)
	return nil
}

func EncodeBytesItem(val interface{}) ([]byte, error) {

	switch val.(type) {
	case *Item:
		return EncodeBytesNodeItem(val.(*Item))
	case *OrderItem:
		return EncodeBytesOrderItem(val.(*OrderItem))
	case *OrderListItem:
		return EncodeBytesOrderListItem(val.(*OrderListItem))
	case *OrderTreeItem:
		return EncodeBytesOrderTreeItem(val.(*OrderTreeItem))
	case *OrderBookItem:
		return EncodeBytesOrderBookItem(val.(*OrderBookItem))
	default:
		return rlp.EncodeToBytes(val)
	}
}

func DecodeBytesItem(bytes []byte, val interface{}) error {

	switch val.(type) {
	case *Item:
		return DecodeBytesNodeItem(bytes, val.(*Item))
	case *OrderItem:
		return DecodeBytesOrderItem(bytes, val.(*OrderItem))
	case *OrderListItem:
		return DecodeBytesOrderListItem(bytes, val.(*OrderListItem))
	case *OrderTreeItem:
		return DecodeBytesOrderTreeItem(bytes, val.(*OrderTreeItem))
	case *OrderBookItem:
		return DecodeBytesOrderBookItem(bytes, val.(*OrderBookItem))
	default:
		return rlp.DecodeBytes(bytes, val)
	}

}
