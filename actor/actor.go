package actor

//Actor producer
type Producer func() Actor

func ExchangeProducer(obj ActorFunc) Producer { return func() Actor { return obj } }

//Actor func
type ActorFunc func(ActorContext)

func (f ActorFunc) Receive(ctx ActorContext) { f(ctx) }

//actor
type Actor interface {
	Receive(ActorContext)
}

//proxy
type ActorRef interface {
	//tell nothing
	Tell(interface{}) error
	//tell self with sender
	Request(interface{}, ActorRef) error
	//stop
	Close() error
}
