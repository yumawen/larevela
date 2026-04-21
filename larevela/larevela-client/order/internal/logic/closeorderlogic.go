package logic

import (
	"context"

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
	return l.svcCtx.Store.CloseOrder(in)
}
