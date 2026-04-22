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

type GetReceiptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReceiptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReceiptLogic {
	return &GetReceiptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetReceiptLogic) GetReceipt(in *chain.GetReceiptReq) (*chain.GetReceiptResp, error) {
	txID := strings.TrimSpace(in.TxId)
	if txID == "" {
		return &chain.GetReceiptResp{Found: false}, nil
	}
	l.Infof("chain.GetReceipt start txId=%s", txID)
	sig, err := solana.SignatureFromBase58(txID)
	if err != nil {
		return &chain.GetReceiptResp{TxId: txID, Found: false}, nil
	}

	rpcURL := strings.TrimSpace(l.svcCtx.Config.SolanaRpcURL)
	if rpcURL == "" {
		rpcURL = tradechain.DefaultDevnetSolanaRPCURL
	}
	client := rpc.New(rpcURL)
	statuses, err := client.GetSignatureStatuses(l.ctx, false, sig)
	if err != nil || statuses == nil || statuses.Value == nil || len(statuses.Value) == 0 || statuses.Value[0] == nil {
		l.Infof("chain.GetReceipt not found txId=%s rpc=%s err=%v", txID, rpcURL, err)
		return &chain.GetReceiptResp{TxId: txID, Found: false}, nil
	}
	status := statuses.Value[0]
	slot := int64(status.Slot)
	confirmationStatus := string(status.ConfirmationStatus)
	if confirmationStatus == "" {
		confirmationStatus = "processed"
	}
	txStatus := int64(1)
	if status.Err != nil {
		txStatus = 0
	}
	l.Infof("chain.GetReceipt success txId=%s status=%d confirmation=%s", txID, txStatus, confirmationStatus)

	return &chain.GetReceiptResp{
		TxId:               txID,
		Status:             txStatus,
		BlockNumber:        slot,
		Slot:               slot,
		FeeAmount:          "0",
		ConfirmationStatus: confirmationStatus,
		Finalized:          status.ConfirmationStatus == rpc.ConfirmationStatusFinalized,
		Found:              true,
	}, nil
}
