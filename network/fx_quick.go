package network

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

/*
*close web func
 */
type CloseFunc func()

/*
* quick start web server
 */
func StartWebServer(addr string) (cancel CloseFunc, err error) {
	var (
		handle Handler
	)
	if handle, err = GetHandler(addr); err == nil {
		//http server
		ln := &http.Server{
			ReadHeaderTimeout: time.Second * 3,
			Addr:              addr,
			Handler:           websocket.Handler(func(conn *websocket.Conn) { handle(conn).Run() }),
		}
		//close func
		cancel = func() { ln.Close() }
		go func() {
			defer ln.Close()
			fmt.Println("start web server: addr =", addr)
			err := ln.ListenAndServe()
			fmt.Println("stop web server: err=", err)
		}()
	}
	return
}

/*
* quick start tcp server
 */
func StartTcpServer(addr string) (cancel CloseFunc, err error) {
	var (
		handle Handler
		ln     net.Listener
	)
	if handle, err = GetHandler(addr); err == nil {
		if ln, err = net.Listen("tcp", addr); err == nil {
			cancel = func() { ln.Close() }
			go func() {
				defer ln.Close()
				fmt.Println("start tcp server: addr =", addr)
				for {
					conn, err := ln.Accept()
					if err == nil {
						go func(link net.Conn) {
							defer link.Close()
							handle(link).Run()
						}(conn)
					} else {
						if temp, ok := err.(net.Error); ok && temp.Temporary() {
							fmt.Println("warning tcp server: err=", err, "[wait 200ms]")
							time.Sleep(time.Millisecond * 200)
						} else {
							fmt.Println("stop tcp server: err=", err)
							break
						}
					}
				}
			}()
		}
	}
	return
}
