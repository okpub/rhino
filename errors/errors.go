package errors

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

func New(err string) error {
	return fmt.Errorf(err)
}

func Newf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

//no used PS: if directly using,which leads to the Stack function can not capture right
func Throwf(format string, args ...interface{}) {
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

//get to the location of the panic
func Stack() (str string) {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(4, pc[:])
	callers := pc[:n]
	frames := runtime.CallersFrames(callers)
	for {
		frame, more := frames.Next()
		file = frame.File
		line = frame.Line
		name = frame.Function
		if !strings.HasPrefix(name, "runtime.") || !more {
			break
		}
	}

	switch {
	case name != "":
		str = fmt.Sprintf("%v:%v", name, line)
	case file != "":
		str = fmt.Sprintf("%v:%v", file, line)
	default:
		str = fmt.Sprintf("pc:%x", pc)
	}
	return
}

func PrintStack() {
	debug.PrintStack()
}
