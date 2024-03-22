package entity

import (
	"errors"
	"fmt"
	"sync"

	"github.com/arravoco/hackathon_backend/cache"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/events"
	eventsdtos "github.com/arravoco/hackathon_backend/events_dtos"
	"github.com/arravoco/hackathon_backend/utils"
)

var wg sync.WaitGroup
var hackathon_id string = "hackathonID"

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
	events.EmitParticipantAccountCreated(&eventsdtos.ParticipantAccountCreatedEventData{
		ParticipantEmail: dataResponse.Email,
		LastName:         dataResponse.LastName,
		FirstName:        dataResponse.FirstName,
		EventData:        eventsdtos.EventData{EventName: "ParticipantAccountCreated"},
		ParticipantType:  "INDIVIDUAL",
	})
	// emit created event

	return dataResponse, err
}

func (p *Participant) RegisterTeam(input dtos.RegisterNewTeamParticipantDTO) (interface{}, error) {
	teamMembers := []eventsdtos.TeamParticipantInfo{}
	dataInput := &data.CreateTeamParticipantRecordData{
		TeamLeadEmail:       input.TeamLeadEmail,
		HackathonId:         hackathon_id,
		TeamName:            input.TeamName,
		CoParticipantEmails: input.CoParticipantEmails,
	}
	insertId, err := data.CreateTeamParticipantRecord(dataInput)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", insertId)
	for _, v := range append(input.CoParticipantEmails, input.TeamLeadEmail) {
		wg.Add(1)
		v := v
		go func() {
			member, err := p.CreateCoParticipantAccount(v)
			if err != nil {
				wg.Done()
				return
			}
			teamMembers = append(teamMembers, eventsdtos.TeamParticipantInfo{
				Email:    member.Email,
				Password: member.Password,
			})
			wg.Done()
		}()
	}
	wg.Wait()
	// emit created event

	events.EmitParticipantAccountCreated(&eventsdtos.ParticipantAccountCreatedEventData{
		TeamParticipants: teamMembers,
		TeamLeadEmail:    input.TeamLeadEmail,
		EventData:        eventsdtos.EventData{EventName: "ParticipantAccountCreated"},
		TeamName:         input.TeamName,
		ParticipantType:  "TEAM",
	})

	return dataInput, err
}

type CoParticipantCreatedData struct {
	Email    string
	Password string
}

func (p *Participant) CreateCoParticipantAccount(email string) (*CoParticipantCreatedData, error) {
	isFound := cache.FindEmailInCache(email)
	if isFound == true {
		utils.MySugarLogger.Error("Email already exists")
		return nil, errors.New("email already exists")
	}
	password := utils.GeneratePassword()
	passwordHash, _ := utils.GenerateHashPassword(password)
	createDataInput :=
		&data.CreateAccountData{
			Email:        email,
			PasswordHash: passwordHash,
			Role:         "PARTICIPANT"}
	_, err := data.CreateAccount(createDataInput)
	if err != nil {
		return nil, err
	}
	cache.AddEmailToCache(email)
	return &CoParticipantCreatedData{
		Email:    createDataInput.Email,
		Password: password,
	}, nil
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
