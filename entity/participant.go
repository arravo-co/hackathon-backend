package entity

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/arravoco/hackathon_backend/cache"
	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/events"
	"github.com/arravoco/hackathon_backend/exports"
)

var wg sync.WaitGroup

// AddMemberToParticipatingTeam
type Participant struct {
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	Email               string `json:"email"`
	passwordHash        string
	Gender              string                   `json:"gender"`
	State               string                   `json:"state"`
	Age                 int                      `json:"age"`
	DOB                 time.Time                `json:"dob"`
	AccountRole         string                   `json:"role"`
	ParticipantId       string                   `json:"participant_id"`
	TeamLeadEmail       string                   `json:"team_lead_email"`
	TeamName            string                   `json:"team_name"`
	TeamRole            string                   `json:"team_role"`
	HackathonId         string                   `json:"hackathon_id"`
	ParticipantType     string                   `json:"type"`
	CoParticipantEmails []string                 `json:"co_participant_emails"`
	ParticipantEmail    string                   `json:"participant_email"`
	InviteList          []exports.InviteInfo     `json:"invite_list"`
	AccountStatus       string                   `json:"account_status"`
	ParticipationStatus string                   `json:"participation_status"`
	Skillset            []string                 `json:"skillset"`
	PhoneNumber         string                   `json:"phone_number"`
	Solution            exports.SolutionDocument `json:"solution"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           time.Time                `json:"updated_at"`
}

type CompleteNewTeamMemberRegistrationEntityData struct {
	FirstName       string
	LastName        string
	Email           string
	Password        string
	PhoneNumber     string
	ConfirmPassword string
	Gender          string
	Skillset        []string
	State           string
	DOB             string
	ParticipantId   string
	TeamLeadEmail   string
	HackathonId     string
	TeamRole        string
}

type TeamMemberAccount struct {
	Email           string   `json:"email"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Gender          string   `json:"gender,omitempty"`
	LinkedInAddress string   `json:"linkedIn_address,omitempty"`
	PhoneNumber     string   `json:"phone_number,omitempty"`
	Skillset        []string `json:"skillset,omitempty"`
	ParticipantId   string   `json:"participant_id,omitempty"`
	HackathonId     string   `json:"hackathon_id"`

	State             string    `json:"state,omitempty"`
	TeamRole          string    `json:"role,omitempty"`
	DOB               time.Time `json:"dob,omitempty"`
	IsEmailVerified   bool      `json:"is_email_verified,omitempty"`
	IsEmailVerifiedAt time.Time `json:"is_email_verified_at,omitempty"`
	Status            string    `json:"status"`
}

func (p Participant) InviteToTeam(dataInput *exports.AddToTeamInviteListData) (interface{}, error) {
	res, err := data.AddToTeamInviteList(dataInput)
	if err != nil {
		return nil, err
	}
	q, err := data.GetQueue("invite_list")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	queuePayload := exports.InvitelistQueuePayload{
		InviterEmail:       dataInput.InviterEmail,
		InviterName:        p.FirstName,
		InviteeEmail:       dataInput.Email,
		ParticipantId:      dataInput.ParticipantId,
		HackathonId:        dataInput.HackathonId,
		TeamLeadEmailEmail: dataInput.InviterEmail,
		TimeSent:           time.Now(),
	}
	byt, err := json.Marshal(queuePayload)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = q.PublishBytes(byt)
	if err != nil {
		fmt.Println(err.Error())
	}
	return res, nil
}

