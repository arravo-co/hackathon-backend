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

type AdminCreatedByAdminWelcomeEmailTaskConsumer struct {
}

func init() {
	queue, err := queue.GetQueue("send_admin_created_by_admin_welcome_email")
	if err != nil {
		fmt.Println("Error getting queue")
		fmt.Println()
	}
	adminWelcomeEmailQueue = queue
	err = adminWelcomeEmailQueue.StartConsuming(1, time.Second)
	if err != nil {
		fmt.Println(err)
		fmt.Println()
	}
	taskConsumer := &AdminCreatedByAdminWelcomeEmailTaskConsumer{}
	adminWelcomeEmailQueue.AddConsumer("admin_welcome_email_list", taskConsumer)
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
	fmt.Println("email tokens prepared to be sent")
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
	fmt.Println("Preparing to send email to new admin created by %s...", payloadStruct.InviterName)
	err = email.SendAdminCreatedByAdminWelcomeEmail(&email.SendAdminCreatedByAdminWelcomeEmailData{
		Email:       payloadStruct.Email,
		InviterName: payloadStruct.InviterName,
		AdminName:   payloadStruct.AdminName,
		Subject:     "Welcome to Arravo's Hackathon - Confirm Your Email Address",
		Token:       payloadStruct.Token,
		TTL:         payloadStruct.TTL.Minute(),
		Link:        link,
		Password:    payloadStruct.Password,
	})
	if err != nil {
		exports.MySugarLogger.Errorln(err)
		return
	}
	d.Ack()
}
