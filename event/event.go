package event

type (
	Handler func(Event)

	Event interface {
		Publication
		Target() interface{}
	}

	EventBus interface {
		On(Handler, ...int) Subscriber
		Off(Subscriber)
		DispatchEvent(Publication)
	}

	Subscriber interface {
		Topics() []int
		Unsubscribe()
		DispatchEvent(Publication)
	}

	Publication interface {
		Type() int
		Message() interface{}
	}
)

//class ObserSet
type OberSet map[int]ArraySubscription

func (hset OberSet) On(method Handler, topics ...int) Subscriber {
	sub := &Subscription{parent: hset, caller: nil, method: method, topics: topics}
	for _, topic := range topics {
		hset[topic] = append(hset[topic], sub)
	}
	return sub
}

func (hset OberSet) Off(sub Subscriber) {
	for _, topic := range sub.Topics() {
		if arr, ok := hset[topic]; ok {
			hset[topic] = arr.RemoveIndex(arr.IndexOf(sub))
			if len(hset[topic]) == 0 {
				delete(hset, topic)
			}
		}
	}
}

func (hset OberSet) DispatchEvent(pub Publication) {
	if arr, ok := hset[pub.Type()]; ok {
		arr = arr.copy() //拷贝，避免调度的时候删除
		for _, obser := range arr {
			obser.DispatchEvent(pub)
		}
	}
}

//array Subscription
type ArraySubscription []Subscriber

func (arr ArraySubscription) IndexOf(p Subscriber) int {
	for i, v := range arr {
		if v == p {
			return i
		}
	}
	return -1
}

func (arr ArraySubscription) RemoveIndex(i int) ArraySubscription {
	if i != -1 {
		return append(arr[:i], arr[i+1:]...)
	}
	return arr
}

func (arr ArraySubscription) copy() (list ArraySubscription) {
	list = make(ArraySubscription, len(arr))
	copy(list, arr)
	return
}

//class Subscription
type Subscription struct {
	parent OberSet
	topics []int
	caller interface{}
	method func(Event)
}

func (this *Subscription) DispatchEvent(pub Publication) {
	this.method(&untypeEvent{target: this.caller, Publication: pub})
}

func (this *Subscription) Topics() []int {
	return this.topics
}

func (this *Subscription) Unsubscribe() {
	this.parent.Off(this)
}

type SubOption func(*Subscription)

func Caller(caller interface{}) SubOption {
	return func(p *Subscription) {
		p.caller = caller
	}
}

//class Event
type untypeEvent struct {
	target interface{}
	Publication
}

func (this *untypeEvent) Target() interface{} { return this.target }
