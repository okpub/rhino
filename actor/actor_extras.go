package actor

type actorContextExtras struct {
	//childs  map[ActorRef]struct{}
	context ActorContext
}

func newActorContextExtras(context ActorContext) *actorContextExtras {
	return &actorContextExtras{
		context: context,
		//childs:  make(map[ActorRef]struct{}),
	}
}

func (this *actorContextExtras) Context() ActorContext {
	return this.context
}

func (this *actorContextExtras) restartStats() {
	// lazy initialize the child restart stats if this is the first time
	// further mutations are handled within "restart"
}

func (this *actorContextExtras) resetReceiveTimeoutTimer(int) {}

func (this *actorContextExtras) stopReceiveTimeoutTimer() {}

func (this *actorContextExtras) killReceiveTimeoutTimer() {}

func (this *actorContextExtras) stopAllChildren() {}
