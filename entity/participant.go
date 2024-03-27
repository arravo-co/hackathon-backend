package entity

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/arravoco/hackathon_backend/cache"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/events"
	"github.com/arravoco/hackathon_backend/exports"
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

func (p *Participant) RegisterIndividual(input dtos.RegisterNewIndividualParticipantDTO) (*exports.CreateIndividualParticipantAccountData, error) {
	passwordHash, err := exports.GenerateHashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	dataInput := &exports.CreateIndividualParticipantAccountData{}
	dataInput.CreateAccountData =
		exports.CreateAccountData{
			Email:        input.Email,
			PasswordHash: passwordHash,
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			Gender:       input.Gender,
			State:        input.State,
			Role:         "PARTICIPANT",
		}
	dataInput.LinkedInAddress = input.LinkedInAddress
	isEmailInCache := cache.FindEmailInCache(dataInput.Email)
	if isEmailInCache {
		return nil, errors.New("email is already existing")
	}
	dataResponse, err := data.CreateIndividualParticipantAccount(dataInput)
	if err != nil {
		return nil, err
	}
	addedToCache := cache.AddEmailToCache(dataInput.Email)

	if !addedToCache {
		exports.MySugarLogger.Warnln("Email is already in cache")
	}

	particicipantDoc, err := data.CreateParticipantRecord(&exports.CreateParticipantRecordData{
		GithubAddress:    input.GithubAddress,
		HackathonId:      hackathon_id,
		Type:             "INDIVIDUAL",
		ParticipantEmail: input.Email,
	})

	if err != nil {
		return nil, err
	}
	fmt.Printf("\n%+v\n", particicipantDoc)
	// emit created event
	events.EmitParticipantAccountCreated(&exports.ParticipantAccountCreatedEventData{
		ParticipantEmail: dataResponse.Email,
		LastName:         dataResponse.LastName,
		FirstName:        dataResponse.FirstName,
		EventData:        exports.EventData{EventName: "ParticipantAccountCreated"},
		ParticipantType:  "INDIVIDUAL",
	})
	return dataInput, nil
}

func (p *Participant) RegisterTeam(input dtos.RegisterNewTeamParticipantDTO) (*exports.TeamParticipantRecordCreatedData, error) {
	teamMembers := []exports.TeamParticipantInfo{}
	dataInput := &exports.CreateParticipantRecordData{
		TeamLeadEmail:       input.TeamLeadEmail,
		HackathonId:         hackathon_id,
		TeamName:            input.TeamName,
		CoParticipantEmails: input.CoParticipantEmails,
	}
	allEmails := append(input.CoParticipantEmails, input.TeamLeadEmail)
	allEmailInterfaces := []interface{}{}
	for _, v := range allEmails {
		allEmailInterfaces = append(allEmailInterfaces, v)
	}
	areEmailsInCache, err := cache.FindEmailsInCache(allEmailInterfaces)
	if err != nil {
		return nil, errors.New("cannot confirm uniqueness of all participants' emails")
	}
	fmt.Println(areEmailsInCache)
	for index, item := range areEmailsInCache {
		if item {
			return nil, errors.New(strings.Join([]string{"email", allEmails[index], " already exists"}, " "))
		}
	}
	particicipantDoc, err := data.CreateParticipantRecord(dataInput)
	if err != nil {
		return nil, err
	}
	for _, v := range append(input.CoParticipantEmails, input.TeamLeadEmail) {
		wg.Add(1)
		v := v
		go func() {
			member, err := p.CreateCoParticipantAccount(v)
			if err != nil {
				wg.Done()
				return
			}
			teamMembers = append(teamMembers, exports.TeamParticipantInfo{
				Email:    member.Email,
				Password: member.Password,
			})
			wg.Done()
		}()
	}
	wg.Wait()
	// emit created event

	events.EmitParticipantAccountCreated(&exports.ParticipantAccountCreatedEventData{
		TeamParticipants: teamMembers,
		TeamLeadEmail:    input.TeamLeadEmail,
		EventData:        exports.EventData{EventName: "ParticipantAccountCreated"},
		TeamName:         input.TeamName,
		ParticipantType:  "TEAM",
	})

	return &exports.TeamParticipantRecordCreatedData{
		TeamLeadEmail:       dataInput.TeamLeadEmail,
		Role:                "PARTICIPANT",
		ParticipantId:       particicipantDoc.ParticipantId,
		TeamName:            particicipantDoc.TeamName,
		CoParticipantEmails: particicipantDoc.CoParticipantEmails,
		GithubAddress:       particicipantDoc.GithubAddress,
		Type:                particicipantDoc.Type,
		HackathonId:         particicipantDoc.HackathonId,
	}, err
}

type CoParticipantCreatedData struct {
	Email         string
	Password      string
	ParticipantId string
}

func (p *Participant) CreateCoParticipantAccount(email string) (*CoParticipantCreatedData, error) {
	isFound := cache.FindEmailInCache(email)
	if isFound {
		exports.MySugarLogger.Error("Email already exists")
		return nil, errors.New("email already exists")
	}
	password := exports.GeneratePassword()
	passwordHash, _ := exports.GenerateHashPassword(password)
	createDataInput :=
		&exports.CreateAccountData{
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

	return nil
}

func (p *Participant) ReconcileParticipantInfo(accountDataInput *exports.AccountDocument, particicipantDataInput *exports.ParticipantDocument) error {
	_, err := data.UpdateParticipantInfoByEmail(&exports.UpdateAccountFilter{Email: p.Email}, &exports.UpdateAccountDocument{
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
	_, err := data.UpdateParticipantInfoByEmail(&exports.UpdateAccountFilter{Email: p.Email}, &exports.UpdateAccountDocument{
		FirstName:       dataInput.FirstName,
		LastName:        dataInput.LastName,
		Gender:          dataInput.Gender,
		State:           dataInput.State,
		GithubAddress:   dataInput.GithubAddress,
		LinkedInAddress: dataInput.LinkedInAddress,
	})
	return err
}
