package core

import (
	"fmt"
	"sync"
	"time"
)

func TestGo(name string, n int, totals int, f func(int)) {
	var wg sync.WaitGroup
	wg.Add(n)
	tm := time.Now()
	for i := 0; i < n; i++ {
		go func(idx int) {
			f(idx)
			wg.Done()
		}(i)
	}
	wg.Wait()
	d := time.Since(tm)
	suc := d / time.Duration(totals)
	fmt.Println("[", name, "] tatols time=", d, " once time=", suc, " num=", totals, " tps=", int64((time.Second*time.Duration(totals))/d))
}

func TestFunc(name string, n int, totals int, f func(int)) {
	tm := time.Now()
	for i := 0; i < n; i++ {
		f(i)
	}
	d := time.Since(tm)
	suc := d / time.Duration(totals)
	fmt.Println("[", name, "] tatols time=", d, " once time=", suc, " num=", totals, " tps=", int64((time.Second*time.Duration(totals))/d))
}
