package trademodel

type UpsertPaymentViewInput struct {
	PaymentNo          string
	OrderNo            string
	TxID               string
	ChainType          string
	Network            string
	ChainID            int64
	PayerAccount       string
	ReceiverAccount    string
	AssetSymbol        string
	AmountExpected     string
	AmountActual       string
	Status             string
	ConfirmationStatus string
	Confirmations      int64
	LastScannedBlock   int64
	LastScannedSlot    int64
	FailureReason      string
	UpdatedSource      string
}

type PaymentView struct {
	PaymentNo          string
	OrderNo            string
	TxID               string
	ChainType          string
	Network            string
	ChainID            int64
	PayerAccount       string
	ReceiverAccount    string
	AssetSymbol        string
	AmountExpected     string
	AmountActual       string
	Status             string
	ConfirmationStatus string
	Confirmations      int64
	LastScannedBlock   int64
	LastScannedSlot    int64
	FailureReason      string
}
