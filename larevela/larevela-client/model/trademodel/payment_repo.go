package trademodel

import (
	"context"
	"fmt"
	"time"
)

func (m *Model) CreateTrialPaymentTx(ctx context.Context, in CreateTrialPaymentTxInput) error {
	if err := m.ensureReady(); err != nil {
		return err
	}
	if in.PaymentNo == "" || in.TxID == "" || in.FromAccount == "" || in.ToAccount == "" || in.AmountSol == "" {
		return fmt.Errorf("missing required transaction fields")
	}
	if in.ChainType == "" {
		in.ChainType = "solana"
	}
	if in.Network == "" {
		in.Network = "devnet"
	}
	if in.ChainID == 0 {
		in.ChainID = 103
	}
	if in.AssetSymbol == "" {
		in.AssetSymbol = "SOL"
	}

	now := time.Now()
	query := `
INSERT INTO payment_transactions (
  payment_no, tx_id, chain_type, network, chain_id, from_account, to_account,
  asset_address, asset_symbol, amount_actual, tx_status, confirmations,
  confirmation_status, reference_id, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'success', 1, 'confirmed', ?, ?, ?)
ON DUPLICATE KEY UPDATE
  chain_type = VALUES(chain_type),
  network = VALUES(network),
  chain_id = VALUES(chain_id),
  from_account = VALUES(from_account),
  to_account = VALUES(to_account),
  asset_address = VALUES(asset_address),
  asset_symbol = VALUES(asset_symbol),
  amount_actual = VALUES(amount_actual),
  tx_status = VALUES(tx_status),
  confirmations = VALUES(confirmations),
  confirmation_status = VALUES(confirmation_status),
  reference_id = VALUES(reference_id),
  updated_at = VALUES(updated_at)
`

	_, err := m.conn.ExecCtx(
		ctx,
		query,
		in.PaymentNo,
		in.TxID,
		in.ChainType,
		in.Network,
		in.ChainID,
		in.FromAccount,
		in.ToAccount,
		in.AssetAddress,
		in.AssetSymbol,
		in.AmountSol,
		in.PlanType,
		now,
		now,
	)
	return err
}

func (m *Model) CreatePaymentIntent(ctx context.Context, in CreatePaymentIntentInput) error {
	if err := m.ensureReady(); err != nil {
		return err
	}
	if in.PaymentNo == "" || in.OrderNo == "" || in.ChainType == "" || in.Network == "" ||
		in.PayerAccount == "" || in.ReceiverAccount == "" || in.AssetSymbol == "" || in.AmountExpected == "" {
		return fmt.Errorf("missing required payment intent fields")
	}
	if in.Status == "" {
		in.Status = "created"
	}

	query := `
INSERT INTO payment_intents (
  payment_no, order_no, chain_type, network, chain_id, pay_mode, payer_account, receiver_account,
  asset_symbol, asset_address, amount_expected, decimals, reference_id, calldata, quote_expired_at,
  confirmations, confirmation_status, status, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, 'transfer', ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 'pending', ?, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  payer_account = VALUES(payer_account),
  receiver_account = VALUES(receiver_account),
  amount_expected = VALUES(amount_expected),
  reference_id = VALUES(reference_id),
  calldata = VALUES(calldata),
  quote_expired_at = VALUES(quote_expired_at),
  updated_at = NOW()
`

	_, err := m.conn.ExecCtx(
		ctx,
		query,
		in.PaymentNo,
		in.OrderNo,
		in.ChainType,
		in.Network,
		in.ChainID,
		in.PayerAccount,
		in.ReceiverAccount,
		in.AssetSymbol,
		in.AssetAddress,
		in.AmountExpected,
		in.Decimals,
		in.ReferenceID,
		in.SerializedMessage,
		in.QuoteExpiredAt,
		in.Status,
	)
	return err
}

func (m *Model) MarkPaymentSubmitted(ctx context.Context, paymentNo, txID, fromAccount string) error {
	if err := m.ensureReady(); err != nil {
		return err
	}
	if paymentNo == "" || txID == "" {
		return fmt.Errorf("paymentNo and txID are required")
	}

	query := `
UPDATE payment_intents
SET tx_id = ?,
    payer_account = CASE WHEN ? <> '' THEN ? ELSE payer_account END,
    status = 'submitted',
    confirmation_status = 'pending',
    updated_at = NOW()
WHERE payment_no = ?
`
	_, err := m.conn.ExecCtx(ctx, query, txID, fromAccount, fromAccount, paymentNo)
	return err
}

func (m *Model) GetPaymentIntent(ctx context.Context, paymentNo string) (*PaymentIntent, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}
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
		AssetAddress       string `db:"asset_address"`
		AssetSymbol        string `db:"asset_symbol"`
		AmountExpected     string `db:"amount_expected"`
		AmountActual       string `db:"amount_actual"`
		Confirmations      int64  `db:"confirmations"`
		ConfirmationStatus string `db:"confirmation_status"`
		Status             string `db:"status"`
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
       IFNULL(asset_address, '') AS asset_address,
       asset_symbol,
       amount_expected,
       IFNULL(amount_actual, '') AS amount_actual,
       confirmations,
       IFNULL(confirmation_status, '') AS confirmation_status,
       status,
       IFNULL(failure_reason, '') AS failure_reason
FROM payment_intents
WHERE payment_no = ?
LIMIT 1
`
	if err := m.conn.QueryRowCtx(ctx, &row, query, paymentNo); err != nil {
		return nil, err
	}
	return &PaymentIntent{
		PaymentNo:          row.PaymentNo,
		OrderNo:            row.OrderNo,
		TxID:               row.TxID,
		ChainType:          row.ChainType,
		Network:            row.Network,
		ChainID:            row.ChainID,
		PayerAccount:       row.PayerAccount,
		ReceiverAccount:    row.ReceiverAccount,
		AssetAddress:       row.AssetAddress,
		AssetSymbol:        row.AssetSymbol,
		AmountExpected:     row.AmountExpected,
		AmountActual:       row.AmountActual,
		Confirmations:      row.Confirmations,
		ConfirmationStatus: row.ConfirmationStatus,
		Status:             row.Status,
		FailureReason:      row.FailureReason,
	}, nil
}

func (m *Model) UpdatePaymentConfirmation(
	ctx context.Context,
	paymentNo string,
	confirmations int64,
	confirmationStatus, amountActual, status, failureReason string,
) error {
	if err := m.ensureReady(); err != nil {
		return err
	}
	if paymentNo == "" {
		return fmt.Errorf("paymentNo is required")
	}

	query := `
UPDATE payment_intents
SET confirmations = ?,
    confirmation_status = ?,
    amount_actual = ?,
    status = ?,
    failure_reason = ?,
    paid_at = CASE WHEN ? = 'paid' THEN NOW() ELSE paid_at END,
    updated_at = NOW()
WHERE payment_no = ?
`
	_, err := m.conn.ExecCtx(
		ctx,
		query,
		confirmations,
		confirmationStatus,
		amountActual,
		status,
		failureReason,
		status,
		paymentNo,
	)
	return err
}
