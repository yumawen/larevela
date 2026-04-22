package trademodel

import (
	"context"
	"fmt"
	"time"
)

func (m *Model) CreateOrder(ctx context.Context, in CreateOrderInput) error {
	if err := m.ensureReady(); err != nil {
		return err
	}
	if in.OrderNo == "" || in.BizType == "" || in.BizID == "" || in.Currency == "" || in.Amount == "" {
		return fmt.Errorf("missing required order fields")
	}
	if in.Status == "" {
		in.Status = "created"
	}

	query := `
INSERT INTO orders (
  order_no, biz_type, biz_id, user_id, currency, amount, status, expired_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  biz_type = VALUES(biz_type),
  biz_id = VALUES(biz_id),
  user_id = VALUES(user_id),
  currency = VALUES(currency),
  amount = VALUES(amount),
  status = VALUES(status),
  expired_at = VALUES(expired_at),
  updated_at = NOW()
`
	_, err := m.conn.ExecCtx(
		ctx,
		query,
		in.OrderNo, in.BizType, in.BizID, in.UserID, in.Currency, in.Amount, in.Status, in.ExpiredAt,
	)
	return err
}

func (m *Model) GetOrder(ctx context.Context, orderNo string) (*Order, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}
	if orderNo == "" {
		return nil, fmt.Errorf("orderNo is required")
	}

	var row struct {
		OrderNo   string    `db:"order_no"`
		BizType   string    `db:"biz_type"`
		BizID     string    `db:"biz_id"`
		UserID    int64     `db:"user_id"`
		Currency  string    `db:"currency"`
		Amount    string    `db:"amount"`
		Status    string    `db:"status"`
		ExpiredAt time.Time `db:"expired_at"`
		PaidAt    time.Time `db:"paid_at"`
	}
	query := `
SELECT order_no, biz_type, biz_id, user_id, currency, amount, status, expired_at, paid_at
FROM orders
WHERE order_no = ?
LIMIT 1
`
	if err := m.conn.QueryRowCtx(ctx, &row, query, orderNo); err != nil {
		return nil, err
	}
	return &Order{
		OrderNo:   row.OrderNo,
		BizType:   row.BizType,
		BizID:     row.BizID,
		UserID:    row.UserID,
		Currency:  row.Currency,
		Amount:    row.Amount,
		Status:    row.Status,
		ExpiredAt: row.ExpiredAt.Unix(),
		PaidAt:    row.PaidAt.Unix(),
	}, nil
}

func (m *Model) MarkOrderPaid(ctx context.Context, orderNo string, paidAt int64) error {
	if err := m.ensureReady(); err != nil {
		return err
	}
	if orderNo == "" {
		return fmt.Errorf("orderNo is required")
	}

	t := time.Unix(paidAt, 0)
	query := `
UPDATE orders
SET status = 'paid',
    paid_at = ?,
    updated_at = NOW()
WHERE order_no = ?
`
	_, err := m.conn.ExecCtx(ctx, query, t, orderNo)
	return err
}

func (m *Model) CloseOrder(ctx context.Context, orderNo string) error {
	if err := m.ensureReady(); err != nil {
		return err
	}
	if orderNo == "" {
		return fmt.Errorf("orderNo is required")
	}

	_, err := m.conn.ExecCtx(ctx, "UPDATE orders SET status = 'closed', updated_at = NOW() WHERE order_no = ?", orderNo)
	return err
}
