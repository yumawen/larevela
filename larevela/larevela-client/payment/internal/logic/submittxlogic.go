package logic

import (
	"context"

	"payment/internal/svc"
	"payment/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitTxLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitTxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitTxLogic {
	return &SubmitTxLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SubmitTxLogic) SubmitTx(in *payment.SubmitTxReq) (*payment.SubmitTxResp, error) {
	// todo: add your logic here and delete this line

	return &payment.SubmitTxResp{}, nil
}
