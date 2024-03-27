package listeners

import (
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/utils/authutils"
	"github.com/arravoco/hackathon_backend/utils/email"
)

func HandleJudgeCreatedEvent(eventDTOData *exports.JudgeAccountCreatedEventData, otherParams ...interface{}) {

	ttl := time.Now().Add(time.Minute * 15)
	dataToken, err := authutils.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
		Email: eventDTOData.JudgeEmail,
		TTL:   ttl,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return
	}
	email.SendJudgeWelcomeEmail(&email.SendJudgeWelcomeEmailData{
		SendWelcomeEmailData: email.SendWelcomeEmailData{
			Email:     eventDTOData.JudgeEmail,
			LastName:  eventDTOData.LastName,
			FirstName: eventDTOData.FirstName,
			Subject:   " Welcome to Arravo's Hackathon - Confirm Your Email Address",
			Token:     dataToken.Token,
			TTL:       dataToken.TTL.Minute()},
	})

}
