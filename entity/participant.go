package entity

import (
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
)

type Participant struct {
	FirstName       string
	LastName        string
	Email           string
	PasswordHash    string
	Gender          string
	State           string
	GithubAddress   string
	LinkedInAddress string
	Role            string
}

func (p *Participant) Register(input dtos.RegisterParticipantDTO) (*data.CreateParticipantAccountData, error) {
	dataInput := &data.CreateParticipantAccountData{}
	return data.CreateParticipantAccount(dataInput)
}
