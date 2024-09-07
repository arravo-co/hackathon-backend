package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/publish"
	taskmgt "github.com/arravoco/hackathon_backend/task_mgt"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/authutils"
	"github.com/arravoco/hackathon_backend/utils/email"
)

func HandleSendParticipantWelcomeAndVerificationEmailJob(eventDTOData *exports.SendWelcomeAndEmailVerificationTokenJobQueuePayload) error {

	ttl := time.Now().Add(time.Minute * time.Duration(15))
	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	dataToken, err := auth.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
		Email: eventDTOData.Email,
		TTL:   ttl,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	taskmgt.UpdateTaskStatusById(eventDTOData.TaskId, "COMPLETED")

	jb := &exports.CreateEmailTokenJobResponsePayload{
		QueuePayload: exports.QueuePayload{
			RecipientQueueName:  "send.participant.created.welcome_email_verification_email",
			RequestingJobTaskID: eventDTOData.ParentTaskID,
			TaskId:              eventDTOData.TaskId,
			RespondingJobTaskID: eventDTOData.TaskId,
			Direction:           "RESPONSE",
		},
		TokenId:        dataToken.Id,
		Token:          dataToken.Token,
		TokenType:      dataToken.TokenType,
		TokenTypeValue: dataToken.TokenTypeValue,
		Status:         dataToken.Status,
		CreatedAt:      dataToken.CreatedAt,
		UpdatedAt:      dataToken.UpdatedAt,
	}
	byt, err := json.Marshal(jb)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	err = publish.Publish(&exports.PublisherConfig{}, byt)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	return nil
}

func HandleCreateEmailTokenJob(eventDTOData *exports.SendWelcomeAndEmailVerificationTokenJobQueuePayload) error {

	ttl := time.Now().Add(time.Minute * time.Duration(15))
	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	dataToken, err := auth.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
		Email: eventDTOData.Email,
		TTL:   ttl,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	taskmgt.UpdateTaskStatusById(eventDTOData.TaskId, "COMPLETED")

	jb := &exports.CreateEmailTokenJobResponsePayload{
		QueuePayload: exports.QueuePayload{
			RequestingJobTaskID: eventDTOData.ParentTaskID,
			TaskId:              eventDTOData.TaskId,
			RespondingJobTaskID: eventDTOData.TaskId,
			Direction:           "RESPONSE",
		},
		TokenId:        dataToken.Id,
		Token:          dataToken.Token,
		TokenType:      dataToken.TokenType,
		TokenTypeValue: dataToken.TokenTypeValue,
		Status:         dataToken.Status,
		CreatedAt:      dataToken.CreatedAt,
		UpdatedAt:      dataToken.UpdatedAt,
	}
	byt, err := json.Marshal(jb)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	err = publish.Publish(&exports.PublisherConfig{}, byt)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	return nil
}

func HandleSendEmailJob(eventDTOData *exports.ParticipantAccountCreatedEventData, otherParams ...interface{}) {
	participantType := eventDTOData.ParticipantType
	switch participantType {
	case "TEAM":
		ttl := time.Now().Add(time.Minute * 15)
		auth := authutils.GetAuthUtilsWithDefaultRepositories()
		dataToken, err := auth.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
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
		auth := authutils.GetAuthUtilsWithDefaultRepositories()
		dataToken, err := auth.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
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
