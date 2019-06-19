package errors

import (
	"fmt"
)

func todo() error {
	fmt.Println("没错")
	panic(New("我是谁"))
	return nil
}

func init() {
	Try(todo, func(err error) {
		fmt.Println("err=", err, Stack())
	})
}
