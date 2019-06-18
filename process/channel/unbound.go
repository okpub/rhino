package channel

type Producer func() MessageQueue

//默认邮箱的大小
const (
	defaultPendingNum = 10
)

//new
func New(args ...Option) MessageQueue {
	this := &Mailbox{
		pendingNum: defaultPendingNum,
		blocking:   true,
	}
	return this.filler(args...)
}

//producers
func Unbounded(args ...Option) Producer {
	return func() MessageQueue {
		return New(args...)
	}
}
