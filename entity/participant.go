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
	Role                string                   `json:"role"`
	ParticipantId       string                   `json:"participant_id"`
	TeamLeadEmail       string                   `json:"team_lead_email"`
	TeamName            string                   `json:"team_name"`
	TeamRole            string                   `json:"team_role"`
	HackathonId         string                   `json:"hackathon_id"`
	Type                string                   `json:"type"`
	CoParticipantEmails []string                 `json:"co_participant_emails"`
	ParticipantEmail    string                   `json:"participant_email"`
	InviteList          []exports.InviteInfo     `json:"invite_list"`
	AccountStatus       string                   `json:"account_status"`
	ParticipationStatus string                   `json:"participation_status"`
	Skillset            []string                 `json:"skillset"`
	PhoneNumber         string                   `json:"phone_number"`
	Solution            exports.SolutionDocument `json:"solution"`
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
		InviterEmail: dataInput.InviterEmail,
		InviterName:  "",
		InviteeEmail: dataInput.Email,
		TimeSent:     time.Now(),
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

func (p *Participant) RegisterNewTeamMember(input *dtos.RegisterNewTeamMemberDTO) (*Participant, error) {
	passwordHash, err := exports.GenerateHashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	dataInput := &exports.CreateTeamMemberAccountData{}
	dob, err := time.Parse(time.RFC3339, input.DOB)
	if err == nil {
		dataInput.DOB = dob
	}
	dataInput.CreateAccountData =
		exports.CreateAccountData{
			Email:        input.Email,
			PasswordHash: passwordHash,
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			Gender:       input.Gender,
			State:        input.State,
			PhoneNumber:  input.PhoneNumber,
			Role:         "PARTICIPANT",
		}
	dataInput.ParticipantId = input.ParticipantId

	isEmailInCache := cache.FindEmailInCache(dataInput.Email)
	if isEmailInCache {
		return nil, errors.New("email is already existing")
	}
	_, err = data.AddMemberToParticipatingTeam(&exports.AddMemberToParticipatingTeamData{
		HackathonId:   config.GetHackathonId(),
		ParticipantId: dataInput.ParticipantId,
		Email:         dataInput.Email,
		Role:          "PARTICIPANT",
	})
	if err != nil {
		return nil, err
	}
	acc, err := data.CreateTeamMemberAccount(&exports.CreateTeamMemberAccountData{
		CreateAccountData: exports.CreateAccountData{
			Email:        dataInput.Email,
			LastName:     dataInput.LastName,
			FirstName:    dataInput.FirstName,
			PasswordHash: passwordHash,
			Gender:       dataInput.Gender,
			Role:         dataInput.Role,
			PhoneNumber:  dataInput.PhoneNumber,
			DOB:          dataInput.DOB,
			HackathonId:  config.GetHackathonId(),
			State:        dataInput.State,
		},
		ParticipantId: dataInput.ParticipantId,
		Skillset:      dataInput.Skillset,
	})
	if err != nil {
		return nil, err
	}
	addedToCache := cache.AddEmailToCache(dataInput.Email)

	if !addedToCache {
		exports.MySugarLogger.Warnln("Email is already in cache")
	}

	// emit created event
	events.EmitParticipantAccountCreated(&exports.ParticipantAccountCreatedEventData{
		ParticipantEmail: acc.Email,
		LastName:         acc.LastName,
		FirstName:        acc.FirstName,
		EventData:        exports.EventData{EventName: "ParticipantAccountCreated"},
		ParticipantType:  "INDIVIDUAL",
	})
	p.FillParticipantInfo(acc.Email)
	return p, nil
}

func (p *Participant) RegisterIndividual(input dtos.RegisterNewParticipantDTO) (*Participant, error) {
	passwordHash, err := exports.GenerateHashPassword(input.Password)
	participantId, err := GenerateParticipantID([]string{input.Email})
	if err != nil {
		return nil, err
	}
	dob, err := time.Parse(time.RFC3339, input.DOB)
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
	dataResponse, err := data.CreateParticipantAccount(dataInput)
	if err != nil {
		return nil, err
	}
	addedToCache := cache.AddEmailToCache(dataInput.Email)

	if !addedToCache {
		exports.MySugarLogger.Warnln("Email is already in cache")
	}
	particicipantDoc, err := data.CreateParticipantRecord(&exports.CreateParticipantRecordData{
		HackathonId:      config.GetHackathonId(),
		Type:             "INDIVIDUAL",
		ParticipantEmail: input.Email,
		ParticipantId:    participantId,
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
	return &Participant{
		ParticipantId:       participantId,
		Email:               dataInput.Email,
		LastName:            dataInput.LastName,
		FirstName:           dataInput.FirstName,
		Gender:              dataInput.Gender,
		DOB:                 dob,
		Role:                dataInput.Role,
		State:               dataInput.State,
		ParticipationStatus: dataResponse.Status,
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

	dob, err := time.Parse(time.RFC3339, input.DOB)
	if err == nil {
		return nil, err
	}
	data.CreateParticipantAccount(&exports.CreateParticipantAccountData{
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
		},
	})

	dataInput := &exports.CreateParticipantRecordData{
		TeamLeadEmail:       input.Email,
		HackathonId:         config.GetHackathonId(),
		TeamName:            input.TeamName,
		CoParticipantEmails: []string{},
		ParticipantId:       participantId,
		ParticipantEmail:    input.Email,
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
		Email:            input.Email,
		Role:             "PARTICIPANT",
		FirstName:        input.FirstName,
		LastName:         input.LastName,
		Gender:           input.Gender,
		State:            input.State,
		Skillset:         input.Skillset,
		PhoneNumber:      input.PhoneNumber,
		DOB:              dob,
		TeamLeadEmail:    particicipantDoc.TeamLeadEmail,
		TeamName:         particicipantDoc.TeamName,
		ParticipantId:    participantId,
		ParticipantEmail: particicipantDoc.ParticipantEmail,
		TeamRole:         "TEAM_LEAD",
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
	p.ParticipationStatus = accountData.Status
	p.ParticipationStatus = particicipantDocData.Status
	p.passwordHash = accountData.PasswordHash
	p.FirstName = accountData.FirstName
	p.LastName = accountData.LastName
	p.Gender = accountData.Gender
	p.State = accountData.State
	p.Role = accountData.Role
	p.TeamName = particicipantDocData.TeamName
	p.TeamLeadEmail = particicipantDocData.TeamLeadEmail
	p.HackathonId = particicipantDocData.HackathonId
	p.Type = particicipantDocData.Type
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
	}
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
