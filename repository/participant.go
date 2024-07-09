package repository

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
	valueobjects "github.com/arravoco/hackathon_backend/value_objects"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddMemberToParticipatingTeam
type ParticipantRepository struct {
	DB                  *query.Query
	Entity              *entity.Participant
	FirstName           string               `json:"first_name"`
	LastName            string               `json:"last_name"`
	Email               string               `json:"email"`
	Gender              string               `json:"gender"`
	State               string               `json:"state"`
	Age                 int                  `json:"age"`
	DOB                 time.Time            `json:"dob"`
	AccountRole         string               `json:"role"`
	ParticipantId       string               `json:"participant_id"`
	TeamLeadEmail       string               `json:"team_lead_email"`
	TeamName            string               `json:"team_name"`
	TeamRole            string               `json:"team_role"`
	HackathonId         string               `json:"hackathon_id"`
	ParticipantType     string               `json:"type"`
	CoParticipants      []CoParticipantInfo  `json:"co_participants"`
	ParticipantEmail    string               `json:"participant_email"`
	InviteList          []exports.InviteInfo `json:"invite_list"`
	AccountStatus       string               `json:"account_status"`
	ParticipationStatus string               `json:"participation_status"`
	Skillset            []string             `json:"skillset"`
	PhoneNumber         string               `json:"phone_number"`
	EmploymentStatus    string               `json:"employment_status"`
	ExperienceLevel     string               `json:"experience_level"`
	Motivation          string               `json:"motivation"`
	HackathonExperience string               `json:"hackathon_experience"`
	YearsOfExperience   int                  `json:"years_of_experience"`
	FieldOfStudy        string               `json:"field_of_study"`
	PreviousProjects    []string             `json:"previous_projects"`
	Solution            *SolutionRepository  `json:"solution"`
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
}
type ParticipantRecord struct {
	ParticipantId       string               `json:"participant_id"`
	TeamLeadEmail       string               `json:"team_lead_email"`
	TeamName            string               `json:"team_name"`
	TeamRole            string               `json:"team_role"`
	HackathonId         string               `json:"hackathon_id"`
	ParticipantType     string               `json:"type"`
	CoParticipants      []CoParticipantInfo  `json:"co_participants"`
	ParticipantEmail    string               `json:"participant_email"`
	InviteList          []exports.InviteInfo `json:"invite_list"`
	ParticipationStatus string               `json:"participation_status"`
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
}
type CoParticipantInfo struct {
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Email               string    `json:"email"`
	Gender              string    `json:"gender"`
	State               string    `json:"state"`
	Age                 int       `json:"age"`
	DOB                 time.Time `json:"dob"`
	AccountStatus       string    `json:"account_status"`
	AccountRole         string    `json:"account_role"`
	TeamRole            string    `json:"team_role"`
	ParticipantId       string    `json:"participant_id"`
	HackathonId         string    `json:"hackathon_id"`
	Skillset            []string  `json:"skillset"`
	PhoneNumber         string    `json:"phone_number"`
	EmploymentStatus    string    `json:"employment_status"`
	ExperienceLevel     string    `json:"experience_level"`
	Motivation          string    `json:"motivation"`
	HackathonExperience string    `json:"hackathon_experience"`
	YearsOfExperience   int       `json:"years_of_experience"`
	FieldOfStudy        string    `json:"field_of_study"`
	PreviousProjects    []string  `json:"previous_projects"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
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

type CoParticipantCreatedData struct {
	Email         string
	Password      string
	ParticipantId string
}

type RemoveMemberFromTeamData struct {
	HackathonId   string `bson:"hackathon_id"`
	ParticipantId string `bson:"participant_id"`
	MemberEmail   string `bson:"email"`
}

func NewParticipantRepository(q *query.Query) *ParticipantRepository {
	return &ParticipantRepository{
		DB: q,
	}
}

func (repo *ParticipantRepository) SetEntity(ent *entity.Participant) {
	repo.Entity = ent
}

func (p *ParticipantRepository) AddMemberToParticipatingTeam(dataToSave *exports.AddMemberToParticipatingTeamData) (*exports.ParticipantDocument, error) {
	partDoc, err := p.DB.AddMemberToParticipatingTeam(dataToSave)
	if err != nil {
		return nil, err
	}

	return partDoc, err
}

func (p ParticipantRepository) CreateParticipantRecord(dataInput *exports.CreateParticipantRecordData) (*ParticipantRecord, error) {
	partDoc, err := p.DB.CreateParticipantRecord(dataInput)
	if err != nil {
		return nil, err
	}

	return &ParticipantRecord{
		TeamLeadEmail:       partDoc.TeamLeadEmail,
		ParticipantId:       partDoc.ParticipantId,
		ParticipationStatus: partDoc.Status,
		HackathonId:         partDoc.HackathonId,
		ParticipantType:     partDoc.Type,
		TeamName:            partDoc.TeamName,
		CreatedAt:           partDoc.CreatedAt,
		UpdatedAt:           partDoc.UpdatedAt,
	}, nil
}

func (p ParticipantRepository) InviteToTeam(dataInput *exports.AddToTeamInviteListData) (interface{}, error) {
	res, err := p.DB.AddToTeamInviteList(dataInput)
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

func (p *ParticipantRepository) RegisterIndividual(input dtos.RegisterNewParticipantDTO) (*ParticipantRepository, error) {
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
	return &ParticipantRepository{
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
		Motivation:          dataInput.Motivation,
		YearsOfExperience:   dataInput.YearsOfExperience,
		HackathonExperience: dataInput.HackathonExperience,
		FieldOfStudy:        dataInput.FieldOfStudy,
		PreviousProjects:    dataInput.PreviousProjects,
		ExperienceLevel:     dataInput.ExperienceLevel,
		EmploymentStatus:    dataInput.EmploymentStatus,
		CreatedAt:           accCreated.CreatedAt,
		UpdatedAt:           accCreated.UpdatedAt,
	}, nil
}

func (p *ParticipantRepository) RegisterTeamLead(input dtos.RegisterNewParticipantDTO) (*ParticipantRepository, error) {
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

	return &ParticipantRepository{
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
		CoParticipants:      []CoParticipantInfo{},
		Motivation:          input.Motivation,
		YearsOfExperience:   input.YearsOfExperience,
		FieldOfStudy:        input.FieldOfStudy,
		HackathonExperience: input.HackathonExperience,
		PreviousProjects:    input.PreviousProjects,
		ExperienceLevel:     input.ExperienceLevel,
		EmploymentStatus:    input.EmploymentStatus,
		CreatedAt:           acc.CreatedAt,
		UpdatedAt:           acc.UpdatedAt,
	}, err
}

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
}

func (p *ParticipantRepository) UpdateParticipantInfo(dataInput *dtos.AuthParticipantInfoUpdateDTO) error {
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

func (s *ParticipantRepository) GetParticipantInfo(participantId string) (*entity.Participant, error) {
	participant, err := s.DB.GetParticipantRecord(participantId)
	if err != nil {
		return nil, err
	}
	partIds := []string{participant.ParticipantId}

	accs, err := data.GetAccountsByParticipantIds(partIds)
	if err != nil {
		return nil, err
	}
	pE := entity.Participant{}
	cs := []CoParticipantInfo{}

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
			pE.Motivation = a.Motivation
			pE.YearsOfExperience = a.YearsOfExperience
			pE.HackathonExperience = a.HackathonExperience
			pE.PreviousProjects = a.PreviousProjects
			pE.CreatedAt = a.CreatedAt
			pE.UpdatedAt = a.UpdatedAt
			if participant.Type == "TEAM" {
				pE.TeamRole = "TEAM_LEAD"
			}
		} else {
			for _, c := range participant.CoParticipants {
				if c.Email == a.Email {
					cs = append(cs, CoParticipantInfo{
						FirstName:           a.FirstName,
						Email:               c.Email,
						LastName:            a.LastName,
						PhoneNumber:         a.PhoneNumber,
						Gender:              a.Gender,
						State:               a.State,
						DOB:                 a.DOB,
						Age:                 time.Now().Year() - a.DOB.Year(),
						AccountStatus:       a.Status,
						ParticipantId:       a.ParticipantId,
						HackathonId:         a.HackathonId,
						TeamRole:            c.Role,
						AccountRole:         a.Role,
						Skillset:            a.Skillset,
						Motivation:          a.Motivation,
						YearsOfExperience:   a.YearsOfExperience,
						FieldOfStudy:        a.FieldOfStudy,
						HackathonExperience: a.HackathonExperience,
						PreviousProjects:    a.PreviousProjects,
						ExperienceLevel:     a.ExperienceLevel,
						EmploymentStatus:    a.EmploymentStatus,
						CreatedAt:           a.CreatedAt,
						UpdatedAt:           a.UpdatedAt,
					})
				}
			}
		}
	}

	return &pE, err
}

func (s *ParticipantRepository) GetParticipantsInfo() ([]entity.Participant, error) {
	partDocs, err := s.DB.GetParticipantsRecordsAggregate()

	var participants []entity.Participant
	for _, part := range partDocs {
		participants = append(participants, entity.Participant{
			Email:               part.Email,
			FirstName:           part.FirstName,
			LastName:            part.LastName,
			Gender:              part.Gender,
			DOB:                 part.DOB,
			TeamLeadEmail:       part.TeamName,
			AccountRole:         part.Role,
			ParticipantId:       part.ParticipantId,
			ParticipantEmail:    part.ParticipantEmail,
			TeamName:            part.TeamName,
			TeamRole:            "TEAM_LEAD",
			PhoneNumber:         part.PhoneNumber,
			ParticipantType:     part.Type,
			Motivation:          part.Motivation,
			YearsOfExperience:   part.YearsOfExperience,
			HackathonExperience: part.HackathonExperience,
			PreviousProjects:    part.PreviousProjects,
			FieldOfStudy:        part.FieldOfStudy,
			CoParticipants: func(arr []exports.CoParticipantAggregateDocument) []valueobjects.CoParticipantInfo {
				var items []valueobjects.CoParticipantInfo
				for _, ar := range arr {
					items = append(items, valueobjects.CoParticipantInfo{
						HackathonId:         ar.HackathonId,
						Email:               ar.Email,
						FirstName:           ar.FirstName,
						LastName:            ar.LastName,
						Gender:              ar.Gender,
						State:               ar.State,
						AccountRole:         ar.AccountRole,
						TeamRole:            ar.TeamRole,
						Motivation:          ar.Motivation,
						DOB:                 ar.DOB,
						Skillset:            ar.Skillset,
						YearsOfExperience:   ar.YearsOfExperience,
						FieldOfStudy:        ar.FieldOfStudy,
						HackathonExperience: ar.HackathonExperience,
						PreviousProjects:    ar.PreviousProjects,
						EmploymentStatus:    ar.EmploymentStatus,
						ExperienceLevel:     ar.ExperienceLevel,
						PhoneNumber:         ar.PhoneNumber,
						CreatedAt:           ar.CreatedAt,
						UpdatedAt:           ar.UpdatedAt,
					})
				}
				return items
			}(part.CoParticipants),
		})
	}
	return participants, err
}

func (s *ParticipantRepository) RemoveMemberFromTeam(dataInput *RemoveMemberFromTeamData) (*entity.TeamMemberAccount, error) {
	_, err := s.DB.RemoveMemberFromParticipatingTeam(&exports.RemoveMemberFromParticipatingTeamData{
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

//

func (s *ParticipantRepository) SelectSolutionForTeam(dataInput *exports.SelectTeamSolutionData) (*entity.Solution, error) {

	partDoc, err := s.DB.SelectSolutionForTeam(&exports.SelectTeamSolutionData{
		HackathonId:   dataInput.HackathonId,
		ParticipantId: dataInput.ParticipantId,
		SolutionId:    dataInput.SolutionId,
	})
	if err != nil {
		return nil, err
	}
	return &entity.Solution{
		Id:          partDoc.Solution.Id.(primitive.ObjectID).Hex(),
		HackathonId: partDoc.Solution.HackathonId,
		Description: partDoc.Solution.Description,
		Objective:   partDoc.Solution.Objective,
		Title:       partDoc.Solution.Title,
		CreatorId:   partDoc.Solution.CreatorId,
	}, err
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
	info.Motivation = account.Motivation
	info.HackathonExperience = account.HackathonExperience
	info.YearsOfExperience = account.YearsOfExperience
	info.FieldOfStudy = account.FieldOfStudy
	info.PreviousProjects = account.PreviousProjects
	// emit created event

	return info
}
