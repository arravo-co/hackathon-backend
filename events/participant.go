package events

import (
	"github.com/arravoco/hackathon_backend/exports"
)

func EmitParticipantAccountCreated(input *exports.ParticipantAccountCreatedEventData) {
	EventEmitter.EmitEvent(ParticipantAccountCreatedEvent, input)
}

func RegisterParticipantCreatedEvent(fn exports.ParticipantAccountCreatedEventHandler) {
	EventEmitter.AddListener(ParticipantAccountCreatedEvent, func(arguments ...interface{}) {
		arg0 := arguments[0].(*exports.ParticipantAccountCreatedEventData)
		args := arguments[1:]
		fn(arg0, args)
	})
}
