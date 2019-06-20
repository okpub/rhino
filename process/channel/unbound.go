package channel

type Producer func() MessageQueue

//new
func New(args ...Option) MessageQueue {
	return MakeBuffer(10, args...) //default value
}

func MakeBuffer(pendingNum int, args ...Option) MessageQueue {
	this := &Mailbox{
		pendingNum: pendingNum, //default blocking
	}
	return this.Init(args...)
}

//producer
func Unbounded(args ...Option) Producer {
	return func() MessageQueue {
		return New(args...)
	}
}
