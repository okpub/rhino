package remote

import (
	"fmt"
	"time"

	"github.com/okpub/rhino/network"
	"github.com/okpub/rhino/process"
)

//server
type testAgent struct {
	process.UntypeBroker
	self SocketProcess
}

func (this *testAgent) DispatchMessage(v interface{}) {
	switch d := v.(type) {
	case error:
		//read timeout
	case []byte:
		fmt.Println("服务端收到:", network.ReadBegin(d))
		this.self.Send(network.WriteBegin(0x012).Flush())
	}
}

//client
type testClient struct {
	process.UntypeBroker
	self SocketProcess
}

func (this *testClient) PostStop() {
	fmt.Println("客户端关闭:", this.self)
}

func (this *testClient) DispatchMessage(v interface{}) {
	switch d := v.(type) {
	case error:
	case []byte:
		fmt.Println("客户端收到:", network.ReadBegin(d))
		//this.self.Send(network.WriteBegin(0x012).Flush())
	}
	//this.self.Post(network.WriteBegin(0x01).Flush())
}

//test
func init() {
	//注册地址代理
	network.OnHandler(func(conn network.Link) network.Runnable {
		//new ref
		self := NewKeepActive(OptionStream(network.WithLink(conn)))
		self.OnRegister(process.NewSyncDispatcher(0), &testAgent{self: self})
		self.Start()
		return network.EmptyRunner(0)
	}, ":8088")
	//启动服务
	stop, _ := network.StartTcpServer(":8088")

	//客户端制造商
	producer := Unbounded(OptionAddr("127.0.0.1:8088"))
	//todo
	ref := producer()
	ref.OnRegister(process.NewDefaultDispatcher(0), &testClient{self: ref})
	ref.Start()
	ref.Send(network.WriteBegin(0x1).Flush())

	time.Sleep(time.Second * 3)
	stop()
	time.Sleep(time.Millisecond * 1)
}
