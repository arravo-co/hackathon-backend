package events

import (
	eventsdtos "github.com/arravoco/hackathon_backend/events_dtos"
)

func EmitParticipantAccountCreated(input *eventsdtos.ParticipantAccountCreatedEventData) {
	EventEmitter.EmitEvent(ParticipantAccountCreatedEvent, input)
}

func RegisterParticipantCreatedEvent(fn eventsdtos.ParticipantAccountCreatedEventHandler) {
	EventEmitter.AddListener(ParticipantAccountCreatedEvent, func(arguments ...interface{}) {
		arg := arguments[0].(*eventsdtos.ParticipantAccountCreatedEventData)
		args := arguments[1:]
		fn(arg, args)
	})
}
