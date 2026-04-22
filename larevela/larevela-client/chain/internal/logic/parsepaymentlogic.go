package logic

import (
	"context"
	"strings"

	"chain/chain"
	"chain/internal/svc"
	"tradechain"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ParsePaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewParsePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParsePaymentLogic {
	return &ParsePaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ParsePaymentLogic) ParsePayment(in *chain.ParsePaymentReq) (*chain.ParsePaymentResp, error) {
	txID := strings.TrimSpace(in.TxId)
	if txID == "" {
		return &chain.ParsePaymentResp{
			Matched:     false,
			MatchReason: "empty txId",
		}, nil
	}
	l.Infof("chain.ParsePayment start txId=%s receiver=%s expectedAmount=%s",
		txID, strings.TrimSpace(in.ReceiverAccount), strings.TrimSpace(in.ExpectedAmount))
	sig, err := solana.SignatureFromBase58(txID)
	if err != nil {
		return &chain.ParsePaymentResp{
			Matched:     false,
			MatchReason: "invalid txId",
		}, nil
	}

	rpcURL := strings.TrimSpace(l.svcCtx.Config.SolanaRpcURL)
	if rpcURL == "" {
		rpcURL = tradechain.DefaultDevnetSolanaRPCURL
	}
	client := rpc.New(rpcURL)
	statuses, err := client.GetSignatureStatuses(l.ctx, false, sig)
	if err != nil || statuses == nil || statuses.Value == nil || len(statuses.Value) == 0 || statuses.Value[0] == nil {
		l.Infof("chain.ParsePayment tx not found txId=%s rpc=%s err=%v", txID, rpcURL, err)
		return &chain.ParsePaymentResp{
			Matched:     false,
			MatchReason: "tx not found",
		}, nil
	}

	status := statuses.Value[0]
	slot := int64(status.Slot)
	var confirmations int64
	if status.Confirmations != nil {
		confirmations = int64(*status.Confirmations)
	}
	if status.Err != nil {
		l.Infof("chain.ParsePayment tx failed on chain txId=%s confirmation=%s", txID, string(status.ConfirmationStatus))
		return &chain.ParsePaymentResp{
			Matched:            false,
			Confirmations:      confirmations,
			ConfirmationStatus: string(status.ConfirmationStatus),
			Finalized:          status.ConfirmationStatus == rpc.ConfirmationStatusFinalized,
			MatchReason:        "tx failed on chain",
		}, nil
	}

	confirmationStatus := string(status.ConfirmationStatus)
	if confirmationStatus == "" {
		confirmationStatus = "processed"
	}
	matchReason := "tx confirmed"
	matched := strings.TrimSpace(in.ExpectedAmount) != "" && strings.TrimSpace(in.ReceiverAccount) != ""
	if !matched {
		matchReason = "receiverAccount/expectedAmount missing"
	}
	l.Infof("chain.ParsePayment success txId=%s matched=%v finalized=%v confirmations=%d",
		txID, matched, status.ConfirmationStatus == rpc.ConfirmationStatusFinalized, confirmations)

	return &chain.ParsePaymentResp{
		Matched:            matched,
		ToAccount:          strings.TrimSpace(in.ReceiverAccount),
		AssetAddress:       strings.TrimSpace(in.AssetAddress),
		AmountActual:       strings.TrimSpace(in.ExpectedAmount),
		BlockNumber:        slot,
		Slot:               slot,
		Confirmations:      confirmations,
		ConfirmationStatus: confirmationStatus,
		Finalized:          status.ConfirmationStatus == rpc.ConfirmationStatusFinalized,
		MatchReason:        matchReason,
	}, nil
}
