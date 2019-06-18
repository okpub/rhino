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

	//net地址
	NetAddr interface {
		Network() string //网络类型
		Addr() string    //客户端dial地址
		PubAddr() string //服务器lnnr地址
		//Next() NetAddr   //next网络节点，可以用作集群或者分布式节点链
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

/*
*是否为暂时的错误(超时或者其他)
 */
func CheckNetTemporary(err error) bool {
	temp, ok := err.(net.Error)
	return ok && temp.Temporary()
}
