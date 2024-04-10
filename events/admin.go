package events

import (
	"github.com/arravoco/hackathon_backend/exports"
)

func EmitAdminAccountCreated(input *exports.AdminAccountCreatedEventData) {
	EventEmitter.EmitEvent(AdminAccountCreatedEvent, input)
}

func EmitAdminAccountCreatedByAdmin(input *exports.AdminAccountCreatedByAdminEventData) {
	EventEmitter.EmitEvent(AdminAccountCreatedByAdminEvent, input)
}

func RegisterAdminCreatedEvent(fn exports.AdminAccountCreatedEventHandler) {
	EventEmitter.AddListener(AdminAccountCreatedEvent, func(arguments ...interface{}) {
		arg0 := arguments[0].(*exports.AdminAccountCreatedEventData)
		args := arguments[1:]
		fn(arg0, args)
	})
}

func RegisterAdminCreatedByAdminEvent(fn exports.AdminAccountCreatedByAdminEventHandler) {
	EventEmitter.AddListener(AdminAccountCreatedByAdminEvent, func(arguments ...interface{}) {
		arg0 := arguments[0].(*exports.AdminAccountCreatedByAdminEventData)
		args := arguments[1:]
		fn(arg0, args)
	})
}
