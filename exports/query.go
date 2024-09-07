package exports

type DatasourceQueryMethods interface {
	AddMemberToParticipatingTeam(dataToSave *AddMemberToParticipatingTeamData) (*ParticipantDocument, error)
	AddToTeamInviteList(dataToSave *AddToTeamInviteListData) (interface{}, error)
	CreateAccount(dataToSave *CreateAccountData) (interface{}, error)
	CreateAdminAccount(dataToSave *CreateAdminAccountData) (*AccountDocument, error)
	CreateJudgeAccount(dataToSave *CreateJudgeAccountData) (*AccountDocument, error)
	CreateParticipantAccount(dataToSave *CreateParticipantAccountData) (*AccountDocument, error)
	CreateParticipantRecord(dataToSave *CreateParticipantRecordData) (*ParticipantDocument, error)
	CreateSolutionData(dataInput *CreateSolutionData) (*SolutionDocument, error)
	CreateTeamMemberAccount(dataToSave *CreateTeamMemberAccountData) (*AccountDocument, error)
	DeleteAccount(identifier string) (*AccountDocument, error)
	FindAccountIdentifier(identifier string) (*AccountDocument, error)
	GetAccountByEmail(email string) (*AccountDocument, error)
	GetAccountsByEmails(emails []string) ([]AccountDocument, error)
	GetAccounts(FilterGetManyAccountDocuments) ([]AccountDocument, error)
	GetAccountsByParticipantIds(participantIds []string) ([]AccountDocument, error)
	GetAccountsOfJudges() ([]AccountDocument, error)
	GetManySolutionData(filterInput interface{}) ([]SolutionDocument, error)
	GetParticipantRecord(participantId string) (*ParticipantDocument, error)
	GetParticipantsRecords() ([]ParticipantDocument, error)
	GetSolutionDataById(id string) (*SolutionDocument, error)
	RemoveMemberFromParticipatingTeam(dataToSave *RemoveMemberFromParticipatingTeamData) (interface{}, error)
	RemoveTeamMemberAccount(dataToSave *RemoveTeamMemberAccountData) (*AccountDocument, error)
	SelectSolutionForTeam(dataToSave *SelectTeamSolutionData) (*ParticipantDocument, error)
	UpdateAccountInfoByEmail(filter *UpdateAccountDocumentFilter, dataInput *UpdateAccountDocument) (*AccountDocument, error)
	UpdateParticipantInfoByEmail(filter *UpdateAccountDocumentFilter, dataInput *UpdateAccountDocument) (*AccountDocument, error)
	UpdatePasswordByEmail(filter *UpdateAccountDocumentFilter, newPasswordHash string) (*AccountDocument, error)
	UpdateSolutionData(id string, dataInput *UpdateSolutionData) (*SolutionDocument, error)
	UpsertToken(dataInput *UpsertTokenData) (*TokenData, error)
	VerifyToken(dataInput *VerifyTokenData) error
}

type JudgeDatasourceQueryMethods interface {
	CreateAccount(dataToSave *CreateAccountData) (*AccountDocument, error)
	CreateJudgeAccount(dataToSave *CreateJudgeAccountData) (*AccountDocument, error)
	DeleteAccount(identifier string) (*AccountDocument, error)
	FindAccountIdentifier(identifier string) (*AccountDocument, error)
	GetAccountByEmail(email string) (*AccountDocument, error)
	GetAccountsByEmails(emails []string) ([]AccountDocument, error)
	GetAccountsOfJudges() ([]AccountDocument, error)
	UpdateAccountInfoByEmail(filter *UpdateAccountDocumentFilter, dataInput *UpdateAccountDocument) (*AccountDocument, error)
	UpdatePasswordByEmail(filter *UpdateAccountDocumentFilter, newPasswordHash string) (*AccountDocument, error)
}

type AdminDatasourceQueryMethods interface {
	CreateAdminAccount(dataToSave *CreateAdminAccountData) (*AccountDocument, error)
	DeleteAccount(identifier string) (*AccountDocument, error)
	FindAccountIdentifier(identifier string) (*AccountDocument, error)
	GetAccounts(FilterGetManyAccountDocuments) ([]AccountDocument, error)
	GetAccountByEmail(email string) (*AccountDocument, error)
	GetAccountsByEmails(emails []string) ([]AccountDocument, error)
	GetAccountsOfJudges() ([]AccountDocument, error)
	UpdateAccountInfoByEmail(filter *UpdateAccountDocumentFilter, dataInput *UpdateAccountDocument) (*AccountDocument, error)
	UpdatePasswordByEmail(filter *UpdateAccountDocumentFilter, newPasswordHash string) (*AccountDocument, error)
}

type TokenDatasourceQueryMethods interface {
	UpsertToken(dataInput *UpsertTokenData) (*TokenData, error)
	VerifyToken(dataInput *VerifyTokenData) error
}