func CompleteNewTeamMemberRegistration(input *CompleteNewTeamMemberRegistrationEntityData) (*Participant, error) {
	passwordHash, err := exports.GenerateHashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	dob, err := time.Parse("2006-01-02", input.DOB)
	if err != nil {
		return nil, err
	}
	isEmailInCache := cache.FindEmailInCache(input.Email)
	if isEmailInCache {
		//return nil, errors.New("email is already existing")
	}
	partDoc, err := data.AddMemberToParticipatingTeam(&exports.AddMemberToParticipatingTeamData{
		HackathonId:   input.HackathonId,
		ParticipantId: input.ParticipantId,
		Email:         input.Email,
		Role:          "PARTICIPANT",
		TeamRole:      input.TeamRole,
	})
	if err != nil {
		return nil, err
	}
	acc, err := data.CreateTeamMemberAccount(&exports.CreateTeamMemberAccountData{
		CreateAccountData: exports.CreateAccountData{
			Email:        input.Email,
			LastName:     input.LastName,
			FirstName:    input.FirstName,
			PasswordHash: passwordHash,
			Gender:       input.Gender,
			Role:         "PARTICIPANT",
			PhoneNumber:  input.PhoneNumber,
			DOB:          dob,
			HackathonId:  input.HackathonId,
			State:        input.State,
			Status:       "REVIEWED",
		},
		ParticipantId: input.ParticipantId,
		Skillset:      input.Skillset,
	})
	if err != nil {
		return nil, err
	}
	addedToCache := cache.AddEmailToCache(input.Email)

	if !addedToCache {
		exports.MySugarLogger.Warnln("Email is already in cache")
	}

	// emit created event
	events.EmitParticipantAccountCreated(&exports.ParticipantAccountCreatedEventData{
		ParticipantEmail: acc.Email,
		LastName:         acc.LastName,
		FirstName:        acc.FirstName,
		EventData:        exports.EventData{EventName: "ParticipantAccountCreated"},
		ParticipantType:  "TEAM",
	})
	return &Participant{
		AccountStatus:       acc.Status,
		ParticipationStatus: partDoc.Status,
		ParticipantType:     partDoc.Type,
		TeamRole:            input.TeamRole,
		AccountRole:         acc.Role,
		Email:               input.Email,
		FirstName:           input.FirstName,
		LastName:            input.LastName,
		DOB:                 dob,
		Gender:              input.Gender,
		HackathonId:         input.HackathonId,
		State:               input.State,
		Skillset:            input.Skillset,
		Age:                 time.Now().Year() - dob.Year(),
		PhoneNumber:         input.PhoneNumber,
		ParticipantId:       input.ParticipantId,
		CreatedAt:           acc.CreatedAt,
		UpdatedAt:           acc.UpdatedAt,
	}, nil
}

func (p *Participant) RegisterIndividual(input dtos.RegisterNewParticipantDTO) (*Participant, error) {
	passwordHash, _ := exports.GenerateHashPassword(input.Password)
	participantId, err := GenerateParticipantID([]string{input.Email})
	if err != nil {
		return nil, err
	}
	dob, err := time.Parse("2006-Jan-02", input.DOB)
	if err == nil {
		return nil, err
	}
	dataInput := &exports.CreateParticipantAccountData{
		CreateAccountData: exports.CreateAccountData{
			Email:        input.Email,
			PasswordHash: passwordHash,
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			Gender:       input.Gender,
			State:        input.State,
			Role:         "PARTICIPANT",
			DOB:          dob,
		},
		Skillset:      input.Skillset,
		ParticipantId: participantId}
	isEmailInCache := cache.FindEmailInCache(dataInput.Email)
	if isEmailInCache {
		//return nil, errors.New("email is already existing")
	}
	accCreated, err := data.CreateParticipantAccount(dataInput)
	if err != nil {
		return nil, err
	}
	addedToCache := cache.AddEmailToCache(dataInput.Email)

	if !addedToCache {
		exports.MySugarLogger.Warnln("Email is already in cache")
	}
	partDoc, err := data.CreateParticipantRecord(&exports.CreateParticipantRecordData{
		HackathonId:      config.GetHackathonId(),
		Type:             "INDIVIDUAL",
		ParticipantEmail: input.Email,
		ParticipantId:    participantId,
	})

	if err != nil {
		return nil, err
	}
	events.EmitParticipantAccountCreated(&exports.ParticipantAccountCreatedEventData{
		ParticipantEmail: accCreated.Email,
		LastName:         accCreated.LastName,
		FirstName:        accCreated.FirstName,
		EventData:        exports.EventData{EventName: "ParticipantAccountCreated"},
		ParticipantType:  "INDIVIDUAL",
	})
	/*
		dt := struct {
			FirstName string `json:"message"`
		}{FirstName: dataResponse.FirstName}
		by, _ := json.Marshal(&dt)
		err = producer.Publish("participant_register", by)
		if err != nil {
			fmt.Println(err.Error())
		}
	*/
	return &Participant{
		HackathonId:         dataInput.HackathonId,
		ParticipantType:     partDoc.Type,
		AccountStatus:       accCreated.Status,
		ParticipantId:       participantId,
		Email:               dataInput.Email,
		LastName:            dataInput.LastName,
		FirstName:           dataInput.FirstName,
		Gender:              dataInput.Gender,
		DOB:                 dob,
		Age:                 time.Now().Year() - dataInput.DOB.Year(),
		AccountRole:         dataInput.Role,
		State:               dataInput.State,
		ParticipationStatus: partDoc.Status,
		Skillset:            accCreated.Skillset,
		PhoneNumber:         accCreated.PhoneNumber,
		CreatedAt:           accCreated.CreatedAt,
		UpdatedAt:           accCreated.UpdatedAt,
	}, nil
}

