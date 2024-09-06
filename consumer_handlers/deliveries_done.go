package consumerhandlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/jobs"
	taskmgt "github.com/arravoco/hackathon_backend/task_mgt"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/authutils"
	"github.com/arravoco/hackathon_backend/utils/email"
	"github.com/jaevor/go-nanoid"
)

func HandleSendWelcomeAndEmailVerificationEmailToJudgeConsumption(by []byte) error {
	payloadStruct := exports.JudgeRegisteredPayload{}
	err := json.Unmarshal([]byte(by), &payloadStruct)
	if err != nil {
		return err
	}
	fmt.Println("HandleSendWelcomeAndEmailVerificationEmailToJudgeConsumption payload: \n\n", payloadStruct)

	ttl := time.Now().Add(time.Minute * 15)
	tokenFunc, _ := nanoid.Custom("1234567890", 6)
	token := tokenFunc()
	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	tokenData, err := auth.GenerateToken(&authutils.GenerateTokenArgs{
		Token:          token,
		TokenType:      "EMAIL",
		TokenTypeValue: payloadStruct.Email,
		TTL:            ttl,
		Status:         "PENDING",
		Scope:          "EMAIL_VERIFICATION",
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}

	link, err := utils.GenerateEmailVerificationLink(&exports.EmailVerificationLinkPayload{
		Token: tokenData.Token,
		TTL:   ttl,
		Email: payloadStruct.Email,
	})
	err = email.SendJudgeCreatedByAdminWelcomeEmail(&email.SendJudgeCreatedByAdminWelcomeEmailData{
		Email:       payloadStruct.Email,
		Name:        payloadStruct.Name,
		Subject:     "Welcome to Arravo's Hackathon - Confirm Your Email Address",
		Token:       tokenData.Token,
		TTL:         ttl.Minute(),
		Link:        link,
		InviterName: payloadStruct.InviterName,
		Password:    payloadStruct.Password,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func HandleSendParticipantWelcomeAndVerificationEmailConsumption(by []byte) error {
	payloadStruct := exports.ParticipantRegisteredPublishPayload{}
	err := json.Unmarshal([]byte(by), &payloadStruct)
	if err != nil {
		return err
	}
	fmt.Println("HandleSendWelcomeAndEmailVerificationEmailToJudgeConsumption payload: \n\n", payloadStruct)

	ttl := time.Now().Add(time.Minute * 15)
	tokenFunc, _ := nanoid.Custom("1234567890", 6)
	token := tokenFunc()
	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	tokenData, err := auth.GenerateToken(&authutils.GenerateTokenArgs{
		Token:          token,
		TokenType:      "EMAIL",
		TokenTypeValue: payloadStruct.Email,
		TTL:            ttl,
		Status:         "PENDING",
		Scope:          "EMAIL_VERIFICATION",
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}

	link, err := utils.GenerateEmailVerificationLink(&exports.EmailVerificationLinkPayload{
		Token: tokenData.Token,
		TTL:   ttl,
		Email: payloadStruct.Email,
	})

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if payloadStruct.TeamRole == "TEAM_LEAD" {

		err = email.SendTeamLeadWelcomeEmail(&email.SendTeamLeadWelcomeEmailData{
			SendWelcomeEmailData: email.SendWelcomeEmailData{
				Email:     payloadStruct.Email,
				FirstName: payloadStruct.FirstName,
				LastName:  payloadStruct.LastName,
				Subject:   "Welcome to Arravo's Hackathon - Confirm Your Email Address",
				Token:     tokenData.Token,
				TTL:       int(time.Until(tokenData.TTL)),
				Link:      link,
			},
		})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

func HandleUploadJudgeProfilePicConsumption(by []byte) error {
	dt := &exports.UploadJudgeProfilePicQueuePayload{}
	err := json.Unmarshal([]byte(by), &dt)
	if err != nil {
		return err
	}
	fmt.Println("HandleUploadJudgeProfilePicConsumption: ")
	tsk, err := taskmgt.GetTaskById(dt.TaskId)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if tsk.Status == "COMPLETED" {
		return nil
	}
	err = jobs.ConsumeUploadJudgeProfilePicQueue(dt)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	taskmgt.UpdateTaskStatusById(tsk.Id, "COMPLETED")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
