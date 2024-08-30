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

type ParticipantService struct {
	ParticipantRepository *repository.ParticipantRepository
	AccountRepository     *repository.AccountRepository
	SolutionRepository    *repository.SolutionRepository
}

func NewParticipantService() *ParticipantService {
	q := query.GetDefaultQuery()
	part := repository.NewParticipantRepository(q)
	acc := repository.NewAccountRepository(q)
	sol := repository.NewSolutionRepository(q)
	return &ParticipantService{
		ParticipantRepository: part,
		AccountRepository:     acc,
		SolutionRepository:    sol,
	}
}

func (s *ParticipantService) CompleteNewTeamMemberRegistration(input *CompleteNewTeamMemberRegistrationEntityData) (*entity.Participant, error) {
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
	acc, err := s.AccountRepository.CreateTeamMemberAccount(&exports.CreateTeamMemberAccountData{
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

func (s *ParticipantService) RegisterTeamLead(input dtos.RegisterNewParticipantDTO) (*entity.Participant, error) {
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
	if err != nil {
		fmt.Println("\nhere: ", err.Error())

		return nil, err
	}
	fmt.Println("\nhere: ", input)
	acc, err := s.AccountRepository.CreateParticipantAccount(&exports.CreateParticipantAccountData{
		ParticipantId:       participantId,
		Skillset:            input.Skillset,
		DOB:                 dob,
		Motivation:          input.Motivation,
		EmploymentStatus:    input.EmploymentStatus,
		ExperienceLevel:     input.ExperienceLevel,
		HackathonExperience: input.HackathonExperience,
		YearsOfExperience:   input.YearsOfExperience,
		FieldOfStudy:        input.FieldOfStudy,
		PreviousProjects:    input.PreviousProjects,
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
		fmt.Println(err)
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

	particicipantDoc, err := s.ParticipantRepository.CreateParticipantRecord(dataInput)
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
		ParticipationStatus: particicipantDoc.ParticipationStatus,
		AccountRole:         "PARTICIPANT",
		ParticipantType:     particicipantDoc.ParticipantType,
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
		PreviousProjects:    input.PreviousProjects,
		HackathonExperience: input.HackathonExperience,
		FieldOfStudy:        input.FieldOfStudy,
		YearsOfExperience:   input.YearsOfExperience,
		CreatedAt:           acc.CreatedAt,
		UpdatedAt:           acc.UpdatedAt,
	}, err
}

func (s *ParticipantService) RemoveMemberFromTeam(dataInput *repository.RemoveMemberFromTeamData) (*entity.TeamMemberAccount, error) {
	_, err := s.ParticipantRepository.RemoveMemberFromTeam(dataInput)
	if err != nil {
		return nil, err
	}
	acc, err := s.AccountRepository.DB.DeleteAccount(dataInput.MemberEmail)
	info := FillTeamMemberInfo(acc)
	return info, err
}

func (s *ParticipantService) SelectionTeamSolution(dataInput *exports.SelectTeamSolutionData) (*entity.Solution, error) {
	fmt.Println(dataInput)
	solEnt, err := s.ParticipantRepository.SelectSolutionForTeam(dataInput)
	if err != nil {
		return nil, err
	}
	return solEnt, err
}

func (s *ParticipantService) FillParticipantInfo(id string) (*entity.Participant, error) {
	ent, err := s.ParticipantRepository.FillParticipantInfo(id)
	// emit created event

	return ent, err
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

func (s *ParticipantService) GetTeamMembersInfo(p *entity.Participant) ([]entity.TeamMemberAccount, error) {
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

func (s *ParticipantService) GetParticipantInfo(participantId string) (*entity.Participant, error) {
	participant, err := s.ParticipantRepository.GetParticipantInfo(participantId)

	return participant, err
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

func (s *ParticipantService) CreateTeamMemberAccount(dataToSave *exports.CreateTeamMemberAccountData) (*entity.TeamMemberAccount, error) {
	tmAccRepo, err := s.AccountRepository.CreateTeamMemberAccount(dataToSave)
	if err != nil {
		return nil, err
	}
	return tmAccRepo, nil
}

func (s *ParticipantService) GetParticipantsInfo() ([]entity.Participant, error) {
	return s.ParticipantRepository.GetParticipantsInfo()
}

func (s *ParticipantService) InviteToTeam(dataInput *exports.AddToTeamInviteListData) (interface{}, error) {
	return s.ParticipantRepository.InviteToTeam(dataInput)
}
