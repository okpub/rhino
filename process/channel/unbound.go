package channel

type Producer func() MessageQueue

//new
func New(args ...Option) MessageQueue {
	this := &myBuffer{}
	this.init(args...)
	return this
}

//producers
func Unbounded(args ...Option) Producer {
	return func() MessageQueue {
		return New(args...)
	}
}
