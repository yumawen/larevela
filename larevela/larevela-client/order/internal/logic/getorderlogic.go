package logic

import (
	"context"
	"fmt"
	"strings"

	"order/internal/svc"
	"order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderLogic) GetOrder(in *order.GetOrderReq) (*order.GetOrderResp, error) {
	if l.svcCtx.TradeModel == nil {
		return nil, fmt.Errorf("database is not configured")
	}
	orderNo := strings.TrimSpace(in.OrderNo)
	if orderNo == "" {
		return nil, fmt.Errorf("orderNo is required")
	}
	l.Infof("order.GetOrder start orderNo=%s", orderNo)

	o, err := l.svcCtx.TradeModel.GetOrder(l.ctx, orderNo)
	if err != nil {
		l.Errorf("order.GetOrder failed orderNo=%s err=%v", orderNo, err)
		return nil, err
	}
	l.Infof("order.GetOrder success orderNo=%s status=%s", o.OrderNo, o.Status)

	return &order.GetOrderResp{
		OrderNo:   o.OrderNo,
		BizType:   o.BizType,
		BizId:     o.BizID,
		UserId:    o.UserID,
		Currency:  o.Currency,
		Amount:    o.Amount,
		Status:    o.Status,
		PaymentNo: "",
		ExpiredAt: o.ExpiredAt,
		PaidAt:    o.PaidAt,
		Found:     true,
	}, nil
}
