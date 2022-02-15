package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	productFieldNames          = builder.RawFieldNames(&Product{})
	productRows                = strings.Join(productFieldNames, ",")
	productRowsExpectAutoSet   = strings.Join(stringx.Remove(productFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	productRowsWithPlaceHolder = strings.Join(stringx.Remove(productFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheProductIdPrefix = "cache:product:id:"
)

type (
	ProductModel interface {
		Insert(data *Product) (sql.Result, error)
		FindOne(id int64) (*Product, error)
		Update(data *Product) error
		Delete(id int64) error
	}

	defaultProductModel struct {
		sqlc.CachedConn
		table string
	}

	Product struct {
		Id         int64     `db:"id"`
		Name       string    `db:"name"`   // 产品名称
		Desc       string    `db:"desc"`   // 产品描述
		Stock      int64     `db:"stock"`  // 产品库存
		Amount     int64     `db:"amount"` // 产品金额
		Status     int64     `db:"status"` // 产品状态
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewProductModel(conn sqlx.SqlConn, c cache.CacheConf) ProductModel {
	return &defaultProductModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`product`",
	}
}

func (m *defaultProductModel) Insert(data *Product) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, productRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.Name, data.Desc, data.Stock, data.Amount, data.Status)

	return ret, err
}

func (m *defaultProductModel) FindOne(id int64) (*Product, error) {
	productIdKey := fmt.Sprintf("%s%v", cacheProductIdPrefix, id)
	var resp Product
	err := m.QueryRow(&resp, productIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", productRows, m.table)
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

func (m *defaultProductModel) Update(data *Product) error {
	productIdKey := fmt.Sprintf("%s%v", cacheProductIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, productRowsWithPlaceHolder)
		return conn.Exec(query, data.Name, data.Desc, data.Stock, data.Amount, data.Status, data.Id)
	}, productIdKey)
	return err
}

func (m *defaultProductModel) Delete(id int64) error {

	productIdKey := fmt.Sprintf("%s%v", cacheProductIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, productIdKey)
	return err
}

func (m *defaultProductModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheProductIdPrefix, primary)
}

func (m *defaultProductModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", productRows, m.table)
	return conn.QueryRow(v, query, primary)
}
