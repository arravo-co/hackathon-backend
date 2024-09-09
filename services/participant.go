package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/repository"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/rabbitmq/amqp091-go"
)

func (s *Service) GetTeamParticipantInfo(email string) (*entity.TeamMemberWithParticipantRecord, error) {

	fmt.Printf("\n\n \n\n %+v \n\n\n", email)
	partAcc, err := s.ParticipantAccountRepository.GetParticipantAccountByEmail(email)
	if err != nil {
		return nil, err
	}
	if partAcc == nil {
		return nil, fmt.Errorf("unable to get participant info")
	}
	//fmt.PrintF("\n\n \n\n%es %+v \n\n\n"),email, partAcc))
	if s.ParticipantRecordRepository == nil {
		return nil, fmt.Errorf("no participant record repository")
	}
	fmt.Printf("\n\n  %+v \n\n\n", partAcc)
	part, err := s.ParticipantRecordRepository.GetSingleParticipantRecordAndMemberAccountsInfo(partAcc.ParticipantId)
	if err != nil {
		return nil, err
	}
	if part == nil {
		return nil, fmt.Errorf("participant with id %v not found", partAcc.ParticipantId)
	}
	var co_parts []entity.ParticipantEntityCoParticipantInfo
	if part.CoParticipants != nil {

		for _, v := range part.CoParticipants {
			co_parts = append(co_parts, entity.ParticipantEntityCoParticipantInfo{
				AccountStatus:       v.AccountStatus,
				Email:               v.Email,
				FirstName:           v.FirstName,
				LastName:            v.LastName,
				Gender:              v.Gender,
				PhoneNumber:         v.PhoneNumber,
				State:               v.State,
				EmploymentStatus:    v.EmploymentStatus,
				ExperienceLevel:     v.ExperienceLevel,
				HackathonExperience: v.HackathonExperience,
				YearsOfExperience:   v.YearsOfExperience,
				FieldOfStudy:        v.FieldOfStudy,
				Skillset:            v.Skillset,
				Motivation:          v.Motivation,
				TeamRole:            v.TeamRole,
				AccountRole:         v.AccountRole,
				AccountId:           v.AccountId,
				HackathonId:         v.HackathonId,
				DOB:                 v.DOB,
				PreviousProjects:    v.PreviousProjects,
				ParticipantId:       v.ParticipantId,
				CreatedAt:           v.CreatedAt,
				UpdatedAt:           v.UpdateAt,
			})
		}
	}
	team_lead := entity.ParticipantEntityTeamLeadInfo{
		HackathonId:         part.TeamLeadInfo.HackathonId,
		FirstName:           part.TeamLeadInfo.FirstName,
		LastName:            part.TeamLeadInfo.LastName,
		Email:               part.TeamLeadInfo.Email,
		Skillset:            part.TeamLeadInfo.Skillset,
		Gender:              part.TeamLeadInfo.Gender,
		AccountId:           part.TeamLeadInfo.AccountId,
		State:               part.TeamLeadInfo.State,
		TeamRole:            part.TeamLeadInfo.TeamRole,
		PhoneNumber:         part.TeamLeadInfo.PhoneNumber,
		AccountStatus:       part.TeamLeadInfo.AccountStatus,
		AccountRole:         part.TeamLeadInfo.AccountRole,
		EmploymentStatus:    part.TeamLeadInfo.EmploymentStatus,
		ExperienceLevel:     part.TeamLeadInfo.ExperienceLevel,
		HackathonExperience: part.TeamLeadInfo.HackathonExperience,
		YearsOfExperience:   part.TeamLeadInfo.YearsOfExperience,
		IsEmailVerified:     part.TeamLeadInfo.IsEmailVerified,
		IsEmailVerifiedAt:   part.TeamLeadInfo.IsEmailVerifiedAt,
		LinkedInAddress:     part.TeamLeadInfo.LinkedInAddress,
		ParticipantId:       part.TeamLeadInfo.ParticipantId,
		PreviousProjects:    part.TeamLeadInfo.PreviousProjects,
		FieldOfStudy:        part.TeamLeadInfo.FieldOfStudy,
		UpdateAt:            part.TeamLeadInfo.UpdateAt,
	}
	sol := &entity.ParticipantEntitySelectedSolution{
		Title:            part.Solution.Title,
		HackathonId:      part.Solution.HackathonId,
		Id:               part.SolutionId,
		Description:      part.Solution.Description,
		SolutionImageUrl: part.Solution.SolutionImageUrl,
		Objective:        part.Solution.Objective,
	}
	_ = &entity.Participant{
		ParticipantId:       part.ParticipantId,
		ParticipantEmail:    part.ParticipantEmail,
		ParticipantType:     part.Type,
		ParticipatantStatus: part.Status,
		CoParticipants:      co_parts,
		TeamLeadEmail:       part.TeamLeadEmail,
		TeamName:            part.TeamName,
		TeamLeadInfo:        team_lead,
		Solution:            sol,
	}
	var team_role string = "TEAM_MEMBER"
	if team_lead.Email == partAcc.Email {
		team_role = team_lead.TeamRole
	}
	return &entity.TeamMemberWithParticipantRecord{
		Email:               partAcc.Email,
		FirstName:           partAcc.FirstName,
		LastName:            partAcc.LastName,
		AccountId:           partAcc.Id,
		PhoneNumber:         partAcc.PhoneNumber,
		HackathonExperience: partAcc.HackathonExperience,
		Skillset:            partAcc.Skillset,
		PreviousProjects:    partAcc.PreviousProjects,
		ParticipantStatus:   part.Status,
		YearsOfExperience:   partAcc.YearsOfExperience,
		ReviewRanking:       part.ReviewRanking,
		IsEmailVerified:     partAcc.IsEmailVerified,
		IsEmailVerifiedAt:   partAcc.IsEmailVerifiedAt,
		DOB:                 partAcc.DOB,
		ParticipantId:       part.ParticipantId,
		CoParticipants:      co_parts,
		TeamLeadInfo:        team_lead,
		EmploymentStatus:    partAcc.EmploymentStatus,
		ExperienceLevel:     partAcc.ExperienceLevel,
		FieldOfStudy:        partAcc.FieldOfStudy,
		LinkedInAddress:     partAcc.LinkedInAddress,
		HackathonId:         partAcc.HackathonId,
		State:               partAcc.State,
		AccountStatus:       partAcc.Status,
		TeamRole:            team_role,
	}, err
}

