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

//new
func New(opts ...Option) SocketProcess {
	this := &mySocket{
		opts: Options{
			SendTimeout: defaultSendTimeout, // default send timeout
		},
	}
	this.init(opts...)
	return this
}

func NewKeepActive(opts ...Option) SocketProcess {
	this := &mySocket{
		opts: Options{
			ReadTimeout: defaultReadTimeout, // default read timeout
			SendTimeout: defaultSendTimeout, // default send timeout
			PingTimeout: defaultPingTimeout, // default heartbeat
		},
	}
	this.init(opts...)
	return this
}

func NewBlocking(opts ...Option) SocketProcess {
	this := &mySocket{}
	this.init(opts...)
	return this
}

//Producer
func Unbounded(opts ...Option) Producer {
	return func() SocketProcess {
		return NewBlocking(opts...)
	}
}
