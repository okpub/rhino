package mysql

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Bool(v interface{}) bool {
	return Int64(v) != 0
}

func Int64(v interface{}) int64 {
	switch d := v.(type) {
	case int64:
		return d
	case string:
		n, _ := strconv.ParseInt(d, 10, 64)
		return n
	case []byte:
		return Int64(string(d))
	}
	return 0
}

func String(v interface{}) string {
	switch d := v.(type) {
	case []byte:
		return string(d)
	case string:
		return d
	case int64:
		return strconv.FormatInt(d, 10)
	}
	return ""
}

func Date(v interface{}) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", String(v), time.Local)
}

//private
func unpacking(fd reflect.Value, v interface{}) {
	switch fd.Kind() {
	case reflect.Bool:
		fd.SetBool(Bool(v))
	case reflect.String:
		fd.SetString(String(v))

	case reflect.Int:
		fd.SetInt(Int64(v))
	case reflect.Int8:
		fd.SetInt(Int64(v))
	case reflect.Int16:
		fd.SetInt(Int64(v))
	case reflect.Int32:
		fd.SetInt(Int64(v))
	case reflect.Int64:
		fd.SetInt(Int64(v))

	case reflect.Uint:
		fd.SetInt(Int64(v))
	case reflect.Uint8:
		fd.SetInt(Int64(v))
	case reflect.Uint16:
		fd.SetInt(Int64(v))
	case reflect.Uint32:
		fd.SetInt(Int64(v))
	case reflect.Uint64:
		fd.SetUint(uint64(Int64(v)))
	default:
		setPointer(fd.Addr().Interface(), v)
	}
}

//其他指针
func setPointer(ptr interface{}, v interface{}) {
	switch p := ptr.(type) {
	case *time.Time:
		if tm, err := Date(v); err == nil {
			*p = tm
		}
	default:
		panic(fmt.Errorf("mysql unpacking err: type=%s value=%v", reflect.TypeOf(ptr).String(), v))
	}
}
