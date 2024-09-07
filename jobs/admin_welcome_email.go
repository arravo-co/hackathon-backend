package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/rmqUtils"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/authutils"
	"github.com/arravoco/hackathon_backend/utils/email"
)

type AdminWelcomeEmailTaskConsumer struct {
	Ch chan interface{}
}

func StartAdminWelcomeEmailQueue() (*AdminWelcomeEmailTaskConsumer, error) {
	queue, err := rmqUtils.GetQueue("send_admin_welcome_email")
	if err != nil {
		fmt.Println("Error getting queue")
		return nil, err
	}
	var adminWelcomeEmailQueue rmq.Queue = queue
	err = adminWelcomeEmailQueue.StartConsuming(1, time.Second)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	taskConsumer := &AdminWelcomeEmailTaskConsumer{}
	str, err := adminWelcomeEmailQueue.AddConsumer("admin_welcome_email_list", taskConsumer)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println("send_admin_welcome_email queue consuming")
	fmt.Println(str)
	return taskConsumer, nil
}

func (c *AdminWelcomeEmailTaskConsumer) Consume(d rmq.Delivery) {
	payload := d.Payload()

	payloadStruct := exports.AdminWelcomeEmailQueuePayload{}
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
	ttl := time.Now().AddDate(0, 0, 7)
	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	dataToken, err := auth.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
		Email: payloadStruct.Email,
		TTL:   ttl,
	})
	if err != nil {
		fmt.Println(err.Error())
		if err := d.Reject(); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
	}
	url, err := config.GetRemoteServerURL()
	if err != nil {

		return
	}
	link, err := utils.GenerateEmailVerificationLink(&exports.EmailVerificationLinkPayload{
		Token:     dataToken.Token,
		TTL:       dataToken.TTL,
		Email:     dataToken.TokenTypeValue,
		ServerUrl: url,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = email.SendAdminWelcomeEmail(&email.SendAdminWelcomeEmailData{
		FirstName: payloadStruct.FirstName,
		LastName:  payloadStruct.LastName,
		Email:     payloadStruct.Email,
		Subject:   "Invitation to Join Arravo Hackathon Link As An Administrator",
		TTL:       7,
		Token:     dataToken.Token,
		Link:      link,
	})
	if err != nil {
		exports.MySugarLogger.Errorln(err)
		return
	}
	d.Ack()
	c.Ch <- struct{}{}
}
