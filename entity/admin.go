package entity

import (
	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/events"
	"github.com/arravoco/hackathon_backend/exports"
)

type Admin struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	passwordHash string
	Gender       string   `json:"gender"`
	Role         string   `json:"role"`
	HackathonId  string   `json:"hackathon_id"`
	Status       string   `json:"status"`
	Skillset     []string `json:"skillset"`
	PhoneNumber  string   `json:"phone_number"`
}

func (ad *Admin) RegisterNewAdmin(dataInput *dtos.CreateNewAdminDTO) error {
	acc, err := data.CreateAdminAccount(&exports.CreateAdminAccountData{
		Role:        "ADMIN",
		Email:       dataInput.Email,
		FirstName:   dataInput.FirstName,
		LastName:    dataInput.LastName,
		Gender:      dataInput.Gender,
		PhoneNumber: dataInput.PhoneNumber,
		HackathonId: config.GetHackathonId(),
	})
	if err != nil {
		return err
	}
	// raise event
	events.EmitAdminAccountCreated(&exports.AdminAccountCreatedEventData{
		Email:     acc.Email,
		LastName:  acc.LastName,
		FirstName: acc.FirstName,
		EventData: exports.EventData{EventName: string(events.AdminAccountCreatedEvent)},
	})

	ad.passwordHash = acc.PasswordHash
	ad.Email = acc.Email
	ad.FirstName = acc.FirstName
	ad.LastName = acc.LastName
	ad.HackathonId = acc.HackathonId
	ad.Gender = acc.Gender
	ad.Role = acc.Role

	return nil
}

func (ad *Admin) AdminCreateNewAdminProlife(dataInput *dtos.CreateNewAdminDTO) error {
	password := exports.GeneratePassword()
	passwordHash, err := exports.GenerateHashPassword(password)
	if err != nil {
		return err
	}
	acc, err := data.CreateAdminAccount(&exports.CreateAdminAccountData{
		Role:         "ADMIN",
		Email:        dataInput.Email,
		FirstName:    dataInput.FirstName,
		LastName:     dataInput.LastName,
		Gender:       dataInput.Gender,
		PhoneNumber:  dataInput.PhoneNumber,
		HackathonId:  config.GetHackathonId(),
		PasswordHash: passwordHash,
	})
	if err != nil {
		return err
	}
	// raise event
	events.EmitAdminAccountCreatedByAdmin(&exports.AdminAccountCreatedByAdminEventData{
		Email:     acc.Email,
		LastName:  acc.LastName,
		FirstName: acc.FirstName,
		EventData: exports.EventData{EventName: string(events.AdminAccountCreatedEvent)},
		Password:  password,
	})

	ad.passwordHash = acc.PasswordHash
	ad.Email = acc.Email
	ad.FirstName = acc.FirstName
	ad.LastName = acc.LastName
	ad.HackathonId = acc.HackathonId
	ad.Gender = acc.Gender
	ad.Role = acc.Role

	return nil
}
