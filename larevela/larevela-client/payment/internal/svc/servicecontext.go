package svc

import (
	"chain/chainclient"
	"ledger/ledgerclient"
	"order/orderclient"
	"payment/internal/config"
	"trademodel"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	TradeModel *trademodel.Model
	OrderRpc   orderclient.Order
	ChainRpc   chainclient.Chain
	LedgerRpc  ledgerclient.Ledger
}

func NewServiceContext(c config.Config) *ServiceContext {
	var tradeModel *trademodel.Model
	dsn := c.Mysql.DSN()
	if dsn != "" {
		conn := sqlx.NewMysql(dsn)
		tradeModel = trademodel.NewModel(conn)
	}

	return &ServiceContext{
		Config:     c,
		TradeModel: tradeModel,
		OrderRpc:   orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
		ChainRpc:   chainclient.NewChain(zrpc.MustNewClient(c.ChainRpc)),
		LedgerRpc:  ledgerclient.NewLedger(zrpc.MustNewClient(c.LedgerRpc)),
	}
}
