package mini

import (
	"context"
	"fmt"
	"sync/atomic"
)

//ref
type ActorRef struct {
	opts      Options
	closeFlag int32
	taskCh    chan interface{}
}

func (this *ActorRef) Init(args ...Option) {
	//this.opts = DefaultActor //default options
	for _, o := range args {
		o(&this.opts)
	}

	this.taskCh = make(chan interface{}, this.opts.PendingNum)
}

func (this *ActorRef) Options() Options {
	return this.opts
}

func (this *ActorRef) With() (child context.Context, cancel context.CancelFunc) {
	child, cancel = context.WithCancel(this.opts.Context)
	return
}

func (this *ActorRef) Run() (err error) {
	err = this.run(this.opts.Context)
	return
}

func (this *ActorRef) run(ctx context.Context) (err error) {
	this.opts.OnStart()
Loop:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break Loop
		case body, ok := <-this.taskCh:
			if ok {
				this.handleMessage(body)
			} else {
				goto End
			}
		}
	}
	//flush
	this.Close()
	for body := range this.taskCh {
		this.handleMessage(body)
	}
End:
	{
		this.opts.OnStop(err)
	}
	return
}

func (this *ActorRef) handleMessage(data interface{}) {
	this.opts.Received(data) //handle data
}

func (this *ActorRef) Send(data interface{}, args ...PublishOption) (err error) {
	defer func() {
		if code := recover(); code != nil {
			err = fmt.Errorf("actor error: %v", code)
		}
	}()

	var (
		options = DefaultPublisher //default send options
	)
	for _, o := range args {
		o(&options)
	}

	select {
	case <-options.Done():
		err = options.Err()
	case this.taskCh <- data:
	}
	return
}

func (this *ActorRef) Close() (err error) {
	if atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		close(this.taskCh)
	} else {
		err = fmt.Errorf("closed context")
	}
	return
}