func (p *Participant) RegisterTeamLead(input dtos.RegisterNewParticipantDTO) (*Participant, error) {
	teamMembers := []exports.TeamParticipantInfo{}
	passwordHash, err := exports.GenerateHashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	participantId, err := GenerateParticipantID([]string{input.Email})
	if err != nil {
		return nil, err
	}

	areEmailsInCache := cache.FindEmailInCache(input.Email)
	if !areEmailsInCache {
		//return nil, errors.New("email already exists")
	}

	dob, err := time.Parse(time.RubyDate, input.DOB)
	if err == nil {
		return nil, err
	}
	acc, err := data.CreateParticipantAccount(&exports.CreateParticipantAccountData{
		ParticipantId: participantId,
		Skillset:      input.Skillset,
		CreateAccountData: exports.CreateAccountData{
			Email:        input.Email,
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			Role:         "PARTICIPANT",
			Gender:       input.Gender,
			PasswordHash: passwordHash,
			State:        input.State,
			DOB:          dob,
			PhoneNumber:  input.PhoneNumber,
			HackathonId:  config.GetHackathonId(),
		},
	})
	if err != nil {
		return nil, err
	}
	dataInput := &exports.CreateParticipantRecordData{
		TeamLeadEmail:       input.Email,
		HackathonId:         config.GetHackathonId(),
		TeamName:            input.TeamName,
		CoParticipantEmails: []string{},
		ParticipantId:       participantId,
		ParticipantEmail:    input.Email,
		Type:                "TEAM",
	}

	particicipantDoc, err := data.CreateParticipantRecord(dataInput)
	if err != nil {
		return nil, err
	}
	// emit created event

	events.EmitParticipantAccountCreated(&exports.ParticipantAccountCreatedEventData{
		TeamParticipants: teamMembers,
		TeamLeadEmail:    input.Email,
		EventData:        exports.EventData{EventName: "ParticipantAccountCreated"},
		ParticipantEmail: input.Email,
		TeamName:         input.TeamName,
		TeamRole:         "TEAM_LEAD",
		ParticipantType:  "TEAM",
	})

	return &Participant{
		Email:               input.Email,
		HackathonId:         dataInput.HackathonId,
		ParticipationStatus: particicipantDoc.Status,
		AccountRole:         "PARTICIPANT",
		ParticipantType:     particicipantDoc.Type,
		FirstName:           input.FirstName,
		LastName:            input.LastName,
		Gender:              input.Gender,
		State:               input.State,
		Skillset:            input.Skillset,
		PhoneNumber:         input.PhoneNumber,
		DOB:                 dob,
		Age:                 time.Now().Year() - dob.Year(),
		TeamLeadEmail:       particicipantDoc.TeamLeadEmail,
		TeamName:            particicipantDoc.TeamName,
		ParticipantId:       participantId,
		ParticipantEmail:    particicipantDoc.ParticipantEmail,
		TeamRole:            "TEAM_LEAD",
		CoParticipantEmails: []string{},
		CreatedAt:           acc.CreatedAt,
		UpdatedAt:           acc.UpdatedAt,
	}, err
}

type CoParticipantCreatedData struct {
	Email         string
	Password      string
	ParticipantId string
}

