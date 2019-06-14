package network

import (
	"fmt"
	"net"
	"sync"
	"time"
)

/*
*用于管理所有的link
*可以开启多个服务, 但是不能单独关闭，如果你想单独关闭，使用join
 */
func NewServer() Server {
	return &myServer{}
}

/*
*class 用来管理端口(同一监听同一处理方式)
* 需要统计
 */
type myServer struct {
	mu  sync.Mutex
	acs ConnSet
	lns ListenSet
}

//快速启动(默认为tcp服务)
func (this *myServer) Start(addr string) error {
	return this.Join(NewListener(addr))
}

func (this *myServer) Group(addr string) (manager Dialer) {
	manager = NewManager(addr)
	this.Join(manager)
	return
}

//手动加入(加入失败会被关闭)
func (this *myServer) Join(ln Listener) (err error) {
	var newAgent Handler
	if newAgent, err = GetHandler(ln.Address()); err == nil {
		this.Wrap(func() {
			this.Serve(ln, newAgent)
		})
	} else {
		ln.Close()
	}
	return
}

//最重要的接口
func (this *myServer) Serve(ln Listener, newAgent Handler) (err error) {
	defer ln.Close()
	this.trackListener(ln, true)
	defer this.trackListener(ln, false)
	fmt.Println("open server: ", ln.Address())
	var conn Link
	for {
		if conn, err = ln.Accept(); err == nil {
			this.joinAndRunner(conn, newAgent)
		} else {
			if temp, ok := err.(net.Error); ok && temp.Temporary() {
				fmt.Println("WARN: server err =", err.Error(), ", addr =", ln.Address(), "[wait 200ms]")
				time.Sleep(time.Millisecond * 200)
			} else {
				fmt.Println("ERROR: server err =", err.Error(), ", addr =", ln.Address(), "[stop serve]")
				break
			}
		}
	}
	return
}

func (this *myServer) Close() (err error) {
	this.mu.Lock()
	err = this.closeListeners()
	this.closeIdelConns()
	this.mu.Unlock()
	return
}

//private
func (this *myServer) joinAndRunner(conn Link, newAgent Handler) {
	this.trackConn(conn, true)
	this.Wrap(func() {
		newAgent(conn).Run()
		this.trackConn(conn, false)
		conn.Close()
	})
}

func (this *myServer) closeListeners() (err error) {
	for ln := range this.lns {
		if cerr := ln.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}
	this.lns = nil
	return
}

func (this *myServer) trackListener(ln Listener, add bool) {
	this.mu.Lock()
	if this.lns == nil {
		this.lns = make(ListenSet)
	}
	if add {
		this.lns[ln] = struct{}{}
	} else {
		delete(this.lns, ln)
	}
	this.mu.Unlock()
}

func (this *myServer) closeIdelConns() {
	for conn := range this.acs {
		conn.Close()
	}
	this.acs = nil
}

func (this *myServer) trackConn(conn Link, add bool) {
	this.mu.Lock()
	if this.acs == nil {
		this.acs = make(ConnSet)
	}
	if add {
		this.acs[conn] = struct{}{}
	} else {
		delete(this.acs, conn)
	}
	this.mu.Unlock()
}

//go func
func (this *myServer) Wrap(fn func()) {
	go fn()
}
