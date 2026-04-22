package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"chain/chainclient"
	"ledger/ledgerclient"
	"order/orderclient"
	"payment/internal/svc"
	"payment/payment"
	"trademodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmPaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmPaymentLogic {
	return &ConfirmPaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfirmPaymentLogic) ConfirmPayment(in *payment.ConfirmPaymentReq) (*payment.ConfirmPaymentResp, error) {
	if l.svcCtx.TradeModel == nil {
		return nil, fmt.Errorf("database is not configured")
	}
	paymentNo := strings.TrimSpace(in.PaymentNo)
	if paymentNo == "" {
		return nil, fmt.Errorf("paymentNo is required")
	}
	l.Infof("payment.ConfirmPayment start paymentNo=%s", paymentNo)

	pi, err := l.svcCtx.TradeModel.GetPaymentIntent(l.ctx, paymentNo)
	if err != nil {
		l.Errorf("payment.ConfirmPayment load payment intent failed paymentNo=%s err=%v", paymentNo, err)
		return nil, err
	}
	if strings.TrimSpace(pi.TxID) == "" {
		if viewErr := l.svcCtx.TradeModel.UpsertPaymentView(l.ctx, trademodel.UpsertPaymentViewInput{
			PaymentNo:          paymentNo,
			OrderNo:            pi.OrderNo,
			TxID:               "",
			ChainType:          pi.ChainType,
			Network:            pi.Network,
			ChainID:            pi.ChainID,
			PayerAccount:       pi.PayerAccount,
			ReceiverAccount:    pi.ReceiverAccount,
			AssetSymbol:        pi.AssetSymbol,
			AmountExpected:     pi.AmountExpected,
			AmountActual:       "",
			Status:             "submitted",
			ConfirmationStatus: "pending",
			Confirmations:      0,
			LastScannedBlock:   0,
			LastScannedSlot:    0,
			FailureReason:      "",
			UpdatedSource:      "confirm",
		}); viewErr != nil {
			logx.WithContext(l.ctx).Errorf("upsert payment_view no txId failed, paymentNo=%s err=%v", paymentNo, viewErr)
		}
		l.Infof("payment.ConfirmPayment pending no txId paymentNo=%s", paymentNo)
		return &payment.ConfirmPaymentResp{
			PaymentNo:          paymentNo,
			TxId:               "",
			Confirmations:      0,
			ConfirmationStatus: "pending",
			AmountActual:       "",
			Status:             "submitted",
			FailureReason:      "",
		}, nil
	}

	parseResp, err := l.svcCtx.ChainRpc.ParsePayment(l.ctx, &chainclient.ParsePaymentReq{
		ChainType:            pi.ChainType,
		Network:              pi.Network,
		ChainId:              pi.ChainID,
		TxId:                 pi.TxID,
		ReceiverAccount:      pi.ReceiverAccount,
		ReceiverTokenAccount: pi.ReceiverTokenAccount,
		AssetAddress:         pi.AssetAddress,
		ExpectedAmount:       pi.AmountExpected,
		ReferenceId:          "",
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("parse payment rpc failed, paymentNo=%s err=%v", paymentNo, err)
		return nil, fmt.Errorf("confirm payment failed")
	}
	if cursorErr := l.svcCtx.TradeModel.UpsertChainScanCursor(l.ctx, trademodel.UpsertChainScanCursorInput{
		ChainType:        pi.ChainType,
		Network:          pi.Network,
		ChainID:          pi.ChainID,
		CursorType:       "payment_confirm",
		LastScannedBlock: parseResp.BlockNumber,
		LastScannedSlot:  parseResp.Slot,
		LastScannedTxID:  pi.TxID,
		Remark:           parseResp.ConfirmationStatus,
	}); cursorErr != nil {
		logx.WithContext(l.ctx).Errorf("upsert chain scan cursor failed, paymentNo=%s txId=%s err=%v", paymentNo, pi.TxID, cursorErr)
	}
	l.Infof("payment.ConfirmPayment chain parsed paymentNo=%s matched=%v finalized=%v confirmations=%d",
		paymentNo, parseResp.Matched, parseResp.Finalized, parseResp.Confirmations)

	status := "confirming"
	failureReason := ""
	if parseResp.Matched {
		status = "paid"
	} else if parseResp.Finalized {
		status = "failed"
		failureReason = parseResp.MatchReason
	}
	amountActual := parseResp.AmountActual
	if strings.TrimSpace(amountActual) == "" {
		amountActual = pi.AmountExpected
	}

	if err = l.svcCtx.TradeModel.UpdatePaymentConfirmation(
		l.ctx,
		paymentNo,
		parseResp.Confirmations,
		parseResp.ConfirmationStatus,
		amountActual,
		status,
		failureReason,
	); err != nil {
		l.Errorf("payment.ConfirmPayment update status failed paymentNo=%s err=%v", paymentNo, err)
		return nil, err
	}
	l.Infof("payment.ConfirmPayment updated paymentNo=%s status=%s failureReason=%s", paymentNo, status, failureReason)
	if viewErr := l.svcCtx.TradeModel.UpsertPaymentView(l.ctx, trademodel.UpsertPaymentViewInput{
		PaymentNo:          paymentNo,
		OrderNo:            pi.OrderNo,
		TxID:               pi.TxID,
		ChainType:          pi.ChainType,
		Network:            pi.Network,
		ChainID:            pi.ChainID,
		PayerAccount:       pi.PayerAccount,
		ReceiverAccount:    pi.ReceiverAccount,
		AssetSymbol:        pi.AssetSymbol,
		AmountExpected:     pi.AmountExpected,
		AmountActual:       amountActual,
		Status:             status,
		ConfirmationStatus: parseResp.ConfirmationStatus,
		Confirmations:      parseResp.Confirmations,
		LastScannedBlock:   parseResp.BlockNumber,
		LastScannedSlot:    parseResp.Slot,
		FailureReason:      failureReason,
		UpdatedSource:      "confirm",
	}); viewErr != nil {
		logx.WithContext(l.ctx).Errorf("upsert payment_view confirm failed, paymentNo=%s txId=%s err=%v", paymentNo, pi.TxID, viewErr)
	}

	if status == "paid" {
		entryNo := fmt.Sprintf("entry-%d", time.Now().UnixNano())
		_, ledgerErr := l.svcCtx.LedgerRpc.CreateEntry(l.ctx, &ledgerclient.CreateEntryReq{
			EntryNo:      entryNo,
			PaymentNo:    paymentNo,
			OrderNo:      pi.OrderNo,
			UserId:       0,
			ChainType:    pi.ChainType,
			Network:      pi.Network,
			ChainId:      pi.ChainID,
			EntryType:    "payment",
			AssetSymbol:  pi.AssetSymbol,
			AssetAddress: pi.AssetAddress,
			Amount:       amountActual,
			Direction:    "credit",
			Remark:       "auto-post from payment confirm",
		})
		if ledgerErr != nil {
			logx.WithContext(l.ctx).Errorf("create ledger entry failed, paymentNo=%s err=%v", paymentNo, ledgerErr)
		} else {
			l.Infof("payment.ConfirmPayment ledger posted paymentNo=%s entryNo=%s", paymentNo, entryNo)
		}

		_, orderErr := l.svcCtx.OrderRpc.MarkOrderPaid(l.ctx, &orderclient.MarkOrderPaidReq{
			OrderNo:   pi.OrderNo,
			PaymentNo: paymentNo,
			PaidAt:    time.Now().UTC().Unix(),
		})
		if orderErr != nil {
			logx.WithContext(l.ctx).Errorf("mark order paid failed, paymentNo=%s orderNo=%s err=%v", paymentNo, pi.OrderNo, orderErr)
		} else {
			l.Infof("payment.ConfirmPayment order marked paid paymentNo=%s orderNo=%s", paymentNo, pi.OrderNo)
		}
	}

	return &payment.ConfirmPaymentResp{
		PaymentNo:          paymentNo,
		TxId:               pi.TxID,
		Confirmations:      parseResp.Confirmations,
		ConfirmationStatus: parseResp.ConfirmationStatus,
		AmountActual:       amountActual,
		Status:             status,
		FailureReason:      failureReason,
	}, nil
}
