package network

import (
	"fmt"
	"net"
	"time"
)

type testAgent struct {
	Conn
}

func newAgent(conn Link) Runnable {
	return &testAgent{Conn: With(conn.(net.Conn))}
}

func (this *testAgent) Run() {
	fmt.Println("new connect:", this.Address())
	for {
		this.SetReadTimeout(time.Second)
		body, err := this.Read()
		if err == nil {
			fmt.Println("读取:", ReadBegin(body))
		} else {
			if CheckNetTemporary(err) {
				fmt.Println("超时或者其他错误:", err)
				this.Write(WriteBegin(0x11).Flush())
			} else {
				fmt.Println("其他错误", err)
				return
			}
		}
	}
}

type testClient struct {
	Conn
}

func newClient(conn Link) Runnable {
	//start
	return &testClient{Conn: With(conn.(net.Conn))}
}

func (this *testClient) Run() {
	fmt.Println("new link:", this.Address())
	for {
		body, err := this.Read()
		if err == nil {
			fmt.Println("客户端:", ReadBegin(body))
		} else {
			if CheckNetTemporary(err) {
				fmt.Println("客户端超时: ", err)
			} else {
				fmt.Println("客户端错误: ", err)
				return
			}
		}
	}
}

func init() {
	//server
	OnHandler(newAgent, ":8084", ":8085", ":7088", ":7077")

	//开启多个服务，如果你只想开启一个，单独new
	tcp := NewServer()
	tcp.Start(":8084")
	tcp.Start(":8085")

	_, err := StartWebServer(":7088")
	fmt.Println("web err1:", err)
	_, err = StartWebServer(":7077")
	fmt.Println("web err1:", err)
	//dail
	OnHandler(newClient, "router", "gate")

	dialer := NewManager("router")
	tcp.Join(dialer)
	//
	conn, err := DialScan(TCP_LINK, "localhost:8084")
	//这里可以将conn包装后发给dialer
	if err == nil {
		dialer.Join(conn)
	}
	//
	dialer = NewManager("gate")
	tcp.Join(dialer)
	conn, err = DialScan(TCP_LINK, "localhost:8084")
	if err == nil {
		dialer.Join(conn)
	}

	time.Sleep(time.Second * 3)
	tcp.Close()
	time.Sleep(time.Second)
	fmt.Println("剧终")
}
