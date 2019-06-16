package network

type (
	//通过连接获得运行器
	Handler   func(Link) Runnable
	ConnSet   map[Link]struct{}
	ListenSet map[Listener]struct{}

	//运行器
	Runnable interface {
		Run()
	}

	//连接
	Link interface {
		Close() error
	}

	//服务
	Server interface {
		Start(string) error
		Group(string) Dialer
		Join(Listener) error
		Close() error
	}

	//监听器
	Listener interface {
		Accept() (Link, error)
		Address() string
		Close() error
	}

	//连接器管理
	Dialer interface {
		Listener
		Join(Link) error
		//Tcp(string) error
		//Web(string) error
	}
)

//空运行
type EmptyRunner int

func (EmptyRunner) Run() { /*不执行任何*/ }
