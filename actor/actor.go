package actor

//Actor生产商
type Producer func() Actor

//actor func
type ActorFunc func(ActorContext)

func (f ActorFunc) Receive(ctx ActorContext) { f(ctx) }

//actor
type Actor interface {
	Receive(ActorContext)
}

//ref
type ActorRef interface {
	//tell nothing
	Tell(interface{}) error
	//tell self with sender
	Request(interface{}, ActorRef) error
	//stop
	Close() error
}
