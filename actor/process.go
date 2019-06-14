package actor

import (
	"fmt"

	"github.com/rhino/process"
	"github.com/rhino/process/channel"
	"github.com/rhino/process/remote"
)

type ProcessProducer func() ActorProcess

type ActorProcess interface {
	process.Process
	SendMessage(ActorRef, interface{}) error
	Stop(ActorRef) error
}

func NewLocalProcess(opts ...channel.Option) ActorProcess {
	return &LocalProcess{MessageQueue: channel.New(opts...)}
}

func LocalUnbounded(opts ...channel.Option) ProcessProducer {
	return func() ActorProcess {
		return NewLocalProcess(opts...)
	}
}

//本地
type LocalProcess struct {
	channel.MessageQueue
}

func (this *LocalProcess) SendMessage(sender ActorRef, data interface{}) error {
	return this.Post(data)
}

func (this *LocalProcess) Stop(sender ActorRef) error {
	return this.Close()
}

//远程
func NewRemoteProcess(opts ...remote.Option) ActorProcess {
	return &RemoteProcess{SocketProcess: remote.New(opts...)}
}

func RemoteUnbounded(opts ...remote.Option) ProcessProducer {
	return func() ActorProcess {
		return &RemoteProcess{SocketProcess: remote.New(opts...)}
	}
}

type RemoteProcess struct {
	remote.SocketProcess
}

func (this *RemoteProcess) SendMessage(sender ActorRef, data interface{}) error {
	switch p := data.(type) {
	case []byte:
		return this.Send(p)
	case MessageEnvelope:
		return this.SendMessage(sender, p.Any())
	default:
		panic(fmt.Errorf("remote tell type Unable to resolve %+v", data))
	}
}

func (this *RemoteProcess) Stop(sender ActorRef) error {
	return this.Close()
}
