package actor

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	Failed           = errors.New("failed canceled")           //添加错误
	Canceled         = errors.New("self canceled")             //自身取消
	ParentCanceled   = errors.New("parent canceled")           //父亲取消
	SystemCanceled   = errors.New("system canceled")           //系统取消
	DeadlineExceeded = errors.New("context deadline exceeded") //超时
	UnrealizedErr    = errors.New("unrealized function")       //未实现
	OverfullErr      = errors.New("channel overfull")          //未实现
)

type (
	//只需要知道退出原因即可(不关心退出时候的信号)
	Context interface {
		Err() error
	}

	//节点(所有关闭操作在removeSelf里面执行)
	Node interface {
		removeSelf(bool, error) error
	}

	NodeSet map[Node]struct{}

	//节点树(本身可以不为节点)
	Tree interface {
		addChild(Node) error
		removeChild(Node)
	}
)

func NewTree(parent Context) PIDGroup {
	return PIDGroup{Context: parent}
}

//actor内核(谁来管理内核 用来继承)
type PIDGroup struct {
	Context
	childs NodeSet
	mu     sync.Mutex
	err    error
}

//interface context
func (this *PIDGroup) Err() (err error) {
	this.mu.Lock()
	err = this.err
	this.mu.Unlock()
	return
}

//interface tree
func (this *PIDGroup) addChild(child Node) (err error) {
	this.mu.Lock()
	if err = this.err; err == nil {
		if this.childs == nil {
			this.childs = make(map[Node]struct{})
		}
		this.childs[child] = struct{}{}
	}
	this.mu.Unlock()
	//Add the fall will be cancelled
	if err != nil {
		child.removeSelf(false, Failed)
	}
	return
}

func (this *PIDGroup) removeChild(child Node) {
	this.mu.Lock()
	delete(this.childs, child)
	this.mu.Unlock()
	return
}

//interface node
func (this *PIDGroup) removeSelf(removeFromParent bool, code error) (err error) {
	if code == nil {
		panic("remove self code=nil")
	}
	this.mu.Lock()
	if err = this.err; err == nil {
		//setter close
		this.err = code
		//remove child
		for child := range this.childs {
			child.removeSelf(false, ParentCanceled)
		}
		this.childs = nil
	}
	this.mu.Unlock()
	//delete
	if removeFromParent {
		Fire(this.Context, this)
	}
	return
}

//static function
func Join(parent Context, child Node) (err error) {
	if p, ok := parentTransformTree(parent); ok {
		err = p.addChild(child)
	} else {
		fmt.Println("WARNING: root tree", Typeof(parent))
	}
	return
}

func Fire(parent Context, child Node) {
	if p, ok := parentTransformTree(parent); ok {
		p.removeChild(child)
	} else {
		fmt.Println("INFO: fire untree parent=", Typeof(parent), " child=", Typeof(child))
	}
}

func parentTransformTree(parent Context) (p Tree, ok bool) {
	p, ok = parent.(Tree)
	return
}

func Typeof(p interface{}) string {
	if p == nil {
		return "<nil>"
	}
	return reflect.TypeOf(p).String()
}
