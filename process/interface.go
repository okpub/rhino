package process

//bbb
type (

	//处理器
	Process interface {
		OnRegister(Dispatcher, Broker, ...Statistics)
		Start() error
		Close() error
	}

	//调度(异步/同步)
	Dispatcher interface {
		Schedule(func())
		Throughput() int
	}

	//代替处理(内部环境)
	Broker interface {
		PreStart()                       //调度准备
		DispatchMessage(interface{})     //派送消息
		ThrowFailure(error, interface{}) //系统抛出错误
		PostStop()                       //调度停止
	}

	//统计
	Statistics interface {
		OnStarted()                     //start ready
		OnPosted(interface{})           //Messages sent (not necessarily successful)
		OnDiscarded(error, interface{}) //The lost message(post success)
		OnReceived(interface{})         //The message received(not necessarily processed)
		OnFree()                        //In general for the heart rate response or no news of free time
	}
)
