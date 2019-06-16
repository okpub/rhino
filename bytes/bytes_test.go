package bytes

import (
	gbytes "bytes"
	"fmt"
	"testing"
)

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
	a := 10

	b := New()
	b.Wobj(a)

	var c int
	b.SeekBegin()
	b.Robj(&c)
	fmt.Println("读到:", c)

	gbytes.NewBuffer([]byte{1})
}

func BenchmarkCopy1(b *testing.B) {
	m := New()
	t := &TestObj{}
	for i := 0; i < 20; i++ {
		m.Wobj(t)
	}
	fmt.Println("read size:", m.Len())
	for i := 0; i < b.N; i++ {
		m.SeekBegin()
		for k := 0; k < 10; k++ {
			t2 := &TestObj{}
			m.Robj(t2)
		}
	}
}

func cBenchmarkCopy2(b *testing.B) {
	m := New()
	t := &TestObj{}
	m.Wobj(t)
	for i := 0; i < b.N; i++ {
		m.SeekBegin()
		//t2 := &TestObj{}

	}
}
