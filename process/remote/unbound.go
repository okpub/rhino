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
	return this.Filler(opts...)
}

func NewAddr(addr string, opts ...Option) SocketProcess {
	this := &Socket{
		conn: WithAddr(addr), //direct connection
	}
	return this.Filler(opts...)
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
	return this.Filler(opts...)
}

//Producer(默认无超时)
func Unbounded(opts ...Option) Producer {
	return func() SocketProcess {
		return New(opts...)
	}
}
