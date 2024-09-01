package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/cache"
	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/repository"
	"github.com/arravoco/hackathon_backend/utils"
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
	/*isEmailInCache := cache.FindEmailInCache(input.Email)
	if isEmailInCache {
		//return nil, errors.New("email is already existing")
	}*/
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
	_, err = s.ParticipantAccountRepository.CreateParticipantAccount(&exports.CreateParticipantAccountData{
		CreateAccountData: exports.CreateAccountData{
			Email:        input.Email,
			LastName:     input.LastName,
			FirstName:    input.FirstName,
			Gender:       input.Gender,
			Role:         "PARTICIPANT",
			PhoneNumber:  input.PhoneNumber,
			HackathonId:  input.HackathonId,
			PasswordHash: passwordHash,
			State:        input.State,
			Status:       "EMAIL_VERIFIED"},
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

	/* emit created event
	events.EmitParticipantAccountCreated(&exports.ParticipantAccountCreatedEventData{
		ParticipantEmail: acc.Email,
		LastName:         acc.LastName,
		FirstName:        acc.FirstName,
		EventData:        exports.EventData{EventName: "ParticipantAccountCreated"},
		ParticipantType:  "TEAM",
	})*/
	partEnt, err := s.GetParticipantInfo(partDoc.ParticipantId)
	return partEnt, nil
}

func (s *Service) InviteToTeam(dataInput *AddToTeamInviteListData) (interface{}, error) {

	res, err := s.ParticipantRecordRepository.AddToTeamInviteList(&exports.AddToTeamInviteListData{
		Email:            dataInput.Email,
		ParticipantId:    dataInput.ParticipantId,
		InviterEmail:     dataInput.InviterEmail,
		InviterFirstName: dataInput.InviterFirstName,
		HackathonId:      dataInput.HackathonId,
	})
	if err != nil {
		return nil, err
	}
	/*
		queue, err := data.GetQueue("invite_list")
		if err != nil {
			//fmt.Printf("%s\n", err.Error())
		} else {}*/
	queuePayload := exports.InvitelistQueuePayload{
		InviterEmail:       dataInput.InviterEmail,
		InviterName:        dataInput.InviterFirstName,
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
	err = s.Publisher.Publish(exports.PublisherConfig{}, byt)
	if err != nil {
		fmt.Println(err.Error())
	}

	//partEnt, err := s.GetParticipantInfo(partDoc.ParticipantId)
	return res, nil
}

func (s *Service) RegisterTeamLead(input *RegisterNewParticipantDTO) (*entity.Participant, error) {
	err := validate.Struct(input)
	if err != nil {
		return nil, err
	}
	teamMembers := []exports.TeamParticipantInfo{}
	passwordHash, err := exports.GenerateHashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	participantId, err := utils.GenerateParticipantID([]string{input.Email})
	if err != nil {
		return nil, err
	}
	dob, err := time.Parse("2006-04-02", input.DOB)
	if err != nil {
		fmt.Println("\nhere: ", err.Error())
		return nil, err
	}
	_, err = s.ParticipantAccountRepository.CreateParticipantAccount(&exports.CreateParticipantAccountData{
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

	pubData := &exports.ParticipantAccountCreatedEventData{
		TeamParticipants: teamMembers,
		TeamLeadEmail:    input.Email,
		EventData:        exports.EventData{EventName: "ParticipantAccountCreated"},
		ParticipantEmail: input.Email,
		TeamName:         input.TeamName,
		TeamRole:         "TEAM_LEAD",
		ParticipantType:  "TEAM",
	}
	by, err := json.Marshal(pubData)
	if err != nil {
		fmt.Println("failed to marshal to publish")
	}
	//panic("About to pusblish")
	err = s.Publisher.Publish(exports.PublisherConfig{}, by)
	if err != nil {
		fmt.Println("failed to publish")
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

	_, err = s.ParticipantRecordRepository.CreateParticipantRecord(dataInput)
	if err != nil {
		return nil, err
	}

	partEnt, err := s.GetParticipantInfo(participantId)
	return partEnt, err
}

func (s *Service) RemoveMemberFromTeam(dataInput *repository.RemoveMemberFromTeamData) (*entity.TeamMemberWithParticipantRecord, error) {
	_, err := s.ParticipantRecordRepository.RemoveCoparticipantFromParticipantRecord(&exports.RemoveMemberFromTeamData{
		HackathonId:   dataInput.HackathonId,
		ParticipantId: dataInput.HackathonId,
		MemberEmail:   dataInput.MemberEmail,
	})
	if err != nil {
		return nil, err
	}
	acc, err := s.ParticipantAccountRepository.MarkParticipantAccountAsDeleted(dataInput.MemberEmail)
	info := s.FillTeamMemberInfo(acc)
	return info, err
}

func (s *Service) SelectTeamSolution(dataInput *exports.SelectTeamSolutionData) (*entity.Participant, error) {
	fmt.Println(dataInput)
	_, err := s.ParticipantRecordRepository.AddSolutionIdToParticipantRecord(dataInput)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func (s *Service) FillTeamMemberInfo(account *exports.ParticipantAccountRepository) *entity.TeamMemberWithParticipantRecord {
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
	accounts, err := s.ParticipantAccountRepository.GetParticipantAccountsByEmail(emails)
	fmt.Println(accounts)
	for _, acc := range accounts {
		oo := s.FillTeamMemberInfo(&exports.ParticipantAccountRepository{
			Email: acc.Email,
		})
		team = append(team, *oo)
	}
	return team, err
}

func (s *Service) GetParticipantInfo(participantId string) (*entity.Participant, error) {
	part, err := s.ParticipantRecordRepository.GetSingleParticipantRecordAndMemberAccountsInfo(participantId)
	var co_parts []entity.ParticipantEntityCoParticipantInfo
	for _, v := range part.CoParticipants {
		co_parts = append(co_parts, entity.ParticipantEntityCoParticipantInfo{
			AccountStatus: v.AccountStatus,
		})
	}
	partEnt := &entity.Participant{
		ParticipantId:       part.ParticipantId,
		ParticipantEmail:    part.ParticipantEmail,
		ParticipantType:     part.Type,
		ParticipationStatus: part.Status,
		CoParticipants:      co_parts,
		TeamLeadEmail:       part.TeamLeadEmail,
		TeamName:            part.TeamName,
		TeamLeadInfo: entity.ParticipantEntityTeamLeadInfo{
			HackathonId: part.TeamLeadInfo.HackathonId,
			FirstName:   part.TeamLeadInfo.FirstName,
			LastName:    part.TeamLeadInfo.LastName,
			Email:       part.TeamLeadInfo.Email,
			Skillset:    part.TeamLeadInfo.Skillset,
			Gender:      part.TeamLeadInfo.Gender,
			AccountId:   part.TeamLeadInfo.AccountId,
			State:       part.TeamLeadInfo.State,
			TeamRole:    part.TeamLeadInfo.TeamRole,
			PhoneNumber: part.TeamLeadInfo.PhoneNumber,
		},
		AccountStatus: part.Status,
		Solution: &entity.ParticipantEntitySelectedSolution{
			Title:            part.Solution.Title,
			Id:               part.SolutionId,
			Description:      part.Solution.Description,
			SolutionImageUrl: part.Solution.SolutionImageUrl,
			Objective:        part.Solution.Objective,
		},
	}
	return partEnt, err
}

func (s *Service) CreateTeamMemberAccount(dataToSave *exports.CreateTeamMemberAccountData) (*entity.TeamMemberWithParticipantRecord, error) {
	tmAccRepo, err := s.ParticipantAccountRepository.CreateParticipantAccount(&exports.CreateParticipantAccountData{})
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

func (s *Service) UpdateParticipantInfo(input *AuthParticipantInfoUpdateDTO) (interface{}, error) {
	return nil, nil
}
