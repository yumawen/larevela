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

	queryCursors := `
INSERT INTO chain_scan (
  chain_type, network, chain_id, cursor_type, last_scanned_block, last_scanned_slot,
  last_scanned_tx_id, remark, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())
ON DUPLICATE KEY UPDATE
  last_scanned_block = GREATEST(last_scanned_block, VALUES(last_scanned_block)),
  last_scanned_slot = GREATEST(last_scanned_slot, VALUES(last_scanned_slot)),
  last_scanned_tx_id = VALUES(last_scanned_tx_id),
  remark = VALUES(remark),
  updated_at = UTC_TIMESTAMP()
`
	queryLegacy := `
INSERT INTO chain_scan (
  chain_type, network, chain_id, cursor_type, last_scanned_block, last_scanned_slot,
  last_scanned_tx_id, remark, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())
ON DUPLICATE KEY UPDATE
  last_scanned_block = GREATEST(last_scanned_block, VALUES(last_scanned_block)),
  last_scanned_slot = GREATEST(last_scanned_slot, VALUES(last_scanned_slot)),
  last_scanned_tx_id = VALUES(last_scanned_tx_id),
  remark = VALUES(remark),
  updated_at = UTC_TIMESTAMP()
`
	args := []any{
		in.ChainType,
		in.Network,
		in.ChainID,
		in.CursorType,
		in.LastScannedBlock,
		in.LastScannedSlot,
		strings.TrimSpace(in.LastScannedTxID),
		strings.TrimSpace(in.Remark),
	}
	_, err := m.conn.ExecCtx(
		ctx,
		queryCursors,
		args...,
	)
	if err == nil {
		return nil
	}
	// Backward compatibility for environments that still use the legacy `chain_scan` table.
	if strings.Contains(err.Error(), "1146") || strings.Contains(err.Error(), "doesn't exist") {
		_, fallbackErr := m.conn.ExecCtx(ctx, queryLegacy, args...)
		return fallbackErr
	}
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
) VALUES (?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())
ON DUPLICATE KEY UPDATE
  status = VALUES(status),
  response_snapshot = VALUES(response_snapshot),
  updated_at = UTC_TIMESTAMP()
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
