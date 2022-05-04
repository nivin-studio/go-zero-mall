package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderModel = (*customOrderModel)(nil)

type (
	// OrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderModel.
	OrderModel interface {
		orderModel

		FindAllByUid(ctx context.Context, uid int64) ([]*Order, error)
		FindOneByUid(ctx context.Context, uid int64) (*Order, error)
		TxInsert(ctx context.Context, tx *sql.Tx, data *Order) (sql.Result, error)
		TxUpdate(ctx context.Context, tx *sql.Tx, data *Order) error
	}

	customOrderModel struct {
		*defaultOrderModel
	}
)

// NewOrderModel returns a model for the database table.
func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf) OrderModel {
	return &customOrderModel{
		defaultOrderModel: newOrderModel(conn, c),
	}
}

func (m *customOrderModel) FindAllByUid(ctx context.Context, uid int64) ([]*Order, error) {
	var resp []*Order

	query := fmt.Sprintf("select %s from %s where `uid` = ?", orderRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, uid)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultOrderModel) FindOneByUid(ctx context.Context, uid int64) (*Order, error) {
	var resp Order

	query := fmt.Sprintf("select %s from %s where `uid` = ? order by create_time desc limit 1", orderRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, uid)

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultOrderModel) TxInsert(ctx context.Context, tx *sql.Tx, data *Order) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, orderRowsExpectAutoSet)
	ret, err := tx.ExecContext(ctx, query, data.Uid, data.Pid, data.Amount, data.Status)

	return ret, err
}

func (m *defaultOrderModel) TxUpdate(ctx context.Context, tx *sql.Tx, data *Order) error {
	productIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, orderRowsWithPlaceHolder)
		return tx.ExecContext(ctx, query, data.Uid, data.Pid, data.Amount, data.Status, data.Id)
	}, productIdKey)
	return err
}
