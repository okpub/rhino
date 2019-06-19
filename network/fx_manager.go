package network

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/okpub/rhino/errors"
)

//可以看作统一类型的连接管理器
func NewManager(addr string) Manager {
	return &myManager{addr: addr, ch: make(chan net.Conn)}
}

//class manager
type myManager struct {
	addr string
	code int32
	ch   chan net.Conn
}

func (this *myManager) Dial(addr string) error {
	if atomic.LoadInt32(&this.code) == 1 {
		return fmt.Errorf("close of dial addr=%s", this.addr)
	}
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err == nil {
		if err = this.joinAndConn(conn); err != nil {
			conn.Close()
		}
	}
	return err
}

//join
func (this *myManager) joinAndConn(conn net.Conn) (err error) {
	defer func() { err = errors.Catch(recover()) }()
	this.ch <- conn
	return
}

//Listener
func (this *myManager) Accept() (net.Conn, error) {
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

func (this *myManager) Addr() net.Addr  { return this }
func (this *myManager) Network() string { return "tcp" }
func (this *myManager) String() string  { return this.addr }
