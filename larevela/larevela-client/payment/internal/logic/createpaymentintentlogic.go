package logic

import (
	"context"

	"payment/internal/svc"
	"payment/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentIntentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePaymentIntentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentIntentLogic {
	return &CreatePaymentIntentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePaymentIntentLogic) CreatePaymentIntent(in *payment.CreatePaymentIntentReq) (*payment.CreatePaymentIntentResp, error) {
	// todo: add your logic here and delete this line

	return &payment.CreatePaymentIntentResp{}, nil
}
