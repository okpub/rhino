package actor

import (
	"github.com/okpub/rhino/process"
	"github.com/okpub/rhino/process/channel"
	"github.com/okpub/rhino/process/remote"
)

func WithActor(producer Producer) *Options {
	return &Options{
		producer:  producer,
		processer: LocalUnbounded(channel.OptionPendingNum(100)),
	}
}

func ActorWithFunc(actor ActorFunc) *Options {
	return WithActor(func() Actor { return actor })
}

func ActorWithStream(producer Producer, stream remote.Stream) *Options {
	return &Options{
		producer:   producer,
		dispatcher: process.NewSyncDispatcher(0),
		processer:  RemoteUnbounded(remote.OptionWithStream(stream)),
	}
}

func ActorWithRemoteFunc(actor ActorFunc, dial func() remote.Stream) *Options {
	return &Options{
		producer:  func() Actor { return actor },
		processer: RemoteUnbounded(remote.OptionWithFunc(dial)),
	}
}
