package event

import (
	"fmt"
	"time"
)

type intMsg int

func (n intMsg) Type() int            { return int(n) }
func (n intMsg) Message() interface{} { return n }

func init() {
	event := make(OberSet)

	event.On(func(data Event) {
		fmt.Println("消息", data.Message())
	}, 1)
	event.On(func(data Event) {
		fmt.Println("消息", data.Message())
	}, 1)
	sub1 := event.On(func(data Event) {
		fmt.Println("消息", data.Message())
	}, 1, 2, 4, 5)
	sub1.Unsubscribe()
	event.DispatchEvent(intMsg(4))
	time.Sleep(time.Second)
}
