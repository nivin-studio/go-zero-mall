package svc

import (
	"mall/service/order/api/internal/config"
	"mall/service/order/rpc/order"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc   order.Order
	ProductRpc product.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		OrderRpc:   order.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
		ProductRpc: product.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
	}
}
