package remote

import (
	"time"
)

/*
* 拷贝出来的无法影响之前的,只能修改
 */
type Options struct {
	//network
	Stream

	//The read timeout, no Settings, is blocking
	ReadTimeout time.Duration

	//The send timeout, no Settings, is blocking
	SendTimeout time.Duration

	//The heartbeat timeout, no Settings, no heartbeat
	PingTimeout time.Duration

	//The delay closing
	DeathDelay time.Duration
}

func (this *Options) Fill(args ...Option) {
	for _, o := range args {
		o(this)
	}
}

func (this *Options) Send(b []byte) error {
	this.SetSendTimeout(this.SendTimeout)
	return this.Write(b)
}

//func (this *Options) copy() Options {
//	return *this
//}

//socket选项
func OptionWithStream(conn Stream) Option {
	return func(p *Options) {
		p.Stream = conn
	}
}

func OptionWithFunc(dial func() Stream) Option {
	return func(p *Options) {
		p.Stream = dial()
	}
}

func OptionReadTimeout(d time.Duration) Option {
	return func(p *Options) {
		p.ReadTimeout = d
	}
}

func OptionSendTimeout(d time.Duration) Option {
	return func(p *Options) {
		p.SendTimeout = d
	}
}

func OptionPingTimeout(d time.Duration) Option {
	return func(p *Options) {
		p.PingTimeout = d
	}
}

func OptionDeathDelay(d time.Duration) Option {
	return func(p *Options) {
		p.DeathDelay = d
	}
}
