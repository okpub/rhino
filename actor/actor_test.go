package actor

//import (
//	"fmt"
//	"time"

//	"github.com/rhino/core"
//)

//type TestDecorator struct {
//	ActorContext
//}

//func (this *TestDecorator) String() string {
//	return "我不告诉你"
//}

//func init() {
//	ctx1 := func(next ContextDecoratorFunc) ContextDecoratorFunc {
//		return func(ctx ActorContext) ActorContext {
//			return next(&TestDecorator{ctx})
//		}
//	}

//	m1 := func(next SenderFunc) SenderFunc {
//		return func(ctx SenderContext, sender ActorRef, message MessageEnvelope) error {
//			fmt.Println("来啊", ctx)
//			//sender.Tell(message)
//			return next(ctx, sender, message)
//		}
//	}

//	ref := WithActor(Stage(), OptionFromFunc(func(ctx ActorContext) {
//		//		fmt.Println(ctx.Any())

//		//		childRef := WithActor(ctx, OptionFromFunc(func(child ActorContext) {
//		//			fmt.Println("child:", child.Any())
//		//		}))
//		//		//
//		//		ctx.Forward(childRef)
//	}), OptionSenderMiddlewareChain(m1), OptionContextMiddlewareChain(ctx1))
//	ref.Tell("我谁是")
//	fmt.Println("over", core.SizeTypeof(ref))
//	time.Sleep(time.Second)
//}

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/okpub/rhino/core"
	"github.com/okpub/rhino/network"
	"github.com/okpub/rhino/process/remote"
)

type BeginObj struct {
	add  bool
	id   int64
	self ActorRef
}

type TellObj struct {
	id   int64
	body []byte
}

type CloseObj struct {
}

type GateActor struct {
	fd_id  int64
	refs   map[int64]ActorRef
	router interface{}
}

func (this *GateActor) NextID() int64 {
	return atomic.AddInt64(&this.fd_id, 1)
}

func (this *GateActor) PreStart(ctx ActorContext) {
	this.refs = make(map[int64]ActorRef)
	this.start_server()
}

func (this *GateActor) Receive(ctx ActorContext) {
	switch o := ctx.Any().(type) {
	case *Started:
		this.PreStart(ctx)
	case *BeginObj:
		this.trackUser(ctx, o)
	case *TellObj:
		this.tellUser(ctx, o.id, o.body)
	case *network.SocketPacket:
		fmt.Println("网关接收:", o)
	default:
		fmt.Println("not handle:", ctx.Any())
	}
}

func (this *GateActor) trackUser(ctx ActorContext, obj *BeginObj) {
	if ref, ok := this.refs[obj.id]; ok {
		if obj.add {
			if ok {
				ctx.Refuse()
			} else {
				this.refs[obj.id] = obj.self
			}
		} else {
			if ref == obj.self {
				delete(this.refs, obj.id)
				ref.Close()
			}
		}
	} else {
		//not id
	}
}

func (this *GateActor) tellUser(ctx ActorContext, id int64, body []byte) (err error) {
	if ref, ok := this.refs[id]; ok {
		//blocking the closed after
		if err = ref.Tell(body); err == nil {
			ref.Close()
		}
	}
	return
}

func (this *GateActor) start_server() {
	network.OnHandler(func(conn network.Link) network.Runnable {
		Stage().ActorOf(WithStream(func() Actor {
			return &testAgent{}
		}, network.WithLink(conn)))
		return network.EmptyRunner(0)
	}, ":8088")
	//open
	network.StartTcpServer(":8088")
}

//server agent
type testAgent struct {
	router  interface{} //路由
	id      int64
	sendRef ActorRef
}

func (this *testAgent) init(ctx ActorContext) {
	this.sendRef = ctx.ActorOf(WithFunc(func(child ActorContext) {
		switch body := child.Any().(type) {
		case []byte:
			child.Respond(body)
			//ctx.Self().Tell(body) //给socket
		case *Stopped:
			ctx.Stop(ctx.Self()) //write close
		case *Started:
			//todo
		default:
			fmt.Println("can't handle:", body)
		}
	}))
	//注册网关自己的发送通道
	ctx.Bubble(&BeginObj{add: true, id: this.id, self: this.sendRef})
}

func (this *testAgent) Receive(ctx ActorContext) {
	switch b := ctx.Any().(type) {
	case []byte:
		fmt.Println("服务端收到:", network.ReadBegin(b))
		ctx.Respond(network.WriteBegin(0x03).Flush())
	case error:
		fmt.Println("心跳", b)
		ctx.Respond(network.WriteBegin(0x01).Flush())
	case *Started:
		fmt.Println("开始")
		this.init(ctx)
	case *Stopped:
		fmt.Println("结束")
		ctx.Bubble(&BeginObj{id: this.id, self: this.sendRef})
	}
}

func init() {
	gateRef := Stage().ActorOf(WithActor(func() Actor {
		return &GateActor{}
	}))
	fmt.Println("尺寸:", core.SizeTypeof(gateRef))
	//闭包
	func() {
		var (
			client remote.SocketProcess
			addr   = "localhost:8088"
		)
		cliRef := Stage().ActorOf(WithFunc(func(ctx ActorContext) {
			switch b := ctx.Any().(type) {
			case *Started:
			case *CloseObj:
				client = nil
			//reset
			case []byte:
				client = remote.New(remote.OptionWithFunc(func() remote.Stream {
					return network.WithErr(network.DialScan("", addr))
				}))
				client.OnRegister(defaultDispatcher, FuncBroker(func(data interface{}) {
					switch body := data.(type) {
					case []byte:
						fmt.Println("body", network.ReadBegin(body))
					}
				}))
				client.Start()
				client.Send(b)
			}
		}))
		//客户端(唯一丢包的可能就是socket断线重连)
		cliRef.Tell(network.WriteBegin(0x2).Flush())
	}()

	time.Sleep(time.Second * 1)
	//gateRef.Tell(&TellObj{id: 1, body: network.WriteBegin(0x03).Flush()})
	time.Sleep(time.Second * 2)
	Stage().Shutdown()

	Stage().ActorOf(WithFunc(func(ctx ActorContext) {
		fmt.Println("what")
	}))
	time.Sleep(1000)
}
