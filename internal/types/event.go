package types

type EventType int

const (
	EventPlaceOrder EventType = iota
	EventCancelOrder
)

type Event struct {
	Type     EventType
	Order    *Order
	OrderID  int64
	Response chan any
}
