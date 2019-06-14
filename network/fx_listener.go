package network

import (
	"net"
)

//默认tcp
func NewListener(addr string) Listener {
	ln, err := net.Listen("tcp", addr)
	return &myListener{addr: addr, err: err, ln: ln}
}

//class listener
type myListener struct {
	addr string
	err  error
	ln   net.Listener
}

func (this *myListener) Accept() (Link, error) {
	if this.err == nil {
		return this.ln.Accept()
	}
	return nil, this.err
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
