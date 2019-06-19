package actor

import (
	"fmt"

	"github.com/okpub/rhino/process"
	"github.com/okpub/rhino/process/channel"
	"github.com/okpub/rhino/process/remote"
)

//本地ref
type LocalProcess struct {
	channel.MessageQueue
}

func (this *LocalProcess) SendMessage(sender ActorRef, data interface{}) error {
	return this.Post(data)
}

func (this *LocalProcess) Stop(sender ActorRef) error {
	return this.Close()
}

//远程ref
type RemoteProcess struct {
	remote.SocketProcess
}

func (this *RemoteProcess) SendMessage(sender ActorRef, data interface{}) error {
	switch b := data.(type) {
	case []byte:
		return this.Send(b)
	case MessageEnvelope:
		return this.SendMessage(sender, b.Any())
	default:
		panic(fmt.Errorf("remote tell type unable to resolve %T", data))
	}
}

func (this *RemoteProcess) Stop(sender ActorRef) error {
	return this.Close()
}

//邮箱
func WithActor(producer Producer, opts ...channel.Option) Option {
	return func(p *Options) {
		p.producer = producer
		p.processer = func() ActorProcess {
			return &LocalProcess{MessageQueue: channel.New(opts...)}
		}
	}
}

func WithFunc(fn func(ActorContext), opts ...channel.Option) Option {
	return WithActor(ExchangeProducer(fn), opts...)
}

//一般为客户端连接
func WithRemoteAddr(fn func(ActorContext), addr string, args ...remote.Option) Option {
	return func(p *Options) {
		p.producer = ExchangeProducer(fn)
		p.processer = func() ActorProcess {
			return &RemoteProcess{SocketProcess: remote.NewAddr(addr, args...)}
		}
	}
}

//一般为服务端连接(会阻塞当前线程)
func WithRemoteStream(fn func(ActorContext), conn interface{}, args ...remote.Option) Option {
	return func(p *Options) {
		p.producer = ExchangeProducer(fn)
		p.dispatcher = process.NewSyncDispatcher(0) //默认同步
		p.processer = func() ActorProcess {
			return &RemoteProcess{SocketProcess: remote.NewKeepActive(conn, args...)}
		}
	}
}
