package repository

import (
	"time"

	"github.com/arravoco/hackathon_backend/exports"
)

type AdminAccountRepository struct {
	DB exports.AdminDatasourceQueryMethods
}

func NewAdminAccountRepository(datasource exports.AdminDatasourceQueryMethods) *AdminAccountRepository {
	return &AdminAccountRepository{
		DB: datasource,
	}
}

//exports.AdminRepositoryInterface

func (ad *AdminAccountRepository) CreateAdminAccount(dataInput *exports.CreateAdminAccountRepositoryDTO) (*exports.AdminAccountRepository, error) {

	acc, err := ad.DB.CreateAdminAccount(&exports.CreateAdminAccountData{
		FirstName:    dataInput.FirstName,
		LastName:     dataInput.LastName,
		Email:        dataInput.Email,
		PasswordHash: dataInput.PasswordHash,
		PhoneNumber:  dataInput.PhoneNumber,
		Gender:       dataInput.Gender,
		HackathonId:  dataInput.HackathonId,
		Status:       dataInput.Status,
		Role:         "ADMIN",
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return nil, err
	}
	/* raise event
	events.EmitAdminAccountCreated(&exports.AdminAccountCreatedEventData{
		Email:     acc.Email,
		LastName:  acc.LastName,
		FirstName: acc.FirstName,
		EventData: exports.EventData{EventName: string(events.AdminAccountCreatedEvent)},
	})*/

	return &exports.AdminAccountRepository{
		Id:           acc.Id.Hex(),
		FirstName:    acc.FirstName,
		LastName:     acc.LastName,
		Email:        acc.Email,
		Role:         acc.Role,
		Gender:       acc.Gender,
		HackathonId:  acc.HackathonId,
		Status:       acc.Status,
		PhoneNumber:  acc.PhoneNumber,
		PasswordHash: acc.PasswordHash,
	}, nil
}

func (ad *AdminAccountRepository) GetAdminAccountByEmail(email string) (*exports.AdminAccountRepository, error) {
	acc, err := ad.DB.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	return &exports.AdminAccountRepository{
		Id:           acc.Id.Hex(),
		FirstName:    acc.FirstName,
		LastName:     acc.LastName,
		Email:        acc.Email,
		Role:         acc.Role,
		Gender:       acc.Gender,
		HackathonId:  acc.HackathonId,
		Status:       acc.Status,
		PhoneNumber:  acc.PhoneNumber,
		PasswordHash: acc.PasswordHash,
	}, nil
}

func (ad *AdminAccountRepository) GetAdminAccounts(opts exports.FilterGetManyAccountRepositories) ([]exports.AdminAccountRepository, error) {
	accs, err := ad.DB.GetAccounts(exports.FilterGetManyAccountDocuments{})
	if err != nil {
		return nil, err
	}
	var accRepos []exports.AdminAccountRepository
	for _, acc := range accs {
		accRepos = append(accRepos, exports.AdminAccountRepository{
			Id:           acc.Id.Hex(),
			FirstName:    acc.FirstName,
			LastName:     acc.LastName,
			Email:        acc.Email,
			Role:         acc.Role,
			Gender:       acc.Gender,
			HackathonId:  acc.HackathonId,
			Status:       acc.Status,
			PhoneNumber:  acc.PhoneNumber,
			PasswordHash: acc.PasswordHash,
		})
	}
	return accRepos, nil
}

func (ad *AdminAccountRepository) UpdateAdminAccount(filter *exports.UpdateAccountRepositoryFilter, dataInput *exports.UpdateAdminAccountRepository) (*exports.AdminAccountRepository, error) {
	acc, err := ad.DB.UpdateAccountInfoByEmail(&exports.UpdateAccountDocumentFilter{
		Email:       filter.Email,
		PhoneNumber: filter.PhoneNumber,
	}, &exports.UpdateAccountDocument{
		FirstName:         dataInput.FirstName,
		LastName:          dataInput.LastName,
		Gender:            dataInput.Gender,
		State:             dataInput.State,
		IsEmailVerified:   dataInput.IsEmailVerified,
		IsEmailVerifiedAt: dataInput.IsEmailVerifiedAt,
	})
	if err != nil {
		return nil, err
	}

	return &exports.AdminAccountRepository{
		Id:           acc.Id.Hex(),
		FirstName:    acc.FirstName,
		LastName:     acc.LastName,
		Email:        acc.Email,
		Role:         acc.Role,
		Gender:       acc.Gender,
		HackathonId:  acc.HackathonId,
		Status:       acc.Status,
		PhoneNumber:  acc.PhoneNumber,
		PasswordHash: acc.PasswordHash,
	}, nil
}
