package config

import (
	"mysqlconf"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql mysqlconf.Conf `json:"mysql"`
}
