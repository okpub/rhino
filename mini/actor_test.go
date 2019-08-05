package mini

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var rootRef *ActorRef

func init() {
	fmt.Println("init")
	time.Sleep(time.Millisecond)

	root, cancel := context.WithCancel(context.Background())
	ref := &ActorRef{
		opts: DefaultActor,
	}
	rootRef = ref
	ref.Init(
		WithContext(root),
		Name("root"),
		OnStart(func() {
			fmt.Println("start :", ref.Options().Name)
		}),
		OnStop(func(err error) {
			fmt.Println("stop :", ref.Options().Name, "[", err, "]")
		}),
		Received(func(data interface{}) {
			fmt.Println("root recv:", data)
		}),
	)

	go func() {
		defer cancel()
		ref.Run()
	}()

	ref.Send("顶级")
	ref.Send("顶级1")
	ref.Send("顶级2")
	ref.Send("顶级3")

	test(ref)
	ref.Close()
	time.Sleep(time.Second * 3)
}

func test(parent *ActorRef) {
	child, cancel := parent.With()

	ref := &ActorRef{
		opts: DefaultActor,
	}
	ref.Init(
		WithContext(child),
		Name("child"),
		OnStart(func() {
			fmt.Println("start :", ref.Options().Name)
		}),
		OnStop(func(err error) {
			fmt.Println("stop :", ref.Options().Name, "[", err, "]")
		}),
		Received(func(data interface{}) {
			fmt.Println("child recv:", data)
		}),
	)
	go func() {
		defer cancel()
		ref.Run()
	}()
	ref.Send("我是谁1")
	ref.Send("我是谁2")
	ref.Send("我是谁3")
	ref.Send("我是谁4")
}

func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rootRef.Send("who's your daddy!")
	}
}
