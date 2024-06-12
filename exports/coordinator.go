package exports

import "time"

type SendParticipantWelcomeAndVerificationEmailCoordinatorState struct {
	CoordinatorId  string    `json:"coordinator_id,omitempty"`
	CurrentStateId string    `json:"current_state_id"` // GENERATE_EMAIL_TOKEN EMAIL_TOKEN_GENERATED
	Email          string    `json:"email,omitempty"`
	Token          string    `json:"token,omitempty"`
	TTL            time.Time `json:"ttl,omitempty"`
	TeamName       string    `json:"team_name"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	TeamRole       string    `json:"team_role"` // TEAM_LEAD TEAM_MEMBER
	TeamLeadName   string    `json:"team_leader_name"`
}

type GenerateEmailTokenState struct {
	StateId string `json:"state_id,omitempty"`
}

type EmailTokenGeneratedState struct {
	StateId string    `json:"state_id,omitempty"`
	Token   string    `json:"token,omitempty"`
	Email   string    `json:"email,omitempty"`
	TTL     time.Time `json:"ttl,omitempty"`
}
