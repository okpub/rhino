package network

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

//默认数据
const (
	NET_Paylen = 4
	NET_Minlen = 4
	NET_Maxlen = 1024 * 1024 * 5
)

/*通用的socket*/
type Stream interface {
	//io
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
	//other
	SetSendTimeout(time.Duration) error
	SetReadTimeout(time.Duration) error
	//addr
	Address() string      //远端地址
	LocalAddress() string //本地的地址
}

//new
func With(conn net.Conn) Stream {
	return &netStream{
		rwc:      conn,
		newBytes: func(n int) []byte { return make([]byte, n) },
	}
}

func WithAddr(addr string) Stream {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err == nil {
		return With(conn)
	}
	return &errorConn{err: fmt.Errorf("Dial Err:" + err.Error()), addr: addr}
}

//class conn
type netStream struct {
	rwc        net.Conn
	ok         bool
	header     [NET_Paylen]byte
	readBuffer []byte //body
	readLen    int    //payload size
	readPos    int    //read pos
	//options
	newBytes func(int) []byte
}

//if you set readtimeout : unsafe read for other thread
func (this *netStream) Read() (body []byte, err error) {
	var n int
process:
	if this.ok {
		if n = this.readLen; n < NET_Minlen || n > NET_Maxlen {
			err = fmt.Errorf("read big pack: body=%d min=%d max=%d", n, NET_Minlen, NET_Maxlen)
		} else {
			n, err = this.rwc.Read(this.readBuffer[this.readPos:])
			if bodySize := len(this.readBuffer); err == nil {
				if this.readPos += n; this.readPos == bodySize {
					this.ok = false
					this.readPos = 0
					this.readLen = 0
					copy(this.readBuffer, this.header[0:])
					body, err = this.readBuffer, nil
				} else {
					if this.readPos > bodySize {
						panic(fmt.Errorf("An unknown error body %d>%d", this.readPos, bodySize))
					}
					goto process //read body
				}
			}
		}
	} else {
		n, err = this.rwc.Read(this.header[this.readPos:])
		if err == nil {
			if this.readPos += n; this.readPos == NET_Paylen {
				this.ok = true
				//big endian
				this.readLen = int(binary.BigEndian.Uint32(this.header[0:]))
				//empty or max full : throw error
				if n = this.readLen; n < NET_Minlen || n > NET_Maxlen {
					err = fmt.Errorf("read big pack: body=%d min=%d max=%d", n, NET_Minlen, NET_Maxlen)
				} else {
					this.readBuffer = this.newBytes(this.readPos + this.readLen)
					goto process //read body
				}
			} else {
				if this.readPos > NET_Paylen {
					panic(fmt.Errorf("An unknown error header %d>%d", this.readPos, NET_Paylen))
				}
				goto process //read header
			}
		}
	}
	return
}

func (this *netStream) Write(b []byte) (err error) {
	_, err = this.rwc.Write(b)
	return
}

func (this *netStream) SetReadTimeout(timeout time.Duration) (err error) {
	if timeout > 0 {
		err = this.rwc.SetReadDeadline(time.Now().Add(timeout))
	} else {
		err = this.rwc.SetReadDeadline(time.Time{})
	}
	return
}

func (this *netStream) SetSendTimeout(timeout time.Duration) (err error) {
	if timeout > 0 {
		err = this.rwc.SetWriteDeadline(time.Now().Add(timeout))
	} else {
		err = this.rwc.SetWriteDeadline(time.Time{})
	}
	return
}
func (this *netStream) Close() error         { return this.rwc.Close() }
func (this *netStream) Address() string      { return this.rwc.RemoteAddr().String() }
func (this *netStream) LocalAddress() string { return this.rwc.LocalAddr().String() }

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
