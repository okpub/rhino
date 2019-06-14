package actor

func NewBroker(ctx *actorContext, p ActorProcess) *ActorBroker {
	return &ActorBroker{actorContext: ctx, ActorProcess: p}
}

type ActorBroker struct {
	*actorContext
	ActorProcess
}

func (this *ActorBroker) ref() ActorProcess {
	return this.ActorProcess
}

func (this *ActorBroker) Tell(data interface{}) error {
	return this.ref().SendMessage(this, data)
}

func (this *ActorBroker) Request(data interface{}, sender ActorRef) error {
	return this.Tell(MSG(sender, data))
}

func (this *ActorBroker) Close() error {
	return this.removeSelf(true, Canceled)
}

//private
func (this *ActorBroker) removeSelf(removeFromParent bool, code error) (err error) {
	if err = this.PIDGroup.removeSelf(false, code); err == nil {
		this.ref().Stop(this)
	}
	if removeFromParent {
		Fire(this.Context, this)
	}
	return
}
