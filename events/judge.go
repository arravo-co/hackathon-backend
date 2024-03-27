package events

import (
	"github.com/arravoco/hackathon_backend/exports"
)

func EmitJudgeAccountCreated(input *exports.JudgeAccountCreatedEventData) {
	EventEmitter.EmitEvent(JudgeAccountCreatedEvent, input)
}

func RegisterJudgeCreatedEvent(fn exports.JudgeAccountCreatedEventHandler) {
	EventEmitter.AddListener(JudgeAccountCreatedEvent, func(arguments ...interface{}) {
		arg0 := arguments[0].(*exports.JudgeAccountCreatedEventData)
		args := arguments[1:]
		fn(arg0, args)
	})
}
