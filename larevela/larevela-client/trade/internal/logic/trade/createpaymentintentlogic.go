// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package trade

import (
	"context"

	"trade/internal/svc"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentIntentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建支付意图
func NewCreatePaymentIntentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentIntentLogic {
	return &CreatePaymentIntentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePaymentIntentLogic) CreatePaymentIntent(req *types.CreatePaymentIntentReq) (resp *types.CreatePaymentIntentResp, err error) {
	// todo: add your logic here and delete this line

	return
}
