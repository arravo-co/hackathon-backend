package events

import (
	eventsdtos "github.com/arravoco/hackathon_backend/events_dtos"
	"github.com/arravoco/hackathon_backend/listeners"
	eventemitter "github.com/vansante/go-event-emitter"
)

var EventEmitter *eventemitter.Emitter
var ListenersToParticipantsCreatedEvent map[string]eventsdtos.ParticipantAccountCreatedEventHandler

var ParticipantAccountCreatedEvent eventemitter.EventType

func init() {
	EventEmitter = eventemitter.NewEmitter(true)
	RegisterParticipantCreatedEvent(listeners.HandleParticipantCreatedEvent)
}

func EmitParticipantAccountCreated(input *eventsdtos.ParticipantAccountCreatedEventData) {
	EventEmitter.EmitEvent(ParticipantAccountCreatedEvent, input)
}

func RegisterParticipantCreatedEvent(fn eventsdtos.ParticipantAccountCreatedEventHandler) {
	EventEmitter.AddListener(ParticipantAccountCreatedEvent, func(arguments ...interface{}) {
		arg := arguments[0].(*eventsdtos.ParticipantAccountCreatedEventData)
		fn(arg)
	})
}
