package process

/*
* untype prosser
 */
type UntypeProcess struct {
	pc  Dispatcher
	do  Broker
	arr []Statistics
}

//interface Prosser
func (this *UntypeProcess) OnRegister(pc Dispatcher, do Broker, args ...Statistics) {
	this.pc = pc
	this.do = do
	this.arr = args
}

func (this *UntypeProcess) Start() (err error) {
	return
}

func (this *UntypeProcess) Close() (err error) {
	return
}

//interface Statistics
func (this *UntypeProcess) OnStarted() {
	for _, f := range this.arr {
		f.OnStarted()
	}
}

func (this *UntypeProcess) OnPosted(v interface{}) {
	for _, f := range this.arr {
		f.OnPosted(v)
	}
}

func (this *UntypeProcess) OnReceived(v interface{}) {
	for _, f := range this.arr {
		f.OnPosted(v)
	}
}

func (this *UntypeProcess) OnDiscarded(err error, v interface{}) {
	for _, f := range this.arr {
		f.OnDiscarded(err, v)
	}
}

func (this *UntypeProcess) OnFree() {
	for _, f := range this.arr {
		f.OnFree()
	}
}

//interface Broker
func (this *UntypeProcess) PreStart()                                { this.do.PreStart() }
func (this *UntypeProcess) DispatchMessage(body interface{})         { this.do.DispatchMessage(body) }
func (this *UntypeProcess) PostStop()                                { this.do.PostStop() }
func (this *UntypeProcess) ThrowFailure(err error, body interface{}) { this.do.ThrowFailure(err, body) }

//interface Dispatcher
func (this *UntypeProcess) Schedule(fn func()) { this.pc.Schedule(fn) }
func (this *UntypeProcess) Throughput() int    { return this.pc.Throughput() }
