// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package trade

import (
	"context"
	"fmt"
	"strings"
	"time"

	"order/orderclient"
	"payment/paymentclient"
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
	payerAccount := strings.TrimSpace(req.PayerAccount)
	if payerAccount == "" {
		return nil, fmt.Errorf("payerAccount is required")
	}
	l.Infof("trade.CreatePaymentIntent start payer=%s orderNo=%s referenceId=%s",
		payerAccount, strings.TrimSpace(req.OrderNo), strings.TrimSpace(req.ReferenceId))

	orderNo := strings.TrimSpace(req.OrderNo)
	if orderNo == "" {
		orderResp, createErr := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &orderclient.CreateOrderReq{
			BizType:   "trial",
			BizId:     fmt.Sprintf("%s-%d", strings.TrimSpace(req.ReferenceId), time.Now().UnixNano()),
			UserId:    0,
			Currency:  "USD",
			Amount:    "0",
			ExpiredAt: time.Now().Add(15 * time.Minute).Unix(),
		})
		if createErr != nil {
			l.Errorf("trade.CreatePaymentIntent create order rpc failed payer=%s err=%v", payerAccount, createErr)
			return nil, createErr
		}
		orderNo = orderResp.OrderNo
		l.Infof("trade.CreatePaymentIntent auto created orderNo=%s", orderNo)
	}
	chainType := strings.TrimSpace(req.ChainType)
	if chainType == "" {
		chainType = "solana"
	}
	network := strings.TrimSpace(req.Network)
	if network == "" {
		network = "devnet"
	}
	chainID := req.ChainId
	if chainID == 0 {
		chainID = 103
	}

	rpcResp, err := l.svcCtx.PaymentRpc.CreatePaymentIntent(l.ctx, &paymentclient.CreatePaymentIntentReq{
		OrderNo:           orderNo,
		ChainType:         chainType,
		Network:           network,
		ChainId:           chainID,
		AssetSymbol:       strings.TrimSpace(req.AssetSymbol),
		AssetAddress:      strings.TrimSpace(req.AssetAddress),
		PayerAccount:      payerAccount,
		PayerTokenAccount: strings.TrimSpace(req.PayerTokenAccount),
		ReferenceId:       strings.TrimSpace(req.ReferenceId),
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("create payment intent rpc failed, orderNo=%s err=%v", orderNo, err)
		return nil, err
	}
	l.Infof("trade.CreatePaymentIntent success paymentNo=%s orderNo=%s status=%s",
		rpcResp.PaymentNo, rpcResp.OrderNo, rpcResp.Status)

	return &types.CreatePaymentIntentResp{
		PaymentNo:            rpcResp.PaymentNo,
		OrderNo:              rpcResp.OrderNo,
		ChainType:            rpcResp.ChainType,
		Network:              rpcResp.Network,
		ChainId:              rpcResp.ChainId,
		PayMode:              rpcResp.PayMode,
		ReceiverAccount:      rpcResp.ReceiverAccount,
		ReceiverTokenAccount: rpcResp.ReceiverTokenAccount,
		AssetAddress:         rpcResp.AssetAddress,
		AssetSymbol:          rpcResp.AssetSymbol,
		AmountExpected:       rpcResp.AmountExpected,
		Decimals:             rpcResp.Decimals,
		ReferenceId:          rpcResp.ReferenceId,
		ContractAddress:      rpcResp.ContractAddress,
		Method:               rpcResp.Method,
		Calldata:             rpcResp.Calldata,
		Value:                rpcResp.Value,
		SerializedMessage:    rpcResp.SerializedMessage,
		QuoteExpiredAt:       rpcResp.QuoteExpiredAt,
		Status:               rpcResp.Status,
	}, nil
}
