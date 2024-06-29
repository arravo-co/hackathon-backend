package services

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/arravoco/hackathon_backend/cache"
	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/events"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/repository"
	valueobjects "github.com/arravoco/hackathon_backend/value_objects"
)

type CompleteNewTeamMemberRegistrationEntityData struct {
	FirstName        string
	LastName         string
	Email            string
	Password         string
	PhoneNumber      string
	ConfirmPassword  string
	Gender           string
	Skillset         []string
	State            string
	DOB              string
	ParticipantId    string
	TeamLeadEmail    string
	HackathonId      string
	TeamRole         string
	EmploymentStatus string
	ExperienceLevel  string
	Motivation       string
}

type RemoveMemberFromTeamData struct {
	HackathonId   string `bson:"hackathon_id"`
	ParticipantId string `bson:"participant_id"`
	MemberEmail   string `bson:"email"`
}

type ParticipantService struct {
	ParticipantRepository *repository.ParticipantRepository
	SolutionRepository    *repository.SolutionRepository
}

func (s *ParticipantService) CompleteNewTeamMemberRegistration(q *query.Query, input *CompleteNewTeamMemberRegistrationEntityData) (*entity.Participant, error) {
	passwordHash, err := exports.GenerateHashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	dob, err := time.Parse("2006-04-02", input.DOB)
	if err != nil {
		return nil, err
	}
	isEmailInCache := cache.FindEmailInCache(input.Email)
	if isEmailInCache {
		//return nil, errors.New("email is already existing")
	}
	partDoc, err := s.ParticipantRepository.AddMemberToParticipatingTeam(&exports.AddMemberToParticipatingTeamData{
		HackathonId:   input.HackathonId,
		ParticipantId: input.ParticipantId,
		Email:         input.Email,
		Role:          "PARTICIPANT",
		TeamRole:      input.TeamRole,
	})
	if err != nil {
		return nil, err
	}
	acc, err := CreateTeamMemberAccount(q, &exports.CreateTeamMemberAccountData{
		CreateAccountData: exports.CreateAccountData{
			Email:        input.Email,
			LastName:     input.LastName,
			FirstName:    input.FirstName,
			PasswordHash: passwordHash,
			Gender:       input.Gender,
			Role:         "PARTICIPANT",
			PhoneNumber:  input.PhoneNumber,
			HackathonId:  input.HackathonId,
			State:        input.State,
			Status:       "EMAIL_VERIFIED",
		},
		TeamRole:         "TEAM_MEMBER",
		DOB:              dob,
		ParticipantId:    input.ParticipantId,
		Skillset:         input.Skillset,
		Motivation:       input.Motivation,
		EmploymentStatus: input.EmploymentStatus,
		ExperienceLevel:  input.ExperienceLevel,
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
	return &entity.Participant{
		AccountStatus:       acc.Status,
		ParticipationStatus: partDoc.Status,
		ParticipantType:     partDoc.Type,
		TeamRole:            input.TeamRole,
		AccountRole:         acc.AccountRole,
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

func InviteToTeam(q query.Query, dataInput *exports.AddToTeamInviteListData) (interface{}, error) {
	partRepo := repository.NewParticipantRepository(&q)
	partEnt, err := partRepo.FillParticipantInfo(dataInput.InviterEmail)
	if err != nil {
		return nil, err
	}
	partRepo.SetEntity(partEnt)
	res, err := partRepo.InviteToTeam(dataInput)
	if err != nil {
		return nil, err
	}
	queue, err := data.GetQueue("invite_list")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	queuePayload := exports.InvitelistQueuePayload{
		InviterEmail:       dataInput.InviterEmail,
		InviterName:        partEnt.FirstName,
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
	err = queue.PublishBytes(byt)
	if err != nil {
		fmt.Println(err.Error())
	}
	return res, nil
}

func RegisterTeamLead(q *query.Query, input dtos.RegisterNewParticipantDTO) (*entity.Participant, error) {
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

	dob, err := time.Parse("2006-04-02", input.DOB)
	if err == nil {
		return nil, err
	}
	acc, err := q.CreateParticipantAccount(&exports.CreateParticipantAccountData{
		ParticipantId:    participantId,
		Skillset:         input.Skillset,
		DOB:              dob,
		Motivation:       input.Motivation,
		EmploymentStatus: input.EmploymentStatus,
		ExperienceLevel:  input.ExperienceLevel,
		CreateAccountData: exports.CreateAccountData{
			Email:        input.Email,
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			Role:         "PARTICIPANT",
			Gender:       input.Gender,
			PasswordHash: passwordHash,
			State:        input.State,
			PhoneNumber:  input.PhoneNumber,
			HackathonId:  config.GetHackathonId(),
			Status:       "EMAIL_UNVERIFIED",
		},
	})
	if err != nil {
		return nil, err
	}
	dataInput := &exports.CreateParticipantRecordData{
		TeamLeadEmail:    input.Email,
		HackathonId:      config.GetHackathonId(),
		TeamName:         input.TeamName,
		CoParticipants:   []exports.CoParticipant{},
		ParticipantId:    participantId,
		ParticipantEmail: input.Email,
		Type:             "TEAM",
	}

	particicipantDoc, err := q.CreateParticipantRecord(dataInput)
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

	return &entity.Participant{
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
		CoParticipants:      []valueobjects.CoParticipantInfo{},
		Motivation:          input.Motivation,
		ExperienceLevel:     input.ExperienceLevel,
		EmploymentStatus:    input.EmploymentStatus,
		CreatedAt:           acc.CreatedAt,
		UpdatedAt:           acc.UpdatedAt,
	}, err
}

func RemoveMemberFromTeam(dataInput *RemoveMemberFromTeamData) (*entity.TeamMemberAccount, error) {
	q := query.GetDefaultQuery()
	_, err := q.RemoveMemberFromParticipatingTeam(&exports.RemoveMemberFromParticipatingTeamData{
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

func FillTeamMemberInfo(account *exports.AccountDocument) *entity.TeamMemberAccount {
	info := &entity.TeamMemberAccount{}
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

func GetTeamMembersInfo(p *entity.Participant) ([]entity.TeamMemberAccount, error) {
	team := []entity.TeamMemberAccount{}
	fmt.Println(p.CoParticipants)
	var emails []string
	emails = []string{}
	for _, obj := range p.CoParticipants {
		emails = append(emails, obj.Email)
	}
	accounts, err := data.GetAccountsByEmails(emails)
	fmt.Println(accounts)
	for _, acc := range accounts {
		oo := FillTeamMemberInfo(&acc)
		team = append(team, *oo)
	}
	return team, err
}

func GetParticipantsInfo() ([]entity.Participant, error) {
	participants, err := data.GetParticipantsRecords()
	if err != nil {
		return nil, err
	}
	partIds := []string{}

	for _, part := range participants {
		partIds = append(partIds, part.ParticipantId)
	}

	accs, err := data.GetAccountsByParticipantIds(partIds)
	if err != nil {
		return nil, err
	}
	pEs := []entity.Participant{}
	for _, p := range participants {
		pE := entity.Participant{}
		cs := []valueobjects.CoParticipantInfo{}
		for _, a := range accs {
			if p.TeamLeadEmail == a.Email || p.ParticipantEmail == a.Email {
				pE.FirstName = a.FirstName
				pE.LastName = a.LastName
				pE.Email = a.Email
				pE.ParticipationStatus = p.Status
				pE.AccountRole = a.Role
				pE.DOB = a.DOB
				pE.Age = time.Now().Year() - a.DOB.Year()
				pE.Gender = a.Gender
				pE.HackathonId = a.HackathonId
				pE.TeamName = p.TeamName
				pE.ParticipantType = p.Type
				pE.InviteList = p.InviteList
				pE.State = a.State
				pE.ParticipantId = a.ParticipantId
				pE.TeamLeadEmail = p.TeamLeadEmail
				pE.AccountStatus = a.Status
				pE.PhoneNumber = a.PhoneNumber
				pE.Skillset = a.Skillset
				if p.Type == "TEAM" {
					pE.TeamRole = "TEAM_LEAD"
				}
			} else {
				for _, c := range p.CoParticipants {
					if c.Email == a.Email {
						cs = append(cs, valueobjects.CoParticipantInfo{
							FirstName:     a.FirstName,
							LastName:      a.LastName,
							Email:         a.Email,
							AccountRole:   a.Role,
							Gender:        a.Gender,
							State:         a.State,
							DOB:           a.DOB,
							Age:           time.Now().Year() - a.DOB.Year(),
							AccountStatus: a.Status,
							Skillset:      a.Skillset,
							PhoneNumber:   a.PhoneNumber,
							HackathonId:   a.HackathonId,
							ParticipantId: a.ParticipantId,
							TeamRole:      c.Email,
							CreatedAt:     a.CreatedAt,
							UpdatedAt:     a.UpdatedAt,
						})
					}
				}
			}
		}
		fmt.Println(p.ParticipantId)
		pE.CoParticipants = cs
		pEs = append(pEs, pE)
	}
	return pEs, err
}

func ReconcileParticipantInfo(p *entity.Participant, accountDataInput *exports.AccountDocument, particicipantDataInput *exports.ParticipantDocument) error {
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

func GetParticipantInfo(q *query.Query, participantId string) (*entity.Participant, error) {
	participant, err := q.GetParticipantRecord(participantId)
	if err != nil {
		return nil, err
	}
	partIds := []string{participant.ParticipantId}

	accs, err := q.GetAccountsByParticipantIds(partIds)
	if err != nil {
		return nil, err
	}
	pE := entity.Participant{}
	cs := []valueobjects.CoParticipantInfo{}

	for _, a := range accs {
		if participant.TeamLeadEmail == a.Email || participant.ParticipantEmail == a.Email {
			pE.FirstName = a.FirstName
			pE.LastName = a.LastName
			pE.Email = a.Email
			pE.ParticipationStatus = participant.Status
			pE.AccountRole = a.Role
			pE.TeamLeadEmail = participant.TeamLeadEmail
			pE.TeamName = participant.TeamName
			pE.State = a.State
			pE.DOB = a.DOB
			pE.Age = time.Now().Year() - a.DOB.Year()
			pE.Gender = a.Gender
			pE.HackathonId = a.HackathonId
			pE.AccountStatus = a.Status
			pE.InviteList = participant.InviteList
			pE.Skillset = a.Skillset
			pE.PhoneNumber = a.PhoneNumber
			pE.ParticipantId = participant.ParticipantId
			pE.ParticipantType = participant.Type
			pE.Motivation = a.Motivation
			pE.ExperienceLevel = a.ExperienceLevel
			pE.EmploymentStatus = a.EmploymentStatus
			pE.CreatedAt = a.CreatedAt
			pE.UpdatedAt = a.UpdatedAt
			if participant.Type == "TEAM" {
				pE.TeamRole = "TEAM_LEAD"
			}
		} else {
			for _, c := range participant.CoParticipants {
				if c.Email == a.Email {
					cs = append(cs, valueobjects.CoParticipantInfo{
						FirstName:        a.FirstName,
						Email:            c.Email,
						LastName:         a.LastName,
						PhoneNumber:      a.PhoneNumber,
						Gender:           a.Gender,
						State:            a.State,
						DOB:              a.DOB,
						Age:              time.Now().Year() - a.DOB.Year(),
						AccountStatus:    a.Status,
						ParticipantId:    a.ParticipantId,
						HackathonId:      a.HackathonId,
						TeamRole:         c.Role,
						AccountRole:      a.Role,
						Skillset:         a.Skillset,
						Motivation:       a.Motivation,
						ExperienceLevel:  a.ExperienceLevel,
						EmploymentStatus: a.EmploymentStatus,
						CreatedAt:        a.CreatedAt,
						UpdatedAt:        a.UpdatedAt,
					})
				}
			}
		}
	}
	pE.CoParticipants = cs

	return &pE, err
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
	slicesOfHash := strings.Split(hashedString, "")
	prefixSlices := slicesOfHash[0:5]
	postFix := slicesOfHash[len(slicesOfHash)-5:]
	sub := strings.Join(append(prefixSlices, postFix...), "")
	return sub, nil
}

func CreateTeamMemberAccount(q *query.Query, dataToSave *exports.CreateTeamMemberAccountData) (*entity.TeamMemberAccount, error) {
	accRepo := repository.NewAccountRepository(q)
	tmAccRepo, err := accRepo.CreateTeamMemberAccount(dataToSave)
	if err != nil {
		return nil, err
	}
	return &entity.TeamMemberAccount{
		Email:             tmAccRepo.Email,
		FirstName:         tmAccRepo.FirstName,
		LastName:          tmAccRepo.LastName,
		Gender:            tmAccRepo.Gender,
		LinkedInAddress:   tmAccRepo.LinkedInAddress,
		PhoneNumber:       tmAccRepo.PhoneNumber,
		ParticipantId:     tmAccRepo.ParticipantId,
		Skillset:          tmAccRepo.Skillset,
		State:             tmAccRepo.State,
		HackathonId:       tmAccRepo.HackathonId,
		TeamRole:          tmAccRepo.TeamRole,
		DOB:               tmAccRepo.DOB,
		IsEmailVerified:   tmAccRepo.IsEmailVerified,
		IsEmailVerifiedAt: tmAccRepo.IsEmailVerifiedAt,
	}, nil
}
