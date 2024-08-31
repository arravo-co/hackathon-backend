package repository

import (
	"time"

	//"gitee.com/golang-module/carbon"
	"github.com/arravoco/hackathon_backend/cache"
	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/events"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/golang-module/carbon"
	//"github.com/golang-module/carbon"
)

type ParticipantAccountRepository struct {
	DB              *query.Query
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

func NewParticipantAccountRepository(q *query.Query) *ParticipantAccountRepository {
	return &ParticipantAccountRepository{
		DB: q,
	}
}

// GetParticipantAccountByEmail
func (p ParticipantRepository) CreateTeamMemberAccount(dataInput *exports.CreateTeamMemberAccountData) (*exports.ParticipantAccountRepository, error) {
	partDoc, err := p.DB.CreateTeamMemberAccount(dataInput)
	if err != nil {
		return nil, err
	}

	return &exports.ParticipantAccountRepository{
		Email:               partDoc.Email,
		LastName:            partDoc.LastName,
		FirstName:           partDoc.FirstName,
		Gender:              partDoc.Gender,
		PasswordHash:        partDoc.PasswordHash,
		Status:              partDoc.Status,
		State:               partDoc.State,
		LinkedInAddress:     partDoc.LinkedInAddress,
		DOB:                 partDoc.DOB,
		EmploymentStatus:    partDoc.EmploymentStatus,
		IsEmailVerified:     partDoc.IsEmailVerified,
		IsEmailVerifiedAt:   partDoc.IsEmailVerifiedAt,
		ExperienceLevel:     partDoc.ExperienceLevel,
		HackathonExperience: partDoc.HackathonExperience,
		YearsOfExperience:   partDoc.YearsOfExperience,
		Motivation:          partDoc.Motivation,
		CreatedAt:           partDoc.CreatedAt,
		UpdatedAt:           partDoc.UpdatedAt,
	}, nil
}

func (p *ParticipantRepository) RegisterTeamLead(input dtos.RegisterNewParticipantDTO) (*exports.ParticipantRecordRepository, error) {
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
	accRepo := NewAccountRepository(p.DB)
	acc, err := accRepo.CreateParticipantAccount(&exports.CreateParticipantAccountData{
		ParticipantId:       participantId,
		Skillset:            input.Skillset,
		DOB:                 dob,
		Motivation:          input.Motivation,
		EmploymentStatus:    input.EmploymentStatus,
		ExperienceLevel:     input.ExperienceLevel,
		YearsOfExperience:   input.YearsOfExperience,
		HackathonExperience: input.HackathonExperience,
		PreviousProjects:    input.PreviousProjects,
		FieldOfStudy:        input.FieldOfStudy,
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

	particicipantDoc, err := p.DB.CreateParticipantRecord(dataInput)
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

	return &exports.ParticipantRecordRepository{
		HackathonId:      dataInput.HackathonId,
		ParticipantType:  particicipantDoc.Type,
		TeamLeadEmail:    particicipantDoc.TeamLeadEmail,
		TeamName:         particicipantDoc.TeamName,
		ParticipantId:    participantId,
		ParticipantEmail: particicipantDoc.ParticipantEmail,
		TeamRole:         "TEAM_LEAD",
		CreatedAt:        acc.CreatedAt,
		UpdatedAt:        acc.UpdatedAt,
	}, err
}

func (p *ParticipantRepository) RegisterIndividual(input dtos.RegisterNewParticipantDTO) (*exports.ParticipantRecordRepository, error) {
	passwordHash, _ := exports.GenerateHashPassword(input.Password)
	participantId, err := GenerateParticipantID([]string{input.Email})
	if err != nil {
		return nil, err
	}
	dob, err := time.Parse("2006-01-02", input.DOB)
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
			HackathonId:  config.GetHackathonId(),
			Status:       "EMAIL_UNVERIFIED",
		},
		DOB:                 dob,
		Skillset:            input.Skillset,
		ParticipantId:       participantId,
		Motivation:          input.Motivation,
		YearsOfExperience:   input.YearsOfExperience,
		PreviousProjects:    input.PreviousProjects,
		HackathonExperience: input.HackathonExperience,
		FieldOfStudy:        input.FieldOfStudy,
		ExperienceLevel:     input.ExperienceLevel,
		EmploymentStatus:    input.EmploymentStatus,
	}
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
	partDoc, err := p.DB.CreateParticipantRecord(&exports.CreateParticipantRecordData{
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
	//age := time.Now().Sub(time.Now())
	carbon.Now()
	return &exports.ParticipantRecordRepository{
		HackathonId:     dataInput.HackathonId,
		ParticipantType: partDoc.Type,
		Status:          accCreated.Status,
		ParticipantId:   participantId,
		CreatedAt:       accCreated.CreatedAt,
		UpdatedAt:       accCreated.UpdatedAt,
	}, nil
}

/*
func (repo *ParticipantRepository) FillParticipantInfo(idOrEmail string) (*entity.Participant, error) {
	p := &entity.Participant{}
	accountData, err := data.GetAccountByEmail(idOrEmail)
	if err != nil {
		return nil, err
	}
	particicipantDocData, err := data.GetParticipantRecord(accountData.ParticipantId)
	if err != nil {
		return nil, err
	}
	p.Email = accountData.Email
	p.AccountStatus = accountData.Status
	p.ParticipationStatus = particicipantDocData.Status
	p.AccountRole = accountData.Role
	p.FirstName = accountData.FirstName
	p.LastName = accountData.LastName
	p.Gender = accountData.Gender
	p.State = accountData.State
	p.TeamName = particicipantDocData.TeamName
	p.Skillset = accountData.Skillset
	p.TeamLeadEmail = particicipantDocData.TeamLeadEmail
	p.HackathonId = particicipantDocData.HackathonId
	p.ParticipantType = particicipantDocData.Type
	p.ParticipantEmail = particicipantDocData.ParticipantEmail
	p.ParticipantId = particicipantDocData.ParticipantId
	p.ReviewRanking = particicipantDocData.ReviewRanking
	p.Age = time.Now().Year() - accountData.DOB.Year()
	p.Motivation = accountData.Motivation
	p.FieldOfStudy = accountData.FieldOfStudy
	p.HackathonExperience = accountData.HackathonExperience
	p.PreviousProjects = accountData.PreviousProjects
	p.YearsOfExperience = accountData.YearsOfExperience
	p.ExperienceLevel = accountData.ExperienceLevel
	p.EmploymentStatus = accountData.EmploymentStatus
	if particicipantDocData.Type == "TEAM" {
		if particicipantDocData.TeamLeadEmail == accountData.Email {
			p.TeamRole = "TEAM_LEAD"
		} else {
			p.TeamRole = "TEAM_MEMBER"
		}
		partAccs, err := data.GetAccountsByParticipantIds([]string{p.ParticipantId})
		if err != nil {
			fmt.Println(err)
		}
		for _, part := range particicipantDocData.CoParticipants {
			for _, acc := range partAccs {
				if acc.Email == part.Email {
					p.CoParticipants = append(p.CoParticipants, valueobjects.CoParticipantInfo{
						Email:               part.Email,
						TeamRole:            part.Role,
						FirstName:           acc.FirstName,
						LastName:            acc.LastName,
						DOB:                 acc.DOB,
						Age:                 time.Now().Year() - acc.DOB.Year(),
						Gender:              acc.Gender,
						AccountStatus:       acc.Status,
						PhoneNumber:         acc.PhoneNumber,
						AccountRole:         acc.Role,
						State:               acc.Role,
						Skillset:            acc.Skillset,
						HackathonId:         acc.HackathonId,
						Motivation:          acc.Motivation,
						YearsOfExperience:   acc.YearsOfExperience,
						FieldOfStudy:        acc.FieldOfStudy,
						PreviousProjects:    acc.PreviousProjects,
						HackathonExperience: acc.HackathonExperience,
						ExperienceLevel:     acc.ExperienceLevel,
						EmploymentStatus:    acc.EmploymentStatus,
					})

				}
			}
		}
	}
	solRepo := NewSolutionRepository(repo.DB)
	sol, err := solRepo.GetSolutionDataById(particicipantDocData.SolutionId)
	if err != nil {
		fmt.Println(err)
		return p, nil
	}
	p.Solution = &entity.Solution{
		Id:          sol.Id,
		HackathonId: sol.HackathonId,
		CreatorId:   sol.CreatorId,
		Title:       sol.Title,
		Description: sol.Description,
		CreatedAt:   sol.CreatedAt,
		UpdatedAt:   sol.UpdatedAt,
	}
	// emit created event

	return p, nil
}*/
