package channel

/*
*公开的选项
 */
type Options struct {
	NonBlocking bool             //非阻塞模式(默认阻塞)
	Buffer      chan interface{} //消息通道
}

//填充
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
func OptionPendingNum(peningNum int) Option {
	return func(p *Options) {
		p.Buffer = make(chan interface{}, peningNum)
	}
}

func OptionBlocking(blocking bool) Option {
	return func(p *Options) {
		p.NonBlocking = !blocking
	}
}

func OptionWithFunc(fn func() chan interface{}) Option {
	return func(p *Options) {
		p.Buffer = fn()
	}
}
