package trademodel

type CreateLedgerEntryInput struct {
	EntryNo      string
	PaymentNo    string
	OrderNo      string
	UserID       int64
	ChainType    string
	Network      string
	ChainID      int64
	EntryType    string
	AssetSymbol  string
	AssetAddress string
	Amount       string
	Direction    string
	Remark       string
}

type LedgerEntry struct {
	EntryNo      string
	PaymentNo    string
	OrderNo      string
	UserID       int64
	ChainType    string
	Network      string
	ChainID      int64
	EntryType    string
	AssetSymbol  string
	AssetAddress string
	Amount       string
	Direction    string
	Status       string
	Remark       string
}
