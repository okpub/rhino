package event

import (
	"fmt"
	"time"
)

func init() {
	event := make(OberSet)

	event.Subscribe(func(data Publication) {
		fmt.Println(data.Message())
	}, 1)
	event.Publish(1)
	fmt.Println("123")
	time.Sleep(time.Second)
}
