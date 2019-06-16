package errors

import (
	"fmt"
	"runtime/debug"
)

func New(err string) error {
	return fmt.Errorf(err)
}

func Newf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func Throw(format string, args ...interface{}) {
	panic(fmt.Errorf("Throw Err: "+format, args...))
}

func Try(method func() error, catch func(error), finals ...func()) (err error) {
	defer func() {
		if err == nil {
			err = Catch(recover())
		}
		if err != nil {
			catch(err)
		}
		for _, die := range finals {
			die()
		}
	}()
	err = method()
	return
}

func Catch(code interface{}) (err error) {
	//defer func() { err = Catch(recover()) }()
	if code != nil {
		switch errc := code.(type) {
		case error:
			err = errc
		default:
			err = fmt.Errorf("catch err: %v", errc)
		}
	}
	return
}

func PrintStack() {
	debug.PrintStack()
}
