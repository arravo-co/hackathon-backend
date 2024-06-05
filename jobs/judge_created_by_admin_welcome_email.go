package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/queue"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/email"
)

type JudgeCreatedByAdminWelcomeEmailTaskConsumer struct {
	Ch chan interface{}
}

func StartConsumingJudgeCreatedByAdminWelcomeEmailQueue() (*JudgeCreatedByAdminWelcomeEmailTaskConsumer, error) {
	fmt.Println("send_judge_created_by_admin_welcome_email started")
	queue, err := queue.GetQueue("send_judge_created_by_admin_welcome_email")
	if err != nil {
		fmt.Println("Error getting queue")
		return nil, err
	}
	err = queue.StartConsuming(1, time.Second)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	taskConsumer := &JudgeCreatedByAdminWelcomeEmailTaskConsumer{}
	str, err := queue.AddConsumer("judge_created_by_admin_welcome_email_list", taskConsumer)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	fmt.Println(str)
	return taskConsumer, nil
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
		exports.MySugarLogger.Errorln(err)
		return
	}
	d.Ack()
	c.Ch <- struct{}{}
}
