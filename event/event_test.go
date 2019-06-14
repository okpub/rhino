package event

import (
	"fmt"
	"time"
)

type TestEvent int

func (n TestEvent) Type() int {
	return int(n)
}

func (n TestEvent) Body() interface{} { return nil }

func init() {
	cb := New()

	cb.OnFunc(2, func() {
		fmt.Println("深恶;:")
	})

	cb.OnFunc(3, func() {
		fmt.Println("我是谁")
	})
	//cb.Off(3)
	cb.DispatchEvent(TestEvent(3))

	time.Sleep(time.Second)
}
