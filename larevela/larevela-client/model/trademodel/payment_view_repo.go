package trademodel

import (
	"context"
	"fmt"
	"strings"
)

func (m *Model) UpsertPaymentView(ctx context.Context, in UpsertPaymentViewInput) error {
	if err := m.ensureReady(); err != nil {
		return err
	}
	if strings.TrimSpace(in.PaymentNo) == "" || strings.TrimSpace(in.OrderNo) == "" || strings.TrimSpace(in.AssetSymbol) == "" {
		return fmt.Errorf("missing required payment view fields")
	}
	if strings.TrimSpace(in.ChainType) == "" {
		in.ChainType = "solana"
	}
	if strings.TrimSpace(in.Network) == "" {
		in.Network = "devnet"
	}
	if strings.TrimSpace(in.Status) == "" {
		in.Status = "created"
	}
	if strings.TrimSpace(in.UpdatedSource) == "" {
		in.UpdatedSource = "payment"
	}

	query := `
INSERT INTO payment_view (
  payment_no, order_no, tx_id, chain_type, network, chain_id, payer_account, receiver_account,
  asset_symbol, amount_expected, amount_actual, status, confirmation_status, confirmations,
  last_scanned_block, last_scanned_slot, failure_reason, updated_source, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())
ON DUPLICATE KEY UPDATE
  order_no = VALUES(order_no),
  tx_id = CASE WHEN VALUES(tx_id) <> '' THEN VALUES(tx_id) ELSE tx_id END,
  chain_type = VALUES(chain_type),
  network = VALUES(network),
  chain_id = VALUES(chain_id),
  payer_account = CASE WHEN VALUES(payer_account) <> '' THEN VALUES(payer_account) ELSE payer_account END,
  receiver_account = CASE WHEN VALUES(receiver_account) <> '' THEN VALUES(receiver_account) ELSE receiver_account END,
  asset_symbol = VALUES(asset_symbol),
  amount_expected = VALUES(amount_expected),
  amount_actual = CASE WHEN VALUES(amount_actual) <> '' THEN VALUES(amount_actual) ELSE amount_actual END,
  status = VALUES(status),
  confirmation_status = VALUES(confirmation_status),
  confirmations = GREATEST(confirmations, VALUES(confirmations)),
  last_scanned_block = GREATEST(last_scanned_block, VALUES(last_scanned_block)),
  last_scanned_slot = GREATEST(last_scanned_slot, VALUES(last_scanned_slot)),
  failure_reason = VALUES(failure_reason),
  updated_source = VALUES(updated_source),
  updated_at = UTC_TIMESTAMP()
`

	_, err := m.conn.ExecCtx(
		ctx,
		query,
		strings.TrimSpace(in.PaymentNo),
		strings.TrimSpace(in.OrderNo),
		strings.TrimSpace(in.TxID),
		strings.TrimSpace(in.ChainType),
		strings.TrimSpace(in.Network),
		in.ChainID,
		strings.TrimSpace(in.PayerAccount),
		strings.TrimSpace(in.ReceiverAccount),
		strings.TrimSpace(in.AssetSymbol),
		strings.TrimSpace(in.AmountExpected),
		strings.TrimSpace(in.AmountActual),
		strings.TrimSpace(in.Status),
		strings.TrimSpace(in.ConfirmationStatus),
		in.Confirmations,
		in.LastScannedBlock,
		in.LastScannedSlot,
		strings.TrimSpace(in.FailureReason),
		strings.TrimSpace(in.UpdatedSource),
	)
	return err
}

func (m *Model) GetPaymentView(ctx context.Context, paymentNo string) (*PaymentView, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}
	paymentNo = strings.TrimSpace(paymentNo)
	if paymentNo == "" {
		return nil, fmt.Errorf("paymentNo is required")
	}

	var row struct {
		PaymentNo          string `db:"payment_no"`
		OrderNo            string `db:"order_no"`
		TxID               string `db:"tx_id"`
		ChainType          string `db:"chain_type"`
		Network            string `db:"network"`
		ChainID            int64  `db:"chain_id"`
		PayerAccount       string `db:"payer_account"`
		ReceiverAccount    string `db:"receiver_account"`
		AssetSymbol        string `db:"asset_symbol"`
		AmountExpected     string `db:"amount_expected"`
		AmountActual       string `db:"amount_actual"`
		Status             string `db:"status"`
		ConfirmationStatus string `db:"confirmation_status"`
		Confirmations      int64  `db:"confirmations"`
		LastScannedBlock   int64  `db:"last_scanned_block"`
		LastScannedSlot    int64  `db:"last_scanned_slot"`
		FailureReason      string `db:"failure_reason"`
	}

	query := `
SELECT payment_no,
       order_no,
       IFNULL(tx_id, '') AS tx_id,
       chain_type,
       network,
       chain_id,
       IFNULL(payer_account, '') AS payer_account,
       IFNULL(receiver_account, '') AS receiver_account,
       asset_symbol,
       amount_expected,
       IFNULL(amount_actual, '') AS amount_actual,
       status,
       IFNULL(confirmation_status, '') AS confirmation_status,
       confirmations,
       last_scanned_block,
       last_scanned_slot,
       IFNULL(failure_reason, '') AS failure_reason
FROM payment_view
WHERE payment_no = ?
LIMIT 1
`

	if err := m.conn.QueryRowCtx(ctx, &row, query, paymentNo); err != nil {
		return nil, err
	}

	return &PaymentView{
		PaymentNo:          row.PaymentNo,
		OrderNo:            row.OrderNo,
		TxID:               row.TxID,
		ChainType:          row.ChainType,
		Network:            row.Network,
		ChainID:            row.ChainID,
		PayerAccount:       row.PayerAccount,
		ReceiverAccount:    row.ReceiverAccount,
		AssetSymbol:        row.AssetSymbol,
		AmountExpected:     row.AmountExpected,
		AmountActual:       row.AmountActual,
		Status:             row.Status,
		ConfirmationStatus: row.ConfirmationStatus,
		Confirmations:      row.Confirmations,
		LastScannedBlock:   row.LastScannedBlock,
		LastScannedSlot:    row.LastScannedSlot,
		FailureReason:      row.FailureReason,
	}, nil
}
