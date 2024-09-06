package consumerhandlers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/jobs"
	"github.com/arravoco/hackathon_backend/publish"
	taskmgt "github.com/arravoco/hackathon_backend/task_mgt"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/email"
	amqp "github.com/rabbitmq/amqp091-go"
)

func HandleSendInviteEmailConsumption(by []byte) error {

	payloadStruct := exports.AddedToInvitelistPublishPayload{}
	err := json.Unmarshal([]byte(by), &payloadStruct)
	if err != nil {
		return err
	}

	ttl := time.Now().Add(time.Hour * 24 * 7)

	linkPayload, err := utils.GenerateTeamInviteLinkPayload(&exports.TeamInviteLinkPayload{
		ParticipantId:      payloadStruct.ParticipantId,
		TeamLeadEmailEmail: payloadStruct.TeamLeadEmailEmail,
		InviteeEmail:       payloadStruct.InviteeEmail,
		HackathonId:        payloadStruct.HackathonId,
		TTL:                ttl.Unix(),
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = email.SendInviteTeamMemberEmail(&email.SendTeamInviteEmailData{
		InviterName:  payloadStruct.InviterName,
		InviterEmail: payloadStruct.InviterEmail,
		InviteeEmail: payloadStruct.InviteeEmail,
		InviteeName:  payloadStruct.InviterName,
		Subject:      "Invitation to Join Arravo Hackathon Link",
		TTL:          ttl.Day(),
		Link: strings.Join(
			[]string{
				strings.Join(
					[]string{
						config.GetServerURL(), "api/v1/auth/team/invite"}, "/"), linkPayload}, "?token="),
	})
	if err != nil {
		exports.MySugarLogger.Errorln(err)
		return err
	}
	return nil
}

func HandleSolutionPicUploadConsumption(response *amqp.Delivery) {
	fmt.Println("HandleSolutionPicUploadConsumption: ", response.Body)
	dt := &exports.UploadSolutionPicQueuePayload{}
	err := json.Unmarshal(response.Body, dt)
	if err != nil {
		fmt.Println(err.Error())
		err = response.Ack(false)
		if err != nil {
			fmt.Println("Failed to reject delivery. Error: ", err.Error())
			return
		}
		fmt.Println("Should never get here.")
		return
	}
	tsk, err := taskmgt.GetTaskById(dt.TaskId)
	if err != nil {
		fmt.Println(err.Error())
		if err := response.Reject(false); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
	}
	if tsk.Status == "COMPLETED" {
		response.Ack(false)
		return
	}
	err = jobs.ConsumeUploadSolutionPicQueue(dt)
	if err != nil {
		fmt.Println(err.Error())
		if err := response.Reject(false); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
	}
	taskmgt.UpdateTaskStatusById(tsk.Id, "COMPLETED")
	if err != nil {
		fmt.Println(err.Error())
		if err := response.Reject(true); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
	}
	response.Ack(false)
}

func HandleCoordinateParticipantWelcomeVerificationConsumption(response *amqp.Delivery) {
	response.Ack(false)
	//coordState:=exports.SendParticipantWelcomeAndVerificationEmailCoordinatorState{}
	var blankType interface{}
	json.Unmarshal(response.Body, blankType)
	switch typedPayload := blankType.(type) {
	case exports.CreateEmailTokenJobResponsePayload:
		fmt.Println(typedPayload)
		if typedPayload.RespondingJobTaskName == "generate_email_token" {
			//:=exports.SendWelcomeAndEmailVerificationQueueRequestPayload{}
			publish.Publish(&exports.PublisherConfig{}, nil)
		} else {
			fmt.Println("Unusual name in this phase:", typedPayload.RespondingJobTaskName)
		}
	}
}
