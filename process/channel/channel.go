package channel

import (
	"fmt"

	"github.com/okpub/rhino/process"
)

var (
	OverfullErr = fmt.Errorf("channel overfull")
)

type Option func(*Mailbox)

type MessageQueue interface {
	process.Process
	Post(interface{}) error
}

//选项(选项可以提取出来，因为选项比较少，所以不单独出来)
func OptionPendingNum(pendingNum int) Option {
	return func(p *Mailbox) {
		p.pendingNum = pendingNum
	}
}

func OptionBuffer(buffer chan interface{}) Option {
	return func(p *Mailbox) {
		p.buffer = buffer //多次使用,可能会被替换
	}
}

func OptionBlocking() Option { //阻塞模式
	return func(p *Mailbox) {
		p.blocking = true
	}
}

func OptionNonBlocking() Option { //非阻塞模式
	return func(p *Mailbox) {
		p.blocking = false
	}
}

func OptionFiller(opts ...Option) Option {
	return func(p *Mailbox) {
		p.filler(opts...)
	}
}
