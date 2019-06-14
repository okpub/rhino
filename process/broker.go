package process

import (
	"fmt"
)

type ProcessType int

const (
	NilFlag ProcessType = iota
	StartFlag
	StopFlag
)

//class untype invoker
type UntypeBroker struct{}

func (*UntypeBroker) PreStart() {
	fmt.Println("Warning: UntypeBroker ignore start!")
}

func (*UntypeBroker) DispatchMessage(data interface{}) {
	fmt.Println("Warning: UntypeBroker ignore message =", data)
}

func (*UntypeBroker) ThrowFailure(err error, body interface{}) {
	fmt.Println("Warning: UntypeBroker ignore failure { err:", err, "body:", body, "}")
}

func (*UntypeBroker) PostStop() {
	fmt.Println("Warning: UntypeBroker ignore stop!")
}

//default func invoker
type funcBroker func(interface{})

func (f funcBroker) DispatchMessage(body interface{})         { f(body) }
func (f funcBroker) PreStart()                                { f(StartFlag) }
func (f funcBroker) PostStop()                                { f(StopFlag) }
func (f funcBroker) ThrowFailure(err error, body interface{}) { f(err) }

func DoWithFunc(f func(interface{})) Broker {
	return funcBroker(f)
}
