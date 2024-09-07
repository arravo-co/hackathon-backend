package exports

import (
	"time"
)

type AccountRepository struct {
	DB                *AccountRepositoryInterface
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

type AccountRepositoryInterface interface {
	GetAccountByEmail(email string) (*AccountRepository, error)
	GetAccountsByEmail(emails []string) ([]*AccountRepository, error)
	CreateAccount(input *CreateAccountData) (*AccountRepository, error)
	//CreateCoParticipantAccount(input *RegisterNewParticipantAccountDTO) (*ParticipantAccountRepository, error)
	UpdateAccount(email string, input *UpdateAccountDTO) error
	GetAccounts() ([]*AccountRepository, error)
	DeleteAccount(identifier string) (*AccountRepository, error)
	MarkAccountAsDeleted(identifier string) (*AccountRepository, error)
	//GetJudgeAccountByEmail(email string) (*JudgeAccountRepository, error)
	//UpdateJudgeAccount(filter *UpdateAccountFilter, dataInput *UpdateAccountDocument) (*JudgeAccountRepository, error)
	FindAccountIdentifier(identifier string) (*AccountRepository, error)
	UpdateParticipantPassword(filter *UpdateAccountDocumentFilter, newPasswordHash string) (*AccountRepository, error)
	ChangePassword(dataInput *PasswordChangeData) (*AccountRepository, error)
}

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
	UpdateJudgePassword(filter *UpdateAccountDocumentFilter, newPasswordHash string) (*JudgeAccountRepository, error)
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
	Id               string    `bson:"_id"`
	HackathonId      string    `bson:"hackathon_id"`
	Title            string    `bson:"title"`
	Description      string    `bson:"description"`
	Objective        string    `bson:"objective"`
	CreatorId        string    `bson:"creator_id"`
	SolutionImageUrl string    `bson:"solution_image_url"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
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
	HackathonId         string
	AccountId           string
	Email               string
	FirstName           string
	LastName            string
	Gender              string
	PhoneNumber         string
	Skillset            []string
	AccountStatus       string
	AccountRole         string
	TeamRole            string
	passwordHash        string
	State               string
	LinkedInAddress     string
	ParticipantId       string
	DOB                 time.Time `bson:"dob,omitempty"`
	EmploymentStatus    string
	ExperienceLevel     string
	Motivation          string
	HackathonExperience string
	YearsOfExperience   int
	FieldOfStudy        string
	PreviousProjects    []string
	IsEmailVerified     bool
	IsEmailVerifiedAt   time.Time
	CreatedAt           time.Time
	UpdateAt            time.Time
}

// CreateTeamMemberAccount
type CoParticipantAggregateData struct {
	HackathonId         string
	AccountId           string
	Email               string
	FirstName           string
	LastName            string
	Gender              string
	PhoneNumber         string
	Skillset            []string
	AccountStatus       string
	AccountRole         string
	TeamRole            string
	passwordHash        string
	State               string
	LinkedInAddress     string
	ParticipantId       string
	DOB                 time.Time `bson:"dob,omitempty"`
	EmploymentStatus    string
	ExperienceLevel     string
	Motivation          string
	HackathonExperience string
	YearsOfExperience   int
	FieldOfStudy        string
	PreviousProjects    []string
	IsEmailVerified     bool
	IsEmailVerifiedAt   time.Time
	CreatedAt           time.Time
	UpdateAt            time.Time
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
	//UpdateJudgeAccount(filter *UpdateAccountFilter, dataInput *UpdateAccountDocument) (*JudgeAccountRepository, error)
	FindAccountIdentifier(identifier string) (*ParticipantAccountRepository, error)
	UpdateParticipantPassword(filter *UpdateAccountDocumentFilter, newPasswordHash string) (*ParticipantAccountRepository, error)
}

type ParticipantRepositoryInterface interface {
	AddMemberInfoToParticipatingTeamRecord(dataToSave *AddMemberToParticipatingTeamData) (*ParticipantRecordRepository, error)
	AddSolutionIdToParticipantRecord(dataInput *SelectTeamSolutionData) (*ParticipantRecordRepository, error)
	AddToTeamInviteList(dataInput *AddToTeamInviteListData) (interface{}, error)
	CreateParticipantRecord(dataInput *CreateParticipantRecordData) (*ParticipantRecordRepository, error)
	GetParticipantRecord(participantId string) (*ParticipantRecordRepository, error)
	GetParticipantsRecords() ([]ParticipantRecordRepository, error)
	AdminUpdateParticipantRecord(filterOpts *UpdateSingleParticipantRecordFilter, dataInput *AdminParticipantInfoUpdateDTO) (*ParticipantRecordRepository, error)
	//RegisterIndividual(input RegisterNewParticipantAccountDTO) (*ParticipantRecordRepository, error)
	//RegisterTeamLead(input RegisterNewParticipantAccountDTO) (*ParticipantRecordRepository, error)
	RemoveCoparticipantFromParticipantRecord(dataInput *RemoveMemberFromTeamData) (*ParticipantRecordRepository, error)
	GetSingleParticipantRecordAndMemberAccountsInfo(participant_id string) (*ParticipantTeamMembersWithAccountsAggregate, error)
	GetMultipleParticipantRecordAndMemberAccountsInfo(dataInput GetParticipantsWithAccountsAggregateFilterOpts) ([]*ParticipantTeamMembersWithAccountsAggregate, error)
}

type TokenDataRepository struct {
	Id             string
	Token          string
	TokenType      string
	TokenTypeValue string
	Scope          string
	TTL            time.Time
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type TokenRepositoryInterface interface {
	UpsertToken(dataInput *UpsertTokenData) (*TokenDataRepository, error)
	VerifyToken(dataInput *VerifyTokenData) error
}

type AdminRepositoryInterface interface {
	CreateAdminAccount(dataToSave *CreateAdminAccountRepositoryDTO) (*AdminAccountRepository, error)
	//DeleteAdminAccount(identifier string) (*AdminAccountRepository, error)
	GetAdminAccountByEmail(email string) (*AdminAccountRepository, error)
	GetAdminAccounts(FilterGetManyAccountRepositories) ([]AdminAccountRepository, error)
	UpdateAdminAccount(filter *UpdateAccountRepositoryFilter, dataInput *UpdateAdminAccountRepository) (*AdminAccountRepository, error)
}

type AdminAccountRepository struct {
	Id              string
	FirstName       string
	LastName        string
	Email           string
	PasswordHash    string
	Gender          string
	Role            string
	HackathonId     string
	Status          string
	PhoneNumber     string
	IsEmailVerified bool      `json:"is_email_verified"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
