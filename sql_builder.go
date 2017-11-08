package database

import (
	"fmt"
	"strings"
	"reflect"
)

type sqlBuilder struct {
  table string
  props []*Property
  orders []string
  conditions []Condition
  limit, offset int64
}

func newSQLBuilder(model Model) *sqlBuilder {
	return &sqlBuilder{
		table: model.TableName(),
	}
}

func (b *sqlBuilder) AddProperty(prop *Property) {
	b.props = append(b.props, prop)
}

func (b *sqlBuilder) cols() []string {
	var cols []string
	for _, prop := range b.props {
		cols = append(cols, prop.Name)
	}

	return cols
}

func (b *sqlBuilder) SelectSQL() (string, []interface{}) {
	return b.SelectSQLCols(b.cols()...)
}

func (b *sqlBuilder) SelectSQLCols(cols ...string) (string, []interface{}) {
	var conds []string
	var values []interface{}
	for _, cond := range b.conditions {
		conds = append(conds, cond.SQL())
		values = append(values, cond.Values()...)
	}

	sql := fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(cols, ", "), b.table)

	if len(conds) > 0 {
		sql = fmt.Sprintf("%s WHERE %s", sql, strings.Join(conds, " AND "))
	}
	if len(b.orders) > 0 {
		sql = fmt.Sprintf("%s ORDER BY %s", sql, strings.Join(b.orders, ", "))
	}
	if b.limit > 0 {
		sql = fmt.Sprintf("%s LIMIT %d,%d", sql, b.offset, b.limit)
	}

	return sql, values
}

func (b *sqlBuilder) Condition(cond Condition) {
	b.conditions = append(b.conditions, cond)
}

func (b *sqlBuilder) UpdateSQL() (string, []interface{}) {
	var values []interface{}

	var updates []string
	for _, prop := range b.props {
		updates = append(updates, fmt.Sprintf("%s = ?", prop.Name))
		values = append(values, prop.Value)
	}
	
	var conds []string
	for _, cond := range b.conditions {
		conds = append(conds, cond.SQL())
		values = append(values, cond.Values()...)
	}

	sql := fmt.Sprintf(`UPDATE %s SET %s WHERE %s`, b.table, strings.Join(updates, ", "), strings.Join(conds, " AND "))

	return sql, values
}

func (b *sqlBuilder) InsertSQL() (string, []interface{}) {
	var values []interface{}

	var placeholders []string
	for _, prop := range b.props {
		placeholders = append(placeholders, "?")
		values = append(values, prop.Value)
	}

	sql := fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s)`, b.table, strings.Join(b.cols(), ", "), strings.Join(placeholders, ", "))

	return sql, values
}

func (b *sqlBuilder) DeleteSQL() (string, []interface{}) {
	var conds []string
	var values []interface{}
	for _, cond := range b.conditions {
		conds = append(conds, cond.SQL())
		values = append(values, cond.Values()...)
	}

	sql := fmt.Sprintf(`DELETE FROM %s WHERE %s`, b.table, strings.Join(conds, " AND "))

	return sql, values
}

func (b *sqlBuilder) Hydrate() {
	for _, prop := range b.props {
		prop.Value = reflect.ValueOf(prop.Pointer).Elem().Interface()
	}
}