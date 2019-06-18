package actor

import (
	"github.com/okpub/rhino/process/channel"
	"github.com/okpub/rhino/process/remote"
)

func WithActor(producer Producer, args ...channel.Option) *Options {
	return &Options{
		producer: producer,
		processer: func() ActorProcess {
			return &LocalProcess{MessageQueue: channel.New(args...)}
		},
	}
}

func WithFunc(actor ActorFunc, args ...channel.Option) *Options {
	return WithActor(func() Actor { return actor }, args...)
}

//一般为客户端连接
func WithRemoteFunc(actor ActorFunc, dial func() remote.Stream) *Options {
	return &Options{
		producer: func() Actor { return actor },
		processer: func() ActorProcess {
			return &RemoteProcess{SocketProcess: remote.New(remote.OptionFunc(dial))}
		},
	}
}

func WithRemoteAddr(actor ActorFunc, addr string) *Options {
	return &Options{
		producer: func() Actor { return actor },
		processer: func() ActorProcess {
			return &RemoteProcess{SocketProcess: remote.New(remote.OptionAddr(addr))}
		},
	}
}

//一般为服务端连接(同步阻塞)
func WithStream(producer Producer, stream remote.Stream) *Options {
	return &Options{
		producer:   producer,
		dispatcher: SyncDispatcher,
		processer: func() ActorProcess {
			return &RemoteProcess{SocketProcess: remote.NewKeepActive(remote.OptionStream(stream))}
		},
	}
}
