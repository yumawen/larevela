package logic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"payment/internal/svc"
	"payment/payment"
	"tradechain"
	"trademodel"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	defaultDevnetUSDCMint = "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU"
	defaultUSDCReceiver   = "7xRr4GRzw5aw441Btum3Zxot6RUVBEUGihMdkwFb17zc"
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

	assetSymbol := strings.ToUpper(strings.TrimSpace(in.AssetSymbol))
	if assetSymbol == "" {
		assetSymbol = "SOL"
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

	var (
		receiverAccount      string
		receiverTokenAccount string
		assetAddress         string
		amountExpected       string
		decimals             int64
		payMode              = "transfer"
		instructions         []solana.Instruction
	)

	switch assetSymbol {
	case "SOL":
		receiverAccount = strings.TrimSpace(l.svcCtx.Config.TrialTransfer.ReceiverAccount)
		if receiverAccount == "" {
			return nil, fmt.Errorf("trialTransfer.receiverAccount is required")
		}
		receiverPubKey, parseErr := solana.PublicKeyFromBase58(receiverAccount)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid trialTransfer.receiverAccount")
		}
		amountExpected = strings.TrimSpace(l.svcCtx.Config.TrialTransfer.AmountSol)
		if amountExpected == "" {
			amountExpected = "0.1"
		}
		amountFloat, parseErr := strconv.ParseFloat(amountExpected, 64)
		if parseErr != nil || amountFloat <= 0 {
			return nil, fmt.Errorf("invalid trialTransfer.amountSol")
		}
		lamports := uint64(math.Round(amountFloat * float64(solana.LAMPORTS_PER_SOL)))
		if lamports == 0 {
			return nil, fmt.Errorf("amountSol too small")
		}
		instructions = append(instructions, system.NewTransferInstruction(lamports, payerPubKey, receiverPubKey).Build())
		assetAddress = ""
		decimals = 9
	case "USDC":
		receiverAccount = strings.TrimSpace(l.svcCtx.Config.TrialTransfer.UsdcReceiver)
		if receiverAccount == "" {
			receiverAccount = defaultUSDCReceiver
		}
		receiverPubKey, parseErr := solana.PublicKeyFromBase58(receiverAccount)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid trialTransfer.usdcReceiver")
		}
		assetAddress = strings.TrimSpace(in.AssetAddress)
		if assetAddress == "" {
			assetAddress = strings.TrimSpace(l.svcCtx.Config.TrialTransfer.UsdcMint)
		}
		if assetAddress == "" {
			assetAddress = defaultDevnetUSDCMint
		}
		usdcMint, parseErr := solana.PublicKeyFromBase58(assetAddress)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid usdc mint")
		}
		amountExpected = strings.TrimSpace(l.svcCtx.Config.TrialTransfer.AmountUsdc)
		if amountExpected == "" {
			amountExpected = "1"
		}
		baseUnits, parseErr := parseDecimalToBaseUnits(amountExpected, 6)
		if parseErr != nil || baseUnits == 0 {
			return nil, fmt.Errorf("invalid trialTransfer.amountUsdc")
		}
		decimals = 6
		sourceTokenAccount := strings.TrimSpace(in.PayerTokenAccount)
		var sourceTokenPubKey solana.PublicKey
		if sourceTokenAccount == "" {
			var ataErr error
			sourceTokenPubKey, _, ataErr = solana.FindAssociatedTokenAddress(payerPubKey, usdcMint)
			if ataErr != nil {
				return nil, fmt.Errorf("failed to derive payer token account")
			}
		} else {
			var ataErr error
			sourceTokenPubKey, ataErr = solana.PublicKeyFromBase58(sourceTokenAccount)
			if ataErr != nil {
				return nil, fmt.Errorf("invalid payerTokenAccount")
			}
		}
		receiverTokenPubKey, _, ataErr := solana.FindAssociatedTokenAddress(receiverPubKey, usdcMint)
		if ataErr != nil {
			return nil, fmt.Errorf("failed to derive receiver token account")
		}
		receiverTokenAccount = receiverTokenPubKey.String()
		createReceiverAtaIx, ataBuildErr := associatedtokenaccount.NewCreateIdempotentInstruction(
			payerPubKey,
			receiverPubKey,
			usdcMint,
		).ValidateAndBuild()
		if ataBuildErr != nil {
			return nil, fmt.Errorf("failed to build receiver token account instruction")
		}
		transferIx, transferBuildErr := token.NewTransferCheckedInstruction(
			baseUnits,
			6,
			sourceTokenPubKey,
			usdcMint,
			receiverTokenPubKey,
			payerPubKey,
			nil,
		).ValidateAndBuild()
		if transferBuildErr != nil {
			return nil, fmt.Errorf("failed to build usdc transfer instruction")
		}
		instructions = append(instructions, createReceiverAtaIx, transferIx)
	default:
		return nil, fmt.Errorf("unsupported assetSymbol: %s", assetSymbol)
	}

	tx, err := solana.NewTransaction(
		instructions,
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
		AssetSymbol:       assetSymbol,
		AssetAddress:      assetAddress,
		AmountExpected:    amountExpected,
		Decimals:          decimals,
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
		AssetSymbol:        assetSymbol,
		AmountExpected:     amountExpected,
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
		PaymentNo:            paymentNo,
		OrderNo:              orderNo,
		ChainType:            chainType,
		Network:              network,
		ChainId:              chainID,
		PayMode:              payMode,
		ReceiverAccount:      receiverAccount,
		ReceiverTokenAccount: receiverTokenAccount,
		AssetAddress:         assetAddress,
		AssetSymbol:          assetSymbol,
		AmountExpected:       amountExpected,
		Decimals:             decimals,
		ReferenceId:          strings.TrimSpace(in.ReferenceId),
		SerializedMessage:    serializedMessage,
		QuoteExpiredAt:       quoteExpiredAt.Unix(),
		Status:               "created",
	}, nil
}

func parseDecimalToBaseUnits(amount string, decimals int) (uint64, error) {
	v := strings.TrimSpace(amount)
	if v == "" {
		return 0, fmt.Errorf("empty amount")
	}
	v = strings.TrimPrefix(v, "+")
	if strings.HasPrefix(v, "-") {
		return 0, fmt.Errorf("negative amount")
	}
	parts := strings.Split(v, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid decimal amount")
	}
	intPart := parts[0]
	if intPart == "" {
		intPart = "0"
	}
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}
	if len(fracPart) > decimals {
		return 0, fmt.Errorf("too many decimal places")
	}
	pow10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	intValue, ok := new(big.Int).SetString(intPart, 10)
	if !ok {
		return 0, fmt.Errorf("invalid integer part")
	}
	intValue.Mul(intValue, pow10)
	if fracPart != "" {
		paddedFrac := fracPart + strings.Repeat("0", decimals-len(fracPart))
		fracValue, fracOK := new(big.Int).SetString(paddedFrac, 10)
		if !fracOK {
			return 0, fmt.Errorf("invalid fraction part")
		}
		intValue.Add(intValue, fracValue)
	}
	if intValue.Sign() <= 0 {
		return 0, fmt.Errorf("amount must be greater than zero")
	}
	if !intValue.IsUint64() {
		return 0, fmt.Errorf("amount exceeds uint64")
	}
	return intValue.Uint64(), nil
}
