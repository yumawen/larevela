package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"order/internal/svc"
	"order/order"
	"trademodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderReq) (*order.CreateOrderResp, error) {
	if l.svcCtx.TradeModel == nil {
		return nil, fmt.Errorf("database is not configured")
	}

	orderNo := fmt.Sprintf("order-%d", time.Now().UnixNano())
	if strings.TrimSpace(in.BizType) == "" || strings.TrimSpace(in.BizId) == "" || strings.TrimSpace(in.Amount) == "" {
		return nil, fmt.Errorf("bizType/bizId/amount are required")
	}
	l.Infof("order.CreateOrder start orderNo=%s bizType=%s bizId=%s userId=%d",
		orderNo, strings.TrimSpace(in.BizType), strings.TrimSpace(in.BizId), in.UserId)

	expiredAt := time.Now().UTC().Add(15 * time.Minute)
	if in.ExpiredAt > 0 {
		expiredAt = time.Unix(normalizeUnixTimestampSeconds(in.ExpiredAt), 0).UTC()
	}

	err := l.svcCtx.TradeModel.CreateOrder(l.ctx, trademodel.CreateOrderInput{
		OrderNo:   orderNo,
		BizType:   strings.TrimSpace(in.BizType),
		BizID:     strings.TrimSpace(in.BizId),
		UserID:    in.UserId,
		Currency:  strings.TrimSpace(in.Currency),
		Amount:    strings.TrimSpace(in.Amount),
		Status:    "created",
		ExpiredAt: expiredAt,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("create order failed, orderNo=%s err=%v", orderNo, err)
		return nil, fmt.Errorf("create order failed")
	}
	l.Infof("order.CreateOrder success orderNo=%s status=created expiredAt=%d", orderNo, expiredAt.Unix())

	return &order.CreateOrderResp{
		OrderNo:   orderNo,
		Status:    "created",
		ExpiredAt: expiredAt.Unix(),
	}, nil
}
