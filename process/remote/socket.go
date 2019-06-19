package remote

import (
	"fmt"
	"net"
	"time"

	"github.com/okpub/rhino/errors"
	"github.com/okpub/rhino/process"
)

//net actor prosser
type Socket struct {
	process.UntypeProcess

	//net io
	conn Stream

	//The read timeout, no Settings, is blocking
	rTimeout time.Duration

	//The send timeout, no Settings, is blocking
	sTimeout time.Duration

	//The heartbeat timeout, no Settings, no heartbeat
	pTimeout time.Duration

	//The delay closing
	dieDelay time.Duration
}

//private 填充，初始化使用
func (this *Socket) Filler(opts ...Option) *Socket {
	for _, o := range opts {
		o(this)
	}
	return this
}

func (this Socket) Copy(opts ...Option) *Socket {
	return this.Filler(opts...)
}

//procsess
func (this *Socket) Start() (err error) {
	this.OnStarted()
	this.Schedule(this.run)
	return
}

func (this *Socket) Close() (err error) {
	err = this.conn.Close()
	return
}

//private
func (this *Socket) run() {
	var (
		err   error
		body  []byte
		debug = false
		conn  = this.conn
	)
	//dead why
	defer func() {
		if debug {
			if err == nil {
				err = errors.Catch(recover())
			}
		}
		if err != nil {
			this.ThrowFailure(err, body)
		}
		this.Close() //must close
		this.Sleep() //wait closing
		this.PostStop()
	}()
	//start call
	this.PreStart()
	//read once time
	conn.SetReadTimeout(this.rTimeout)
	var pong bool
	for {
		body, err = conn.Read()
		//ping timeout
		conn.SetReadTimeout(this.pTimeout)
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
					conn.SetReadTimeout(this.rTimeout)
					//heartbeat notice
					this.DispatchMessage(err)
					//this.OnFree()
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

func (this *Socket) Sleep() {
	if this.dieDelay > 0 {
		fmt.Println("wait closing :", this.dieDelay)
		time.Sleep(this.dieDelay)
	}
}

func (this *Socket) Send(b []byte) (err error) {
	var (
		conn = this.conn
	)
	this.OnPosted(b)
	//If you not handle return value, please do not set up, Avoid sending part only
	conn.SetSendTimeout(this.sTimeout)
	if err = conn.Write(b); err != nil {
		this.OnDiscarded(err, b)
	}
	return
}
