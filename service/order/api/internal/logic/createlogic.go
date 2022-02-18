package logic

import (
	"context"

	"mall/service/order/api/internal/svc"
	"mall/service/order/api/internal/types"
	"mall/service/order/rpc/order"
	"mall/service/product/rpc/product"

	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateLogic {
	return CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req types.CreateRequest) (resp *types.CreateResponse, err error) {
	// 获取 OrderRpc BuildTarget
	orderRpcBusiServer, err := l.svcCtx.Config.OrderRpc.BuildTarget()
	if err != nil {
		return nil, status.Error(100, "订单创建异常")
	}

	// 获取 ProductRpc BuildTarget
	productRpcBusiServer, err := l.svcCtx.Config.ProductRpc.BuildTarget()
	if err != nil {
		return nil, status.Error(100, "订单创建异常")
	}

	// dtm 服务的 etcd 注册地址
	var dtmServer = "etcd://etcd:2379/dtmservice"
	// 创建一个gid
	gid := dtmgrpc.MustGenGid(dtmServer)
	// 创建一个saga协议的事务
	saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
		Add(orderRpcBusiServer+"/orderclient.Order/Create", orderRpcBusiServer+"/orderclient.Order/CreateRevert", &order.CreateRequest{
			Uid:    req.Uid,
			Pid:    req.Pid,
			Amount: req.Amount,
			Status: 0,
		}).
		Add(productRpcBusiServer+"/productclient.Product/DecrStock", productRpcBusiServer+"/productclient.Product/DecrStockRevert", &product.DecrStockRequest{
			Id:  req.Pid,
			Num: 1,
		})

	// 事务提交
	err = saga.Submit()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	return &types.CreateResponse{}, nil
}
