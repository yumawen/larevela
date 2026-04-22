// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"order/orderclient"
	"payment/paymentclient"
	"trade/internal/config"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	PaymentRpc paymentclient.Payment
	OrderRpc   orderclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		PaymentRpc: paymentclient.NewPayment(zrpc.MustNewClient(c.PaymentRpc)),
		OrderRpc:   orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
	}
}
