package mysql

import (
	"database/sql"
	"reflect"
)

//mysql唯一标志
const TagKey = "mysql"

//mysql array row
type RowArray []interface{}

func (arr RowArray) Obj(ptr interface{}) (err error) {
	vals := reflect.ValueOf(ptr).Elem()
	for i := 0; i < vals.NumField(); i++ {
		if item := vals.Field(i); item.CanSet() {
			unpacking(item, arr[i])
		}
	}
	return
}

//mysql map row
type RowObject map[string]interface{}

func (obj RowObject) Obj(ptr interface{}) (err error) {
	tags := reflect.TypeOf(ptr).Elem()
	//values
	vals := reflect.ValueOf(ptr).Elem()
	for i := 0; i < vals.NumField(); i++ {
		if key, ok := tags.Field(i).Tag.Lookup(TagKey); ok {
			if item := vals.Field(i); item.CanSet() {
				unpacking(item, obj[key])
			}
		}
	}
	return
}

/*查询结果*/
func Table(rows *sql.Rows, code error) (arr []RowObject, err error) {
	if err = code; err != nil {
		return
	}
	defer rows.Close()
	var columns []string
	if columns, err = rows.Columns(); err != nil {
		return
	}
	var (
		count = len(columns)
		vals  = make([]interface{}, count)
		ptrs  = make([]interface{}, count)
	)
	for rows.Next() {
		for i := 0; i < count; i++ {
			ptrs[i] = &vals[i]
		}
		rows.Scan(ptrs...)
		obj := make(RowObject)
		for i, k := range columns {
			obj[k] = vals[i]
			//common.Trace("%s=%v", k, vals[i])
		}
		arr = append(arr, obj)
	}
	return
}

func Array(rows *sql.Rows, code error) (arr []RowArray, err error) {
	if err = code; err != nil {
		return
	}
	defer rows.Close()
	var columns []string
	if columns, err = rows.Columns(); err != nil {
		return
	}
	var (
		count = len(columns)
		vals  = make([]interface{}, count)
		ptrs  = make([]interface{}, count)
	)
	for rows.Next() {
		for i := 0; i < count; i++ {
			ptrs[i] = &vals[i]
		}
		rows.Scan(ptrs...)
		obj := make(RowArray, count)
		for i := range columns {
			obj[i] = vals[i]
			//common.Trace("%s=%v", columns[i], vals[i])
		}
		arr = append(arr, obj)
	}
	return
}
