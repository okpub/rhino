package remote

import (
	"time"

	"github.com/okpub/rhino/process"
)

type Option func(*Options)

type SocketProcess interface {
	process.Process
	Options() Options
	Read() ([]byte, error)
	Send([]byte) error
}

/*适用于tcp*/
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
