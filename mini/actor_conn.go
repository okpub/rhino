package mini

import (
	"context"
	"net"
)

type ActorConn struct {
	*ActorRef
	readCh chan []byte
	conn   net.Conn
}

func (this *ActorConn) Init(args ...Option) {
	this.ActorRef.Init(args...)
	this.readCh = make(chan []byte)
	//socket run
	go func(conn net.Conn) {
		var (
			//ctx = this.opts.Context
			b [1024]byte
		)
		for {
			n, err := conn.Read(b[0:])
			if err == nil {
				body := b[:n]
				this.readCh <- body //这里要避免堵塞
			} else {
				this.Close()
				break
			}
		}
	}(this.conn)
}

func (this *ActorConn) Run() (err error) {
	err = this.run(this.opts.Context, this.conn)
	return
}

func (this *ActorConn) run(ctx context.Context, conn net.Conn) (err error) {
	this.opts.OnStart()
Loop:
	for {
		select {
		case <-ctx.Done():
			break Loop
		case body := <-this.readCh:
			this.OnRead(body) //remote message
		case body, ok := <-this.taskCh:
			if ok {
				this.OnSend(body) //local message
			} else {
				goto End
			}
		}
	}
	this.Close()
	for body := range this.taskCh {
		this.OnSend(body)
	}
End:
	{
		//close socket
		conn.Close()
		//event stop
		this.opts.OnStop(err)
	}
	return
}

//message
func (this *ActorConn) OnRead(body []byte) {
	this.opts.Received(body)
}

func (this *ActorConn) OnSend(data interface{}) {
	//this.opts.Received(pack)
}
