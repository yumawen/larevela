package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"payment/internal/svc"
	"payment/payment"
	"trademodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitTxLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitTxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitTxLogic {
	return &SubmitTxLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SubmitTxLogic) SubmitTx(in *payment.SubmitTxReq) (*payment.SubmitTxResp, error) {
	if l.svcCtx.TradeModel == nil {
		return nil, fmt.Errorf("database is not configured")
	}

	paymentNo := strings.TrimSpace(in.PaymentNo)
	if paymentNo == "" {
		paymentNo = fmt.Sprintf("trial-%d", time.Now().UnixNano())
	}

	txID := strings.TrimSpace(in.TxId)
	fromAccount := strings.TrimSpace(in.FromAccount)
	if txID == "" || fromAccount == "" {
		return nil, fmt.Errorf("paymentNo/txId/fromAccount are required")
	}
	idemKey := fmt.Sprintf("submit_tx:%s:%s", paymentNo, txID)
	processingSnapshot, _ := json.Marshal(map[string]any{
		"paymentNo": paymentNo,
		"txId":      txID,
		"status":    "processing",
	})
	if idemErr := l.svcCtx.TradeModel.UpsertIdempotencyRecord(l.ctx, trademodel.UpsertIdempotencyRecordInput{
		IdemKey:          idemKey,
		BizType:          "submit_tx",
		BizNo:            paymentNo,
		Status:           "processing",
		ResponseSnapshot: string(processingSnapshot),
	}); idemErr != nil {
		logx.WithContext(l.ctx).Errorf("persist idempotency submit_tx processing failed, paymentNo=%s txId=%s err=%v", paymentNo, txID, idemErr)
	}
	l.Infof("payment.SubmitTx start paymentNo=%s txId=%s from=%s", paymentNo, txID, fromAccount)

	receiverAccount := strings.TrimSpace(l.svcCtx.Config.TrialTransfer.ReceiverAccount)
	pi, piErr := l.svcCtx.TradeModel.GetPaymentIntent(l.ctx, paymentNo)
	if piErr != nil {
		logx.WithContext(l.ctx).Errorf("load payment intent before persist tx failed, paymentNo=%s err=%v", paymentNo, piErr)
		return nil, fmt.Errorf("payment intent not found")
	}
	if strings.TrimSpace(pi.ReceiverAccount) != "" {
		receiverAccount = strings.TrimSpace(pi.ReceiverAccount)
	}
	if receiverAccount == "" {
		return nil, fmt.Errorf("receiverAccount is required")
	}
	amountExpected := strings.TrimSpace(pi.AmountExpected)
	if amountExpected == "" {
		amountExpected = strings.TrimSpace(l.svcCtx.Config.TrialTransfer.AmountSol)
		if amountExpected == "" {
			amountExpected = "0.1"
		}
	}

	fromTokenAccount := strings.TrimSpace(in.FromTokenAccount)
	if fromTokenAccount == "" {
		fromTokenAccount = strings.TrimSpace(pi.PayerTokenAccount)
	}
	toTokenAccount := strings.TrimSpace(pi.ReceiverTokenAccount)

	err := l.svcCtx.TradeModel.CreateTrialPaymentTx(l.ctx, trademodel.CreateTrialPaymentTxInput{
		PaymentNo:    paymentNo,
		TxID:         txID,
		FromAccount:  fromAccount,
		ToAccount:    receiverAccount,
		FromTokenAccount: fromTokenAccount,
		ToTokenAccount:   toTokenAccount,
		AmountSol:    amountExpected,
		ChainType:    pi.ChainType,
		Network:      pi.Network,
		ChainID:      pi.ChainID,
		AssetSymbol:  pi.AssetSymbol,
		AssetAddress: pi.AssetAddress,
		PlanType:     "free_trial",
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("persist submit tx failed, paymentNo=%s txId=%s err=%v", paymentNo, txID, err)
		failSnapshot, _ := json.Marshal(map[string]any{
			"paymentNo": paymentNo,
			"txId":      txID,
			"status":    "failed",
			"reason":    "persist transaction failed",
		})
		_ = l.svcCtx.TradeModel.UpsertIdempotencyRecord(l.ctx, trademodel.UpsertIdempotencyRecordInput{
			IdemKey:          idemKey,
			BizType:          "submit_tx",
			BizNo:            paymentNo,
			Status:           "failed",
			ResponseSnapshot: string(failSnapshot),
		})
		return nil, fmt.Errorf("persist transaction failed")
	}
	l.Infof("payment.SubmitTx persisted transaction paymentNo=%s txId=%s", paymentNo, txID)

	if err = l.svcCtx.TradeModel.MarkPaymentSubmitted(l.ctx, paymentNo, txID, fromAccount); err != nil {
		logx.WithContext(l.ctx).Errorf("mark payment submitted failed, paymentNo=%s txId=%s err=%v", paymentNo, txID, err)
		failSnapshot, _ := json.Marshal(map[string]any{
			"paymentNo": paymentNo,
			"txId":      txID,
			"status":    "failed",
			"reason":    "update payment status failed",
		})
		_ = l.svcCtx.TradeModel.UpsertIdempotencyRecord(l.ctx, trademodel.UpsertIdempotencyRecordInput{
			IdemKey:          idemKey,
			BizType:          "submit_tx",
			BizNo:            paymentNo,
			Status:           "failed",
			ResponseSnapshot: string(failSnapshot),
		})
		return nil, fmt.Errorf("update payment status failed")
	}
	l.Infof("payment.SubmitTx marked submitted paymentNo=%s", paymentNo)
	if viewErr := l.svcCtx.TradeModel.UpsertPaymentView(l.ctx, trademodel.UpsertPaymentViewInput{
		PaymentNo:          paymentNo,
		OrderNo:            pi.OrderNo,
		TxID:               txID,
		ChainType:          pi.ChainType,
		Network:            pi.Network,
		ChainID:            pi.ChainID,
		PayerAccount:       fromAccount,
		ReceiverAccount:    receiverAccount,
		AssetSymbol:        pi.AssetSymbol,
		AmountExpected:     pi.AmountExpected,
		AmountActual:       "",
		Status:             "submitted",
		ConfirmationStatus: "pending",
		Confirmations:      0,
		LastScannedBlock:   0,
		LastScannedSlot:    0,
		FailureReason:      "",
		UpdatedSource:      "payment",
	}); viewErr != nil {
		logx.WithContext(l.ctx).Errorf("upsert payment_view on submit failed, paymentNo=%s txId=%s err=%v", paymentNo, txID, viewErr)
	}

	confirmResp, err := NewConfirmPaymentLogic(l.ctx, l.svcCtx).ConfirmPayment(&payment.ConfirmPaymentReq{
		PaymentNo: paymentNo,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("confirm payment after submit failed, paymentNo=%s txId=%s err=%v", paymentNo, txID, err)
		failSnapshot, _ := json.Marshal(map[string]any{
			"paymentNo": paymentNo,
			"txId":      txID,
			"status":    "failed",
			"reason":    "confirm payment failed",
		})
		_ = l.svcCtx.TradeModel.UpsertIdempotencyRecord(l.ctx, trademodel.UpsertIdempotencyRecordInput{
			IdemKey:          idemKey,
			BizType:          "submit_tx",
			BizNo:            paymentNo,
			Status:           "failed",
			ResponseSnapshot: string(failSnapshot),
		})
		return nil, err
	}
	successSnapshot, _ := json.Marshal(map[string]any{
		"paymentNo":          paymentNo,
		"txId":               txID,
		"status":             confirmResp.Status,
		"confirmationStatus": confirmResp.ConfirmationStatus,
	})
	if idemErr := l.svcCtx.TradeModel.UpsertIdempotencyRecord(l.ctx, trademodel.UpsertIdempotencyRecordInput{
		IdemKey:          idemKey,
		BizType:          "submit_tx",
		BizNo:            paymentNo,
		Status:           "success",
		ResponseSnapshot: string(successSnapshot),
	}); idemErr != nil {
		logx.WithContext(l.ctx).Errorf("persist idempotency submit_tx success failed, paymentNo=%s txId=%s err=%v", paymentNo, txID, idemErr)
	}
	l.Infof("payment.SubmitTx confirm result paymentNo=%s status=%s confirmation=%s",
		paymentNo, confirmResp.Status, confirmResp.ConfirmationStatus)

	return &payment.SubmitTxResp{
		PaymentNo:          paymentNo,
		TxId:               txID,
		ConfirmationStatus: confirmResp.ConfirmationStatus,
		Status:             confirmResp.Status,
	}, nil
}
