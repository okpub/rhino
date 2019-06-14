package event

/*
*携带参数的回调
 */
func NewHandler(args ...interface{}) IDispatcher {
	return &Function{args: args}
}

type Function struct {
	args []interface{}
}

func (this *Function) DispatchEvent(event Event) (err error) {
	switch args := this.args; fn := args[0].(type) {
	case func():
		fn()
	case func(Event):
		fn(event)
	case func(interface{}):
		fn(args[1])
	case func(int):
		fn(args[1].(int))
	case func(string):
		fn(args[1].(string))
	case func(Event, interface{}):
		fn(event, args[1])
	case func(Event, int):
		fn(event, args[1].(int))
	case func(Event, string):
		fn(event, args[1].(string))
	default:
		panic("The callback function is wrong")
	}
	return
}
