package channel

import (
	"fmt"
	"testing"
	"time"

	"github.com/okpub/rhino/process"
)

type testBroker struct {
	process.UntypeBroker
}

func (*testBroker) DispatchMessage(data interface{}) {

}

//test
func init() {
	fmt.Println("init")
	time.Sleep(time.Millisecond * 100)
}

func BenchmarkTest(b *testing.B) {
	mb := New(OptionPendingNum(100))
	mb.OnRegister(process.NewDefaultDispatcher(0), &testBroker{})
	mb.Start()
	for i := 0; i < b.N; i++ {
		mb.Post("我是谁")
	}
	mb.Close()
}
