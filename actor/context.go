package actor

//actor context
type ActorContext interface {
	Context
	//base
	basePart
	//info
	infoPart
	//tell
	sendPart
	//recv
	recvPart
	//message
	messagePart
	//spawn
	spawnPart
}

//actor system
type ActorSystem interface {
	SpawnerContext
	Shutdown()
	Wait()
}

//context
type (
	SpawnerContext interface {
		Context
		infoPart
		spawnPart
	}

	SenderContext interface {
		Context
		infoPart
		messagePart
		sendPart
	}

	ReceiverContext interface {
		Context
		infoPart
		messagePart
		recvPart
	}
)

//子往父发消息堵塞，父往子发消息不能堵塞(注意环路堵塞)
type (
	basePart interface {
		//Watch(ActorRef)
		//UnWatch(ActorRef)
	}

	infoPart interface {
		Actor() Actor
		Self() ActorRef
		Parent() ActorRef
		System() ActorSystem
	}

	sendPart interface {
		Sender() ActorRef                    //当前发送者
		Forward(ActorRef) error              //转发
		Request(ActorRef, interface{}) error //请求对方
		Respond(interface{}) error           //回复消息(sender不存在的时候失败)
		Bubble(interface{}) error            //冒泡(报告上级)
		Refuse() error                       //关闭当前sender(回绝)
	}

	recvPart interface {
		Receive(MessageEnvelope)
	}

	messagePart interface {
		Any() interface{}
	}

	spawnPart interface {
		ActorOf(*Options) ActorRef
		Stop(ActorRef)
	}
)
