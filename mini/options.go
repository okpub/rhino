package mini

import (
	"context"
)

var (
	DefaultActor = Options{
		Context:    context.Background(),
		Name:       "test",
		PendingNum: 100,
		OnStart:    func() {},
		OnStop:     func(err error) {},
		Received:   func(data interface{}) {}, //fmt.Println("miss handle:", data) },
	}
)

type Option func(*Options)

type Options struct {
	context.Context
	PendingNum int
	Name       string
	OnStart    func()
	OnStop     func(error)
	Received   func(interface{})
}

//options
func WithContext(ctx context.Context) Option {
	return func(p *Options) {
		p.Context = ctx
	}
}

func Name(name string) Option {
	return func(p *Options) {
		p.Name = name
	}
}

func PendingNum(n int) Option {
	return func(p *Options) {
		p.PendingNum = n
	}
}

func OnStart(fn func()) Option {
	return func(p *Options) {
		p.OnStart = fn
	}
}

func Received(fn func(interface{})) Option {
	return func(p *Options) {
		p.Received = fn
	}
}

func OnStop(fn func(error)) Option {
	return func(p *Options) {
		p.OnStop = fn
	}
}

//send
var (
	DefaultPublisher = PublishOptions{
		Context: context.Background(),
	}
)

type PublishOption func(*PublishOptions)

type PublishOptions struct {
	context.Context
	//Sender ActorRef
}
