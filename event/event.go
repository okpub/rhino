package event

type (
	//事件
	Event interface {
		Type() int
		Body() interface{}
	}

	//派送
	IDispatcher interface {
		DispatchEvent(Event) error
	}

	//订阅
	Subscriber interface {
		IDispatcher
		Topics() []int //订阅的事件
	}

	IEventDispatcher interface {
		IDispatcher
		//other
		OnFunc(int, ...interface{}) IDispatcher
		//event
		On(int, IDispatcher) IDispatcher
		Off(int, ...IDispatcher)
		//sub
		AddSubscriber(Subscriber)
		UnSubscriber(Subscriber)
	}
)

//观察同一调度只能被调度一次
type (
	ObserSet map[IDispatcher]struct{}
	EventSet map[int]ObserSet
)

func New() IEventDispatcher {
	return make(EventDispatcher)
}

//消息派发(单线程内)
type EventDispatcher EventSet

func (events EventDispatcher) OnFunc(cmd int, args ...interface{}) IDispatcher {
	return events.On(cmd, NewHandler(args...))
}

func (events EventDispatcher) On(cmd int, target IDispatcher) IDispatcher {
	obsers, ok := events[cmd]
	if !ok {
		obsers = make(ObserSet)
		events[cmd] = obsers
	}
	obsers[target] = struct{}{}
	return target
}

func (events EventDispatcher) Off(cmd int, args ...IDispatcher) {
	if len(args) == 0 {
		delete(events, cmd)
	} else {
		if obser, ok := events[cmd]; ok {
			for _, p := range args {
				delete(obser, p)
			}
			if len(obser) == 0 {
				delete(events, cmd)
			}
		}
	}
}

//订阅消息
func (events EventDispatcher) AddSubscriber(sbr Subscriber) {
	for _, cmd := range sbr.Topics() {
		events.On(cmd, sbr)
	}
}

func (events EventDispatcher) UnSubscriber(sbr Subscriber) {
	for _, cmd := range sbr.Topics() {
		events.Off(cmd, sbr)
	}
}

//由于map遍历是乱序随机的，调度不会顺序执行
func (events EventDispatcher) DispatchEvent(event Event) (err error) {
	if obser, ok := events[event.Type()]; ok {
		for p := range obser {
			p.DispatchEvent(event)
		}
	}
	return
}
