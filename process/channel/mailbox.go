package channel

import (
	"github.com/okpub/rhino/errors"
	"github.com/okpub/rhino/process"
)

//class mailbox(假邮箱，关闭后不能重复使用)
type Mailbox struct {
	process.UntypeProcess
	//message
	buffer     chan interface{} //消息通道
	pendingNum int              //会存在默认通道
	blocking   bool             //非阻塞模式(默认阻塞)
}

func (this *Mailbox) filler(opts ...Option) *Mailbox {
	for _, o := range opts {
		o(this)
	}
	if this.buffer == nil {
		this.buffer = make(chan interface{}, this.pendingNum)
	}
	return this
}

//process
func (this *Mailbox) Start() (err error) {
	this.OnStarted()
	this.Schedule(this.run)
	return
}

func (this *Mailbox) Close() (err error) {
	defer func() { err = errors.Catch(recover()) }()
	close(this.buffer)
	return
}

func (this *Mailbox) run() {
	var (
		body  interface{}
		err   error
		debug = false
	)
	defer func() {
		if debug {
			if err = errors.Catch(recover()); err != nil {
				this.ThrowFailure(err, body)
			}
		}
		this.Close()
		this.PostStop()
	}()
	this.PreStart()
	//run
	for body = range this.buffer {
		//First of all statistics (processing failure will also record)
		this.OnReceived(body)
		//process message
		this.DispatchMessage(body)
	}
}

func (this *Mailbox) Post(v interface{}) (err error) {
	this.OnPosted(v)
	return errors.Try(func() error {
		return this.sendMessage(v)
	}, func(err error) {
		this.OnDiscarded(err, v)
	})
}

//private
func (this *Mailbox) sendMessage(v interface{}) (err error) {
	if this.blocking {
		this.buffer <- v
	} else {
		select {
		case this.buffer <- v:
		default:
			err = OverfullErr
		}
	}
	return
}
