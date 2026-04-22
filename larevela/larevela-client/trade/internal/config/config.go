// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	PaymentRpc zrpc.RpcClientConf `json:"paymentRpc"` // 调用进入payment
	OrderRpc   zrpc.RpcClientConf `json:"orderRpc"`   // 先创建订单
}
