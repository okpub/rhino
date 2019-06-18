package event

import (
	"fmt"
)

type (
	Handler func(Publication)

	EventStream interface {
		Publish(int, ...interface{}) error
		Subscribe(Handler, ...int) Subscriber
	}

	Publisher interface {
		Publish(int, ...interface{}) error
	}

	Subscriber interface {
		Unsubscribe()
		Topics() []int
		Notify(Publication)
	}

	Publication interface {
		Topic() int
		Message() interface{}
	}
)

//class ObserSet
type OberSet map[int]ArraySubscription

func (hset OberSet) Publish(topic int, args ...interface{}) (err error) {
	if arr, ok := hset[topic]; ok {
		for _, sub := range arr {
			sub.Notify(NewEvent(topic, args))
		}
	} else {
		err = fmt.Errorf("can't find topic=%d", topic)
	}
	return
}

func (hset OberSet) Subscribe(fn Handler, args ...int) Subscriber {
	sub := &Subscription{owner: hset, topics: args, fn: fn}
	for _, topic := range args {
		hset[topic] = append(hset[topic], sub)
	}
	return sub
}

func (hset OberSet) Unsubscribe(sub Subscriber) {
	for _, topic := range sub.Topics() {
		if arr, ok := hset[topic]; ok {
			hset[topic] = arr.RemoveIndex(arr.IndexOf(sub))
			if len(hset[topic]) == 0 {
				delete(hset, topic)
			}
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

func (arr ArraySubscription) Copy() (list ArraySubscription) {
	list = make(ArraySubscription, len(arr))
	copy(list, arr)
	return
}

//class Subscription
type Subscription struct {
	owner  OberSet
	topics []int
	fn     Handler
}

func (this *Subscription) Notify(event Publication) {
	this.fn(event)
}

func (this *Subscription) Topics() []int {
	return this.topics
}

func (this *Subscription) Unsubscribe() {
	if this.owner != nil {
		this.owner.Unsubscribe(this)
	}
}

//class Event
type Event struct {
	topic int
	args  []interface{} //message is array
}

func NewEvent(topic int, args []interface{}) Publication {
	return &Event{topic: topic, args: args}
}

func (this *Event) Topic() int           { return this.topic }
func (this *Event) Message() interface{} { return this.args }
