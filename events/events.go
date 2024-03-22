package events

import (
	eventemitter "github.com/vansante/go-event-emitter"
)

var EventEmitter *eventemitter.Emitter

var ParticipantAccountCreatedEvent eventemitter.EventType

func init() {
	EventEmitter = eventemitter.NewEmitter(true)
	//RegisterParticipantCreatedEvent(listeners.HandleParticipantCreatedEvent)
}
