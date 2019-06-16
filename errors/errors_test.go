package errors

import (
	"fmt"
)

func init() {
	Try(func() error {
		fmt.Println("没错")
		Throw("我是谁")
		return nil
	}, func(err error) {
		fmt.Println("err=", err)
	})
}
