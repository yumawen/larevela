// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package trade

import (
	"context"

	"trade/internal/svc"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitTxLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 前端提交链上交易ID
func NewSubmitTxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitTxLogic {
	return &SubmitTxLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitTxLogic) SubmitTx(req *types.SubmitTxReq) (resp *types.SubmitTxResp, err error) {
	// todo: add your logic here and delete this line

	return
}
