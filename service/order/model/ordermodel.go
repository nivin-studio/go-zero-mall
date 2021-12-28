package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/builder"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
)

var (
	orderFieldNames          = builder.RawFieldNames(&Order{})
	orderRows                = strings.Join(orderFieldNames, ",")
	orderRowsExpectAutoSet   = strings.Join(stringx.Remove(orderFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	orderRowsWithPlaceHolder = strings.Join(stringx.Remove(orderFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheOrderIdPrefix = "cache:order:id:"
)

type (
	OrderModel interface {
		Insert(data *Order) (sql.Result, error)
		FindOne(id int64) (*Order, error)
		FindAllByUid(uid int64) ([]*Order, error)
		Update(data *Order) error
		Delete(id int64) error
	}

	defaultOrderModel struct {
		sqlc.CachedConn
		table string
	}

	Order struct {
		Id         int64     `db:"id"`
		Uid        int64     `db:"uid"`    // 用户ID
		Pid        int64     `db:"pid"`    // 产品ID
		Amount     int64     `db:"amount"` // 订单金额
		Status     int64     `db:"status"` // 订单状态
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf) OrderModel {
	return &defaultOrderModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`order`",
	}
}

func (m *defaultOrderModel) Insert(data *Order) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, orderRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.Uid, data.Pid, data.Amount, data.Status)

	return ret, err
}

func (m *defaultOrderModel) FindOne(id int64) (*Order, error) {
	orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, id)
	var resp Order
	err := m.QueryRow(&resp, orderIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", orderRows, m.table)
		return conn.QueryRow(v, query, id)
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

func (m *defaultOrderModel) FindAllByUid(uid int64) ([]*Order, error) {
	var resp []*Order

	query := fmt.Sprintf("select %s from %s where `uid` = ?", orderRows, m.table)
	err := m.QueryRowsNoCache(&resp, query, uid)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultOrderModel) Update(data *Order) error {
	orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, orderRowsWithPlaceHolder)
		return conn.Exec(query, data.Uid, data.Pid, data.Amount, data.Status, data.Id)
	}, orderIdKey)
	return err
}

func (m *defaultOrderModel) Delete(id int64) error {

	orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, orderIdKey)
	return err
}

func (m *defaultOrderModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheOrderIdPrefix, primary)
}

func (m *defaultOrderModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", orderRows, m.table)
	return conn.QueryRow(v, query, primary)
}
