package logic

import (
	"context"
	"fmt"
	"strings"

	"order/internal/svc"
	"order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloseOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCloseOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloseOrderLogic {
	return &CloseOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CloseOrderLogic) CloseOrder(in *order.CloseOrderReq) (*order.CloseOrderResp, error) {
	if l.svcCtx.TradeModel == nil {
		return nil, fmt.Errorf("database is not configured")
	}
	orderNo := strings.TrimSpace(in.OrderNo)
	if orderNo == "" {
		return nil, fmt.Errorf("orderNo is required")
	}
	l.Infof("order.CloseOrder start orderNo=%s", orderNo)
	if err := l.svcCtx.TradeModel.CloseOrder(l.ctx, orderNo); err != nil {
		l.Errorf("order.CloseOrder failed orderNo=%s err=%v", orderNo, err)
		return nil, err
	}
	l.Infof("order.CloseOrder success orderNo=%s", orderNo)
	return &order.CloseOrderResp{
		OrderNo: orderNo,
		Status:  "closed",
	}, nil
}
