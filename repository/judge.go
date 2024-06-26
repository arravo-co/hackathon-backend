package repository

import (
	"time"

	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/exports"
)

type Judge struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	passwordHash      string
	Gender            string    `json:"gender"`
	Role              string    `json:"role"`
	HackathonId       string    `json:"hackathon_id"`
	Status            string    `json:"status"`
	State             string    `json:"state"`
	PhoneNumber       string    `json:"phone_number"`
	Bio               string    `json:"bio"`
	ProfilePictureUrl string    `json:"profile_picture_url"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func GetJudgeEntity(email string) (*Judge, error) {
	acc, err := data.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}
	judge := Judge{}
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
	return &judge, nil
}

func (judge *Judge) FillJudgeEntity(email string) error {
	acc, err := data.GetAccountByEmail(email)
	if err != nil {
		return err
	}

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
	return nil
}

func (p *Judge) Register(input dtos.RegisterNewJudgeDTO) (*exports.CreateJudgeAccountData, error) {

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
			Role:         "JUDGE", Status: "EMAIL_UNVERIFIED"},
		Bio: input.Bio,
	}
	dataResponse, err := data.CreateJudgeAccount(dataInput)
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
	return dataResponse, err
}

func (p *Judge) UpdateJudgeProfile(input dtos.UpdateJudgeDTO) error {

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
	dataResponse, err := data.UpdateAccountInfoByEmail(&exports.UpdateAccountFilter{
		Email: p.Email,
	}, dataInput)
	// emit created event
	if err != nil {
		return err
	}
	p.FirstName = dataResponse.FirstName
	p.LastName = dataResponse.LastName
	p.Gender = dataResponse.Gender
	p.Bio = dataResponse.Bio

	p.ProfilePictureUrl = dataResponse.ProfilePictureUrl
	p.Role = dataResponse.Role
	p.State = dataResponse.State
	/*events.EmitJudgeAccountCreated(&exports.JudgeAccountCreatedByAdminEventData{
		JudgeEmail: dataResponse.Email,
		LastName:   dataResponse.LastName,
		FirstName:  dataResponse.FirstName,
		EventData:  exports.EventData{EventName: "JudgeAccountCreated"},
	})
	*/
	return nil
}

func GetJudges() ([]*Judge, error) {

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
	var ent []*Judge
	for _, acc := range dataResponse {
		ent = append(ent, &Judge{
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
