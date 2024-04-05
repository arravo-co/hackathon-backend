package exports

import "time"

type InvitelistQueuePayload struct {
	InviterEmail string    `json:"inviter_email"`
	InviteeEmail string    `json:"invitee_email"`
	InviterName  string    `json:"inviter_name"`
	TimeSent     time.Time `json:"time_sent"`
}

type AdminWelcomeEmailQueuePayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
