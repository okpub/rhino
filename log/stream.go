package log

//system
var stage = &eventStream{}

func Subscribe(fn func(Event)) Subscriber {
	return stage.Subscribe(OptionEvent(fn))
}

func Unsubscribe(sub Subscriber) {
	stage.Unsubscribe(sub)
}

type (
	Subscriber interface {
		Notify(Event)
		Unsubscribe()
	}
)

//class eventStream
type eventStream struct {
	ArraySubscription
}

func (this *eventStream) Subscribe(args ...Option) Subscriber {
	sub := newOptions(this, args...)
	this.ArraySubscription = append(this.ArraySubscription, sub)
	return sub
}

func (this *eventStream) Unsubscribe(sub Subscriber) {
	this.ArraySubscription = this.ArraySubscription.RemoveIndex(this.ArraySubscription.IndexOf(sub))
}

func (this *eventStream) Publish(evt Event) {
	//If in the process of scheduling is deleted, then it is not safe
	for _, sub := range this.ArraySubscription {
		sub.Notify(evt)
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
