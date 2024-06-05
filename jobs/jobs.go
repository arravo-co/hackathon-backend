package jobs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/queue"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/email"
)

var inviteListQueue rmq.Queue

type InvitelistTaskConsumer struct {
	Ch chan interface{}
}

func StartConsumingInviteTaskQueue() (*InvitelistTaskConsumer, error) {
	queue, err := queue.GetQueue("invite_list")
	if err != nil {
		fmt.Println("Error getting queue")
		return nil, err
	}
	inviteListQueue = queue
	err = inviteListQueue.StartConsuming(1, time.Second)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	taskConsumer := &InvitelistTaskConsumer{}
	str, err := inviteListQueue.AddConsumer("inviteList", taskConsumer)
	if err != nil {
		return nil, err
	}
	fmt.Println(str)
	return taskConsumer, nil
}

func (c *InvitelistTaskConsumer) Consume(d rmq.Delivery) {
	payload := d.Payload()

	payloadStruct := exports.InvitelistQueuePayload{}
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
		if err := d.Reject(); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
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
		return
	}
	d.Ack()
	c.Ch <- struct{}{}
}
