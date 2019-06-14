package channel

/*
*公开的选项
 */
type Options struct {
	NonBlocking bool             //非阻塞模式(默认阻塞)
	PendingNum  int              //通道大小(需要默认值)
	Buffer      chan interface{} //消息通道
}

//reset buffer
func (this *Options) Fill(args ...Option) {
	for _, o := range args {
		o(this)
	}
}

func (this *Options) Close() (err error) {
	close(this.Buffer)
	return
}

func (this *Options) Post(v interface{}) (err error) {
	if this.NonBlocking {
		select {
		case this.Buffer <- v:
		default:
			err = OverfullErr
		}
	} else {
		this.Buffer <- v
	}
	return
}

//选项
func OptionPendingNum(pendingNum int) Option {
	return func(p *Options) {
		p.PendingNum = pendingNum
	}
}

func OptionBuffer(buffer chan interface{}) Option {
	return func(p *Options) {
		p.Buffer = buffer
	}
}

func OptionBlocking() Option { //阻塞模式
	return func(p *Options) {
		p.NonBlocking = false
	}
}

func OptionNonBlocking() Option { //非阻塞模式
	return func(p *Options) {
		p.NonBlocking = true
	}
}
