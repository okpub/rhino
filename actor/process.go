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

func UnboundedLocal(pendingNum int) ProcessProducer {
	return func() ActorProcess { return &LocalProcess{MessageQueue: channel.MakeBuffer(pendingNum)} }
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

//一般为客户端连接
func WithRemoteAddr(addr string, args ...remote.Option) Option {
	return func(p *Options) {
		p.processer = func() ActorProcess {
			return &RemoteProcess{SocketProcess: remote.NewAddr(addr, args...)}
		}
	}
}

//一般为服务端连接
func WithRemoteStream(conn interface{}, args ...remote.Option) Option {
	return func(p *Options) {
		p.dispatcher = process.NewSyncDispatcher(0) //默认同步
		p.processer = func() ActorProcess {
			return &RemoteProcess{SocketProcess: remote.NewKeepActive(conn, args...)}
		}
	}
}
