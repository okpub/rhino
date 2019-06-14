package channel

import (
	"fmt"
	"time"

	"github.com/rhino/process"
)

//test
func init() {
	mb := New(OptionPendingNum(100))
	mb.OnRegister(process.NewDefaultDispatcher(0), process.DoWithFunc(func(v interface{}) {
		fmt.Println("派送消息:", v)
	}))
	mb.Start()
	mb.Post("我是谁")
	mb.Post("我是谁2")
	opts := mb.Options()
	opts.Post("我不是谁")

	mb.Close()
	time.Sleep(time.Millisecond * 100)
}
