package remote

//// Broker is an interface used for asynchronous messaging.
//type Broker interface {
//	Address() string
//	Connect(string) error
//	Disconnect() error
//	Reconnect() error
//	Publish(interface{}) error
//	Subscribe(Handler, interface{}) (Subscriber, error)
//}

//// Handler is used to process messages via a subscription of a topic.
//// The handler is passed a publication interface which contains the
//// message and optional Ack method to acknowledge receipt of the message.
//type Handler func(Publication) error

//// Publication is given to a subscription handler for processing
//type Publication interface {
//	Topic() int          //主题/cmd
//	ServerID() int       //指定id
//	MessageType() int    //指定类型
//	Body() interface{}   //携带消息
//	Target() interface{} //绑定的对象
//}

////订阅者会关注所有的主题，至于订阅者内部是否需要处理，自己决定(订阅者可以是多个)
//type Subscriber interface {
//	Call(Publication) error
//	Unsubscribe() error
//}
