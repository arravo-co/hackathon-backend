package listeners

import (
	eventsdtos "github.com/arravoco/hackathon_backend/events_dtos"
	"github.com/arravoco/hackathon_backend/utils/email"
)

func HandleParticipantCreatedEvent(data *eventsdtos.ParticipantAccountCreatedEventData, otherParams ...interface{}) {
	email.SendWelcomeEmail(&email.SendWelcomeEmailData{
		Email:     data.ParticipantEmail,
		LastName:  data.LastName,
		FirstName: data.FirstName,
		Subject:   "Welcome onboard",
	})
}
