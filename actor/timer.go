package actor

import (
	"time"
)

type Timer interface {
	Once(time.Duration, ...interface{})
	Loop(time.Duration, ...interface{})
	Start(time.Duration, int, ...interface{})
	Running() bool
	Num() uint64
	Stop()
}

//timer
func NewTimer(ctx ActorContext) Timer {
	return &ActorTimer{ActorContext: ctx}
}

type ActorTimer struct {
	ActorContext
	id       uint64
	openFlag bool
	exit     chan struct{}
}

func (this *ActorTimer) next() uint64 {
	this.id++
	return this.id
}

func (this *ActorTimer) Once(delay time.Duration, args ...interface{}) {
	this.Start(delay, 1, args...)
}

func (this *ActorTimer) Loop(delay time.Duration, args ...interface{}) {
	this.Start(delay, -1, args...)
}

func (this *ActorTimer) Start(delay time.Duration, count int, args ...interface{}) {
	if this.openFlag {
		close(this.exit)
	}
	this.openFlag = true
	this.exit = make(chan struct{})
	//go thread
	event := &timerCtx{
		Timer: this,
		id:    this.next(),
		count: count,
		args:  args,
	}
	go setTimeout(this.exit, delay, count, func() { this.Self().Tell(event.schedule) })
}

func (this *ActorTimer) Stop() {
	if this.openFlag {
		this.openFlag = false
		close(this.exit)
	}
}

func (this *ActorTimer) Running() bool { return this.openFlag }
func (this *ActorTimer) Num() uint64   { return this.id }

//class timerCtx
type timerCtx struct {
	Timer
	id    uint64
	count int
	args  []interface{}
}

func (this *timerCtx) schedule() {
	if this.Running() && this.id == this.Num() {
		//stopped
		if !timerActive(&this.count) {
			this.Stop()
		}
		/*注意回调方式, 必须要有回调，不然报错*/
		args := this.args
		switch cb := args[0]; fn := cb.(type) {
		case func():
			fn()
		case func(int):
			fn(args[1].(int))
		case func(string):
			fn(args[1].(string))
		case func(interface{}):
			fn(args[1])
		default:
			panic("actor timer cna't callback")
		}
	}
}

//private static
func setTimeout(exit <-chan struct{}, delay time.Duration, count int, caller func()) {
	timer := time.NewTicker(delay)
	defer timer.Stop()
	for running := true; running; running = timerActive(&count) {
		select {
		case <-exit:
			return
		case <-timer.C:
			caller()
		}
	}
}

func timerActive(count *int) (active bool) {
	if *count == 0 {
		active = false
	} else if *count < 0 {
		active = true
	} else {
		*count--
		active = *count > 0
	}
	return
}
