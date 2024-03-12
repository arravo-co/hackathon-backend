package entity

import (
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/utils"
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

func (p *Participant) Register(input dtos.RegisterNewParticipantDTO) (*data.CreateParticipantAccountData, error) {
	passwordHash, err := utils.GenerateHashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	dataInput := &data.CreateParticipantAccountData{
		Email:           input.Email,
		PasswordHash:    passwordHash,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Gender:          input.Gender,
		GithubAddress:   input.GithubAddress,
		LinkedInAddress: input.LinkedInAddress,
		State:           input.State,
	}
	dataResponse, err := data.CreateParticipantAccount(dataInput)
	// emit created event

	return dataResponse, err
}
