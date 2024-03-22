package entity

import (
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/utils"
)

type Participant struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	passwordHash    string
	Gender          string `json:"gender"`
	State           string `json:"state"`
	GithubAddress   string `json:"github_address"`
	LinkedInAddress string `json:"linkedIn_address"`
	Role            string `json:"role"`
}

func (p *Participant) RegisterIndividual(input dtos.RegisterNewIndividualParticipantDTO) (*data.CreateIndividualParticipantAccountData, error) {
	passwordHash, err := utils.GenerateHashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	dataInput := &data.CreateIndividualParticipantAccountData{}
	dataInput.CreateAccountData =
		data.CreateAccountData{
			Email:        input.Email,
			PasswordHash: passwordHash,
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			Gender:       input.Gender,
			State:        input.State, Role: "PARTICIPANT"}
	dataInput.GithubAddress = input.GithubAddress
	dataInput.LinkedInAddress = input.LinkedInAddress
	dataResponse, err := data.CreateIndividualParticipantAccount(dataInput)
	// emit created event

	return dataResponse, err
}

func (p *Participant) RegisterTeam(input dtos.RegisterNewTeamParticipantDTO) (interface{}, error) {

	dataInput := &data.CreateTeamParticipantRecordData{
		TeamLeadEmail:       input.TeamLeadEmail,
		HackathonId:         "",
		TeamName:            input.TeamName,
		CoParticipantEmails: input.CoParticipantEmails,
	}
	data.CreateTeamParticipantRecord(dataInput)

	for _, v := range input.CoParticipantEmails {
		p.CreateCoParticipantAccount(v)
	}

	password := utils.GeneratePassword()
	passwordHash, err := utils.GenerateHashPassword(password)
	if err != nil {
		return nil, err
	}

	// create account of team lead
	createDataInput :=
		&data.CreateAccountData{
			Email:        input.TeamLeadEmail,
			PasswordHash: passwordHash,
			Role:         "PARTICIPANT"}
	dataResponse, err := data.CreateAccount(createDataInput)

	// emit created event

	return dataResponse, err
}

func (p *Participant) CreateCoParticipantAccount(email string) {
	password := utils.GeneratePassword()
	passwordHash, _ := utils.GenerateHashPassword(password)
	createDataInput :=
		&data.CreateAccountData{
			Email:        email,
			PasswordHash: passwordHash,
			Role:         "PARTICIPANT"}
	_, err = data.CreateAccount(createDataInput)
}

func (p *Participant) GetParticipant(input string) error {
	accountData, err := data.GetAccountByEmail(input)
	if err != nil {
		return err
	}

	p.Email = accountData.Email
	p.passwordHash = accountData.PasswordHash
	p.FirstName = accountData.FirstName
	p.LastName = accountData.LastName
	p.Gender = accountData.Gender
	p.State = accountData.State
	p.Role = accountData.Role
	//p.GithubAddress = accountData.GithubAddress
	p.LinkedInAddress = accountData.LinkedInAddress
	// emit created event

	return err
}

func (p *Participant) ReconcileParticipantInfo(accountDataInput *data.AccountDocument, particicipantDataInput *data.ParticipantDocument) error {
	_, err := data.UpdateParticipantInfoByEmail(&data.UpdateAccountFilter{Email: p.Email}, &data.UpdateAccountDocument{
		FirstName:       accountDataInput.FirstName,
		LastName:        accountDataInput.LastName,
		Gender:          accountDataInput.Gender,
		State:           accountDataInput.State,
		GithubAddress:   particicipantDataInput.GithubAddress,
		LinkedInAddress: accountDataInput.LinkedInAddress,
	})
	return err
}

func (p *Participant) UpdateParticipantInfo(dataInput *dtos.AuthParticipantInfoUpdateDTO) error {
	_, err := data.UpdateParticipantInfoByEmail(&data.UpdateAccountFilter{Email: p.Email}, &data.UpdateAccountDocument{
		FirstName:       dataInput.FirstName,
		LastName:        dataInput.LastName,
		Gender:          dataInput.Gender,
		State:           dataInput.State,
		GithubAddress:   dataInput.GithubAddress,
		LinkedInAddress: dataInput.LinkedInAddress,
	})
	return err
}
