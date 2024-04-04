package exports

import "time"

type InvitelistQueuePayload struct {
	InviterEmail string    `json:"inviter_email"`
	InviteeEmail string    `json:"invitee_email"`
	InviterName  string    `json:"inviter_name"`
	TimeSent     time.Time `json:"time_sent"`
}
