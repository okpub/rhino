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
	this := &mySocket{}
	this.init(opts...)
	return this
}

func NewKeepActive(opts ...Option) SocketProcess {
	this := &mySocket{
		opts: Options{
			ReadTimeout: defaultReadTimeout,
			SendTimeout: defaultSendTimeout,
			PingTimeout: defaultPingTimeout,
		},
	}
	this.init(opts...)
	return this
}

//Producer
func Unbounded(opts ...Option) Producer {
	return func() SocketProcess {
		return New(opts...)
	}
}
