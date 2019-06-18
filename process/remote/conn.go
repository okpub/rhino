package remote

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
	NET_Paylen = 4
	NET_Maxlen = 1024 * 1024 * 5 //读包最大限制5mb linux: cat /proc/sys/net/core/rmem_max
)

/*做一个通用的socket*/
type Stream interface {
	//io
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
	//other
	SetSendTimeout(time.Duration) error
	SetReadTimeout(time.Duration) error
	Address() string      //远端地址
	LocalAddress() string //本地的地址
}

//new
func With(conn net.Conn) Stream {
	return &myConn{
		rwc:  conn,
		wbuf: bufio.NewWriter(conn),
		rbuf: bufio.NewReader(conn),
	}
}

func WithAddr(addr string) Stream {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err == nil {
		return With(conn)
	}
	return &errorConn{err: fmt.Errorf("Dial Err:" + err.Error()), addr: addr}
}

//class conn (注意：这里读写是不安全的，目前考虑是单线去写)
type myConn struct {
	rwc  net.Conn
	wbuf *bufio.Writer
	rbuf *bufio.Reader
}

func (this *myConn) Read() (body []byte, err error) {
	var lenData [NET_Paylen]byte
	_, err = io.ReadFull(this.rbuf, lenData[0:])
	if err == nil {
		//big endian
		n := binary.BigEndian.Uint32(lenData[0:])
		//empty or max full : throw error
		if n > NET_Maxlen || n < 1 {
			err = fmt.Errorf("read body len big: len=%d max=%d", n, NET_Maxlen)
			//this.Close()
		} else {
			//new body (can get pool)
			body = make([]byte, NET_Paylen+n)
			//read body
			if _, err = io.ReadFull(this.rbuf, body[NET_Paylen:]); err == nil {
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
func (this *myConn) Close() error         { return this.rwc.Close() }
func (this *myConn) Address() string      { return this.rwc.RemoteAddr().String() }
func (this *myConn) LocalAddress() string { return this.rwc.LocalAddr().String() }

//class error conn
type errorConn struct {
	err  error
	addr string
}

func (this *errorConn) Read() ([]byte, error)              { return nil, this.err }
func (this *errorConn) Write([]byte) error                 { return this.err }
func (this *errorConn) Close() error                       { return this.err }
func (this *errorConn) Address() string                    { return this.addr }
func (this *errorConn) LocalAddress() string               { return "undefined" }
func (this *errorConn) SetReadTimeout(time.Duration) error { return this.err }
func (this *errorConn) SetSendTimeout(time.Duration) error { return this.err }
