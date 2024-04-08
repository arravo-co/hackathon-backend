package events

import (
	"github.com/arravoco/hackathon_backend/exports"
)

func EmitJudgeAccountCreatedByAdmin(input *exports.JudgeAccountCreatedByAdminEventData) {
	EventEmitter.EmitEvent(JudgeAccountCreatedByAdminEvent, input)
}

func RegisterJudgeCreatedByAdminEvent(fn exports.JudgeAccountCreatedByAdminEventHandler) {
	EventEmitter.AddListener(JudgeAccountCreatedByAdminEvent, func(arguments ...interface{}) {
		arg0 := arguments[0].(*exports.JudgeAccountCreatedByAdminEventData)
		args := arguments[1:]
		fn(arg0, args)
	})
}
