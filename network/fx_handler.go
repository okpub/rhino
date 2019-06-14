package network

import (
	"fmt"
	"sync"
)

/*
*所有服务器建立的注册(地址相同，那么建立build属性就一样)
 */
var (
	DefaultRegister = make(map[string]Handler)
	DefaultMux      sync.Mutex
)

func GetHandler(addr string) (handle Handler, err error) {
	DefaultMux.Lock()
	if cb, ok := DefaultRegister[addr]; ok {
		handle = cb
	} else {
		err = fmt.Errorf("can't find handler addr %s", addr)
	}
	DefaultMux.Unlock()
	return
}

func OnHandler(handle Handler, args ...string) {
	DefaultMux.Lock()
	for _, addr := range args {
		DefaultRegister[addr] = handle
	}
	DefaultMux.Unlock()
}

func OffHandler(args ...string) {
	DefaultMux.Lock()
	for _, addr := range args {
		delete(DefaultRegister, addr)
	}
	DefaultMux.Unlock()
}
