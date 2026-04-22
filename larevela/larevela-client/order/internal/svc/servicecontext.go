package svc

import (
	"order/internal/config"
	"trademodel"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	TradeModel *trademodel.Model
}

func NewServiceContext(c config.Config) *ServiceContext {
	var tradeModel *trademodel.Model
	if dsn := c.Mysql.DSN(); dsn != "" {
		tradeModel = trademodel.NewModel(sqlx.NewMysql(dsn))
	}

	return &ServiceContext{
		Config:     c,
		TradeModel: tradeModel,
	}
}
