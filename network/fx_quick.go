package network

import (
	"fmt"
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
func StartWebServer(addr string, args ...Option) (cancel CloseFunc, err error) {
	options := NewOptions(args...)
	//http server
	ln := &http.Server{
		ReadHeaderTimeout: time.Second * 3,
		Addr:              addr,
		Handler:           websocket.Handler(func(conn *websocket.Conn) { options.Handler(conn) }),
	}
	//stop func
	cancel = func() { ln.Close() }
	go func() {
		defer ln.Close()
		fmt.Println("start web server: addr =", addr)
		err := ln.ListenAndServe()
		fmt.Println("stop web server: err=", err)
		options.OnClose(ln, err)
	}()
	return
}

/*
* quick start tcp server
 */
func StartTcpServer(addr string, args ...Option) (cancel CloseFunc, err error) {
	ln := NewServer(args...)
	err = ln.Start(addr)
	cancel = func() { ln.Close() }
	return
}
