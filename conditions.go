package database

import (
	"fmt"
	"reflect"
	"strings"
)

type Condition interface {
	SQL() string
	Values() []interface{}
}

type simpleCondition struct {
	sql   string
	value interface{}
}

func (cond *simpleCondition) SQL() string {
	if !strings.Contains(cond.sql, " ") {
		return fmt.Sprintf("%s = ?", cond.sql)
	}

	if strings.Contains(cond.sql, " IN") {
		v := reflect.ValueOf(cond.value)
		placeholders := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			placeholders[i] = "?"
		}
		return fmt.Sprintf("%s (%s)", cond.sql, strings.Join(placeholders, ", "))
	}

	if !strings.Contains(cond.sql, "?") {
		return fmt.Sprintf("%s ?", cond.sql)
	}

	return cond.sql
}

func (cond *simpleCondition) Values() []interface{} {
	if strings.Contains(cond.sql, " IN") {
		v := reflect.ValueOf(cond.value)
		var values []interface{}
		for i := 0; i < v.Len(); i++ {
			values = append(values, v.Index(i).Interface())
		}
		return values
	}

	return []interface{}{cond.value}
}

type compareJSONCondition struct {
	column, path string
	value        interface{}
}

// CompareJSON creates a new condition that checks if a value inside a JSON
// object of a column is equal to the provided value.
func CompareJSON(column, path string, value interface{}) Condition {
	return &compareJSONCondition{
		column: column,
		path:   path,
		value:  value,
	}
}

func (cond *compareJSONCondition) SQL() string {
	return fmt.Sprintf("JSON_EXTRACT(%s, '%s') = ?", cond.column, cond.path)
}

func (cond *compareJSONCondition) Values() []interface{} {
	return []interface{}{cond.value}
}

// EscapeLike escapes a value to insert it in a LIKE query without unexpected wildcards.
// After using this function to clean the value you can add the wildcards you need
// to the query.
func EscapeLike(str string) string {
	str = strings.Replace(str, "%", `\%`, -1)
	str = strings.Replace(str, "_", `\_`, -1)
	return str
}
