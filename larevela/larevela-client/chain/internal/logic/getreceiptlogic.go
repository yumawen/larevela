package logic

import (
	"context"

	"chain/chain"
	"chain/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReceiptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReceiptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReceiptLogic {
	return &GetReceiptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetReceiptLogic) GetReceipt(in *chain.GetReceiptReq) (*chain.GetReceiptResp, error) {
	// todo: add your logic here and delete this line

	return &chain.GetReceiptResp{}, nil
}
