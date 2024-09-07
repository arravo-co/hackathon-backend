package consumerhandlers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/email"
)

func SendInviteEmailQueueHandler(by []byte) error {

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
		TTL:          int(time.Until(ttl).Hours() / 24),
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
