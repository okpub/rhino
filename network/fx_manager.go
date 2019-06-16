package network

import (
	"fmt"
	"sync/atomic"

	"github.com/okpub/rhino/errors"
)

//如果你不加入到服务里面是无法快速启动监听的
func NewManager(addr string) Dialer {
	return &myManager{addr: addr, ch: make(chan Link)}
}

//class manager
type myManager struct {
	addr string
	code int32
	ch   chan Link
}

//join
func (this *myManager) Join(conn Link) (err error) {
	defer func() { err = errors.Catch(recover()) }()
	this.ch <- conn
	return
}

//Listener
func (this *myManager) Accept() (Link, error) {
	if conn, ok := <-this.ch; ok {
		return conn, nil
	}
	return nil, fmt.Errorf("accept closed addr=%s", this.addr)
}

func (this *myManager) Close() (err error) {
	if atomic.CompareAndSwapInt32(&this.code, 0, 1) {
		close(this.ch)
	} else {
		err = fmt.Errorf("close nothing addr=%s", this.addr)
	}
	return
}

func (this *myManager) Address() string {
	return this.addr
}
