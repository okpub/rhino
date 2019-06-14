package sync

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-sloth/vodka/context"
)

//同步消息机制
type (

	//同步消息
	MessageChannel interface {
		//生成子通道
		Accept(context.Context) Conn
		//读取
		Read(Conn, time.Duration) (body interface{}, err error)
		//提交 用于跨进程执行
		Commit(int64, interface{}) error
		//撤回 用于跨进程执行
		Rollback(int64) error
	}

	//同步io 多协程不安全(只能依赖于某个actor中)
	Conn interface {
		SyncID() int64
		//外部请不要随意调用
		Read() (interface{}, error)
		ReadTimeout(time.Duration) (interface{}, error)
		Write(interface{}) error
		WriteTimeout(time.Duration, interface{}) error
		Rollback() error //撤回一次read
		Close() error    //用于同进程内执行
	}
)

/*
* 异端同步消息 cli->svr/svr->cli
 */
var (
	backgroundChannel = &eventChannel{mp: make(ConnSet)}
	sync_id           int64
)

func incr() int64 {
	return atomic.AddInt64(&sync_id, 1)
}

func Channel() MessageChannel {
	return backgroundChannel
}

//class
type ConnSet map[int64]Conn

type eventChannel struct {
	id int64
	mp ConnSet
	mu sync.Mutex
}

/*
*用于监听退出
 */
func (this *eventChannel) Accept(ctx context.Context) (conn Conn) {
	conn = With(ctx)
	this.mu.Lock()
	if this.mp == nil {
		this.mp = make(ConnSet)
	}
	this.mp[conn.SyncID()] = conn
	this.mu.Unlock()
	return
}

/*
*读取后关闭
 */
func (this *eventChannel) Read(conn Conn, delay time.Duration) (body interface{}, err error) {
	if delay > 0 {
		body, err = conn.ReadTimeout(delay)
	} else {
		body, err = conn.Read()
	}
	//remove
	this.mu.Lock()
	delete(this.mp, conn.SyncID())
	this.mu.Unlock()
	//close
	conn.Close()
	return
}

/*
* 提交, 只会提交一次
 */
func (this *eventChannel) Commit(id int64, v interface{}) (err error) {
	if conn, ok := this.loadAndRemove(id); ok {
		conn.Write(v)
	} else {
		err = fmt.Errorf("commit err: can't find sync id=%d", id)
	}
	return
}

/*
* 撤回, 只会撤回一次
 */
func (this *eventChannel) Rollback(id int64) (err error) {
	if conn, ok := this.loadAndRemove(id); ok {
		conn.Close()
	} else {
		err = fmt.Errorf("rollback err: can't find sync id=%d", id)
	}
	return
}

//private
func (this *eventChannel) loadAndRemove(id int64) (conn Conn, ok bool) {
	this.mu.Lock()
	conn, ok = this.mp[id]
	if ok {
		delete(this.mp, id)
	}
	this.mu.Unlock()
	return
}

/*
* 同步辅助功能
 */
func With(ctx context.Context) Conn {
	this := &sync_conn{
		PIDSignal: context.WithSignal(ctx),
		id:        incr(),
		blockCh:   make(chan interface{}),
		errorCh:   make(chan error),
	}
	return this
}

type sync_conn struct {
	*context.PIDSignal
	id      int64
	blockCh chan interface{}
	errorCh chan error
}

/*
* blocking id
 */
func (this *sync_conn) SyncID() int64 {
	return this.id
}

/*
* read blocking
 */
func (this *sync_conn) Read() (body interface{}, err error) {
	select {
	case <-this.Done():
		err = this.Err()
	case err = <-this.errorCh:
		//rollback
	case body = <-this.blockCh:
	}
	return
}

/*
* read with timeout
 */
func (this *sync_conn) ReadTimeout(delay time.Duration) (body interface{}, err error) {
	select {
	case <-this.Done():
		err = this.Err()
	case err = <-this.errorCh:
	//rollback
	case <-time.After(delay):
		err = context.DeadlineExceeded
	case body = <-this.blockCh:
	}
	return
}

/*
* write blocking
 */
func (this *sync_conn) Write(body interface{}) (err error) {
	select {
	case <-this.Done():
		err = this.Err()
	case this.blockCh <- body:
		//	default:
		//		err = fmt.Errorf("write error: blocking channel")
	}
	return
}

/*
* write with timeout
 */
func (this *sync_conn) WriteTimeout(delay time.Duration, body interface{}) (err error) {
	select {
	case <-this.Done():
		err = this.Err()
	case <-time.After(delay):
		err = context.DeadlineExceeded
	case this.blockCh <- body:
		//	default:
		//		err = fmt.Errorf("write error: blocking channel")
	}
	return
}

func (this *sync_conn) Rollback() (err error) {
	select {
	case <-this.Done():
		err = this.Err()
	case this.errorCh <- fmt.Errorf("rollback sync id=%d", this.id):
		//ok
	default:
		err = fmt.Errorf("rollback error: not read")
	}
	return
}

/*
* close sign
 */
func (this *sync_conn) Close() error {
	return this.RemoveSelf(true, context.Canceled)
}
