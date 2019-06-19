package network

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

/*协议部分=默认参数*/
const (
	NET_Paylen   = 4
	NET_Maxlen   = 1024 * 1024 * 5 //读包最大限制5mb linux: cat /proc/sys/net/core/rmem_max
	NET_Dialtime = time.Second * 5
)

//封装net.Conn(协议处理后的read和write)
type Conn interface {
	Close() error
	Read() ([]byte, error)
	Write([]byte) error
	//func
	Address() string      //远端地址
	LocalAddress() string //本地地址
	SetReadTimeout(time.Duration) error
	SetSendTimeout(time.Duration) error
}

//new
func With(conn net.Conn) Conn {
	return &myConn{
		rwc:  conn,
		wbuf: bufio.NewWriter(conn),
		rbuf: bufio.NewReader(conn),
	}
}

func WithAddr(addr string) Conn {
	conn, err := net.DialTimeout("tcp", addr, NET_Dialtime)
	if err == nil {
		return With(conn)
	}
	return &errorConn{err: fmt.Errorf("Dial Err:" + err.Error()), addr: addr}
}

//协议封装(不适合UDP)
type myConn struct {
	rwc  net.Conn
	wbuf *bufio.Writer
	rbuf *bufio.Reader
}

func (this *myConn) Read() (body []byte, err error) {
	//read paylen (multithreading unsafe)
	var lenData [NET_Paylen]byte
	_, err = io.ReadFull(this.rbuf, lenData[0:])
	if err == nil {
		//big endian
		n := binary.BigEndian.Uint32(lenData[0:])
		//empty or max full : throw error
		if n > NET_Maxlen || n < 1 {
			err = fmt.Errorf("read body len big: len=%d max=%d", n, NET_Maxlen)
		} else {
			//new body (can get pool)
			body = make([]byte, NET_Paylen+n)
			//read body
			_, err = io.ReadFull(this.rbuf, body[NET_Paylen:])
			if err == nil {
				//write paylen
				copy(body, lenData[0:])
			}
		}
	}
	return
}

func (this *myConn) Write(b []byte) (err error) {
	if _, err = this.wbuf.Write(b); err == nil {
		err = this.wbuf.Flush()
	}
	return
}

//conn
func (this *myConn) Close() error {
	return this.rwc.Close()
}

func (this *myConn) SetReadTimeout(timeout time.Duration) (err error) {
	if timeout > 0 {
		err = this.rwc.SetReadDeadline(time.Now().Add(timeout))
	} else {
		err = this.rwc.SetReadDeadline(time.Time{})
	}
	return
}

func (this *myConn) SetSendTimeout(timeout time.Duration) (err error) {
	if timeout > 0 {
		err = this.rwc.SetWriteDeadline(time.Now().Add(timeout))
	} else {
		err = this.rwc.SetWriteDeadline(time.Time{})
	}
	return
}

func (this *myConn) Address() string {
	return this.rwc.RemoteAddr().String()
}

func (this *myConn) LocalAddress() string {
	return this.rwc.LocalAddr().String()
}

//die conn
type errorConn struct {
	err  error
	addr string
}

//base
func (this *errorConn) Read() ([]byte, error) { return nil, this.err }
func (this *errorConn) Write([]byte) error    { return this.err }
func (this *errorConn) Close() error          { return this.err }

//other
func (this *errorConn) Address() string                    { return this.addr }
func (this *errorConn) LocalAddress() string               { return "undefined" }
func (this *errorConn) SetReadTimeout(time.Duration) error { return this.err }
func (this *errorConn) SetSendTimeout(time.Duration) error { return this.err }
