package actor

import (
	"fmt"
)

var (
	SendNilErr = fmt.Errorf("sender is nil")
)

type UntypeContext struct {
	self    ActorRef
	parent  ActorRef
	system  ActorSystem
	actor   Actor
	message interface{}
}

//infoPart
func (this *UntypeContext) Self() ActorRef      { return this.self }
func (this *UntypeContext) Parent() ActorRef    { return this.parent }
func (this *UntypeContext) System() ActorSystem { return this.system }
func (this *UntypeContext) Actor() Actor        { return this.actor }

//spawnPart
func (this *UntypeContext) Stop(target ActorRef) { this.stopIdelActive(target) }

//private
func (this *UntypeContext) stopIdelActive(target ActorRef) (err error) {
	if target == nil {
		err = SendNilErr
	} else {
		err = target.Close()
	}
	return
}

func (this *UntypeContext) setMessageEnvelope(value interface{}) { this.message = value }
func (this *UntypeContext) getMessageEnvelope() interface{}      { return this.message }

//messagePart
func (this *UntypeContext) Any() interface{} { return UnwrapEnvelopeMessage(this.message) }

//sendPart sender
func (this *UntypeContext) Sender() ActorRef { return UnwrapEnvelopeSender(this.message) }

//signature
func (this *UntypeContext) getSignatureEnvelope(value interface{}) MessageEnvelope {
	return MSG(this.self, value)
}
