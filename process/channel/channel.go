package channel

import (
	"fmt"

	"github.com/okpub/rhino/process"
)

var (
	OverfullErr = fmt.Errorf("channel overfull")
)

type MessageQueue interface {
	process.Process
	Post(interface{}) error
}

type Option func(*Mailbox)

//options
func OptionBuffer(obj interface{}) Option {
	return func(p *Mailbox) { //PS: change buffer
		switch buf := obj.(type) {
		case chan interface{}:
			p.buffer = buf
		case func() chan interface{}:
			p.buffer = buf()
		case func(int) chan interface{}:
			p.buffer = buf(p.pendingNum)
		default:
			panic(fmt.Errorf("buffer paramer error: [class %T]", obj))
		}
	}
}

func OptionNum(n int) Option {
	return func(p *Mailbox) {
		p.pendingNum = n
	}
}

func OptionBlocking() Option { //blocking mode
	return func(p *Mailbox) {
		p.nonblocking = false
	}
}

func OptionNonBlocking() Option { //nonblocking mode
	return func(p *Mailbox) {
		p.nonblocking = true
	}
}
