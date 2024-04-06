package listeners

import (
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/authutils"
	"github.com/arravoco/hackathon_backend/utils/email"
)

func HandleParticipantCreatedEvent(eventDTOData *exports.ParticipantAccountCreatedEventData, otherParams ...interface{}) {
	participantType := eventDTOData.ParticipantType
	switch participantType {
	case "TEAM":
		ttl := time.Now().Add(time.Minute * 15)
		dataToken, err := authutils.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
			Email: eventDTOData.ParticipantEmail,
			TTL:   ttl,
		})
		if err != nil {
			exports.MySugarLogger.Error(err)
			return
		}
		link, err := utils.GenerateEmailVerificationLink(&exports.EmailVerificationLinkPayload{
			Token: dataToken.Token,
			TTL:   dataToken.TTL,
			Email: dataToken.TokenTypeValue,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if eventDTOData.TeamRole == "TEAM_LEAD" {
			email.SendTeamLeadWelcomeEmail(&email.SendTeamLeadWelcomeEmailData{
				SendWelcomeEmailData: email.SendWelcomeEmailData{Email: eventDTOData.ParticipantEmail,
					LastName:  eventDTOData.LastName,
					FirstName: eventDTOData.FirstName,
					Subject:   " Welcome to Arravo's Hackathon - Confirm Your Email Address",
					Token:     dataToken.Token,
					TTL:       dataToken.TTL.Minute(),
					Link:      link,
				},
			})
		}

	case "INDIVIDUAL":
		ttl := time.Now().Add(time.Minute * 15)
		dataToken, err := authutils.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
			Email: eventDTOData.ParticipantEmail,
			TTL:   ttl,
		})
		if err != nil {
			exports.MySugarLogger.Error(err)
			return
		}
		link, err := utils.GenerateEmailVerificationLink(&exports.EmailVerificationLinkPayload{
			Token: dataToken.Token,
			TTL:   dataToken.TTL,
			Email: dataToken.TokenTypeValue,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		email.SendIndividualParticipantWelcomeEmail(&email.SendIndividualWelcomeEmailData{
			SendWelcomeEmailData: email.SendWelcomeEmailData{
				Email:     eventDTOData.ParticipantEmail,
				LastName:  eventDTOData.LastName,
				FirstName: eventDTOData.FirstName,
				Subject:   " Welcome to Arravo's Hackathon - Confirm Your Email Address",
				Token:     dataToken.Token,
				TTL:       dataToken.TTL.Minute(),
				Link:      link},
		})

	}
}
