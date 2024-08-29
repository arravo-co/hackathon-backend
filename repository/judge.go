package repository

import (
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/exports"
)

type JudgeAccountRepository struct {
	DB exports.JudgeDatasourceQueryMethods
}

func NewJudgeAccountRepository(datasource exports.JudgeDatasourceQueryMethods) *JudgeAccountRepository {
	return &JudgeAccountRepository{
		DB: datasource,
	}
}

func (ac *JudgeAccountRepository) GetJudgeByEmail(email string) (*exports.JudgeAccountRepository, error) {
	acc, err := ac.DB.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}
	judge := &exports.JudgeAccountRepository{}
	judge.Email = acc.Email
	judge.Role = acc.Role
	judge.FirstName = acc.FirstName
	judge.LastName = acc.LastName
	judge.Gender = acc.Gender
	judge.HackathonId = acc.HackathonId
	judge.PhoneNumber = acc.PhoneNumber
	judge.Status = acc.Status
	judge.State = acc.State
	judge.Bio = acc.Bio
	judge.ProfilePictureUrl = acc.ProfilePictureUrl
	judge.CreatedAt = acc.CreatedAt
	judge.UpdatedAt = acc.UpdatedAt
	return judge, nil
}

func (dt *JudgeAccountRepository) CreateJudgeAccount(input *exports.RegisterNewJudgeDTO) (*exports.JudgeAccountRepository, error) {

	passwordHash, err := exports.GenerateHashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	dataInput := &exports.CreateJudgeAccountData{
		CreateAccountData: exports.CreateAccountData{
			Email:        input.Email,
			PasswordHash: passwordHash,
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			Gender:       input.Gender,
			State:        input.State,
			Role:         "JUDGE",
			Status:       "EMAIL_UNVERIFIED",
		},
		Bio: input.Bio,
	}
	dataResponse, err := dt.DB.CreateJudgeAccount(dataInput)
	// emit created event
	if err != nil {
		return nil, err
	}
	judge := &exports.JudgeAccountRepository{
		FirstName:         dataResponse.FirstName,
		LastName:          dataResponse.LastName,
		Email:             dataResponse.Email,
		State:             dataInput.State,
		Gender:            dataInput.Gender,
		HackathonId:       dataInput.HackathonId,
		PhoneNumber:       dataInput.PhoneNumber,
		ProfilePictureUrl: dataInput.ProfilePictureUrl,
		IsEmailVerified:   dataResponse.IsEmailVerified,
		Role:              dataResponse.Role,
		Status:            dataResponse.Status,
		Bio:               dataResponse.Bio,
		CreatedAt:         dataResponse.CreatedAt,
	}
	/*events.EmitJudgeAccountCreated(&exports.JudgeAccountCreatedByAdminEventData{
		JudgeEmail: dataResponse.Email,
		LastName:   dataResponse.LastName,
		FirstName:  dataResponse.FirstName,
		EventData:  exports.EventData{EventName: "JudgeAccountCreated"},
	})
	*/
	return judge, err
}

func (dt *JudgeAccountRepository) UpdateJudgeProfile(email string, input dtos.UpdateJudgeDTO) error {

	dataInput := &exports.UpdateAccountDocument{}
	if input.LastName != "" {
		dataInput.LastName = input.LastName
	}
	if input.State != "" {
		dataInput.State = input.State
	}
	if input.FirstName != "" {
		dataInput.FirstName = input.FirstName
	}
	if input.Bio != "" {
		dataInput.Bio = input.Bio
	}
	if input.Gender != "" {
		dataInput.Gender = input.Gender
	}
	if input.ProfilePictureUrl != "" {
		dataInput.ProfilePictureUrl = input.ProfilePictureUrl
	}
	_, err := dt.DB.UpdateAccountInfoByEmail(&exports.UpdateAccountFilter{
		Email: email,
	}, dataInput)
	// emit created event
	if err != nil {
		return err
	}
	return nil
}

func (dt *JudgeAccountRepository) GetJudges() ([]*exports.JudgeAccountRepository, error) {

	dataResponse, err := data.GetAccountsOfJudges()
	// emit created event
	if err != nil {
		return nil, err
	}
	/*events.EmitJudgeAccountCreated(&exports.JudgeAccountCreatedByAdminEventData{
		JudgeEmail: dataResponse.Email,
		LastName:   dataResponse.LastName,
		FirstName:  dataResponse.FirstName,
		EventData:  exports.EventData{EventName: "JudgeAccountCreated"},
	})
	*/
	var ent []*exports.JudgeAccountRepository
	for _, acc := range dataResponse {
		ent = append(ent, &exports.JudgeAccountRepository{
			FirstName:         acc.FirstName,
			LastName:          acc.LastName,
			Email:             acc.Email,
			Gender:            acc.Gender,
			Role:              acc.Role,
			HackathonId:       acc.HackathonId,
			Status:            acc.Status,
			State:             acc.State,
			PhoneNumber:       acc.PhoneNumber,
			ProfilePictureUrl: acc.ProfilePictureUrl,
			Bio:               acc.Bio,
			CreatedAt:         acc.CreatedAt,
			UpdatedAt:         acc.UpdatedAt,
		})
	}
	return ent, err
}

func (dt *JudgeAccountRepository) DeleteJudgeAccount(identifier string) (*exports.JudgeAccountRepository, error) {
	return nil, nil
}
func (dt *JudgeAccountRepository) GetJudgeAccountByEmail(email string) (*exports.JudgeAccountRepository, error) {
	return nil, nil
}
func (dt *JudgeAccountRepository) UpdateJudgeAccountInfoByEmail(filter *exports.UpdateAccountFilter, dataInput *exports.UpdateAccountDocument) (*exports.JudgeAccountRepository, error) {
	return nil, nil
}
func (dt *JudgeAccountRepository) UpdateJudgePasswordByEmail(filter *exports.UpdateAccountFilter, newPasswordHash string) (*exports.JudgeAccountRepository, error) {
	return nil, nil
}
