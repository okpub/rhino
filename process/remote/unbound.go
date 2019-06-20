package remote

import (
	"time"
)

//default values
var (
	defaultReadTimeout = time.Second * 3
	defaultPingTimeout = time.Minute * 1
	defaultSendTimeout = time.Second * 5 //no used
)

type Producer func() SocketProcess

//new
func New(opts ...Option) SocketProcess {
	this := &Socket{}
	return this.Init(opts...)
}

func NewAddr(addr string, opts ...Option) SocketProcess {
	this := &Socket{}
	//dial
	OptionStream(addr)(this)
	//other
	return this.Init(opts...)
}

func NewKeepActive(obj interface{}, opts ...Option) SocketProcess {
	this := &Socket{
		rTimeout: defaultReadTimeout,
		//sTimeout: defaultSendTimeout,
		pTimeout: defaultPingTimeout,
	}
	//Note that type: net.Conn
	OptionStream(obj)(this)
	//other options
	return this.Init(opts...)
}

//Producer(默认无超时)
func Unbounded(opts ...Option) Producer {
	return func() SocketProcess {
		return New(opts...)
	}
}
