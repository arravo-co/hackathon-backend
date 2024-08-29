package exports

/*type QueryWithDatasource struct {
	Datasource   DBInterface
	QueryMethods DatasourceQueryMethods
}

type Query struct {
	Methods DatasourceQueryMethods
}*/

type DatasourceQueryMethods interface {
	AddMemberToParticipatingTeam(dataToSave *AddMemberToParticipatingTeamData) (*ParticipantDocument, error)
	AddToTeamInviteList(dataToSave *AddToTeamInviteListData) (interface{}, error)
	CreateAccount(dataToSave *CreateAccountData) (interface{}, error)
	CreateAdminAccount(dataToSave *CreateAdminAccountData) (*AccountDocument, error)
	CreateJudgeAccount(dataToSave *CreateJudgeAccountData) (*CreateJudgeAccountData, error)
	CreateParticipantAccount(dataToSave *CreateParticipantAccountData) (*AccountDocument, error)
	CreateParticipantRecord(dataToSave *CreateParticipantRecordData) (*ParticipantDocument, error)
	CreateSolutionData(dataInput *CreateSolutionData) (*SolutionDocument, error)
	CreateTeamMemberAccount(dataToSave *CreateTeamMemberAccountData) (*AccountDocument, error)
	DeleteAccount(identifier string) (*AccountDocument, error)
	FindAccountIdentifier(identifier string) (*AccountDocument, error)
	GetAccountByEmail(email string) (*AccountDocument, error)
	GetAccountsByEmails(emails []string) ([]AccountDocument, error)
	GetAccountsByParticipantIds(participantIds []string) ([]AccountDocument, error)
	GetAccountsOfJudges() ([]AccountDocument, error)
	GetManySolutionData(filterInput interface{}) ([]SolutionDocument, error)
	GetParticipantRecord(participantId string) (*ParticipantDocument, error)
	GetParticipantsRecords() ([]ParticipantDocument, error)
	GetParticipantsRecordsAggregate() ([]ParticipantAccountWithCoParticipantsDocument, error)
	GetSolutionDataById(id string) (*SolutionDocument, error)
	RemoveMemberFromParticipatingTeam(dataToSave *RemoveMemberFromParticipatingTeamData) (interface{}, error)
	RemoveTeamMemberAccount(dataToSave *RemoveTeamMemberAccountData) (*AccountDocument, error)
	SelectSolutionForTeam(dataToSave *SelectTeamSolutionData) (*ParticipantDocument, error)
	UpdateAccountInfoByEmail(filter *UpdateAccountFilter, dataInput *UpdateAccountDocument) (*AccountDocument, error)
	UpdateParticipantInfoByEmail(filter *UpdateAccountFilter, dataInput *UpdateAccountDocument) (*AccountDocument, error)
	UpdatePasswordByEmail(filter *UpdateAccountFilter, newPasswordHash string) (*AccountDocument, error)
	UpdateSolutionData(id string, dataInput *UpdateSolutionData) (*SolutionDocument, error)
}

type JudgeDatasourceQueryMethods interface {
	CreateAccount(dataToSave *CreateAccountData) (interface{}, error)
	CreateJudgeAccount(dataToSave *CreateJudgeAccountData) (*CreateJudgeAccountData, error)
	DeleteAccount(identifier string) (*AccountDocument, error)
	FindAccountIdentifier(identifier string) (*AccountDocument, error)
	GetAccountByEmail(email string) (*AccountDocument, error)
	GetAccountsByEmails(emails []string) ([]AccountDocument, error)
	GetAccountsOfJudges() ([]AccountDocument, error)
	UpdateAccountInfoByEmail(filter *UpdateAccountFilter, dataInput *UpdateAccountDocument) (*AccountDocument, error)
	UpdatePasswordByEmail(filter *UpdateAccountFilter, newPasswordHash string) (*AccountDocument, error)
}
