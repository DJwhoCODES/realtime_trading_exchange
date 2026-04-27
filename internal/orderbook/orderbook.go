package orderbook

import (
	"sort"

	"github.com/djwhocodes/trading-exchange/internal/types"
)

type Orderbook struct {
	bids map[int64]*Level
	asks map[int64]*Level

	bidPrices []int64 // DESC
	askPrices []int64 // ASC

	// orderID → node (for O(1) cancel)
	orderIndex map[int64]*node
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		bids:       make(map[int64]*Level),
		asks:       make(map[int64]*Level),
		orderIndex: make(map[int64]*node),
	}
}

func (ob *Orderbook) AddOrder(o *types.Order) {
	var (
		level *Level
		ok    bool
	)

	if o.Side == types.Bid {
		level, ok = ob.bids[o.Price]
		if !ok {
			level = &Level{Price: o.Price}
			ob.bids[o.Price] = level
			ob.insertBidPrice(o.Price)
		}
	} else {
		level, ok = ob.asks[o.Price]
		if !ok {
			level = &Level{Price: o.Price}
			ob.asks[o.Price] = level
			ob.insertAskPrice(o.Price)
		}
	}

	n := level.Add(o)

	ob.orderIndex[o.ID] = n
}

func (ob *Orderbook) Cancel(orderID int64) {
	n, ok := ob.orderIndex[orderID]
	if !ok {
		return
	}

	o := n.order
	var level *Level

	if o.Side == types.Bid {
		level = ob.bids[o.Price]
	} else {
		level = ob.asks[o.Price]
	}

	level.Remove(n)
	delete(ob.orderIndex, orderID)

	if level.IsEmpty() {
		ob.removePriceLevel(o.Side, o.Price)
	}
}

func (ob *Orderbook) BestBid() *Level {
	if len(ob.bidPrices) == 0 {
		return nil
	}
	return ob.bids[ob.bidPrices[0]]
}

func (ob *Orderbook) BestAsk() *Level {
	if len(ob.askPrices) == 0 {
		return nil
	}
	return ob.asks[ob.askPrices[0]]
}

func (ob *Orderbook) insertBidPrice(price int64) {
	i := sort.Search(len(ob.bidPrices), func(i int) bool {
		return ob.bidPrices[i] <= price
	})

	ob.bidPrices = append(ob.bidPrices, 0)
	copy(ob.bidPrices[i+1:], ob.bidPrices[i:])
	ob.bidPrices[i] = price
}

func (ob *Orderbook) insertAskPrice(price int64) {
	i := sort.Search(len(ob.askPrices), func(i int) bool {
		return ob.askPrices[i] >= price
	})

	ob.askPrices = append(ob.askPrices, 0)
	copy(ob.askPrices[i+1:], ob.askPrices[i:])
	ob.askPrices[i] = price
}

func (ob *Orderbook) removePriceLevel(side types.Side, price int64) {
	if side == types.Bid {
		delete(ob.bids, price)
		for i, p := range ob.bidPrices {
			if p == price {
				ob.bidPrices = append(ob.bidPrices[:i], ob.bidPrices[i+1:]...)
				break
			}
		}
	} else {
		delete(ob.asks, price)
		for i, p := range ob.askPrices {
			if p == price {
				ob.askPrices = append(ob.askPrices[:i], ob.askPrices[i+1:]...)
				break
			}
		}
	}
}
