package remote

import (
	"fmt"
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
}

//socket选项
func OptionStream(obj interface{}) Option {
	return func(p *Socket) {
		switch conn := obj.(type) {
		case net.Conn:
			p.conn = With(conn)
		case Stream:
			p.conn = conn
		case func() Stream:
			p.conn = conn()
		case string:
			p.conn = WithAddr(conn)
		default:
			panic(fmt.Errorf("remote paramer error: [class %T]", obj))
		}
	}
}

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

func OptionNonPing() Option {
	return func(p *Socket) {
		p.pTimeout = 0
	}
}

func OptionDeathDelay(d time.Duration) Option {
	return func(p *Socket) {
		p.dieDelay = d
	}
}
