package listeners

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/queue"
	"github.com/arravoco/hackathon_backend/utils/authutils"
)

func HandleAdminCreatedEvent(eventDTOData *exports.AdminAccountCreatedEventData, otherParams ...interface{}) {
	/*
		ttl := time.Now().AddDate(0, 0, 7)
		dataToken, err := authutils.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
			Email: eventDTOData.Email,
			TTL:   ttl,
		})
		if err != nil {
			exports.MySugarLogger.Error(err)
			return
		}
	*/
	queue, err := queue.GetQueue("send_admin_welcome_email")
	if err != nil {
		return
	}
	payload := exports.AdminWelcomeEmailQueuePayload{
		Email:     eventDTOData.Email,
		FirstName: eventDTOData.FirstName,
		LastName:  eventDTOData.LastName,
	}
	byt, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = queue.PublishBytes(byt)
	if err != nil {
		fmt.Println(err)
		return
	}
	/*
		err = email.SendAdminWelcomeEmail(&email.SendAdminWelcomeEmailData{Email: eventDTOData.Email,
			LastName:  eventDTOData.LastName,
			FirstName: eventDTOData.FirstName,
			Subject:   " Welcome to Arravo's Hackathon - Confirm Your Email Address",
			Token:     dataToken.Token,
			TTL:       dataToken.TTL.Minute(),
		})

		if err != nil {
			fmt.Println(err)
			return
		}
	*/
}

func HandleAdminCreatedByAdminEvent(eventDTOData *exports.AdminAccountCreatedByAdminEventData, otherParams ...interface{}) {
	fmt.Println("listener has recieved admin_created_by_admin_event")
	ttl := time.Now().Add(time.Minute * 15)
	dataToken, err := authutils.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
		Email: eventDTOData.Email,
		TTL:   ttl,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return
	}
	queue, err := queue.GetQueue("send_admin_created_by_admin_welcome_email")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Queue created")
	payload := exports.AdminCreatedByAdminWelcomeEmailQueuePayload{
		Email:       eventDTOData.Email,
		AdminName:   eventDTOData.AdminName,
		InviterName: eventDTOData.InviterName,
		Password:    eventDTOData.Password,
		TTL:         dataToken.TTL,
		Token:       dataToken.Token,
	}
	byt, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Go to")
	}
	fmt.Println("Queue payload published")
	queue.PublishBytes(byt)

}
