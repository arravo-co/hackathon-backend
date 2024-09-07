package exports

import "time"

type AddedToInvitelistPublishPayload struct {
	ParticipantId      string    `json:"participant_id"`
	HackathonId        string    `json:"hackathon_id"`
	TeamLeadEmailEmail string    `json:"teamlead_email"`
	InviterEmail       string    `json:"inviter_email"`
	InviteeEmail       string    `json:"invitee_email"`
	InviterName        string    `json:"inviter_name"`
	TimeSent           time.Time `json:"time_sent"`
}

type ParticipantRegisteredPublishPayload struct {
	AccountId        string `json:"account_id"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	ParticipantEmail string `json:"participant_email"`
	ParticipantId    string `json:"participant_id"`
	TeamLeadEmail    string `json:"team_lead_email"`
	TeamName         string `json:"team_name"`
	TeamRole         string `json:"team_role"`
	TeamLeadName     string `json:"team_lead_name"`
	ParticipantType  string `json:"participant_type"`
}

type JudgeRegisteredPublishPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type JudgeRegisteredByAdminPublishPayload struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	InviterName string `json:"inviter_name"`
}

type AdminRegisteredPublishPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AdminRegisteredByAdminPublishPayload struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	InviterName  string `json:"inviter_name"`
	InviterEmail string `json:"inviter_email"`
}

type AddedToInviteListPublishPayload struct {
	ParticipantId      string    `json:"participant_id"`
	HackathonId        string    `json:"hackathon_id"`
	TeamLeadEmailEmail string    `json:"teamlead_email"`
	InviterEmail       string    `json:"inviter_email"`
	InviteeEmail       string    `json:"invitee_email"`
	InviterName        string    `json:"inviter_name"`
	TimeSent           time.Time `json:"time_sent"`
}
