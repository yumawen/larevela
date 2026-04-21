// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package trade

import (
	"context"

	"trade/internal/svc"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询支付状态
func NewQueryPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPaymentLogic {
	return &QueryPaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryPaymentLogic) QueryPayment() (resp *types.QueryPaymentResp, err error) {
	// todo: add your logic here and delete this line

	return
}
