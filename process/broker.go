package process

import (
	"fmt"
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
