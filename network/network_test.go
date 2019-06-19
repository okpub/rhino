package network

import (
	"fmt"
	"net"
	"time"

	"github.com/okpub/rhino/bytes"
)

func init() {
	//开启多个服务，如果你只想开启一个，单独new
	tcp := NewServer(OptionHandler(func(conn net.Conn) (err error) {
		fd := With(conn)
		fd.SetReadTimeout(time.Millisecond)
		for {
			body, err := fd.Read()
			fd.SetReadTimeout(time.Second * 1)
			if err == nil {
				fmt.Println("完整包:", len(body))
			} else {
				if temp, ok := err.(net.Error); ok && temp.Temporary() {
					fmt.Println("time out:", err)
				} else {
					fmt.Println("close err:", err)
					break
				}
			}
		}
		return
	}))
	tcp.Start(":8084")
	//tcp.Start(":8085")

	//StartWebServer(":7088")
	//fmt.Println("web err1:", err)
	//_, err = StartWebServer(":7077")
	//fmt.Println("web err1:", err)

	time.Sleep(time.Millisecond)
	conn, _ := DialScan("", "localhost:8084")
	time.Sleep(time.Millisecond * 10)
	for i := 0; i < 10; i++ {
		b := bytes.New()
		b.Wint(100)
		b.Wint(3)
		conn.Write(b.Bytes())
		var bb [100 - 4]byte
		conn.Write(bb[0:])
	}
	time.Sleep(time.Second * 3)
	tcp.Close()
	time.Sleep(time.Second)
	fmt.Println("剧终")
}
