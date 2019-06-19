package actor

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/okpub/rhino/core"
	"github.com/okpub/rhino/network"
	"github.com/okpub/rhino/process"
	"github.com/okpub/rhino/process/remote"
)

type BeginObj struct {
	Add bool
}

var fd_id int64

func NextID() int64 {
	return atomic.AddInt64(&fd_id, 1)
}

type GateActor struct {
	refs   map[int64]ActorRef
	router interface{}
}

func (this *GateActor) PreStart(ctx ActorContext) {
	this.refs = make(map[int64]ActorRef)
}

func (this *GateActor) Receive(ctx ActorContext) {
	switch o := ctx.Any().(type) {
	case *Started:
		this.PreStart(ctx)
	case *BeginObj:
		this.trackUser(ctx, o)
	case *network.SocketPacket:
		//this.tellUser(ctx, o.id, o.body)
		fmt.Println("网关接收:", o)
	default:
		fmt.Println("not handle:", ctx.Any())
	}
}

func (this *GateActor) trackUser(ctx ActorContext, obj *BeginObj) {
	var (
		session = ctx.Sender().(*Session)
		id      = session.id
		added   = obj.Add
	)

	if ref, ok := this.refs[id]; ok {
		if added {
			if ok {
				ctx.Refuse()
			} else {
				this.refs[id] = session
			}
		} else {
			if ref == session {
				delete(this.refs, id)
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

//server agent
type Session struct {
	id int64
	ActorRef
}

func NewSession(conn remote.SocketProcess) *Session {
	return &Session{ //会话自身是无法修改自己的
		id: NextID(),
		ActorRef: Stage().ActorOf(WithFunc(func(ctx ActorContext) {
			switch body := ctx.Any().(type) {
			case []byte:
				conn.Send(body)
			case *Started:
				//on
			case *Stopped:
				conn.Close()
			default:
				fmt.Printf("miss session handle [class %T] \n", body)
			}
		})),
	}
}

func init() {
	gateRef := Stage().ActorOf(WithActor(func() Actor {
		return &GateActor{}
	}))
	fmt.Println("尺寸:", core.SizeTypeof(gateRef))
	network.OnHandler(func(conn net.Conn) error {
		var (
			session *Session
			ref     = remote.NewKeepActive(conn, remote.OptionNonPing()) //去掉心跳
		)
		ref.OnRegister(process.NewSyncDispatcher(0), DoFunc(func(any interface{}) {
			switch body := any.(type) {
			case []byte:
				gateRef.Request(network.ReadBegin(body), session)
				//fmt.Println("服务端收到:", network.ReadBegin(body))
			case error:
				ref.Send(network.WriteBegin(0x01).Flush())
			case *Started:
				//建立会话
				session = NewSession(ref)
				gateRef.Request(&BeginObj{Add: true}, session)
			case *Stopped:
				//移除会话
				session.Close()
				gateRef.Request(&BeginObj{}, session)
			default:
				fmt.Printf("miss handle [class %T] \n", any)
			}
		}))
		return ref.Start()
	}, ":8088")
	//open
	network.StartTcpServer(":8088")
	//闭包
	func() {
		var (
			client ActorRef
			addr   = "localhost:8088"
		)
		cliRef := Stage().ActorOf(WithFunc(func(ctx ActorContext) {
			switch b := ctx.Any().(type) {
			case *Started:
			case bool:
				client = nil
			//reset
			case []byte:
				client = ctx.ActorOf(WithRemoteAddr(addr), WithFunc(func(child ActorContext) {
					switch body := child.Any().(type) {
					case []byte:
						fmt.Println("body", network.ReadBegin(body))
					}
				}))
				client.Tell(b)
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
