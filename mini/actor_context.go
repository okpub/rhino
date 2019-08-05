package mini

import (
	"context"
)

type actorContextExtras struct {
}

//ctx
type ActorContext struct {
	parent *ActorContext
	*ActorRef
	//The real context
	extras *actorContextExtras
}

func newActorContext(parent *ActorContext, self *ActorRef) *ActorContext {
	return &ActorContext{
		parent:   parent,
		ActorRef: self,
	}
}

func (this *ActorContext) Self() *ActorRef {
	return this.ActorRef
}

func (this *ActorContext) Parent() *ActorRef {
	return this.parent.Self()
}

//handle message
func (this *ActorContext) handleMessage(ctx context.Context, data interface{}) {
	this.opts.Received(context.WithValue(ctx, "message", data))
}
