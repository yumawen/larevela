package config

import (
	"mysqlconf"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql         mysqlconf.Conf     `json:"mysql"`
	TrialTransfer TrialTransferConf  `json:"trialTransfer"`
	OrderRpc      zrpc.RpcClientConf `json:"orderRpc"`
	ChainRpc      zrpc.RpcClientConf `json:"chainRpc"`
	LedgerRpc     zrpc.RpcClientConf `json:"ledgerRpc"`
}

type TrialTransferConf struct {
	ReceiverAccount string `json:"receiverAccount"`
	AmountSol       string `json:"amountSol"`
	UsdcMint        string `json:"usdcMint"`
	AmountUsdc      string `json:"amountUsdc"`
	UsdcReceiver    string `json:"usdcReceiver"`
}
