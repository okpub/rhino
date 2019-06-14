package actor

type actorContextExtras struct {
	context ActorContext
}

func newActorContextExtras(context ActorContext) *actorContextExtras {
	return &actorContextExtras{
		context: context,
	}
}

func (this *actorContextExtras) Context() ActorContext { return this.context }

// lazy initialize the child restart stats if this is the first time
// further mutations are handled within "restart"
func (this *actorContextExtras) restartStats() {}

func (this *actorContextExtras) resetReceiveTimeoutTimer(int) {}

func (this *actorContextExtras) stopReceiveTimeoutTimer() {}

func (this *actorContextExtras) killReceiveTimeoutTimer() {}
