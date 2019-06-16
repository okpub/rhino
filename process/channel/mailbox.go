package channel

import (
	"fmt"

	"github.com/okpub/rhino/errors"
	"github.com/okpub/rhino/process"
)

//class mailbox(假邮箱，关闭后不能重复使用)
type myBuffer struct {
	process.UntypeProcess
	//options
	opts Options
}

func (this *myBuffer) init(args ...Option) {
	this.opts.Filler(args...)
	//uninitialized
	if pendingNum := this.opts.PendingNum; this.opts.Buffer == nil {
		if pendingNum < 5 {
			fmt.Printf("Warning: The email is too small, easy to jam. [size=%d] \n", pendingNum)
		} else if pendingNum > 1000 {
			fmt.Printf("Info: The email is too much. [size=%d] \n", pendingNum)
		}
		this.opts.Buffer = make(chan interface{}, pendingNum)
	}
}

//process
func (this *myBuffer) Start() (err error) {
	this.OnStarted()
	this.Schedule(this.run)
	return
}

func (this *myBuffer) Close() (err error) {
	defer func() { err = errors.Catch(recover()) }()
	this.opts.Close()
	return
}

func (this *myBuffer) run() {
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
	for body = range this.opts.Buffer {
		//先统计(处理失败也会记录)
		this.OnReceived(body)
		//后处理
		this.DispatchMessage(body)
	}
}

func (this *myBuffer) Post(v interface{}) (err error) {
	this.OnPosted(v)
	return errors.Try(func() error {
		return this.opts.Post(v)
	}, func(err error) {
		this.OnDiscarded(err, v)
	})
}

//options无法改变内部选项
func (this *myBuffer) Options() Options {
	return this.opts
}
