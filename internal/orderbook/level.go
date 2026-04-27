package orderbook

import "github.com/djwhocodes/trading-exchange/internal/types"

type node struct {
	order *types.Order
	prev  *node
	next  *node
}

type Level struct {
	Price       int64
	head        *node
	tail        *node
	TotalVolume int64
	length      int
}

func (l *Level) Add(o *types.Order) {
	n := &node{
		order: o,
		prev:  nil,
		next:  nil,
	}

	if l.tail == nil {
		l.head = n
		l.tail = n
	} else {
		l.tail.next = n
		n.prev = l.tail
		l.tail = n
	}

	l.TotalVolume += o.Size
	l.length++
}

func (l *Level) Remove(n *node) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		l.head = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		l.tail = n.prev
	}

	l.TotalVolume -= n.order.Size
	l.length--
}

func (l *Level) Head() *types.Order {
	if l.head == nil {
		return nil
	}
	return l.head.order
}

func (l *Level) Pop() *types.Order {
	if l.head == nil {
		return nil
	}

	n := l.head
	l.Remove(n)
	return n.order
}

func (l *Level) IsEmpty() bool {
	return l.length == 0
}

func (l *Level) Len() int {
	return l.length
}
