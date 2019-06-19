package network

import (
	"fmt"
	"net"
	"time"
)

//默认tcp
func Listen(addr string) net.Listener {
	ln, err := net.Listen("tcp", addr)
	return &myListener{addr: addr, err: err, ln: ln}
}

//class listener
type myListener struct {
	addr string
	err  error
	ln   net.Listener
}

func (this *myListener) Accept() (conn net.Conn, err error) {
	if err = this.err; err != nil {
		return
	}
	var (
		ln           = this.ln
		tempDuration = time.Millisecond * 100
		maxDuration  = time.Second * 3
	)
process:
	conn, err = ln.Accept()
	if temp, ok := err.(net.Error); ok && temp.Temporary() {
		if tempDuration *= 2; tempDuration > maxDuration {
			tempDuration = maxDuration
		}
		fmt.Println("INFO: ", err.Error(), ", addr =", ln.Addr().String(), "[wait ", tempDuration, "]")
		time.Sleep(tempDuration)
		goto process
	}
	return
}

func (this *myListener) Close() (err error) {
	if err = this.err; err == nil {
		err = this.ln.Close()
	}
	return
}

func (this *myListener) Addr() net.Addr  { return this }
func (this *myListener) Network() string { return "tcp" }
func (this *myListener) String() string  { return this.addr }
