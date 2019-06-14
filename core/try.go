package core

import (
	"fmt"
	"runtime/debug"
)

func Throw(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

/*
 * 捕获错误(不发生错误，性能不会下降)
 */
func Try(f func() error) (err error) {
	defer func() {
		if code, ok := Catch(recover()); ok {
			err = code
		}
	}()
	err = f()
	return
}

/*
 * 未发送错误
 */
func TryOk(err interface{}) bool {
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

/*
 * 发送错误false
 */
func PrintStack() {
	debug.PrintStack()
}

/*
* 捕获，不一定产生错误
 */
func Catch(flag interface{}) (err error, ok bool) {
	//defer func() { Catch(recover()) }()
	if flag != nil {
		ok = true
		switch p := flag.(type) {
		case error:
			err = p
		default:
			err = fmt.Errorf("catch err: %v", p)
		}
	}
	return
}

func CatchErr(flag interface{}) (err error) {
	//defer func() { Catch(recover()) }()
	if flag != nil {
		switch p := flag.(type) {
		case error:
			err = p
		default:
			err = fmt.Errorf("catch err: %v", p)
		}
	}
	return
}
