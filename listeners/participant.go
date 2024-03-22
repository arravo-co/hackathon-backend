package listeners

import (
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/utils/email"
)

type AuthUtil interface {
	InitiateEmailVerification(dataInput *exports.AuthUtilsConfigTokenData) (*exports.TokenData, error)
}

func HandleParticipantCreatedEvent(eventDTOData *exports.ParticipantAccountCreatedEventData, authutils AuthUtil, otherParams ...interface{}) {
	participantType := eventDTOData.ParticipantType
	switch participantType {
	case "TEAM":

	case "INDIVIDUAL":
		ttl := time.Now().Add(time.Minute * 15)
		dataToken, err := authutils.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
			Email: eventDTOData.ParticipantEmail,
			TTL:   ttl,
		})
		if err != nil {

		}
		email.SendIndividualParticipantWelcomeEmail(&email.SendIndividualWelcomeEmailData{
			Email:     eventDTOData.ParticipantEmail,
			LastName:  eventDTOData.LastName,
			FirstName: eventDTOData.FirstName,
			Subject:   "Welcome onboard",
			Token:     dataToken.Token,
		})

	}
}
