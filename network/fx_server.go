package network

import (
	"fmt"
	"net"
	"sync"
)

/*
*用于管理所有的link
*可以开启多个服务, 但是不能单独关闭，如果你想单独关闭，使用join
 */
func NewServer(opts ...Option) Symbiote {
	this := &Server{
		options: NewOptions(opts...),
	}
	return this
}

/*
* class 用来管理端口(同一监听同一处理方式)
* 需要统计
 */
type Server struct {
	options *Options
	mu      sync.Mutex
	acs     ConnSet
	lns     ListenSet
}

//快速启动(默认为tcp服务)
func (this *Server) Start(addr string) (err error) {
	this.Wrap(func() {
		this.Serve(Listen(addr), this.options.Exchange(addr))
	})
	return
}

func (this *Server) Serve(ln net.Listener, handle Handler) (err error) {
	defer ln.Close()
	this.trackListener(ln, true)
	defer this.trackListener(ln, false)
	fmt.Println("open server: ", ln.Addr().String())
	var conn net.Conn
	for {
		if conn, err = ln.Accept(); err == nil {
			this.joinAndRunner(conn, handle)
		} else {
			fmt.Println("ERROR: server err =", err.Error(), ", addr =", ln.Addr().String(), "[stop serve]")
			break
		}
	}
	this.options.OnClose(ln, err)
	return
}

func (this *Server) Close() (err error) {
	this.mu.Lock()
	err = this.closeListeners()
	this.closeIdleConns()
	this.mu.Unlock()
	return
}

//private
func (this *Server) joinAndRunner(conn net.Conn, handle Handler) {
	this.trackConn(conn, true)
	this.Wrap(func() {
		handle(conn)
		this.trackConn(conn, false)
		conn.Close()
		//this.OnClose(conn, err)
	})
}

func (this *Server) closeListeners() (err error) {
	for ln := range this.lns {
		if cerr := ln.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}
	this.lns = nil
	return
}

func (this *Server) trackListener(ln net.Listener, add bool) {
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

func (this *Server) closeIdleConns() {
	for conn := range this.acs {
		conn.Close()
	}
	this.acs = nil
}

func (this *Server) trackConn(conn net.Conn, add bool) {
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
func (this *Server) Wrap(fn func()) {
	go fn()
}
