package network

import (
	"fmt"
	"net"
	"reflect"
	"sync"
)

var (
	DefaultOptions = Options{
		OnClose: func(obj interface{}, err error) {
			fmt.Println("miss close: [ class", reflect.TypeOf(obj).String(), "] err =", err)
		},
	}
)

func DefaultHandler(conn net.Conn) error {
	return fmt.Errorf("miss handle [addr %s]", conn.RemoteAddr().String())
}

type Option func(*Options)

type Options struct {
	Handler
	OnClose func(interface{}, error)
}

func NewOptions(opts ...Option) *Options {
	return DefaultOptions.Copy(opts...)
}

func (this *Options) Exchange(addr string) (handle Handler) {
	if handle = this.Handler; handle == nil {
		if handle = GetHandler(addr); handle == nil {
			handle = DefaultHandler
		}
	}
	return
}

func (this *Options) Init(opts ...Option) *Options {
	for _, o := range opts {
		o(this)
	}
	return this
}

func (this Options) Copy(opts ...Option) *Options {
	return this.Init(opts...)
}

//options
func OptionHandler(fn func(net.Conn) error) Option {
	return func(p *Options) {
		p.Handler = fn
	}
}

func OptionOnClose(fn func(interface{}, error)) Option {
	return func(p *Options) {
		p.OnClose = fn
	}
}

//global handler
var (
	globalMap = make(map[string]Handler)
	globalMux = new(sync.Mutex)
)

func OnHandler(handle Handler, addrs ...string) {
	globalMux.Lock()
	for _, addr := range addrs {
		globalMap[addr] = handle
	}
	globalMux.Unlock()
}

func UnHandler(addrs ...string) {
	globalMux.Lock()
	for _, addr := range addrs {
		delete(globalMap, addr)
	}
	globalMux.Unlock()
}

func GetHandler(addr string) (handle Handler) {
	globalMux.Lock()
	handle = globalMap[addr]
	globalMux.Unlock()
	return
}
