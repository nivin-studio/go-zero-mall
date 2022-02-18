package logic

import (
	"context"

	"mall/service/pay/api/internal/svc"
	"mall/service/pay/api/internal/types"
	"mall/service/pay/rpc/pay"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) CallbackLogic {
	return CallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CallbackLogic) Callback(req types.CallbackRequest) (resp *types.CallbackResponse, err error) {
	_, err = l.svcCtx.PayRpc.Callback(l.ctx, &pay.CallbackRequest{
		Id:     req.Id,
		Uid:    req.Uid,
		Oid:    req.Oid,
		Amount: req.Amount,
		Source: req.Source,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &types.CallbackResponse{}, nil
}
