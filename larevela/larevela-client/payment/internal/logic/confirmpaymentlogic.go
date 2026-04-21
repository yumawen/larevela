package logic

import (
	"context"

	"payment/internal/svc"
	"payment/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmPaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmPaymentLogic {
	return &ConfirmPaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfirmPaymentLogic) ConfirmPayment(in *payment.ConfirmPaymentReq) (*payment.ConfirmPaymentResp, error) {
	// todo: add your logic here and delete this line

	return &payment.ConfirmPaymentResp{}, nil
}
