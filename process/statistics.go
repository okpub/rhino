package process

//class untype Statistics
type UntypeStatistics struct{}

func (*UntypeStatistics) OnStarted()                     {}
func (*UntypeStatistics) OnPosted(interface{})           {}
func (*UntypeStatistics) OnReceived(interface{})         {}
func (*UntypeStatistics) OnDiscarded(error, interface{}) {}
func (*UntypeStatistics) OnFree()                        {}
