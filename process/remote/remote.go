package remote

import (
	"net"
	"time"

	"github.com/okpub/rhino/process"
)

type Option func(*Socket)

/*
* 准确来说，选项和具体实例需要分开，这里未分开
 */
type SocketProcess interface {
	process.Process
	Send([]byte) error
	//Copy(...Option) SocketProcess //改变参数，得到新副本(无法修改已经运行的参数部分)
}

//socket选项
func OptionStream(conn Stream) Option {
	return func(p *Socket) {
		p.conn = conn
	}
}

func OptionConn(conn interface{}) Option {
	return func(p *Socket) {
		p.conn = With(conn.(net.Conn))
	}
}

func OptionFunc(dial func() Stream) Option {
	return func(p *Socket) {
		p.conn = dial()
	}
}

func OptionAddr(addr string) Option {
	return func(p *Socket) {
		p.conn = WithAddr(addr)
	}
}

//timeout
func OptionReadTimeout(d time.Duration) Option {
	return func(p *Socket) {
		p.rTimeout = d
	}
}

func OptionSendTimeout(d time.Duration) Option {
	return func(p *Socket) {
		p.sTimeout = d
	}
}

func OptionPingTimeout(d time.Duration) Option {
	return func(p *Socket) {
		p.pTimeout = d
	}
}

func OptionDeathDelay(d time.Duration) Option {
	return func(p *Socket) {
		p.dieDelay = d
	}
}

//filler
func OptionFiller(opts ...Option) Option {
	return func(p *Socket) {
		p.filler(opts...)
	}
}
