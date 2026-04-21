package logic

import (
	"context"

	"payment/internal/svc"
	"payment/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentLogic {
	return &GetPaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPaymentLogic) GetPayment(in *payment.GetPaymentReq) (*payment.GetPaymentResp, error) {
	// todo: add your logic here and delete this line

	return &payment.GetPaymentResp{}, nil
}
