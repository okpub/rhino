# rhino  (QQ技术群：75205017)
go网络基础框架
# 使用领域
rhino励志打造一个可靠的游戏服务器框架，其核心内容在于

· actor 主要使用在应用层的逻辑部分（actor内部是单线程）例如：日志服务，db服务，游戏服务，网关等等

· network 用于tcp，http，udp的通信，封装读写包

· 其他部分是封装了一些重复劳动的工具

· rpc模块，目前还没有，以后会加入，考虑一个高可用的方案

# acotr模型
1 核心api在actor中，其主要actor处理逻辑部分在于process

2 process包括邮件(channel)和远程(retome)他们都提供对actor的支持

3 所以一切皆actor


# network
net.Conn有个问题需要说下，就是SetReadDeadline/SetWriteDeadline这两个是读写超时的绝对时间

设置读写会超时，那么就需要处理read和write返回值n，查看源码发现net.Conn中write其实是一个write full，所以一般情况不需要设置SetWriteDeadline。

network.stream使用len+body的读取方法，一般处理read/write在单线程比较安全可靠，借助actor模型 由上级actor对其write，来更好的解决和应对分布式中环路堵塞和超时的问题



# 编写一个简单的服务器
```go
func main(){
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
}
```
![75205017.jpg](image/75205017.jpg)
