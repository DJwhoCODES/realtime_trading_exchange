package types

type Side int

const (
	Bid Side = iota
	Ask
)

func (s Side) String() string {
	if s == Bid {
		return "BID"
	}
	return "ASK"
}

type Order struct {
	ID     int64
	UserID int64
	Price  int64
	Size   int64
	Side   Side
	Time   int64
}

func (o *Order) IsFilled() bool {
	return o.Size == 0
}

type Trade struct {
	Price int64
	Size  int64
	BidID int64
	AskID int64
	Time  int64
}
