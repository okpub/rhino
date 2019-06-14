package actor

//actor message
type MessageEnvelope interface {
	Replace(interface{})
	Any() interface{}
	Sender() ActorRef
}

//static message
func MSG(self ActorRef, body interface{}) MessageEnvelope {
	return &actorMessage{self, body}
}

func WrapEnvelope(message interface{}) MessageEnvelope {
	if env, ok := message.(MessageEnvelope); ok {
		return env
	}
	return MSG(nil, message)
}

func UnwrapEnvelopeMessage(message interface{}) interface{} {
	if env, ok := message.(MessageEnvelope); ok {
		return env.Any()
	}
	return message
}

func UnwrapEnvelopeSender(message interface{}) ActorRef {
	if env, ok := message.(MessageEnvelope); ok {
		return env.Sender()
	}
	return nil
}

//class
type actorMessage struct {
	sender ActorRef
	body   interface{}
}

func (this *actorMessage) Replace(data interface{}) { this.body = data }
func (this *actorMessage) Any() interface{}         { return this.body }
func (this *actorMessage) Sender() ActorRef         { return this.sender }

//class message
type Started struct{}
type Stopped struct{}
type Restart struct{}

func (*Started) String() string { return "start" }
func (*Stopped) String() string { return "stop" }
func (*Restart) String() string { return "restart" }

var (
	started = &Started{}
	stoped  = &Stopped{}
	restart = &Restart{}
)

//error
type Error struct {
	err  error
	body interface{}
}

func (this *Error) Error() string     { return this.err.Error() }
func (this *Error) Body() interface{} { return this.body }

func Errof(err error, body interface{}) *Error {
	return &Error{err, body}
}
