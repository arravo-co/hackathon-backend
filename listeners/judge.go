package listeners

import (
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/authutils"
	"github.com/arravoco/hackathon_backend/utils/email"
)

func HandleJudgeCreatedEvent(eventDTOData *exports.JudgeAccountCreatedByAdminEventData, otherParams ...interface{}) {

	ttl := time.Now().Add(time.Minute * 15)
	dataToken, err := authutils.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
		Email: eventDTOData.JudgeEmail,
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
	email.SendJudgeWelcomeEmail(&email.SendJudgeWelcomeEmailData{
		Email:       eventDTOData.JudgeEmail,
		InviterName: eventDTOData.InviterName,
		Subject:     " Welcome to Arravo's Hackathon - Confirm Your Email Address",
		Token:       dataToken.Token,
		TTL:         dataToken.TTL.Minute(),
		Link:        link,
	})

}
