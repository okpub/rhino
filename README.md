# rhino
go网络基础框架
# 前沿


# network
服务端: 通过设置SetReadTimeout来获得心跳，而无须在上层做timer，可以通过test查看例子

客户端：保持默认设置即可

network.stream封装了len+body的读取方法，由于跨协程read不安全所以socket一般不会直接面向应用层逻辑，借助actor模型 由actor对其控制，来更好的解决和应对分布式中环路堵塞和短暂丢包的处理

# acotr模型
目前使用chan，所以在关闭后未设置重启机制
acotr内部实现了链路机制，mailbox可以自定义流量统计
基于本actor的功能，后期可能会做调整，上层api不会改变

# 编写一个简单的服务器
```go
  network.StartTcpServer(":8088", network.OptionHandler(func(conn net.Conn) (err error) {
		Stage().ActorOf(WithRemoteStream(func(ctx ActorContext) {
			switch body := ctx.Any().(type) {
			case *Started:
				fmt.Println("connect addr:", conn.RemoteAddr().String())
			case *Stopped:
			case []byte:
				fmt.Println("message: ", network.ReadBegin(body))
				ctx.Send(ctx.Self(), network.WriteBegin(0x102).Flush())
			case error: //设置了心跳会被通知
			default:
				fmt.Printf("untreated type %T \n", body)
			}
		}, conn))
		return
	}))

	cli := Stage().ActorOf(WithRemoteAddr(func(ctx ActorContext) {
		switch body := ctx.Any().(type) {
		case *Started:
			fmt.Println("open ok")
		case *Stopped:
		case []byte:
			fmt.Println("cli respond: ", network.ReadBegin(body))
		case Failure: //最终错误退出的原因

		default:
			fmt.Printf("other object miss handle %T \n", body)
		}
	}, "localhost:8088"))

	psend := network.WriteBegin(0x101, "who's your daddy!")
	cli.Tell(psend.Flush())
```
