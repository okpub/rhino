package network

import (
	"net"
	"time"

	"golang.org/x/net/websocket"
)

/*
* 目前使用actor框架，基本network可以休眠
 */

const (
	WEB_LINK = "web"
	TCP_LINK = "tcp"
	UDP_LINK = "udp"
)

type (
	ConnSet   map[net.Conn]struct{}
	ListenSet map[net.Listener]struct{}
	Handler   func(net.Conn) error

	//服务
	Symbiote interface {
		Start(string) error
		Serve(net.Listener, Handler) error
		Close() error
	}

	//连接器管理(目前没啥用,actor里面已经实现)
	Manager interface {
		net.Listener
		Dial(string) error
	}
)

//自选类型
func DialScan(kind string, addr string) (net.Conn, error) {
	switch kind {
	case WEB_LINK:
		return websocket.Dial("ws://"+addr, "", "http://"+addr)
	case TCP_LINK:
		return net.DialTimeout("tcp", addr, time.Second*5)
	case UDP_LINK:
		return net.DialTimeout("udp", addr, time.Second*5)
	default:
		return DialScan(TCP_LINK, addr)
	}
}
