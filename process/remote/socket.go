package remote

import (
	"fmt"
	"net"
	"time"

	"github.com/rhino/process"
)

//net actor prosser
type mySocket struct {
	process.UntypeProcess
	//opts
	opts Options
}

func (this *mySocket) init(args ...Option) {
	this.opts.Fill(args...)
}

//procsess
func (this *mySocket) Start() (err error) {
	this.OnStarted()
	this.Schedule(this.run)
	return
}

func (this *mySocket) Close() (err error) {
	err = this.opts.Close()
	return
}

//private
func (this *mySocket) run() {
	var (
		err   error
		body  []byte
		debug = false
	)
	//dead why
	defer func() {
		if debug {
			if err == nil {
				err = process.CatchError(recover())
			} else {
				process.CatchError(recover()) //Ignore other error
			}
		}
		if err != nil {
			this.ThrowFailure(err, body)
		}
		this.Close()
		this.wait() //wait closing
		this.PostStop()
	}()
	//start call
	this.PreStart()
	//read once time
	this.opts.SetReadTimeout(this.opts.ReadTimeout)
	var pong bool
	for {
		body, err = this.opts.Read()
		//ping timeout
		this.opts.SetReadTimeout(this.opts.PingTimeout)
		if err == nil {
			pong = true
			//record read
			this.OnReceived(body)
			//call message
			this.DispatchMessage(body)
		} else {
			//check error
			if temp, ok := err.(net.Error); ok && temp.Temporary() {
				if pong {
					pong = false
					this.opts.SetReadTimeout(this.opts.ReadTimeout)
					//heartbeat notice
					this.DispatchMessage(err)
					this.OnFree()
				} else {
					this.Close() //close by timeout
					break
				}
			} else {
				break
			}
		}
	}
}

func (this *mySocket) wait() {
	if this.opts.DeathDelay > 0 {
		fmt.Println("wait closing :", this.opts.DeathDelay)
		time.Sleep(this.opts.DeathDelay)
	}
}

//If set will write timeout (default send)
func (this *mySocket) Send(b []byte) (err error) {
	this.OnPosted(b)
	if err = this.opts.Send(b); err != nil {
		this.OnDiscarded(err, b)
	}
	return
}

//options无法改变内部选项
func (this *mySocket) Options() Options {
	return this.opts
}
