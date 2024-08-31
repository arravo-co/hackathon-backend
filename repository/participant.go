package repository

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	//"gitee.com/golang-module/carbon"

	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	//"github.com/golang-module/carbon"
)

// AddMemberToParticipatingTeam
type ParticipantRepository struct {
	DB                  *query.Query
	Entity              *entity.Participant
	FirstName           string                                      `json:"first_name"`
	LastName            string                                      `json:"last_name"`
	Email               string                                      `json:"email"`
	Gender              string                                      `json:"gender"`
	State               string                                      `json:"state"`
	Age                 int                                         `json:"age"`
	DOB                 time.Time                                   `json:"dob"`
	AccountRole         string                                      `json:"role"`
	ParticipantId       string                                      `json:"participant_id"`
	TeamLeadEmail       string                                      `json:"team_lead_email"`
	TeamName            string                                      `json:"team_name"`
	TeamRole            string                                      `json:"team_role"`
	HackathonId         string                                      `json:"hackathon_id"`
	ParticipantType     string                                      `json:"type"`
	CoParticipants      []CoParticipantInfo                         `json:"co_participants"`
	ParticipantEmail    string                                      `json:"participant_email"`
	InviteList          []exports.ParticipantDocumentTeamInviteInfo `json:"invite_list"`
	AccountStatus       string                                      `json:"account_status"`
	ParticipationStatus string                                      `json:"participation_status"`
	Skillset            []string                                    `json:"skillset"`
	PhoneNumber         string                                      `json:"phone_number"`
	EmploymentStatus    string                                      `json:"employment_status"`
	ExperienceLevel     string                                      `json:"experience_level"`
	Motivation          string                                      `json:"motivation"`
	HackathonExperience string                                      `json:"hackathon_experience"`
	YearsOfExperience   int                                         `json:"years_of_experience"`
	FieldOfStudy        string                                      `json:"field_of_study"`
	PreviousProjects    []string                                    `json:"previous_projects"`
	Solution            *SolutionRepository                         `json:"solution"`
	CreatedAt           time.Time                                   `json:"created_at"`
	UpdatedAt           time.Time                                   `json:"updated_at"`
}
type ParticipantRecord struct {
	ParticipantId       string                                      `json:"participant_id"`
	TeamLeadEmail       string                                      `json:"team_lead_email"`
	TeamName            string                                      `json:"team_name"`
	TeamRole            string                                      `json:"team_role"`
	HackathonId         string                                      `json:"hackathon_id"`
	ParticipantType     string                                      `json:"type"`
	CoParticipants      []CoParticipantInfo                         `json:"co_participants"`
	ParticipantEmail    string                                      `json:"participant_email"`
	InviteList          []exports.ParticipantDocumentTeamInviteInfo `json:"invite_list"`
	ParticipationStatus string                                      `json:"participation_status"`
	CreatedAt           time.Time                                   `json:"created_at"`
	UpdatedAt           time.Time                                   `json:"updated_at"`
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

func (p *ParticipantRepository) AddMemberInfoToParticipatingTeamRecord(dataToSave *exports.AddMemberToParticipatingTeamData) (*exports.ParticipantDocument, error) {
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

func (p ParticipantRepository) AddToTeamInviteList(dataInput *exports.AddToTeamInviteListData) (interface{}, error) {
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

func (p *ParticipantRepository) UpdateParticipantRecord(dataInput *dtos.AuthParticipantInfoUpdateDTO) error {
	_, err := p.DB.UpdateParticipantInfoByEmail(&exports.UpdateAccountFilter{Email: p.Email}, &exports.UpdateAccountDocument{
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
	slicesOfHash := strings.Split(hashedString, "")
	prefixSlices := slicesOfHash[0:5]
	postFix := slicesOfHash[len(slicesOfHash)-5:]
	sub := strings.Join([]string{"PARTICIPANT_ID_", strings.Join(append(prefixSlices, postFix...), "")}, "")
	return sub, nil
}

func (s *ParticipantRepository) GetParticipantRecord(participantId string) (*exports.ParticipantRecordRepository, error) {
	participant, err := s.DB.GetParticipantRecord(participantId)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	pE := &exports.ParticipantRecordRepository{
		ParticipantId:    participant.ParticipantId,
		TeamName:         participant.TeamName,
		ParticipantEmail: participant.ParticipantEmail,
		TeamLeadEmail:    participant.TeamLeadEmail,
		Id:               participant.Id.Hex(),
		ParticipantType:  participant.Type,
		SolutionId:       participant.SolutionId,
		HackathonId:      participant.HackathonId,
		InviteList:       participant.InviteList,
		Status:           participant.Status,
		ReviewRanking:    participant.ReviewRanking,
		CoParticipants:   participant.CoParticipants,
		CreatedAt:        participant.CreatedAt,
		UpdatedAt:        participant.UpdatedAt,
	}

	return pE, err
}

func (s *ParticipantRepository) GetParticipantsRecords() ([]exports.ParticipantRecordRepository, error) {
	partDocs, err := s.DB.GetParticipantsRecords()

	var participants []exports.ParticipantRecordRepository
	for _, participant := range partDocs {
		participants = append(participants, exports.ParticipantRecordRepository{

			ParticipantId:    participant.ParticipantId,
			TeamName:         participant.TeamName,
			ParticipantEmail: participant.ParticipantEmail,
			TeamLeadEmail:    participant.TeamLeadEmail,
			Id:               participant.Id.Hex(),
			ParticipantType:  participant.Type,
			SolutionId:       participant.SolutionId,
			HackathonId:      participant.HackathonId,
			InviteList:       participant.InviteList,
			Status:           participant.Status,
			ReviewRanking:    participant.ReviewRanking,
			CoParticipants:   participant.CoParticipants,
			CreatedAt:        participant.CreatedAt,
			UpdatedAt:        participant.UpdatedAt,
		})
	}
	return participants, err
}

func (s *ParticipantRepository) RemoveCoparticipantFromParticipantRecord(dataInput *RemoveMemberFromTeamData) (*entity.TeamMemberWithParticipantRecord, error) {
	_, err := s.DB.RemoveMemberFromParticipatingTeam(&exports.RemoveMemberFromParticipatingTeamData{
		HackathonId:   dataInput.HackathonId,
		MemberEmail:   dataInput.MemberEmail,
		ParticipantId: dataInput.ParticipantId,
	})
	if err != nil {
		return nil, err
	}
	acc, err := s.DB.DeleteAccount(dataInput.MemberEmail)
	info := FillTeamMemberInfo(acc)
	return info, err
}

//

func (s *ParticipantRepository) AddSolutionIdToParticipantRecord(dataInput *exports.SelectTeamSolutionData) (*entity.Solution, error) {

	partDoc, err := s.DB.SelectSolutionForTeam(&exports.SelectTeamSolutionData{
		HackathonId:   dataInput.HackathonId,
		ParticipantId: dataInput.ParticipantId,
		SolutionId:    dataInput.SolutionId,
	})
	if err != nil {
		return nil, err
	}
	return &entity.Solution{
		Id:          partDoc.SolutionId,
		HackathonId: partDoc.Solution.HackathonId,
		Description: partDoc.Solution.Description,
		Objective:   partDoc.Solution.Objective,
		Title:       partDoc.Solution.Title,
		CreatorId:   partDoc.Solution.CreatorId,
	}, nil
}

func FillTeamMemberInfo(account *exports.AccountDocument) *entity.TeamMemberWithParticipantRecord {
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
	info.Motivation = account.Motivation
	info.HackathonExperience = account.HackathonExperience
	info.YearsOfExperience = account.YearsOfExperience
	info.FieldOfStudy = account.FieldOfStudy
	info.PreviousProjects = account.PreviousProjects
	// emit created event

	return info
}

func (s *ParticipantRepository) GetSingleParticipantRecordAndMemberAccountsInfo(participant_id string) (*exports.ParticipantTeamMembersWithAccountsAggregate, error) {
	arggs, err := s.DB.GetParticipantsWithAccountsAggregate(nil)
	if err != nil {
		return nil, err
	}
	if arggs != nil {
		return nil, nil
	}
	arg := arggs[0]

	team_lead_info := exports.TeamLeadInfoParticipantRecordRepositoryAggregate{
		Email:         arg.TeamLeadInfo.Email,
		AccountId:     arg.TeamLeadInfo.AccountId,
		FirstName:     arg.TeamLeadInfo.FirstName,
		LastName:      arg.TeamLeadInfo.LastName,
		Gender:        arg.TeamLeadInfo.Gender,
		CreatedAt:     arg.TeamLeadInfo.CreatedAt,
		UpdateAt:      arg.TeamLeadInfo.UpdateAt,
		Skillset:      arg.TeamLeadInfo.Skillset,
		AccountStatus: arg.TeamLeadInfo.AccountStatus,
		AccountRole:   arg.TeamLeadInfo.AccountRole,
	}
	arr := &exports.ParticipantTeamMembersWithAccountsAggregate{
		Id:                arggs[0].Id.String(),
		TeamLeadEmail:     arggs[0].TeamLeadEmail,
		TeamName:          arggs[0].TeamName,
		TeamLeadFirstName: arggs[0].TeamLeadFirstName,
		TeamLeadLastName:  arggs[0].TeamLeadLastName,
		TeamLeadGender:    arggs[0].TeamLeadGender,
		TeamLeadAccountId: arggs[0].TeamLeadAccountId,
		TeamLeadInfo:      team_lead_info,
	}
	return arr, nil
}
func (s *ParticipantRepository) GetMultipleParticipantRecordAndMemberAccountsInfo(dataInput interface{}) (*exports.ParticipantTeamMembersWithAccountsAggregate, error) {
	arggs, err := s.DB.GetParticipantsWithAccountsAggregate(nil)
	if err != nil {
		return nil, err
	}
	var arr []*exports.ParticipantTeamMembersWithAccountsAggregate
	for _, v := range arggs {

	}
	return arr, nil
}
