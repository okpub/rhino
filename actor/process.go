package actor

import (
	"fmt"

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
		panic(fmt.Errorf("remote tell type unable to resolve %+v", data))
	}
}

func (this *RemoteProcess) Stop(sender ActorRef) error {
	return this.Close()
}
