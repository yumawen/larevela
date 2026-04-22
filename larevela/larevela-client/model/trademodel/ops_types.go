package trademodel

type UpsertChainScanCursorInput struct {
	ChainType        string
	Network          string
	ChainID          int64
	CursorType       string
	LastScannedBlock int64
	LastScannedSlot  int64
	LastScannedTxID  string
	Remark           string
}

type UpsertIdempotencyRecordInput struct {
	IdemKey          string
	BizType          string
	BizNo            string
	Status           string
	ResponseSnapshot string
}
