package exports

import "time"

type InvitelistQueuePayload struct {
	ParticipantId      string    `json:"participant_id"`
	HackathonId        string    `json:"hackathon_id"`
	TeamLeadEmailEmail string    `json:"teamlead_email"`
	InviterEmail       string    `json:"inviter_email"`
	InviteeEmail       string    `json:"invitee_email"`
	InviterName        string    `json:"inviter_name"`
	TimeSent           time.Time `json:"time_sent"`
}

type ParticipantWelcomeEmailQueuePayload struct {
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Password    string    `json:"password"`
	TTL         time.Time `json:"ttl"`
	Link        string    `json:"link"`
	Token       string    `json:"token"`
	InviterName string    `json:"inviter_name"`
}

type AdminWelcomeEmailQueuePayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AdminCreatedByAdminWelcomeEmailQueuePayload struct {
	Email       string `json:"email"`
	AdminName   string `json:"admin_name"`
	InviterName string `json:"last_name"`
	Password    string `json:"password"`
}

type JudgeRegisteredPayload struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	InviterName string `json:"inviter_name"`
}

type PlayQueuePayload struct {
	Time time.Time
}