func (s *Service) CompleteNewTeamMemberRegistration(input *CompleteNewTeamMemberRegistrationDTO) (*entity.Participant, error) {
	err := validate.Struct(input)
	if err != nil {
		return nil, err
	}
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
	//fmt.Println("uyfrrdrrszwswwseaesrsrrerrtdtyfyguguugufydtrsesezzesrdtdtfyfgugugugugyfiyfyfyyyyfyfyffu")
	partAcc, err := s.ParticipantAccountRepository.CreateParticipantAccount(&exports.CreateParticipantAccountData{
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
	partEnt, err := s.GetSingleParticipantWithAccountsInfo(partDoc.ParticipantId)
	if err != nil {
		return nil, err
	}
	/*
		addedToCache := cache.AddEmailToCache(input.Email)

		if !addedToCache {
			exports.MySugarLogger.Warnln("Email is already in cache")
		}
	*/
	/* emit created event
	events.EmitParticipantAccountCreated(&exports.ParticipantAccountCreatedEventData{
		ParticipantEmail: acc.Email,
		LastName:         acc.LastName,
		FirstName:        acc.FirstName,
		EventData:        exports.EventData{EventName: "ParticipantAccountCreated"},
		ParticipantType:  "TEAM",
	})*/

	pubData := &exports.ParticipantRegisteredPublishPayload{
		AccountId:        partAcc.Id,
		Email:            input.Email,
		LastName:         input.LastName,
		FirstName:        input.FirstName,
		ParticipantId:    input.ParticipantId,
		ParticipantEmail: partDoc.ParticipantEmail,
		ParticipantType:  partDoc.Type,
		TeamLeadEmail:    partDoc.TeamLeadEmail,
		TeamRole:         input.TeamRole,
		TeamName:         partDoc.TeamName,
		TeamLeadName:     partEnt.TeamLeadInfo.FirstName,
	}

	by, err := json.Marshal(pubData)
	if err != nil {
		fmt.Printf("Failed to marshal: %v\n", err)
	} else {
		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.ParticipantsExchange,
			KeyName:      exports.ParticipantRegisteredRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.ParticipantsExchange,
			KeyName:      exports.ParticipantTeamMemberSendWelcomeEmailRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	return partEnt, nil
}

func (s *Service) InviteToTeam(dataInput *AddToTeamInviteListData) (*entity.Participant, error) {
	err := validate.Struct(dataInput)
	if err != nil {
		return nil, err
	}
	part, err := s.ParticipantRecordRepository.GetParticipantRecord(dataInput.ParticipantId)
	if err != nil {
		return nil, fmt.Errorf("failed to get participant record")
	}
	if part.TeamLeadEmail != dataInput.InviterEmail {
		return nil, fmt.Errorf("only team leads can invite new members")
	}
	if part.Status == "UNREVIEWED" {
		return nil, fmt.Errorf("an unreviewed participant cannot have invite new members")
	}
	_, err = s.ParticipantRecordRepository.AddToTeamInviteList(&exports.AddToTeamInviteListData{
		Email:            dataInput.Email,
		ParticipantId:    dataInput.ParticipantId,
		InviterEmail:     dataInput.InviterEmail,
		InviterFirstName: dataInput.InviterFirstName,
		HackathonId:      dataInput.HackathonId,
	})
	if err != nil {
		return nil, err
	}
	queuePayload := exports.AddedToInvitelistPublishPayload{
		InviterEmail:       dataInput.InviterEmail,
		InviterName:        dataInput.InviterFirstName,
		InviteeEmail:       dataInput.Email,
		ParticipantId:      dataInput.ParticipantId,
		HackathonId:        dataInput.HackathonId,
		TeamLeadEmailEmail: dataInput.InviterEmail,
		TimeSent:           time.Now(),
	}
	by, err := json.Marshal(queuePayload)
	if err != nil {
		fmt.Printf("Failed to marshal: %v\n", err)
	} else {
		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.InvitationsExchange,
			KeyName:      exports.ParticipantInvitedRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.InvitationsExchange,
			KeyName:      exports.ParticipantSendInvitationEmailRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

	}

	partEnt, err := s.GetSingleParticipantWithAccountsInfo(part.ParticipantId)
	return partEnt, err
}

func (s *Service) RegisterTeamLead(input *RegisterNewParticipantDTO) (*entity.Participant, error) {
	err := validate.Struct(input)
	if err != nil {
		return nil, err
	}
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
	partAcc, err := s.ParticipantAccountRepository.CreateParticipantAccount(&exports.CreateParticipantAccountData{
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

	partDoc, err := s.ParticipantRecordRepository.CreateParticipantRecord(dataInput)
	if err != nil {
		return nil, err
	}

	pubData := &exports.ParticipantRegisteredPublishPayload{
		AccountId:        partAcc.Id,
		Email:            input.Email,
		LastName:         input.LastName,
		FirstName:        input.FirstName,
		ParticipantId:    participantId,
		ParticipantEmail: partDoc.ParticipantEmail,
		ParticipantType:  partDoc.Type,
		TeamLeadEmail:    partDoc.TeamLeadEmail,
		TeamRole:         "TEAM_LEAD",
		TeamName:         partDoc.TeamName,
	}

	by, err := json.Marshal(pubData)
	if err != nil {
		fmt.Printf("Failed to marshal: %v\n", err)
	} else {
		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.ParticipantsExchange,
			KeyName:      exports.ParticipantRegisteredRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.ParticipantsExchange,
			KeyName:      exports.ParticipantTeamLeadSendWelcomeEmailRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

	}

	partEnt, err := s.GetSingleParticipantWithAccountsInfo(participantId)
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

func (s *Service) SelectTeamSolution(dataInput *SelectTeamSolutionData) error {
	fmt.Println(dataInput)
	_, err := s.ParticipantRecordRepository.AddSolutionIdToParticipantRecord(&exports.SelectTeamSolutionData{
		HackathonId:   dataInput.HackathonId,
		ParticipantId: dataInput.ParticipantId,
		SolutionId:    dataInput.SolutionId,
	})
	if err != nil {
		return err
	}
	//partEnt, err := s.GetParticipantInfo(dataInput.ParticipantId)
	return nil
}

func (s *Service) FillTeamMemberInfo(account *exports.ParticipantAccountRepository) *entity.TeamMemberWithParticipantRecord {
	info := &entity.TeamMemberWithParticipantRecord{}
	info.Email = account.Email
	info.AccountStatus = account.Status
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

func (s *Service) GetTeamMembersInfo(participant_id string) ([]entity.TeamMemberWithParticipantRecord, error) {

	team := []entity.TeamMemberWithParticipantRecord{}
	partAggre, err := s.ParticipantRecordRepository.GetSingleParticipantRecordAndMemberAccountsInfo(participant_id)

	for _, acc := range partAggre.CoParticipants {
		team = append(team, entity.TeamMemberWithParticipantRecord{
			Email:               acc.Email,
			LastName:            acc.Email,
			FirstName:           acc.FirstName,
			FieldOfStudy:        acc.FieldOfStudy,
			AccountRole:         acc.AccountRole,
			TeamRole:            acc.TeamRole,
			DOB:                 acc.DOB,
			AccountId:           acc.AccountId,
			IsEmailVerifiedAt:   acc.IsEmailVerifiedAt,
			PhoneNumber:         acc.PhoneNumber,
			ParticipantId:       acc.ParticipantId,
			LinkedInAddress:     acc.LinkedInAddress,
			HackathonExperience: acc.HackathonExperience,
			YearsOfExperience:   acc.YearsOfExperience,
			IsEmailVerified:     acc.IsEmailVerified,
			HackathonId:         acc.HackathonId,
			CreatedAt:           acc.CreatedAt,
			UpdatedAt:           acc.UpdateAt,
		})
	}
	return team, err
}

func (s *Service) GetSingleParticipantWithAccountsInfo(participantId string) (*entity.Participant, error) {
	part, err := s.ParticipantRecordRepository.GetSingleParticipantRecordAndMemberAccountsInfo(participantId)
	if err != nil {
		return nil, err
	}
	if part == nil {
		return nil, fmt.Errorf("participant with id %v not found", participantId)
	}
	var co_parts []entity.ParticipantEntityCoParticipantInfo
	if part.CoParticipants != nil {

		for _, v := range part.CoParticipants {
			co_parts = append(co_parts, entity.ParticipantEntityCoParticipantInfo{
				AccountStatus:       v.AccountStatus,
				Email:               v.Email,
				FirstName:           v.FirstName,
				LastName:            v.LastName,
				Gender:              v.Gender,
				PhoneNumber:         v.PhoneNumber,
				State:               v.State,
				EmploymentStatus:    v.EmploymentStatus,
				ExperienceLevel:     v.ExperienceLevel,
				HackathonExperience: v.HackathonExperience,
				YearsOfExperience:   v.YearsOfExperience,
				FieldOfStudy:        v.FieldOfStudy,
				Skillset:            v.Skillset,
				Motivation:          v.Motivation,
				TeamRole:            v.TeamRole,
				AccountRole:         v.AccountRole,
				AccountId:           v.AccountId,
				HackathonId:         v.HackathonId,
				DOB:                 v.DOB,
				PreviousProjects:    v.PreviousProjects,
				ParticipantId:       v.ParticipantId,
				CreatedAt:           v.CreatedAt,
				UpdatedAt:           v.UpdateAt,
			})
		}
	}
	partEnt := &entity.Participant{
		ParticipantId:       part.ParticipantId,
		ParticipantEmail:    part.ParticipantEmail,
		ParticipantType:     part.Type,
		ParticipatantStatus: part.Status,
		CoParticipants:      co_parts,
		TeamLeadEmail:       part.TeamLeadEmail,
		TeamName:            part.TeamName,
		TeamLeadInfo: entity.ParticipantEntityTeamLeadInfo{
			HackathonId:         part.TeamLeadInfo.HackathonId,
			FirstName:           part.TeamLeadInfo.FirstName,
			LastName:            part.TeamLeadInfo.LastName,
			Email:               part.TeamLeadInfo.Email,
			Skillset:            part.TeamLeadInfo.Skillset,
			Gender:              part.TeamLeadInfo.Gender,
			AccountId:           part.TeamLeadInfo.AccountId,
			State:               part.TeamLeadInfo.State,
			TeamRole:            part.TeamLeadInfo.TeamRole,
			PhoneNumber:         part.TeamLeadInfo.PhoneNumber,
			AccountStatus:       part.TeamLeadInfo.AccountStatus,
			AccountRole:         part.TeamLeadInfo.AccountRole,
			EmploymentStatus:    part.TeamLeadInfo.EmploymentStatus,
			ExperienceLevel:     part.TeamLeadInfo.ExperienceLevel,
			HackathonExperience: part.TeamLeadInfo.HackathonExperience,
			YearsOfExperience:   part.TeamLeadInfo.YearsOfExperience,
			IsEmailVerified:     part.TeamLeadInfo.IsEmailVerified,
			IsEmailVerifiedAt:   part.TeamLeadInfo.IsEmailVerifiedAt,
			LinkedInAddress:     part.TeamLeadInfo.LinkedInAddress,
			ParticipantId:       part.TeamLeadInfo.ParticipantId,
			PreviousProjects:    part.TeamLeadInfo.PreviousProjects,
			FieldOfStudy:        part.TeamLeadInfo.FieldOfStudy,
			UpdateAt:            part.TeamLeadInfo.UpdateAt,
		},
		Solution: &entity.ParticipantEntitySelectedSolution{
			Title:            part.Solution.Title,
			HackathonId:      part.Solution.HackathonId,
			Id:               part.SolutionId,
			Description:      part.Solution.Description,
			SolutionImageUrl: part.Solution.SolutionImageUrl,
			Objective:        part.Solution.Objective,
		},
	}
	return partEnt, err
}

func (s *Service) GetMultipleParticipantsWithAccounts(opts *GetParticipantsWithAccountsAggregateFilterOpts) ([]entity.Participant, error) {
	fmt.Printf("%#v\n", opts)
	parts, err := s.ParticipantRecordRepository.
		GetMultipleParticipantRecordAndMemberAccountsInfo(
			exports.GetParticipantsWithAccountsAggregateFilterOpts(*opts))
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v\n", parts)
	//fmt.Printf("%#v\n", parts)
	var partEnts []entity.Participant
	for _, v := range parts {
		var co_parts []entity.ParticipantEntityCoParticipantInfo
		for _, v := range v.CoParticipants {
			co_parts = append(co_parts, entity.ParticipantEntityCoParticipantInfo{
				ParticipantId:       v.ParticipantId,
				PhoneNumber:         v.PhoneNumber,
				PreviousProjects:    v.PreviousProjects,
				HackathonId:         v.HackathonId,
				HackathonExperience: v.HackathonExperience,
				YearsOfExperience:   v.YearsOfExperience,
				DOB:                 v.DOB,
				Email:               v.Email,
				EmploymentStatus:    v.EmploymentStatus,
				ExperienceLevel:     v.ExperienceLevel,
				LastName:            v.LastName,
				FirstName:           v.FirstName,
				FieldOfStudy:        v.FieldOfStudy,
				State:               v.State,
				Skillset:            v.Skillset,
				AccountStatus:       v.AccountStatus,
				AccountId:           v.AccountId,
				AccountRole:         v.AccountRole,
				Gender:              v.Gender,
				CreatedAt:           v.CreatedAt,
				UpdatedAt:           v.UpdateAt,
				TeamRole:            v.TeamRole,
				Motivation:          v.Motivation,
			})
		}
		var invite_list []entity.InviteInfo
		for _, v := range v.InviteList {
			invite_list = append(invite_list, entity.InviteInfo{
				Email:     v.Email,
				InviterId: v.InviterId,
				Time:      v.Time,
			})
		}
		var sol *entity.ParticipantEntitySelectedSolution
		if v.SolutionId != "" {
			sol = &entity.ParticipantEntitySelectedSolution{
				Id:               v.Solution.Id,
				HackathonId:      v.Solution.HackathonId,
				Title:            v.Solution.Title,
				Description:      v.Solution.Description,
				SolutionImageUrl: v.Solution.SolutionImageUrl,
				Objective:        v.Solution.Objective,
			}
		}
		partEnts = append(partEnts, entity.Participant{
			ParticipantId: v.ParticipantId,
			TeamLeadEmail: v.TeamLeadEmail,
			TeamName:      v.TeamName,
			HackathonId:   v.HackathonId,
			InviteList:    invite_list,
			TeamLeadInfo: entity.ParticipantEntityTeamLeadInfo{
				HackathonId:         v.HackathonId,
				TeamRole:            "TEAM_LEAD",
				FirstName:           v.TeamLeadInfo.FirstName,
				LastName:            v.TeamLeadInfo.LastName,
				Gender:              v.TeamLeadInfo.Gender,
				Skillset:            v.TeamLeadInfo.Skillset,
				PhoneNumber:         v.TeamLeadInfo.PhoneNumber,
				AccountId:           v.TeamLeadInfo.AccountId,
				AccountStatus:       v.TeamLeadInfo.AccountStatus,
				AccountRole:         v.TeamLeadInfo.AccountRole,
				YearsOfExperience:   v.TeamLeadInfo.YearsOfExperience,
				HackathonExperience: v.TeamLeadInfo.HackathonExperience,
				EmploymentStatus:    v.TeamLeadInfo.EmploymentStatus,
				Email:               v.TeamLeadInfo.Email,
				ExperienceLevel:     v.TeamLeadInfo.ExperienceLevel,
				IsEmailVerified:     v.TeamLeadInfo.IsEmailVerified,
				IsEmailVerifiedAt:   v.TeamLeadInfo.IsEmailVerifiedAt,
				CreatedAt:           v.TeamLeadInfo.CreatedAt,
				ParticipantId:       v.TeamLeadInfo.ParticipantId,
				UpdateAt:            v.TeamLeadInfo.UpdateAt,
				Motivation:          v.TeamLeadInfo.Motivation,
				PreviousProjects:    v.TeamLeadInfo.PreviousProjects,
				FieldOfStudy:        v.TeamLeadInfo.FieldOfStudy,
			},
			Solution:            sol,
			CoParticipants:      co_parts,
			ParticipantType:     v.Type,
			ParticipantEmail:    v.ParticipantEmail,
			ParticipatantStatus: v.Status,
			CreatedAt:           v.CreatedAt,
			UpdatedAt:           v.UpdatedAt,
		})
	}
	//fmt.Printf("%#v\n", partEnts)
	fmt.Println("---------------------------------------------------------------------------")
	return partEnts, nil
}

func (s *Service) UpdateParticipantInfo(input *AuthParticipantInfoUpdateDTO) (interface{}, error) {
	return nil, nil
}

func (s *Service) AdminUpdateParticipantInfo(filter *UpdateSingleParticipantRecordFilter, input *AdminParticipantInfoUpdateDTO) (interface{}, error) {
	s.ParticipantRecordRepository.AdminUpdateParticipantRecord(&exports.UpdateSingleParticipantRecordFilter{
		HackathonId:   filter.HackathonId,
		ParticipantId: filter.ParticipantId,
	}, &exports.AdminParticipantInfoUpdateDTO{
		ReviewRanking: input.ReviewRanking,
		Status:        input.Status,
	})
	return nil, nil
}
