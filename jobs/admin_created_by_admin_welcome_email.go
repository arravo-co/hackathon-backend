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

type AdminCreatedByAdminWelcomeEmailTaskConsumer struct {
	Ch chan interface{}
}

func StartConsumingAdminCreatedByAdminWelcomeEmailQueue() (*AdminCreatedByAdminWelcomeEmailTaskConsumer, error) {
	fmt.Println("send_admin_created_by_admin_welcome_email started")
	queue, err := rmqUtils.GetQueue("send_admin_created_by_admin_welcome_email")
	if err != nil {
		fmt.Println("Error getting queue")
		return nil, err
	}
	err = queue.StartConsuming(1, time.Second)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	taskConsumer := &AdminCreatedByAdminWelcomeEmailTaskConsumer{}
	str, err := queue.AddConsumer("admin_welcome_email_list", taskConsumer)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	fmt.Println(str)
	return taskConsumer, nil
}

func (c *AdminCreatedByAdminWelcomeEmailTaskConsumer) Consume(d rmq.Delivery) {
	fmt.Println("job consumer received job")
	payload := d.Payload()

	payloadStruct := exports.AdminCreatedByAdminWelcomeEmailQueuePayload{}
	err := json.Unmarshal([]byte(payload), &payloadStruct)
	if err != nil {
		fmt.Println(err.Error())
		if err := d.Reject(); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
	}
	ttl := time.Now().AddDate(0, 0, 7)
	dataToken, err := authutils.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
		Email: payloadStruct.Email,
		TTL:   ttl,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return
	}
	fmt.Println("email tokens prepared to be sent")
	serverUrl, err := config.GetRemoteServerURL()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	redirectUrl := "https://hackathon-admin-frontend.onrender.com"
	if err != nil {

	}
	link, err := utils.GenerateEmailVerificationLink(&exports.EmailVerificationLinkPayload{
		Token:       dataToken.Token,
		TTL:         dataToken.TTL,
		Email:       payloadStruct.Email,
		ServerUrl:   serverUrl,
		RedirectUrl: redirectUrl,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("\nPreparing to send email to new admin created by %s...\n", payloadStruct.InviterName)
	err = email.SendAdminCreatedByAdminWelcomeEmail(&email.SendAdminCreatedByAdminWelcomeEmailData{
		Email:       payloadStruct.Email,
		InviterName: payloadStruct.InviterName,
		AdminName:   payloadStruct.AdminName,
		Subject:     "Welcome to Arravo's Hackathon - Confirm Your Email Address",
		Token:       dataToken.Token,
		TTL:         7,
		Link:        link,
		Password:    payloadStruct.Password,
	})
	if err != nil {
		exports.MySugarLogger.Errorln(err)
		return
	}
	d.Ack()
	c.Ch <- struct{}{}
}
