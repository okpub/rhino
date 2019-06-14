package process

//异步调度器
type goroutineDispatcher int

func (n goroutineDispatcher) Schedule(fn func()) { go fn() }
func (n goroutineDispatcher) Throughput() int    { return int(n) }

func NewDefaultDispatcher(throughput int) Dispatcher {
	return goroutineDispatcher(throughput)
}

//同步调度器
type synchronizedDispatcher int

func (n synchronizedDispatcher) Schedule(fn func()) { fn() }
func (n synchronizedDispatcher) Throughput() int    { return int(n) }

func NewSyncDispatcher(throughput int) Dispatcher {
	return synchronizedDispatcher(throughput)
}
