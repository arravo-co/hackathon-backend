package exports

import "time"

type EmailVerificationLinkPayload struct {
	Email string    `json:"email"`
	Token string    `json:"token"`
	TTL   time.Time `json:"ttl"`
}

type TeamInviteLinkPayload struct {
	InviteeEmail       string `json:"invitee_email"`
	TeamLeadEmailEmail string `json:"teamlead_email"`
	HackathonId        string `json:"hackathon_id"`
	ParticipantId      string `json:"participant_id"`
	TTL                int64  `json:"ttl"`
}

type PaswordRecoveryPayload struct {
	Token string    `json:"token"`
	Email string    `json:"inviter_email"`
	TTL   time.Time `json:"ttl"`
}
