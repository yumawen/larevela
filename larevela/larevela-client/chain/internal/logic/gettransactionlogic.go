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

type GetTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionLogic {
	return &GetTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTransactionLogic) GetTransaction(in *chain.GetTransactionReq) (*chain.GetTransactionResp, error) {
	txID := strings.TrimSpace(in.TxId)
	if txID == "" {
		return &chain.GetTransactionResp{Found: false}, nil
	}
	l.Infof("chain.GetTransaction start txId=%s", txID)
	sig, err := solana.SignatureFromBase58(txID)
	if err != nil {
		return &chain.GetTransactionResp{TxId: txID, Found: false}, nil
	}

	rpcURL := strings.TrimSpace(l.svcCtx.Config.SolanaRpcURL)
	if rpcURL == "" {
		rpcURL = tradechain.DefaultDevnetSolanaRPCURL
	}
	client := rpc.New(rpcURL)
	statuses, err := client.GetSignatureStatuses(l.ctx, false, sig)
	if err != nil || statuses == nil || statuses.Value == nil || len(statuses.Value) == 0 || statuses.Value[0] == nil {
		l.Infof("chain.GetTransaction not found txId=%s rpc=%s err=%v", txID, rpcURL, err)
		return &chain.GetTransactionResp{TxId: txID, Found: false}, nil
	}
	status := statuses.Value[0]
	slot := int64(status.Slot)

	confirmationStatus := string(status.ConfirmationStatus)
	if confirmationStatus == "" {
		confirmationStatus = "processed"
	}
	l.Infof("chain.GetTransaction success txId=%s slot=%d confirmation=%s", txID, slot, confirmationStatus)
	return &chain.GetTransactionResp{
		TxId:               txID,
		BlockNumber:        slot,
		Slot:               slot,
		ConfirmationStatus: confirmationStatus,
		Finalized:          status.ConfirmationStatus == rpc.ConfirmationStatusFinalized,
		Found:              true,
	}, nil
}
