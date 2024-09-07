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

func SendWelcomeAndEmailVerificationEmailToAdmin(by []byte) error {
	payloadStruct := exports.AdminRegisteredPublishPayload{}
	err := json.Unmarshal([]byte(by), &payloadStruct)
	if err != nil {
		return err
	}
	fmt.Println("SendWelcomeAndEmailVerificationEmailToAdmin payload: \n\n", payloadStruct)

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
		exports.MySugarLogger.Error(err)
		return err
	}
	err = email.SendAdminWelcomeEmail(&email.SendAdminWelcomeEmailData{
		Email:     payloadStruct.Email,
		FirstName: payloadStruct.FirstName,
		LastName:  payloadStruct.LastName,
		Subject:   "Welcome to Arravo's Hackathon - Confirm Your Email Address",
		Token:     tokenData.Token,
		TTL:       int(time.Until(ttl).Minutes()),
		Link:      link,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func SendWelcomeAndEmailVerificationEmailToAdminRegisteredByAdmin(by []byte) error {
	payloadStruct := exports.AdminRegisteredByAdminPublishPayload{}
	err := json.Unmarshal([]byte(by), &payloadStruct)
	if err != nil {
		return err
	}
	fmt.Println("SendWelcomeAndEmailVerificationEmailToAdminRegisteredByAdmin payload: \n\n", payloadStruct)

	ttl := time.Now().Add((time.Hour * 24 * 7) + time.Second*6)
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
		exports.MySugarLogger.Error(err)
		return err
	}
	err = email.SendAdminCreatedByAdminWelcomeEmail(&email.SendAdminCreatedByAdminWelcomeEmailData{
		Email:       payloadStruct.Email,
		AdminName:   payloadStruct.Name,
		Subject:     "Welcome to Arravo's Hackathon - Confirm Your Email Address",
		Token:       tokenData.Token,
		TTL:         int(time.Until(ttl).Hours() / 24),
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

func SendWelcomeAndEmailVerificationEmailToJudge(by []byte) error {
	payloadStruct := exports.JudgeRegisteredPublishPayload{}
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
		exports.MySugarLogger.Error(err)
		return err
	}
	err = email.SendJudgeWelcomeEmail(&email.SendJudgeWelcomeEmailData{
		Email:     payloadStruct.Email,
		JudgeName: payloadStruct.FirstName,
		Subject:   "Welcome to Arravo's Hackathon - Confirm Your Email Address",
		Token:     tokenData.Token,
		TTL:       ttl.Minute(),
		Link:      link,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func SendWelcomeAndEmailVerificationEmailToJudgeRegisteredByAdmin(by []byte) error {
	payloadStruct := exports.JudgeRegisteredByAdminPublishPayload{}
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
		exports.MySugarLogger.Error(err)
		return err
	}
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

func SendTeamLeadWelcomeAndVerificationEmail(by []byte) error {
	payloadStruct := exports.ParticipantRegisteredPublishPayload{}
	err := json.Unmarshal([]byte(by), &payloadStruct)
	if err != nil {
		return err
	}
	fmt.Println("SendTeamLeadWelcomeAndVerificationEmail payload: \n\n", payloadStruct)

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

	return nil
}

func SendTeamMemberWelcomeAndVerificationEmail(by []byte) error {
	payloadStruct := exports.ParticipantRegisteredPublishPayload{}
	err := json.Unmarshal([]byte(by), &payloadStruct)
	if err != nil {
		return err
	}
	fmt.Println("SendTeamMemberWelcomeAndVerificationEmail payload: \n\n", payloadStruct)

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

	err = email.SendTeamMemberWelcomeEmail(&email.SendTeamMemberWelcomeEmailData{
		SendWelcomeEmailData: email.SendWelcomeEmailData{
			Email:     payloadStruct.Email,
			FirstName: payloadStruct.FirstName,
			LastName:  payloadStruct.LastName,
			Subject:   "Welcome to Arravo's Hackathon - Confirm Your Email Address",
			Token:     tokenData.Token,
			TTL:       int(time.Until(tokenData.TTL)),
			Link:      link,
		},
		TeamName:     payloadStruct.TeamName,
		TeamLeadName: payloadStruct.TeamLeadEmail,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
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
