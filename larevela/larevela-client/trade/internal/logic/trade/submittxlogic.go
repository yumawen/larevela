// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package trade

import (
	"context"
	"fmt"
	"strings"
	"time"

	"payment/paymentclient"
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
	paymentNo := strings.TrimSpace(req.PaymentNo)
	txID := strings.TrimSpace(req.TxId)
	fromAccount := strings.TrimSpace(req.FromAccount)
	if paymentNo == "" {
		paymentNo = fmt.Sprintf("trial-%d", time.Now().UnixNano())
	}
	if txID == "" || fromAccount == "" {
		return nil, fmt.Errorf("paymentNo/txId/fromAccount are required")
	}
	l.Infof("trade.SubmitTx start paymentNo=%s txId=%s from=%s", paymentNo, txID, fromAccount)

	rpcResp, err := l.svcCtx.PaymentRpc.SubmitTx(l.ctx, &paymentclient.SubmitTxReq{
		PaymentNo:        paymentNo,
		TxId:             txID,
		FromAccount:      fromAccount,
		FromTokenAccount: strings.TrimSpace(req.FromTokenAccount),
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("submit tx rpc failed, paymentNo=%s txId=%s err=%v", paymentNo, txID, err)
		return nil, err
	}
	l.Infof("trade.SubmitTx success paymentNo=%s txId=%s status=%s confirmation=%s",
		rpcResp.PaymentNo, rpcResp.TxId, rpcResp.Status, rpcResp.ConfirmationStatus)

	return &types.SubmitTxResp{
		PaymentNo:          rpcResp.PaymentNo,
		TxId:               rpcResp.TxId,
		ConfirmationStatus: rpcResp.ConfirmationStatus,
		Status:             rpcResp.Status,
	}, nil
}
