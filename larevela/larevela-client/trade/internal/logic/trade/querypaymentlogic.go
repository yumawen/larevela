// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package trade

import (
	"context"
	"fmt"
	"strings"

	"payment/paymentclient"
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

func (l *QueryPaymentLogic) QueryPayment(req *types.QueryPaymentReq) (resp *types.QueryPaymentResp, err error) {
	paymentNo := strings.TrimSpace(req.PaymentNo)
	if paymentNo == "" {
		return nil, fmt.Errorf("paymentNo is required")
	}
	l.Infof("trade.QueryPayment start paymentNo=%s", paymentNo)

	rpcResp, err := l.svcCtx.PaymentRpc.GetPayment(l.ctx, &paymentclient.GetPaymentReq{
		PaymentNo: paymentNo,
	})
	if err != nil {
		l.Errorf("trade.QueryPayment rpc failed paymentNo=%s err=%v", paymentNo, err)
		return nil, err
	}
	l.Infof("trade.QueryPayment success paymentNo=%s status=%s confirmations=%d",
		rpcResp.PaymentNo, rpcResp.Status, rpcResp.Confirmations)

	return &types.QueryPaymentResp{
		PaymentNo:            rpcResp.PaymentNo,
		OrderNo:              rpcResp.OrderNo,
		TxId:                 rpcResp.TxId,
		ChainType:            rpcResp.ChainType,
		Network:              rpcResp.Network,
		ChainId:              rpcResp.ChainId,
		PayerAccount:         rpcResp.PayerAccount,
		PayerTokenAccount:    rpcResp.PayerTokenAccount,
		ReceiverAccount:      rpcResp.ReceiverAccount,
		ReceiverTokenAccount: rpcResp.ReceiverTokenAccount,
		AssetAddress:         rpcResp.AssetAddress,
		AssetSymbol:          rpcResp.AssetSymbol,
		AmountExpected:       rpcResp.AmountExpected,
		AmountActual:         rpcResp.AmountActual,
		Confirmations:        rpcResp.Confirmations,
		ConfirmationStatus:   rpcResp.ConfirmationStatus,
		Status:               rpcResp.Status,
		FailureReason:        rpcResp.FailureReason,
	}, nil
}
