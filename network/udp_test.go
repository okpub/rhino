package network

import (
	"fmt"
	"net"
	"testing"
)

/*
* udp本来就是为了快速通信
* 为了尽量避免丢包，需要做一些基本处理，但是如果做的处理太多，那和tcp无区别那还不如用tcp
* 1 控制发包速度: 比如1ms控制一次发包 1次最多(1024*10)bit(一般一个addr为2Wbit, 一个端口缓存区为20W+)
		假设1个包大小为50b,那么1ms能发200个包(200个包会分10次发送)，100ms能发200*100个包
* 2 因为1个gate不会多余1w人，所以保证100ms能派送1w*50bit就能畅通50Wbit,1ms需要发送5000bit
*
*/
type UDPServer struct{}

func (this *UDPServer) Start(addr string) error {
	go func() {
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err == nil {
			ln, err := net.ListenUDP("udp", udpAddr)
			if err == nil {
				this.Serve(ln)
			} else {
				fmt.Println("UDP listen err:", err, " addr=", addr)
			}
		} else {
			fmt.Println("UDP failed err:", err, " addr=", addr)
		}
	}()
	return nil
}

//udp就是为了快速发包而存在，所以不要过多的封装
func (this *UDPServer) Serve(ln *net.UDPConn) (err error) {
	defer ln.Close()
	body := make([]byte, 1024)
	for {
		n, addr, err := ln.ReadFromUDP(body)
		if err == nil {
			pack := ReadBegin(body[:n]) //read packing
			ln.WriteToUDP(WriteCopy(pack, 1, 2).Flush(), addr)
		} else {
			fmt.Println("udp read err:", err, addr)
		}
	}
}

func init() {
	udp := &UDPServer{}
	udp.Start(":8189")
}

//test
func BenchmarkUDPTest(b *testing.B) {
	conn, err := DialScan(UDP_LINK, "localhost:8189")
	if err != nil {
		return
	}
	var str string
	for i := 0; i < 1000; i++ {
		str += "a"
	}

	b.RunParallel(func(p *testing.PB) {
		var i int
		body := make([]byte, 1024)
		for p.Next() {
			i++
			_, err := conn.Write(WriteSync(int64(i), 0x01, 1, str).Flush())
			if err == nil {
				n, _ := conn.Read(body)
				ReadBegin(body[:n])
				//fmt.Println("客户端收到消息:", ReadBegin(body[:n]))
			} else {
				fmt.Println("send err: ", err)
			}
		}
	})
}
