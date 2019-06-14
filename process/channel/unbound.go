package channel

type Producer func() MessageQueue

//默认邮箱的大小
const (
	defaultPendingNum = 10
)

//new
func New(args ...Option) MessageQueue {
	this := &myBuffer{
		opts: Options{
			PendingNum: defaultPendingNum,
		},
	}
	this.init(args...)
	return this
}

//producers
func Unbounded(args ...Option) Producer {
	return func() MessageQueue {
		return New(args...)
	}
}
