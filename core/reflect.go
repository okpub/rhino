package core

import (
	"fmt"
	"reflect"
)

func Typeof(p interface{}) string {
	if p == nil {
		return "<nil>"
	}
	return reflect.ValueOf(p).Type().String()
}

/*
*对象构造器
 */
type ObjectModel map[int]reflect.Type

//register(1, Test{})	//注册结构，而非指针
func (hset ObjectModel) Register(cmd int, obj interface{}) {
	hset[cmd] = reflect.TypeOf(obj)
}

func (hset ObjectModel) UnRegister(cmd int) {
	delete(hset, cmd)
}

func (hset ObjectModel) New(cmd int) (val interface{}, err error) {
	if field, ok := hset[cmd]; ok {
		val = reflect.New(field).Interface()
	} else {
		err = fmt.Errorf("undefined class cmd=%d", cmd)
	}
	return
}
