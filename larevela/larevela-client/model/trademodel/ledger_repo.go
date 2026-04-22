package trademodel

import (
	"context"
	"fmt"
)

func (m *Model) CreateLedgerEntry(ctx context.Context, in CreateLedgerEntryInput) error {
	if err := m.ensureReady(); err != nil {
		return err
	}
	if in.EntryNo == "" || in.PaymentNo == "" || in.OrderNo == "" || in.AssetSymbol == "" || in.Amount == "" {
		return fmt.Errorf("missing required ledger fields")
	}

	query := `
INSERT INTO ledger_entries (
  entry_no, payment_no, order_no, user_id, chain_type, network, chain_id, entry_type,
  asset_symbol, asset_address, amount, direction, status, remark, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'posted', ?, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  amount = VALUES(amount),
  direction = VALUES(direction),
  remark = VALUES(remark),
  updated_at = NOW()
`
	_, err := m.conn.ExecCtx(
		ctx,
		query,
		in.EntryNo, in.PaymentNo, in.OrderNo, in.UserID, in.ChainType, in.Network, in.ChainID,
		in.EntryType, in.AssetSymbol, in.AssetAddress, in.Amount, in.Direction, in.Remark,
	)
	return err
}

func (m *Model) GetLedgerEntry(ctx context.Context, entryNo string) (*LedgerEntry, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}
	if entryNo == "" {
		return nil, fmt.Errorf("entryNo is required")
	}

	var row struct {
		EntryNo      string `db:"entry_no"`
		PaymentNo    string `db:"payment_no"`
		OrderNo      string `db:"order_no"`
		UserID       int64  `db:"user_id"`
		ChainType    string `db:"chain_type"`
		Network      string `db:"network"`
		ChainID      int64  `db:"chain_id"`
		EntryType    string `db:"entry_type"`
		AssetSymbol  string `db:"asset_symbol"`
		AssetAddress string `db:"asset_address"`
		Amount       string `db:"amount"`
		Direction    string `db:"direction"`
		Status       string `db:"status"`
		Remark       string `db:"remark"`
	}
	query := `
SELECT entry_no, payment_no, order_no, user_id, chain_type, network, chain_id, entry_type,
       asset_symbol, asset_address, amount, direction, status, IFNULL(remark, '')
FROM ledger_entries
WHERE entry_no = ?
LIMIT 1
`
	if err := m.conn.QueryRowCtx(ctx, &row, query, entryNo); err != nil {
		return nil, err
	}
	return &LedgerEntry{
		EntryNo:      row.EntryNo,
		PaymentNo:    row.PaymentNo,
		OrderNo:      row.OrderNo,
		UserID:       row.UserID,
		ChainType:    row.ChainType,
		Network:      row.Network,
		ChainID:      row.ChainID,
		EntryType:    row.EntryType,
		AssetSymbol:  row.AssetSymbol,
		AssetAddress: row.AssetAddress,
		Amount:       row.Amount,
		Direction:    row.Direction,
		Status:       row.Status,
		Remark:       row.Remark,
	}, nil
}
