package exports

import "time"

type EmailVerificationLinkPayload struct {
	Email string    `json:"email"`
	Token string    `json:"token"`
	TTL   time.Time `json:"ttl"`
}

type TeamInviteLinkPayload struct {
	InviteeEmail string `json:"invitee_email"`
	Token        string `json:"token"`
	InviterEmail string `json:"inviter_email"`
	TTL          int64  `json:"ttl"`
}
