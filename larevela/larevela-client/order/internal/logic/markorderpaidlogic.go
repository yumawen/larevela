package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"order/internal/svc"
	"order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkOrderPaidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkOrderPaidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkOrderPaidLogic {
	return &MarkOrderPaidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkOrderPaidLogic) MarkOrderPaid(in *order.MarkOrderPaidReq) (*order.MarkOrderPaidResp, error) {
	if l.svcCtx.TradeModel == nil {
		return nil, fmt.Errorf("database is not configured")
	}
	orderNo := strings.TrimSpace(in.OrderNo)
	if orderNo == "" {
		return nil, fmt.Errorf("orderNo is required")
	}
	l.Infof("order.MarkOrderPaid start orderNo=%s paymentNo=%s", orderNo, strings.TrimSpace(in.PaymentNo))
	paidAt := in.PaidAt
	if paidAt == 0 {
		paidAt = time.Now().UTC().Unix()
	}
	paidAt = normalizeUnixTimestampSeconds(paidAt)

	if err := l.svcCtx.TradeModel.MarkOrderPaid(l.ctx, orderNo, paidAt); err != nil {
		l.Errorf("order.MarkOrderPaid failed orderNo=%s err=%v", orderNo, err)
		return nil, err
	}
	l.Infof("order.MarkOrderPaid success orderNo=%s paidAt=%d", orderNo, paidAt)

	return &order.MarkOrderPaidResp{
		OrderNo: orderNo,
		Status:  "paid",
	}, nil
}
