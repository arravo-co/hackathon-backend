package exports

import (
	"time"
)

type JudgeAccountRepository struct {
	DB                *JudgeRepositoryInterface
	Id                string
	FirstName         string
	LastName          string
	Email             string
	PasswordHash      string
	Gender            string
	Role              string
	HackathonId       string
	Status            string
	State             string
	PhoneNumber       string
	Bio               string
	IsEmailVerified   bool
	ProfilePictureUrl string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type JudgeRepositoryInterface interface {
	GetJudgeAccountByEmail(email string) (*JudgeAccountRepository, error)
	CreateJudgeAccount(input *RegisterNewJudgeDTO) (*JudgeAccountRepository, error)
	UpdateJudgeAccount(email string, input *UpdateJudgeDTO) error
	GetJudges() ([]*JudgeAccountRepository, error)
	DeleteJudgeAccount(identifier string) (*JudgeAccountRepository, error)
	//GetJudgeAccountByEmail(email string) (*JudgeAccountRepository, error)
	//UpdateJudgeAccount(filter *UpdateAccountFilter, dataInput *UpdateAccountDocument) (*JudgeAccountRepository, error)
	UpdateJudgePassword(filter *UpdateAccountFilter, newPasswordHash string) (*JudgeAccountRepository, error)
}

type ParticipantAccountRepository struct {
	DB                *ParticipantAccountRepositoryInterface
	Id                string
	FirstName         string
	LastName          string
	Email             string
	PasswordHash      string
	Gender            string
	Role              string
	HackathonId       string
	Status            string
	State             string
	PhoneNumber       string
	Bio               string
	IsEmailVerified   bool
	ProfilePictureUrl string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	LinkedInAddress     string
	Skillset            []string
	ParticipantId       string
	DOB                 time.Time `bson:"dob,omitempty"`
	EmploymentStatus    string
	ExperienceLevel     string
	Motivation          string
	HackathonExperience string
	YearsOfExperience   int
	FieldOfStudy        string
	PreviousProjects    []string
	IsEmailVerifiedAt   time.Time
}

type ParticipantDocumentParticipantSelectedSolution struct {
	Id               string    `json:"id"`
	HackathonId      string    `json:"hackathon_id"`
	Title            string    `json:"name"`
	Description      string    `json:"description"`
	Objective        string    `json:"objective"`
	CreatorId        string    `json:"creator_id"`
	SolutionImageUrl string    `json:"solution_image_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ParticipantDocumentTeamCoParticipantInfo struct {
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	AccountId     string    `json:"account_id"`
	Role          string    `json:"team_role"`
	ParticipantId string    `json:"participant_id"`
	HackathonId   string    `json:"hackathon_id"`
	AddedToTeamAt time.Time `json:"added_to_team_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ParticipantRecordRepository struct {
	DB              ParticipantRepositoryInterface
	Id              string
	ParticipantId   string
	ParticipantType string
	SolutionId      string
	GithubAddress   string
	ReviewRanking   int
	//TeamLeadFirstName string
	//TeamLeadLastName  string
	//TeamLeadGender    string
	TeamLeadAccountId string
	TeamLeadEmail     string
	TeamName          string
	//TeamRole          string
	Type             string
	HackathonId      string
	CoParticipants   []ParticipantDocumentTeamCoParticipantInfo
	ParticipantEmail string
	InviteList       []ParticipantDocumentTeamInviteInfo
	Status           string
	Solution         *ParticipantDocumentParticipantSelectedSolution
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ParticipantTeamMembersWithAccountsAggregate struct {
	Id               string
	ParticipantId    string
	HackathonId      string
	Type             string
	TeamLeadEmail    string
	SolutionId       string
	Solution         ParticipantDocumentParticipantSelectedSolution
	TeamName         string
	ParticipantEmail string
	GithubAddress    string
	InviteList       []ParticipantDocumentTeamInviteInfo
	ReviewRanking    int
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time

	TeamLeadFirstName string
	TeamLeadLastName  string
	TeamLeadGender    string
	TeamLeadAccountId string

	TeamLeadInfo   TeamLeadInfoParticipantRecordRepositoryAggregate
	CoParticipants []CoParticipantAggregateData
}

type TeamLeadInfoParticipantRecordRepositoryAggregate struct {
	HackathonId   string
	AccountId     string
	Email         string
	FirstName     string
	LastName      string
	Gender        string
	PhoneNumber   string
	Skillset      []string
	AccountStatus string
	AccountRole   string
	TeamRole      string
	PasswordHash  string
	State         string
	CreatedAt     string
	UpdateAt      string
}

// CreateTeamMemberAccount
type CoParticipantAggregateData struct {
	HackathonId   string
	AccountId     string
	Email         string
	FirstName     string
	LastName      string
	Gender        string
	PhoneNumber   string
	Skillset      []string
	AccountStatus string
	AccountRole   string
	TeamRole      string
	PasswordHash  string
	State         string
	CreatedAt     string
	UpdateAt      string
}
type ParticipantAccountRepositoryInterface interface {
	GetParticipantAccountByEmail(email string) (*ParticipantAccountRepository, error)
	GetParticipantAccountsByEmail(emails []string) ([]*ParticipantAccountRepository, error)
	CreateParticipantAccount(input *CreateParticipantAccountData) (*ParticipantAccountRepository, error)
	//CreateCoParticipantAccount(input *RegisterNewParticipantAccountDTO) (*ParticipantAccountRepository, error)
	UpdateParticipantAccount(email string, input *UpdateParticipantDTO) error
	GetParticipantAccounts() ([]*ParticipantAccountRepository, error)
	DeleteParticipantAccount(identifier string) (*ParticipantAccountRepository, error)
	MarkParticipantAccountAsDeleted(identifier string) (*ParticipantAccountRepository, error)
	//GetJudgeAccountByEmail(email string) (*JudgeAccountRepository, error)
	//UpdateJudgeAccount(filter *UpdateAccountFilter, dataInput *UpdateAccountDocument) (*JudgeAccountRepository, error)
	UpdateParticipantPassword(filter *UpdateAccountFilter, newPasswordHash string) (*ParticipantAccountRepository, error)
}

type ParticipantRepositoryInterface interface {
	AddMemberInfoToParticipatingTeamRecord(dataToSave *AddMemberToParticipatingTeamData) (*ParticipantRecordRepository, error)
	AddSolutionIdToParticipantRecord(dataInput *SelectTeamSolutionData) (*ParticipantRecordRepository, error)
	AddToTeamInviteList(dataInput *AddToTeamInviteListData) (interface{}, error)
	CreateParticipantRecord(dataInput *CreateParticipantRecordData) (*ParticipantRecordRepository, error)
	GetParticipantRecord(participantId string) (*ParticipantRecordRepository, error)
	GetParticipantsRecords() ([]ParticipantRecordRepository, error)
	//RegisterIndividual(input RegisterNewParticipantAccountDTO) (*ParticipantRecordRepository, error)
	//RegisterTeamLead(input RegisterNewParticipantAccountDTO) (*ParticipantRecordRepository, error)
	RemoveCoparticipantFromParticipantRecord(dataInput *RemoveMemberFromTeamData) (*ParticipantRecordRepository, error)
	GetSingleParticipantRecordAndMemberAccountsInfo(participant_id string) (*ParticipantTeamMembersWithAccountsAggregate, error)
}
