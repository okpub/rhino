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

// Default values
var (
	SyncDispatcher          = process.NewSyncDispatcher(0)
	defaultDispatcher       = process.NewDefaultDispatcher(0)
	defaultContextDecorator = func(ctx ActorContext) ActorContext { return ctx }

	defaultSpawner = func(parent SpawnerContext, opts *Options) ActorRef {
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
	}
)

//option
type Option func(*Options)

//options(可以添加refer的一个装饰器)
type Options struct {
	producer   Producer
	dispatcher process.Dispatcher
	processer  ProcessProducer
	//guardianStrategy        SupervisorStrategy
	//supervisionStrategy     SupervisorStrategy
	receiverMiddleware      []ReceiverMiddleware
	receiverMiddlewareChain ReceiverFunc
	senderMiddleware        []SenderMiddleware
	senderMiddlewareChain   SenderFunc
	contextDecorator        []ContextDecorator
	contextDecoratorChain   ContextDecoratorFunc
	//spawner
	spawner              SpawnFunc
	spawnMiddleware      []SpawnMiddleware
	spawnMiddlewareChain SpawnFunc
}

func (this Options) Copy(args ...Option) *Options {
	for _, o := range args {
		o(&this)
	}
	return &this
}

func (this *Options) NewActor() Actor {
	return this.producer()
}

func (this *Options) NewProcess() ActorProcess {
	return this.processer()
}

func (this *Options) GetDispatcher() process.Dispatcher {
	if this.dispatcher == nil {
		return defaultDispatcher
	}
	return this.dispatcher
}

func (this *Options) ContextWrapper(ctx ActorContext) ActorContext {
	if this.contextDecoratorChain == nil {
		return defaultContextDecorator(ctx)
	}
	return this.contextDecoratorChain(ctx)
}

func (this *Options) GetSpawner() (spawner SpawnFunc) {
	if spawner = this.spawnMiddlewareChain; spawner == nil {
		if spawner = this.spawner; spawner == nil {
			spawner = defaultSpawner
		}
	}
	return
}

//private
func (this *Options) spawn(parent SpawnerContext) ActorRef {
	return this.GetSpawner()(parent, this)
}

/*
*选项
 */
func OptionFromFunc(fn func(ActorContext)) Option {
	return OptionProducer(func() Actor { return ActorFunc(fn) })
}

func OptionProducer(producer Producer) Option {
	return func(p *Options) {
		p.producer = producer
	}
}

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
			return defaultContextDecorator(ctx)
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
