package process

import (
	"fmt"
	"testing"
)

type Test struct {
	N   int
	M   int
	C   int
	Str string
	TestObj
}

type TestObj struct {
	N   int
	M   int
	C   int
	Str string
}

func (this Test) copy1() Test {
	return this
}

func (this Test) copy2() *Test {
	return &this
}

func (this *Test) copy3() Test {
	return *this
}

//test
func init() {
	var m = Test{}
	fmt.Printf("%p \n", &m.TestObj)
	m.TestObj.M = 1
	c := m.TestObj
	fmt.Printf("%p \n", &m.TestObj)
	fmt.Printf("%p \n", &m.TestObj)
	fmt.Printf("%p \n", &c)
	fmt.Printf("%+v \n", m.TestObj)
}

func BenchmarkCopy1(b *testing.B) {
	var obj = Test{}
	for i := 0; i < b.N; i++ {
		c := obj.copy1()
		c.M++
	}
}

func BenchmarkCopy2(b *testing.B) {
	var obj = Test{}
	for i := 0; i < b.N; i++ {
		c := obj.copy2()
		c.M++
	}
}

func BenchmarkCopy3(b *testing.B) {
	var obj = Test{}
	for i := 0; i < b.N; i++ {
		c := obj.copy3()
		c.M++
	}
}
