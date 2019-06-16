package core

import (
	"fmt"
	"testing"
)

type Test struct {
	i int
	m int
	c int
	d int
	r int8
}

func init() {
	//SysErr()
	var i int
	var z Test
	fmt.Println(Sizeof(i))
	fmt.Println(Sizeof(z))
	init2()
}

func BenchmarkTestZero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ZeroSec()
	}
}

func init2() {
	cls := make(ObjectCreator)
	cls.Register(1, Test{})
	v, _ := cls.New(1)
	fmt.Println("新对象:", Typeof(v))
}
