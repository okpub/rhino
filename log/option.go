package log

type Option func(*Subscription)

func newOptions(owner *eventStream, args ...Option) Subscriber {
	this := &Subscription{owner: owner}
	for _, o := range args {
		o(this)
	}
	return this
}

//class Subscription
type Subscription struct {
	owner *eventStream
	fn    func(Event)
	lv    Level
}

func (this *Subscription) Notify(event Event) {
	if event.Level >= this.lv {
		this.fn(event)
	}
}

func (this *Subscription) Unsubscribe() {
	this.owner.Unsubscribe(this)
}

//options
func OptionEvent(fn func(Event)) Option {
	return func(p *Subscription) {
		p.fn = fn
	}
}

func OptionLevel(lv Level) Option {
	return func(p *Subscription) {
		p.lv = lv
	}
}
