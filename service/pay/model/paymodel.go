package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PayModel = (*customPayModel)(nil)

var (
	cachePayOidPrefix = "cache:pay:oid:"
)

type (
	// PayModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPayModel.
	PayModel interface {
		payModel

		FindOneByOid(ctx context.Context, oid int64) (*Pay, error)
	}

	customPayModel struct {
		*defaultPayModel
	}
)

// NewPayModel returns a model for the database table.
func NewPayModel(conn sqlx.SqlConn, c cache.CacheConf) PayModel {
	return &customPayModel{
		defaultPayModel: newPayModel(conn, c),
	}
}

func (m *defaultPayModel) FindOneByOid(ctx context.Context, oid int64) (*Pay, error) {
	payOidKey := fmt.Sprintf("%s%v", cachePayOidPrefix, oid)
	var resp Pay
	err := m.QueryRowCtx(ctx, &resp, payOidKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `oid` = ? limit 1", payRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, oid)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
