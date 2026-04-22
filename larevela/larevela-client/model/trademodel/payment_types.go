package trademodel

import "time"

type CreateTrialPaymentTxInput struct {
	PaymentNo   string
	TxID        string
	FromAccount string
	ToAccount   string
	AmountSol   string
	ChainType   string
	Network     string
	ChainID     int64
	AssetSymbol string
	AssetAddress string
	PlanType    string
}

type CreatePaymentIntentInput struct {
	PaymentNo         string
	OrderNo           string
	ChainType         string
	Network           string
	ChainID           int64
	PayerAccount      string
	ReceiverAccount   string
	AssetSymbol       string
	AssetAddress      string
	AmountExpected    string
	Decimals          int64
	ReferenceID       string
	SerializedMessage string
	QuoteExpiredAt    time.Time
	Status            string
}

type PaymentIntent struct {
	PaymentNo          string
	OrderNo            string
	TxID               string
	ChainType          string
	Network            string
	ChainID            int64
	PayerAccount       string
	ReceiverAccount    string
	AssetAddress       string
	AssetSymbol        string
	AmountExpected     string
	AmountActual       string
	Confirmations      int64
	ConfirmationStatus string
	Status             string
	FailureReason      string
}
