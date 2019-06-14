package bit

import (
	"fmt"
	"sync"

	"github.com/okpub/rhino/core"
	//	"reflect"
)

type Mp interface {
	Kind()
}

type T1 struct {
	Ms string
}

func (T1) Kind() {

}

type String string
type INT64 int64

type TestObj struct {
	Age   INT64
	Name  String
	Sex1  int
	Sex2  int8
	Sex3  int16
	Sex4  int32
	Sex5  int64
	Sex7  uint
	Sex8  uint8
	Sex9  uint16
	Sex19 uint32
}

func init() {
	b := New()

	t := &TestObj{}
	b.Wobj(t)
	b.SeekBegin()
	t2 := &TestObj{}
	b.Robj(t2)
	fmt.Println(t2)

	//
	b2 := New()

	q1 := []int32{12, 312312313, 23}
	b2.Wobj(q1)
	b2.SeekBegin()
	q2 := make([]int32, 3)
	b2.Robj(q2)
	fmt.Println(q2)
	//
	test1()
	test2()
}

func test1() {
	m := New()
	t := &TestObj{}
	m.Wobj(t)
	core.TestGo("new", 1000000, 1000000, func(int) {
		mc := New()
		mc.Write(m.Bytes())
		mc.SeekBegin()
		t2 := &TestObj{}
		mc.Robj(t2)
	})
}

var p = sync.Pool{New: func() interface{} {
	return WithSize(512)
}}

func test2() {
	m := New()
	t := &TestObj{}
	m.Wobj(t)
	for i := 0; i < 10000; i++ {
		p.Put(p.New())
	}
	fmt.Println(m.Len())
	core.TestGo("pool", 1000000, 1000000, func(int) {
		mc := p.Get().(IBytes)
		mc.SeekBegin()
		mc.Write(m.Bytes())
		t2 := &TestObj{}
		mc.Robj(t2)
		p.Put(mc)
	})
}
