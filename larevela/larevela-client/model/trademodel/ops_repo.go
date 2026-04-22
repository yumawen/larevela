package trademodel

import (
	"context"
	"fmt"
	"strings"
)

func (m *Model) UpsertChainScanCursor(ctx context.Context, in UpsertChainScanCursorInput) error {
	if err := m.ensureReady(); err != nil {
		return err
	}
	if strings.TrimSpace(in.ChainType) == "" || strings.TrimSpace(in.Network) == "" ||
		in.ChainID == 0 || strings.TrimSpace(in.CursorType) == "" {
		return fmt.Errorf("missing required chain cursor fields")
	}

	query := `
INSERT INTO chain_scan (
  chain_type, network, chain_id, cursor_type, last_scanned_block, last_scanned_slot,
  last_scanned_tx_id, remark, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  last_scanned_block = GREATEST(last_scanned_block, VALUES(last_scanned_block)),
  last_scanned_slot = GREATEST(last_scanned_slot, VALUES(last_scanned_slot)),
  last_scanned_tx_id = VALUES(last_scanned_tx_id),
  remark = VALUES(remark),
  updated_at = NOW()
`
	_, err := m.conn.ExecCtx(
		ctx,
		query,
		in.ChainType,
		in.Network,
		in.ChainID,
		in.CursorType,
		in.LastScannedBlock,
		in.LastScannedSlot,
		strings.TrimSpace(in.LastScannedTxID),
		strings.TrimSpace(in.Remark),
	)
	return err
}

func (m *Model) UpsertIdempotencyRecord(ctx context.Context, in UpsertIdempotencyRecordInput) error {
	if err := m.ensureReady(); err != nil {
		return err
	}
	if strings.TrimSpace(in.IdemKey) == "" || strings.TrimSpace(in.BizType) == "" ||
		strings.TrimSpace(in.BizNo) == "" || strings.TrimSpace(in.Status) == "" {
		return fmt.Errorf("missing required idempotency fields")
	}

	var snapshot any
	if strings.TrimSpace(in.ResponseSnapshot) != "" {
		snapshot = in.ResponseSnapshot
	}

	query := `
INSERT INTO idempotency_records (
  idem_key, biz_type, biz_no, status, response_snapshot, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  status = VALUES(status),
  response_snapshot = VALUES(response_snapshot),
  updated_at = NOW()
`
	_, err := m.conn.ExecCtx(
		ctx,
		query,
		strings.TrimSpace(in.IdemKey),
		strings.TrimSpace(in.BizType),
		strings.TrimSpace(in.BizNo),
		strings.TrimSpace(in.Status),
		snapshot,
	)
	return err
}
