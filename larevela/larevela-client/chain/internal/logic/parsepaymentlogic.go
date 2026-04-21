package logic

import (
	"context"

	"chain/chain"
	"chain/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ParsePaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewParsePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParsePaymentLogic {
	return &ParsePaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ParsePaymentLogic) ParsePayment(in *chain.ParsePaymentReq) (*chain.ParsePaymentResp, error) {
	// todo: add your logic here and delete this line

	return &chain.ParsePaymentResp{}, nil
}
