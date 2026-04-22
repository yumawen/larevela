package logic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"payment/internal/svc"
	"payment/payment"
	"tradechain"
	"trademodel"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentIntentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePaymentIntentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentIntentLogic {
	return &CreatePaymentIntentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePaymentIntentLogic) CreatePaymentIntent(in *payment.CreatePaymentIntentReq) (*payment.CreatePaymentIntentResp, error) {
	if l.svcCtx.TradeModel == nil {
		return nil, fmt.Errorf("database is not configured")
	}

	orderNo := strings.TrimSpace(in.OrderNo)
	if orderNo == "" {
		orderNo = fmt.Sprintf("trial-order-%d", time.Now().UnixNano())
	}
	paymentNo := fmt.Sprintf("pay-%d", time.Now().UnixNano())

	payerAccount := strings.TrimSpace(in.PayerAccount)
	if payerAccount == "" {
		return nil, fmt.Errorf("payerAccount is required")
	}
	l.Infof("payment.CreatePaymentIntent start orderNo=%s payer=%s referenceId=%s",
		orderNo, payerAccount, strings.TrimSpace(in.ReferenceId))
	payerPubKey, err := solana.PublicKeyFromBase58(payerAccount)
	if err != nil {
		return nil, fmt.Errorf("invalid payerAccount")
	}

	receiverAccount := strings.TrimSpace(l.svcCtx.Config.TrialTransfer.ReceiverAccount)
	if receiverAccount == "" {
		return nil, fmt.Errorf("trialTransfer.receiverAccount is required")
	}
	receiverPubKey, err := solana.PublicKeyFromBase58(receiverAccount)
	if err != nil {
		return nil, fmt.Errorf("invalid trialTransfer.receiverAccount")
	}

	amountSol := strings.TrimSpace(l.svcCtx.Config.TrialTransfer.AmountSol)
	if amountSol == "" {
		amountSol = "0.1"
	}
	amountFloat, err := strconv.ParseFloat(amountSol, 64)
	if err != nil || amountFloat <= 0 {
		return nil, fmt.Errorf("invalid trialTransfer.amountSol")
	}
	lamports := uint64(math.Round(amountFloat * float64(solana.LAMPORTS_PER_SOL)))
	if lamports == 0 {
		return nil, fmt.Errorf("amountSol too small")
	}

	network := strings.TrimSpace(in.Network)
	if network == "" {
		network = "devnet"
	}
	chainType := strings.TrimSpace(in.ChainType)
	if chainType == "" {
		chainType = "solana"
	}
	chainID := in.ChainId
	if chainID == 0 {
		chainID = 103
	}

	rpcClient := rpc.New(tradechain.DefaultDevnetSolanaRPCURL)
	latestBlockhash, err := rpcClient.GetLatestBlockhash(l.ctx, rpc.CommitmentFinalized)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("get latest blockhash failed: %v", err)
		return nil, fmt.Errorf("failed to build payment transaction")
	}
	l.Infof("payment.CreatePaymentIntent got blockhash orderNo=%s", orderNo)

	transferIx := system.NewTransferInstruction(lamports, payerPubKey, receiverPubKey).Build()
	tx, err := solana.NewTransaction(
		[]solana.Instruction{transferIx},
		latestBlockhash.Value.Blockhash,
		solana.TransactionPayer(payerPubKey),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create unsigned transaction")
	}

	serializedMessageBytes, err := tx.Message.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize unsigned transaction message")
	}
	serializedMessage := base64.StdEncoding.EncodeToString(serializedMessageBytes)
	quoteExpiredAt := time.Now().Add(10 * time.Minute)

	err = l.svcCtx.TradeModel.CreatePaymentIntent(l.ctx, trademodel.CreatePaymentIntentInput{
		PaymentNo:         paymentNo,
		OrderNo:           orderNo,
		ChainType:         chainType,
		Network:           network,
		ChainID:           chainID,
		PayerAccount:      payerAccount,
		ReceiverAccount:   receiverAccount,
		AssetSymbol:       "SOL",
		AssetAddress:      "",
		AmountExpected:    amountSol,
		Decimals:          9,
		ReferenceID:       strings.TrimSpace(in.ReferenceId),
		SerializedMessage: serializedMessage,
		QuoteExpiredAt:    quoteExpiredAt,
		Status:            "created",
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("create payment intent persist failed, paymentNo=%s err=%v", paymentNo, err)
		return nil, fmt.Errorf("failed to persist payment intent")
	}
	if viewErr := l.svcCtx.TradeModel.UpsertPaymentView(l.ctx, trademodel.UpsertPaymentViewInput{
		PaymentNo:          paymentNo,
		OrderNo:            orderNo,
		TxID:               "",
		ChainType:          chainType,
		Network:            network,
		ChainID:            chainID,
		PayerAccount:       payerAccount,
		ReceiverAccount:    receiverAccount,
		AssetSymbol:        "SOL",
		AmountExpected:     amountSol,
		AmountActual:       "",
		Status:             "created",
		ConfirmationStatus: "pending",
		Confirmations:      0,
		LastScannedBlock:   0,
		LastScannedSlot:    0,
		FailureReason:      "",
		UpdatedSource:      "payment",
	}); viewErr != nil {
		logx.WithContext(l.ctx).Errorf("upsert payment_view on create intent failed, paymentNo=%s err=%v", paymentNo, viewErr)
	}
	idemSnapshot, _ := json.Marshal(map[string]any{
		"paymentNo": paymentNo,
		"orderNo":   orderNo,
		"status":    "created",
	})
	if idemErr := l.svcCtx.TradeModel.UpsertIdempotencyRecord(l.ctx, trademodel.UpsertIdempotencyRecordInput{
		IdemKey:          fmt.Sprintf("create_intent:%s", paymentNo),
		BizType:          "create_intent",
		BizNo:            paymentNo,
		Status:           "success",
		ResponseSnapshot: string(idemSnapshot),
	}); idemErr != nil {
		logx.WithContext(l.ctx).Errorf("persist idempotency create_intent failed, paymentNo=%s err=%v", paymentNo, idemErr)
	}
	l.Infof("payment.CreatePaymentIntent success paymentNo=%s orderNo=%s status=created", paymentNo, orderNo)

	return &payment.CreatePaymentIntentResp{
		PaymentNo:         paymentNo,
		OrderNo:           orderNo,
		ChainType:         chainType,
		Network:           network,
		ChainId:           chainID,
		PayMode:           "transfer",
		ReceiverAccount:   receiverAccount,
		AssetAddress:      "",
		AssetSymbol:       "SOL",
		AmountExpected:    amountSol,
		Decimals:          9,
		ReferenceId:       strings.TrimSpace(in.ReferenceId),
		SerializedMessage: serializedMessage,
		QuoteExpiredAt:    quoteExpiredAt.Unix(),
		Status:            "created",
	}, nil
}
