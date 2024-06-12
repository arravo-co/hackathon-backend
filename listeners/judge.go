package listeners

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/publish"
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

func HandleJudgeCreatedByAdminEvent(eventDTOData *exports.JudgeAccountCreatedByAdminEventData, otherParams ...interface{}) {

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
	payload := exports.JudgeCreatedByAdminWelcomeEmailQueuePayload{
		Email:       eventDTOData.JudgeEmail,
		Name:        eventDTOData.JudgeName,
		InviterName: eventDTOData.InviterName,
		Password:    eventDTOData.Password,
		TTL:         dataToken.TTL,
		Token:       dataToken.Token,
		Link:        link,
	}
	byt, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = publish.Publish(&exports.PublisherConfig{
		RabbitMQExchange: "",
		RabbitMQKey:      "send.judge.created.admin.welcome_email",
	}, byt)
	if err != nil {
		fmt.Println("Failed to send email:  ", err.Error())
		return
	}
	/*
		queue, err := rmqUtils.GetQueue("send_judge_created_by_admin_welcome_email")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		payload := exports.JudgeCreatedByAdminWelcomeEmailQueuePayload{
			Email:       eventDTOData.JudgeEmail,
			Name:        eventDTOData.JudgeName,
			InviterName: eventDTOData.InviterName,
			Password:    eventDTOData.Password,
			TTL:         dataToken.TTL,
			Token:       dataToken.Token,
			Link:        link,
		}
		byt, err := json.Marshal(payload)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = queue.StartConsuming(1, time.Second)
		queue.AddConsumer("judge_created_by_admin_welcome_email_list", &JudgeCreatedByAdminWelcomeEmailTaskConsumer{})
		err = queue.PublishBytes(byt)
		if err != nil {
			fmt.Println(err.Error())
			return
		}*/
}

type JudgeCreatedByAdminWelcomeEmailTaskConsumer struct {
	Ch chan interface{}
}

func (c *JudgeCreatedByAdminWelcomeEmailTaskConsumer) Consume(d rmq.Delivery) {
	payload := d.Payload()

	payloadStruct := exports.JudgeCreatedByAdminWelcomeEmailQueuePayload{}
	err := json.Unmarshal([]byte(payload), &payloadStruct)
	if err != nil {
		fmt.Println(err.Error())
		if err := d.Reject(); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
	}
	fmt.Println(payloadStruct)
	link, err := utils.GenerateEmailVerificationLink(&exports.EmailVerificationLinkPayload{
		Token: payloadStruct.Token,
		TTL:   payloadStruct.TTL,
		Email: payloadStruct.Email,
	})
	if err != nil {
		fmt.Println(err.Error())
		if err := d.Reject(); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
	}
	err = email.SendJudgeCreatedByAdminWelcomeEmail(&email.SendJudgeCreatedByAdminWelcomeEmailData{
		Email:       payloadStruct.Email,
		Name:        payloadStruct.Name,
		Subject:     " Welcome to Arravo's Hackathon - Confirm Your Email Address",
		Token:       payloadStruct.Token,
		TTL:         payloadStruct.TTL.Minute(),
		Link:        link,
		InviterName: payloadStruct.InviterName,
		Password:    payloadStruct.Password,
	})
	if err != nil {
		fmt.Println(err.Error())
		if err := d.Reject(); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
	}
	d.Ack()
	c.Ch <- struct{}{}
}
