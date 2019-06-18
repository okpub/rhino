package network

import (
	"fmt"
	"net"
	"time"
)

//默认tcp
func NewListener(addr string) Listener {
	ln, err := net.Listen("tcp", addr)
	return &myListener{addr: addr, err: err, ln: ln, free: time.Millisecond * 200}
}

//class listener
type myListener struct {
	addr string
	err  error
	free time.Duration
	ln   net.Listener
}

func (this *myListener) Accept() (conn Link, err error) {
process:
	if err = this.err; err == nil {
		conn, err = this.ln.Accept()
		if temp, ok := err.(net.Error); ok && temp.Temporary() {
			fmt.Println("WARN: ", err.Error(), ", addr =", this.addr, "[wait ", this.free, "]")
			time.Sleep(this.free)
			goto process
		}
	}
	return
}

func (this *myListener) Close() (err error) {
	if err = this.err; err == nil {
		err = this.ln.Close()
	}
	return
}

func (this *myListener) Address() string {
	return this.addr
}
