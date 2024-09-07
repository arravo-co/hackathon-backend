package services

import (
	"encoding/json"
	"fmt"

	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/rabbitmq/amqp091-go"
)

func (s *Service) RegisterAdmin(dataInput *CreateNewAdminDTO) (*entity.Admin, error) {

	passwordHash, err := exports.GenerateHashPassword(dataInput.Password)
	if err != nil {
		return nil, err
	}
	accAdmin, err := s.AdminAccountRepository.CreateAdminAccount(&exports.CreateAdminAccountRepositoryDTO{
		FirstName:    dataInput.FirstName,
		LastName:     dataInput.LastName,
		Email:        dataInput.Email,
		HackathonId:  dataInput.HackathonId,
		Role:         "ADMIN",
		PasswordHash: passwordHash,
		PhoneNumber:  dataInput.PhoneNumber,
		Status:       "EMAIL_UNVERIFIED",
	})
	if err != nil {
		return nil, err
	}

	pubData := &exports.AdminRegisteredPublishPayload{
		Email:     dataInput.Email,
		FirstName: dataInput.FirstName,
		LastName:  dataInput.LastName,
	}
	by, err := json.Marshal(pubData)
	if err != nil {
		fmt.Printf("Failed to marshal: %v\n", err)
	} else {
		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.AdminsExchange,
			KeyName:      exports.AdminRegisteredRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.AdminsExchange,
			KeyName:      exports.AdminSendWelcomeEmailRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

	}
	return &entity.Admin{
		AccountId:       accAdmin.Id,
		HackathonId:     accAdmin.HackathonId,
		FirstName:       accAdmin.FirstName,
		LastName:        accAdmin.LastName,
		Gender:          accAdmin.Gender,
		IsEmailVerified: accAdmin.IsEmailVerified,
		Email:           accAdmin.Email,
		EmailVerifiedAt: accAdmin.EmailVerifiedAt,
		Status:          accAdmin.Status,
		State:           accAdmin.Status,
	}, nil
}

func (s *Service) AdminCreateNewAdminProfile(dataInput *CreateNewAdminByAuthAdminDTO) (*entity.Admin, error) {
	password := exports.GeneratePassword()
	passwordHash, err := exports.GenerateHashPassword(password)
	if err != nil {
		return nil, err
	}

	accAdmin, err := s.AdminAccountRepository.CreateAdminAccount(&exports.CreateAdminAccountRepositoryDTO{
		FirstName:    dataInput.FirstName,
		LastName:     dataInput.LastName,
		Email:        dataInput.Email,
		Gender:       dataInput.Gender,
		HackathonId:  dataInput.HackathonId,
		Role:         "ADMIN",
		PasswordHash: passwordHash,
		PhoneNumber:  dataInput.PhoneNumber,
		Status:       "EMAIL_UNVERIFIED",
	})
	if err != nil {
		return nil, err
	}

	pubData := &exports.AdminRegisteredByAdminPublishPayload{
		Email:        dataInput.Email,
		Name:         dataInput.FirstName,
		InviterName:  dataInput.InviterName,
		InviterEmail: dataInput.Email,
		Password:     password,
	}
	by, err := json.Marshal(pubData)
	if err != nil {
		fmt.Printf("Failed to marshal: %v\n", err)
	} else {
		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.AdminsExchange,
			KeyName:      exports.AdminRegisteredByAdminSendWelcomeEmailRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.AdminsExchange,
			KeyName:      exports.AdminRegisteredByAdminRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

	}

	return &entity.Admin{
		AccountId:       accAdmin.Id,
		HackathonId:     accAdmin.HackathonId,
		FirstName:       accAdmin.FirstName,
		LastName:        accAdmin.LastName,
		Gender:          accAdmin.Gender,
		IsEmailVerified: accAdmin.IsEmailVerified,
		Email:           accAdmin.Email,
		EmailVerifiedAt: accAdmin.EmailVerifiedAt,
		Status:          accAdmin.Status,
		State:           accAdmin.Status,
		CreatedAt:       accAdmin.CreatedAt,
		UpdatedAt:       accAdmin.UpdatedAt,
	}, nil
}

func (s *Service) GetAdminInfo(email string) (*entity.Admin, error) {
	accAdmin, err := s.AdminAccountRepository.GetAdminAccountByEmail(email)

	if err != nil {
		return nil, err
	}
	return &entity.Admin{
		AccountId:       accAdmin.Id,
		HackathonId:     accAdmin.HackathonId,
		FirstName:       accAdmin.FirstName,
		LastName:        accAdmin.LastName,
		Gender:          accAdmin.Gender,
		IsEmailVerified: accAdmin.IsEmailVerified,
		Email:           accAdmin.Email,
		EmailVerifiedAt: accAdmin.EmailVerifiedAt,
		Status:          accAdmin.Status,
		State:           accAdmin.Status,
		CreatedAt:       accAdmin.CreatedAt,
		UpdatedAt:       accAdmin.UpdatedAt,
	}, nil
}
