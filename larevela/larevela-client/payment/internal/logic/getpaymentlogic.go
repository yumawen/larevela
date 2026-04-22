package logic

import (
	"context"
	"fmt"
	"strings"

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
	if l.svcCtx.TradeModel == nil {
		return nil, fmt.Errorf("database is not configured")
	}
	paymentNo := strings.TrimSpace(in.PaymentNo)
	if paymentNo == "" {
		return nil, fmt.Errorf("paymentNo is required")
	}
	l.Infof("payment.GetPayment start paymentNo=%s", paymentNo)

	view, viewErr := l.svcCtx.TradeModel.GetPaymentView(l.ctx, paymentNo)
	if viewErr == nil {
		l.Infof("payment.GetPayment hit payment_view paymentNo=%s status=%s", view.PaymentNo, view.Status)
		return &payment.GetPaymentResp{
			PaymentNo:            view.PaymentNo,
			OrderNo:              view.OrderNo,
			TxId:                 view.TxID,
			ChainType:            view.ChainType,
			Network:              view.Network,
			ChainId:              view.ChainID,
			PayerAccount:         view.PayerAccount,
			PayerTokenAccount:    "",
			ReceiverAccount:      view.ReceiverAccount,
			ReceiverTokenAccount: "",
			AssetAddress:         "",
			AssetSymbol:          view.AssetSymbol,
			AmountExpected:       view.AmountExpected,
			AmountActual:         view.AmountActual,
			Confirmations:        view.Confirmations,
			ConfirmationStatus:   view.ConfirmationStatus,
			Status:               view.Status,
			FailureReason:        view.FailureReason,
		}, nil
	}

	pi, err := l.svcCtx.TradeModel.GetPaymentIntent(l.ctx, paymentNo)
	if err != nil {
		l.Errorf("payment.GetPayment failed paymentNo=%s err=%v", paymentNo, err)
		return nil, err
	}
	l.Infof("payment.GetPayment success paymentNo=%s status=%s", pi.PaymentNo, pi.Status)

	return &payment.GetPaymentResp{
		PaymentNo:            pi.PaymentNo,
		OrderNo:              pi.OrderNo,
		TxId:                 pi.TxID,
		ChainType:            pi.ChainType,
		Network:              pi.Network,
		ChainId:              pi.ChainID,
		PayerAccount:         pi.PayerAccount,
		PayerTokenAccount:    "",
		ReceiverAccount:      pi.ReceiverAccount,
		ReceiverTokenAccount: "",
		AssetAddress:         pi.AssetAddress,
		AssetSymbol:          pi.AssetSymbol,
		AmountExpected:       pi.AmountExpected,
		AmountActual:         pi.AmountActual,
		Confirmations:        pi.Confirmations,
		ConfirmationStatus:   pi.ConfirmationStatus,
		Status:               pi.Status,
		FailureReason:        pi.FailureReason,
	}, nil
}
