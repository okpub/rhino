package actor

import (
	"github.com/okpub/rhino/process"
)

//default func invoker
type funcBroker func(interface{})

func (f funcBroker) DispatchMessage(msg interface{})          { f(msg) }
func (f funcBroker) PreStart()                                { f(started) }
func (f funcBroker) PostStop()                                { f(stopped) }
func (f funcBroker) ThrowFailure(err error, body interface{}) { f(Failure{err, body}) }

func DoFunc(fn func(interface{})) process.Broker {
	return funcBroker(fn)
}
