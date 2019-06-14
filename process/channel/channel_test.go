package channel

import (
	"fmt"
	"time"

	"github.com/okpub/rhino/process"
)

type testBroker struct {
	process.UntypeBroker
}

func (this *testBroker) p() {
	fmt.Println(1)
}

//test
func init() {
	mb := New(OptionPendingNum(100))
	mb.OnRegister(process.NewDefaultDispatcher(0), &testBroker{})
	mb.Start()
	mb.Post("我是谁")
	mb.Post("我是谁2")
	opts := mb.Options()
	opts.Post("我不是谁")

	mb.Close()
	time.Sleep(time.Millisecond * 100)
}
