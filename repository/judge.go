package repository

import (
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

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

	err := validate.Struct(input)
	if err != nil {
		return nil, err
	}
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
	id := dataResponse.Id.Hex()
	judge := &exports.JudgeAccountRepository{
		Id:                id,
		FirstName:         dataResponse.FirstName,
		LastName:          dataResponse.LastName,
		Email:             dataResponse.Email,
		State:             dataInput.State,
		Gender:            dataInput.Gender,
		HackathonId:       dataInput.HackathonId,
		PhoneNumber:       dataInput.PhoneNumber,
		ProfilePictureUrl: dataInput.ProfilePictureUrl,
		PasswordHash:      dataInput.PasswordHash,
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

func (dt *JudgeAccountRepository) UpdateJudgeAccount(email string, input *exports.UpdateJudgeDTO) error {

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
	_, err := dt.DB.UpdateAccountInfoByEmail(&exports.UpdateAccountDocumentFilter{
		Email: email,
	}, dataInput)
	// emit created event
	if err != nil {
		return err
	}
	return nil
}

func (dt *JudgeAccountRepository) GetJudges() ([]*exports.JudgeAccountRepository, error) {

	dataResponse, err := dt.DB.GetAccountsOfJudges()
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
			PasswordHash:      acc.PasswordHash,
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
	acc, err := dt.DB.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}
	id := acc.Id.Hex()
	j := &exports.JudgeAccountRepository{
		Id:                id,
		Email:             acc.Email,
		LastName:          acc.LastName,
		FirstName:         acc.FirstName,
		Bio:               acc.Bio,
		PhoneNumber:       acc.PhoneNumber,
		Status:            acc.Status,
		Role:              acc.Role,
		PasswordHash:      acc.PasswordHash,
		State:             acc.State,
		HackathonId:       acc.HackathonId,
		IsEmailVerified:   acc.IsEmailVerified,
		ProfilePictureUrl: acc.ProfilePictureUrl,
		Gender:            acc.Gender,
		CreatedAt:         acc.CreatedAt,
		UpdatedAt:         acc.UpdatedAt,
	}
	return j, nil
}

func (dt *JudgeAccountRepository) UpdateJudgePassword(filter *exports.UpdateAccountDocumentFilter, newPasswordHash string) (*exports.JudgeAccountRepository, error) {
	return nil, nil
}
