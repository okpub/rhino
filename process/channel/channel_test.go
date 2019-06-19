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
	fmt.Println(data)
}

//test
func init() {
	fmt.Println("init")
	BenchmarkTest(&testing.B{N: 10})
	time.Sleep(time.Millisecond * 100)
}

func BenchmarkTest(b *testing.B) {
	mb := New(OptionNonBlocking())
	mb.OnRegister(process.NewDefaultDispatcher(0), &testBroker{})
	mb.Start()
	for i := 0; i < b.N; i++ {
		err := mb.Post("我是谁")
		if err != nil {
			fmt.Println("err:", err)
		}
	}
	mb.Close()
}
