package logic

import (
	"context"

	"chain/chain"
	"chain/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionLogic {
	return &GetTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTransactionLogic) GetTransaction(in *chain.GetTransactionReq) (*chain.GetTransactionResp, error) {
	// todo: add your logic here and delete this line

	return &chain.GetTransactionResp{}, nil
}
