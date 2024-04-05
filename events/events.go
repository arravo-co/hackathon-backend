package events

import (
	"github.com/arravoco/hackathon_backend/listeners"
	eventemitter "github.com/vansante/go-event-emitter"
)

var EventEmitter *eventemitter.Emitter

var ParticipantAccountCreatedEvent eventemitter.EventType = "ParticipantAccountCreated"
var JudgeAccountCreatedEvent eventemitter.EventType = "JudgeAccountCreated"

var AdminAccountCreatedEvent eventemitter.EventType = "AdminAccountCreated"

var AdminAccountCreatedByAdminEvent eventemitter.EventType = "AdminAccountCreatedByAdmin"

func init() {
	EventEmitter = eventemitter.NewEmitter(true)
	RegisterParticipantCreatedEvent(listeners.HandleParticipantCreatedEvent)
	RegisterJudgeCreatedEvent(listeners.HandleJudgeCreatedEvent)
	RegisterAdminCreatedEvent(listeners.HandleAdminCreatedEvent)
	RegisterAdminCreatedByAdminEvent(listeners.HandleAdminCreatedByAdminEvent)
}