func (p *Participant) FillParticipantInfo(input string) error {
	accountData, err := data.GetAccountByEmail(input)
	if err != nil {
		return err
	}
	particicipantDocData, err := data.GetParticipantRecord(accountData.ParticipantId)
	if err != nil {
		return err
	}
	p.Email = accountData.Email
	p.AccountStatus = accountData.Status
	p.ParticipationStatus = particicipantDocData.Status
	p.AccountRole = accountData.Role
	p.passwordHash = accountData.PasswordHash
	p.FirstName = accountData.FirstName
	p.LastName = accountData.LastName
	p.Gender = accountData.Gender
	p.State = accountData.State
	p.TeamName = particicipantDocData.TeamName
	p.TeamLeadEmail = particicipantDocData.TeamLeadEmail
	p.HackathonId = particicipantDocData.HackathonId
	p.ParticipantType = particicipantDocData.Type
	p.ParticipantEmail = particicipantDocData.ParticipantEmail
	p.CoParticipantEmails = particicipantDocData.CoParticipantEmails
	p.Solution = particicipantDocData.Solution
	p.ParticipantId = particicipantDocData.ParticipantId
	p.Age = time.Now().Year() - accountData.DOB.Year()
	if particicipantDocData.Type == "TEAM" {
		if particicipantDocData.TeamLeadEmail == accountData.Email {
			p.TeamRole = "TEAM_LEAD"
		} else {
			p.TeamRole = "TEAM_MEMBER"
		}
		p.CoParticipantEmails = particicipantDocData.CoParticipantEmails
	}
	// emit created event

	return nil
}

func FillTeamMemberInfo(account *exports.AccountDocument) *TeamMemberAccount {
	info := &TeamMemberAccount{}
	info.Email = account.Email
	info.Status = account.Status
	info.FirstName = account.FirstName
	info.LastName = account.LastName
	info.Gender = account.Gender
	info.State = account.State
	info.TeamRole = account.Role
	info.HackathonId = account.HackathonId
	info.DOB = account.DOB

	// emit created event

	return info
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

func GenerateParticipantID(emails []string) (string, error) {
	slices.Sort[[]string](emails)
	joined := strings.Join(emails, ":")
	h := sha256.New()
	_, err := h.Write([]byte(joined))
	if err != nil {
		return "", err
	}
	hashByte := h.Sum(nil)
	hashedString := fmt.Sprintf("%x", hashByte)
	fmt.Println("hashedString")
	fmt.Println(joined)
	fmt.Println(hashedString)
	fmt.Println("hashedString")
	slicesOfHash := strings.Split(hashedString, "")
	prefixSlices := slicesOfHash[0:5]
	postFix := slicesOfHash[len(slicesOfHash)-5:]
	sub := strings.Join(append(prefixSlices, postFix...), "")
	return sub, nil
}

type RemoveMemberFromTeamData struct {
	HackathonId   string `bson:"hackathon_id"`
	ParticipantId string `bson:"participant_id"`
	MemberEmail   string `bson:"email"`
}

func (p *Participant) RemoveMemberFromTeam(dataInput *RemoveMemberFromTeamData) (*TeamMemberAccount, error) {
	if p == nil {
		return nil, errors.New("unable to find participant")
	}
	_, err := data.RemoveMemberFromParticipatingTeam(&exports.RemoveMemberFromParticipatingTeamData{
		HackathonId:   dataInput.HackathonId,
		MemberEmail:   dataInput.MemberEmail,
		ParticipantId: dataInput.ParticipantId,
	})
	if err != nil {
		return nil, err
	}
	acc, err := data.DeleteAccount(dataInput.MemberEmail)
	info := FillTeamMemberInfo(acc)
	return info, err
}

func (p *Participant) GetTeamMembersInfo() ([]TeamMemberAccount, error) {
	team := []TeamMemberAccount{}
	fmt.Println(p.CoParticipantEmails)
	accounts, err := data.GetAccountsByEmails(p.CoParticipantEmails)
	fmt.Println(accounts)
	for _, acc := range accounts {
		oo := FillTeamMemberInfo(&acc)
		team = append(team, *oo)
	}
	return team, err
}
