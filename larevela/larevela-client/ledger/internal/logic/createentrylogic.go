package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ledger/internal/svc"
	"ledger/ledger"
	"trademodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEntryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEntryLogic {
	return &CreateEntryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateEntryLogic) CreateEntry(in *ledger.CreateEntryReq) (*ledger.CreateEntryResp, error) {
	if l.svcCtx.TradeModel == nil {
		return nil, fmt.Errorf("database is not configured")
	}

	entryNo := strings.TrimSpace(in.EntryNo)
	if entryNo == "" {
		entryNo = fmt.Sprintf("entry-%d", time.Now().UnixNano())
	}
	if strings.TrimSpace(in.PaymentNo) == "" || strings.TrimSpace(in.OrderNo) == "" {
		return nil, fmt.Errorf("paymentNo and orderNo are required")
	}
	l.Infof("ledger.CreateEntry start entryNo=%s paymentNo=%s orderNo=%s amount=%s",
		entryNo, strings.TrimSpace(in.PaymentNo), strings.TrimSpace(in.OrderNo), strings.TrimSpace(in.Amount))

	err := l.svcCtx.TradeModel.CreateLedgerEntry(l.ctx, trademodel.CreateLedgerEntryInput{
		EntryNo:      entryNo,
		PaymentNo:    strings.TrimSpace(in.PaymentNo),
		OrderNo:      strings.TrimSpace(in.OrderNo),
		UserID:       in.UserId,
		ChainType:    strings.TrimSpace(in.ChainType),
		Network:      strings.TrimSpace(in.Network),
		ChainID:      in.ChainId,
		EntryType:    strings.TrimSpace(in.EntryType),
		AssetSymbol:  strings.TrimSpace(in.AssetSymbol),
		AssetAddress: strings.TrimSpace(in.AssetAddress),
		Amount:       strings.TrimSpace(in.Amount),
		Direction:    strings.TrimSpace(in.Direction),
		Remark:       strings.TrimSpace(in.Remark),
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("create ledger entry failed, entryNo=%s err=%v", entryNo, err)
		return nil, fmt.Errorf("create ledger entry failed")
	}
	l.Infof("ledger.CreateEntry success entryNo=%s status=posted", entryNo)

	return &ledger.CreateEntryResp{
		EntryNo: entryNo,
		Status:  "posted",
	}, nil
}
