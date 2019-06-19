package actor

import (
	"github.com/okpub/rhino/process"
)

type SpawnFunc func(SpawnerContext, *Options) ActorRef

type ContextDecorator func(next ContextDecoratorFunc) ContextDecoratorFunc
type ReceiverMiddleware func(next ReceiverFunc) ReceiverFunc
type SenderMiddleware func(next SenderFunc) SenderFunc
type SpawnMiddleware func(next SpawnFunc) SpawnFunc

//process interface
type ProcessProducer func() ActorProcess

type ActorProcess interface {
	process.Process
	SendMessage(ActorRef, interface{}) error
	Stop(ActorRef) error
}

// default values
var (
	defaultOptions = Options{
		dispatcher: process.NewDefaultDispatcher(0),
		spawner: func(parent SpawnerContext, opts *Options) ActorRef {
			ctx := newActorContext(parent, opts)
			//new process
			self := NewBroker(&ctx, opts.NewProcess())
			ctx.self = self
			//addchild
			Join(parent, self)
			//init and start
			self.OnRegister(opts.GetDispatcher(), &ctx)
			self.Start()
			return self
		},
	}
)

//option
type Option func(*Options)

//options
type Options struct {
	spawner    SpawnFunc
	producer   Producer
	processer  ProcessProducer
	dispatcher process.Dispatcher
	//guardianStrategy        SupervisorStrategy
	//supervisionStrategy     SupervisorStrategy
	receiverMiddleware      []ReceiverMiddleware
	receiverMiddlewareChain ReceiverFunc
	senderMiddleware        []SenderMiddleware
	senderMiddlewareChain   SenderFunc
	contextDecorator        []ContextDecorator
	contextDecoratorChain   ContextDecoratorFunc
	spawnMiddleware         []SpawnMiddleware
	spawnMiddlewareChain    SpawnFunc
}

//基于默认选项配置来做文章(如果你不用默认，那么自己重新定义)
func NewOptions(opts ...Option) *Options {
	return defaultOptions.Copy(opts...)
}

func (this *Options) Filler(args ...Option) *Options {
	for _, o := range args {
		o(this)
	}
	return this
}

func (this Options) Copy(args ...Option) *Options {
	return this.Filler(args...)
}

func (this *Options) NewActor() Actor {
	return this.producer()
}

func (this *Options) NewProcess() ActorProcess {
	return this.processer()
}

func (this *Options) GetDispatcher() process.Dispatcher {
	return this.dispatcher
}

func (this *Options) ContextWrapper(ctx ActorContext) ActorContext {
	if this.contextDecoratorChain == nil {
		return ctx
	}
	return this.contextDecoratorChain(ctx)
}

//private
func (this *Options) spawn(parent SpawnerContext) ActorRef {
	if spawner := this.spawnMiddlewareChain; spawner != nil {
		return spawner(parent, this)
	}
	return this.spawner(parent, this)
}

/*
*选项
 */
func OptionDispatcher(dispatcher process.Dispatcher) Option {
	return func(p *Options) {
		p.dispatcher = dispatcher
	}
}

func OptionSpawner(spawner SpawnFunc) Option {
	return func(p *Options) {
		p.spawner = spawner
	}
}

//middleware chain
func OptionContextMiddlewareChain(childs ...ContextDecorator) Option {
	return func(p *Options) {
		p.contextDecorator = append(p.contextDecorator, childs...)
		p.contextDecoratorChain = makeContextDecoratorChain(func(ctx ActorContext) ActorContext {
			return ctx
		}, p.contextDecorator...)
	}
}

func OptionReceiverMiddlewareChain(childs ...ReceiverMiddleware) Option {
	return func(p *Options) {
		p.receiverMiddleware = append(p.receiverMiddleware, childs...)
		p.receiverMiddlewareChain = makeReceiverMiddlewareChain(func(target ReceiverContext, message MessageEnvelope) {
			target.Receive(message)
		}, p.receiverMiddleware...)
	}
}

func OptionSenderMiddlewareChain(childs ...SenderMiddleware) Option {
	return func(p *Options) {
		p.senderMiddleware = append(p.senderMiddleware, childs...)
		p.senderMiddlewareChain = makeSenderMiddlewareChain(func(_ SenderContext, sender ActorRef, message MessageEnvelope) error {
			return sender.Tell(message)
		}, p.senderMiddleware...)
	}
}

func OptionSpawnMiddleware(childs ...SpawnMiddleware) Option {
	return func(p *Options) {
		p.spawnMiddleware = append(p.spawnMiddleware, childs...)
		p.spawnMiddlewareChain = makeSpawnMiddlewareChain(func(parent SpawnerContext, opts *Options) ActorRef {
			return opts.spawner(parent, opts)
		}, p.spawnMiddleware...)
	}
}
