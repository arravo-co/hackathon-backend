package services

import (
	"encoding/json"
	"fmt"
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

func (s *Service) CompleteNewTeamMemberRegistration(input *CompleteNewTeamMemberRegistrationDTO) (*entity.Participant, error) {
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
	partDoc, err := s.ParticipantRecordRepository.AddMemberInfoToParticipatingTeamRecord(&exports.AddMemberToParticipatingTeamData{
		HackathonId:   input.HackathonId,
		ParticipantId: input.ParticipantId,
		Email:         input.Email,
		Role:          "PARTICIPANT",
		TeamRole:      input.TeamRole,
	})
	if err != nil {
		return nil, err
	}
	acc, err := s.ParticipantAccountRepository.CreateCoParticipantAccount(&exports.RegisterNewParticipantAccountDTO{

		Email:       input.Email,
		LastName:    input.LastName,
		FirstName:   input.FirstName,
		Gender:      input.Gender,
		Role:        "PARTICIPANT",
		PhoneNumber: input.PhoneNumber,
		HackathonId: input.HackathonId,
		State:       input.State,
		Status:      "EMAIL_VERIFIED",
		//TeamRole:         "TEAM_MEMBER",
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

func (s *Service) InviteToTeam(q query.Query, dataInput *exports.AddToTeamInviteListData) (interface{}, error) {
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

func (s *Service) RegisterTeamLead(input dtos.RegisterNewParticipantDTO) (*entity.Participant, error) {
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

	particicipantDoc, err := s.ParticipantRecordRepository.CreateParticipantRecord(dataInput)
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

func (s *Service) RemoveMemberFromTeam(dataInput *repository.RemoveMemberFromTeamData) (*entity.TeamMemberWithParticipantRecord, error) {
	_, err := s.ParticipantRecordRepository.RemoveCoparticipantFromParticipantRecord(dataInput)
	if err != nil {
		return nil, err
	}
	acc, err := s.AccountRepository.DB.DeleteAccount(dataInput.MemberEmail)
	info := FillTeamMemberInfo(acc)
	return info, err
}

func (s *Service) SelectionTeamSolution(dataInput *exports.SelectTeamSolutionData) (*entity.Solution, error) {
	fmt.Println(dataInput)
	solEnt, err := s.ParticipantRecordRepository.SelectSolutionForTeam(dataInput)
	if err != nil {
		return nil, err
	}
	return solEnt, err
}

func (s *Service) FillParticipantInfo(id string) (*entity.Participant, error) {
	ent, err := s.ParticipantRecordRepository.FillParticipantInfo(id)
	// emit created event

	return ent, err
}

func (s *Service) FillTeamMemberInfo(account *exports.AccountDocument) *entity.TeamMemberWithParticipantRecord {
	info := &entity.TeamMemberWithParticipantRecord{}
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

func (s *Service) GetTeamMembersInfo(p *entity.Participant) ([]entity.TeamMemberWithParticipantRecord, error) {
	team := []entity.TeamMemberWithParticipantRecord{}
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

func (s *Service) GetParticipantInfo(participantId string) (*entity.Participant, error) {
	participant, err := s.ParticipantRecordRepository.GetParticipantInfo(participantId)

	return participant, err
}

func (s *Service) CreateTeamMemberAccount(dataToSave *exports.CreateTeamMemberAccountData) (*entity.TeamMemberWithParticipantRecord, error) {
	tmAccRepo, err := s.ParticipantAccountRepository.CreateCoParticipantAccount(&exports.RegisterNewParticipantAccountDTO{})
	if err != nil {
		return nil, err
	}
	fmt.Println(tmAccRepo)
	return nil, nil
}

func (s *Service) GetParticipantsInfo() ([]entity.Participant, error) {
	parts, err := s.ParticipantRecordRepository.GetParticipantsRecords()
	if err != nil {
		return nil, err
	}
	fmt.Println(parts)
	return nil, nil
}
