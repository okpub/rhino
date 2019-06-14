package actor

//默认系统
var stageContext = &actorSystem{}

func Stage() ActorSystem {
	return stageContext
}

//class system(不选择copy那么不能作为context)
type actorSystem struct {
	UntypeContext
	PIDGroup
}

func (this *actorSystem) System() ActorSystem { return this }
func (this *actorSystem) Shutdown()           { this.removeSelf(false, Canceled) }
func (this *actorSystem) Wait()               { /*not used*/ }

func (this *actorSystem) ActorOf(opts *Options) ActorRef {
	return opts.spawn(this)
}
