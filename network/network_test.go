package network

import (
	"fmt"
	"net"
	"time"
)

func SvrHandler(conn net.Conn) (err error) {
	fd := With(conn)
	var (
		body []byte
	)
	fmt.Println("new connect:", fd.Address())
	for {
		body, err = fd.Read()
		if err == nil {
			fmt.Println("读取:", ReadBegin(body))
		} else {
			//if CheckNetTemporary(err) {
			//	fmt.Println("超时或者其他错误:", err)
			//this.Write(WriteBegin(0x11).Flush())
			fmt.Println("其他错误", err)
			return
		}
	}
}

type testClient struct {
	Conn
}

func (this *testClient) Run() {
	fmt.Println("new link:", this.Address())
	for {
		body, err := this.Read()
		if err == nil {
			fmt.Println("客户端:", ReadBegin(body))
		} else {
			//if CheckNetTemporary(err) {
			//	fmt.Println("客户端超时: ", err)
			//} else {
			fmt.Println("客户端错误: ", err)
			return
			//}
		}
	}
}

func init() {
	//开启多个服务，如果你只想开启一个，单独new
	tcp := NewServer(OptionHandler(SvrHandler))
	tcp.Start(":8084")
	tcp.Start(":8085")

	StartWebServer(":7088")
	//fmt.Println("web err1:", err)
	//_, err = StartWebServer(":7077")
	//fmt.Println("web err1:", err)

	time.Sleep(time.Millisecond)
	WithAddr("localhost:8084")
	WithAddr("localhost:8085")

	time.Sleep(time.Second * 3)
	tcp.Close()
	time.Sleep(time.Second)
	fmt.Println("剧终")
}
