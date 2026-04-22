package trademodel

import "time"

type CreateOrderInput struct {
	OrderNo   string
	BizType   string
	BizID     string
	UserID    int64
	Currency  string
	Amount    string
	Status    string
	ExpiredAt time.Time
}

type Order struct {
	OrderNo   string
	BizType   string
	BizID     string
	UserID    int64
	Currency  string
	Amount    string
	Status    string
	ExpiredAt int64
	PaidAt    int64
}
