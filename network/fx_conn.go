package network

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
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
	return &durConn{
		c:    conn,
		wbuf: bufio.NewWriter(conn),
		rbuf: bufio.NewReader(conn),
	}
}

func WithLink(conn Link) Conn {
	return With(conn.(net.Conn))
}

func WithErr(conn net.Conn, err error) Conn {
	if err == nil {
		return With(conn)
	}
	return &emptyConn{error: err}
}

func WithAddr(addr string) Conn {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err == nil {
		return With(conn)
	}
	return &emptyConn{error: err}
}

/*协议部分=默认参数*/
const (
	NET_Paylen = 4
	NET_Maxlen = 1024 * 1024 * 5 //读包最大限制5mb linux: cat /proc/sys/net/core/rmem_max
)

//协议封装(不适合UDP)
type durConn struct {
	c    net.Conn
	wbuf *bufio.Writer
	rbuf *bufio.Reader
}

func (this *durConn) Read() (body []byte, err error) {
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

func (this *durConn) Write(b []byte) (err error) {
	if _, err = this.wbuf.Write(b); err == nil {
		err = this.wbuf.Flush()
	}
	return
}

//conn
func (this *durConn) Close() error {
	return this.c.Close()
}

func (this *durConn) SetReadTimeout(timeout time.Duration) (err error) {
	if timeout > 0 {
		err = this.c.SetReadDeadline(time.Now().Add(timeout))
	} else {
		err = this.c.SetReadDeadline(time.Time{})
	}
	return
}

func (this *durConn) SetSendTimeout(timeout time.Duration) (err error) {
	if timeout > 0 {
		err = this.c.SetWriteDeadline(time.Now().Add(timeout))
	} else {
		err = this.c.SetWriteDeadline(time.Time{})
	}
	return
}

func (this *durConn) Address() string {
	return this.c.RemoteAddr().String()
}

func (this *durConn) LocalAddress() string {
	return this.c.LocalAddr().String()
}

//die conn
type emptyConn struct {
	error
}

//base
func (this *emptyConn) Read() ([]byte, error) { return nil, this.error }
func (this *emptyConn) Write(_ []byte) error  { return this.error }
func (this *emptyConn) Close() error          { return this.error }

//other
func (this *emptyConn) Address() string                      { return "undefined" }
func (this *emptyConn) LocalAddress() string                 { return "undefined" }
func (this *emptyConn) SetReadTimeout(_ time.Duration) error { return this.error }
func (this *emptyConn) SetSendTimeout(_ time.Duration) error { return this.error }
