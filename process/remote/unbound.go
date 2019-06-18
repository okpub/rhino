package remote

import (
	"time"
)

//default values
const (
	defaultReadTimeout = time.Second * 3
	defaultSendTimeout = time.Second * 5
	defaultPingTimeout = time.Minute * 1
)

type Producer func() SocketProcess

//写超时，无心跳和读超时
func New(opts ...Option) SocketProcess {
	this := &Socket{
		sTimeout: defaultSendTimeout, // default send timeout
	}
	return this.filler(opts...)
}

//存在心跳和读写超时
func NewKeepActive(opts ...Option) SocketProcess {
	this := &Socket{
		rTimeout: defaultReadTimeout, // default read timeout
		sTimeout: defaultSendTimeout, // default send timeout
		pTimeout: defaultPingTimeout, // default heartbeat
	}
	return this.filler(opts...)
}

//无心跳，无读写超时
func NewBlocking(opts ...Option) SocketProcess {
	return new(Socket).filler(opts...)
}

//Producer(默认阻塞)
func Unbounded(opts ...Option) Producer {
	return func() SocketProcess {
		return NewBlocking(opts...)
	}
}
