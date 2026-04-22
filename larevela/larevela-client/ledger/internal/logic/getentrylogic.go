package logic

import (
	"context"
	"fmt"
	"strings"

	"ledger/internal/svc"
	"ledger/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEntryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEntryLogic {
	return &GetEntryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEntryLogic) GetEntry(in *ledger.GetEntryReq) (*ledger.GetEntryResp, error) {
	if l.svcCtx.TradeModel == nil {
		return nil, fmt.Errorf("database is not configured")
	}
	entryNo := strings.TrimSpace(in.EntryNo)
	if entryNo == "" {
		return nil, fmt.Errorf("entryNo is required")
	}
	l.Infof("ledger.GetEntry start entryNo=%s", entryNo)

	entry, err := l.svcCtx.TradeModel.GetLedgerEntry(l.ctx, entryNo)
	if err != nil {
		l.Errorf("ledger.GetEntry failed entryNo=%s err=%v", entryNo, err)
		return nil, err
	}
	l.Infof("ledger.GetEntry success entryNo=%s status=%s", entry.EntryNo, entry.Status)

	return &ledger.GetEntryResp{
		EntryNo:      entry.EntryNo,
		PaymentNo:    entry.PaymentNo,
		OrderNo:      entry.OrderNo,
		UserId:       entry.UserID,
		ChainType:    entry.ChainType,
		Network:      entry.Network,
		ChainId:      entry.ChainID,
		EntryType:    entry.EntryType,
		AssetSymbol:  entry.AssetSymbol,
		AssetAddress: entry.AssetAddress,
		Amount:       entry.Amount,
		Direction:    entry.Direction,
		Status:       entry.Status,
		Remark:       entry.Remark,
	}, nil
}
