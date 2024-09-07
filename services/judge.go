package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/rabbitmq/amqp091-go"
)

type UploadJudgeProfilePictureOpt struct {
	PictureFile *multipart.FileHeader
	Email       string
}

func (s *Service) UploadJudgeProfile(opts UploadJudgeProfilePictureOpt) error {
	profPic := opts.PictureFile
	email := opts.Email
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	opt := utils.UploadOpts{
		Folder:         filepath.Join(dir, "uploads"),
		FileNamePrefix: fmt.Sprintf("%s_", email),
	}
	filePath, err := utils.SaveFile(profPic, []utils.UploadOpts{opt}...)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	payload := exports.UploadJudgeProfilePicQueuePayload{
		Email:    email,
		FilePath: filePath,
	}
	byt, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		err = Publish(&PublishOpts{
			ExchangeName: exports.UploadJobsExchange,
			RMQConn:      s.AppResources.RabbitMQConn,
			KeyName:      exports.UploadJudgeProfilePicRoutingKeyName,
			Data:         byt,
			ExchangeKind: amqp091.ExchangeDirect,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	fmt.Println("Queue payload published")
	return nil
}

func (s *Service) RegisterNewJudge(dataInput *RegisterNewJudgeDTO) (*entity.Judge, error) {
	err := validate.Struct(dataInput)
	if err != nil {
		return nil, err
	}
	dataToSave := &exports.RegisterNewJudgeDTO{
		FirstName: dataInput.FirstName,
		LastName:  dataInput.LastName,
		Email:     dataInput.Email,
		Gender:    dataInput.Gender,
		Password:  dataInput.Password,
		Bio:       dataInput.Bio,
		State:     dataInput.State,
	}
	created, err := s.JudgeAccountRepository.CreateJudgeAccount(dataToSave)
	if err != nil {
		return nil, err
	}

	pubData := &exports.JudgeRegisteredPublishPayload{
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
			ExchangeName: exports.JudgesExchange,
			KeyName:      exports.JudgeRegisteredRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.JudgesExchange,
			KeyName:      exports.JudgeSendWelcomeEmailRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)

		}

	}

	judge := entity.Judge{
		Id:                created.Id,
		HackathonId:       created.HackathonId,
		LastName:          created.LastName,
		FirstName:         created.FirstName,
		Email:             created.Email,
		Gender:            created.Gender,
		PhoneNumber:       created.PhoneNumber,
		ProfilePictureUrl: created.ProfilePictureUrl,
		Role:              created.Role,
		Bio:               created.Bio,
		Status:            created.Status,
		CreatedAt:         created.CreatedAt,
		UpdatedAt:         created.UpdatedAt,
	}

	return &judge, nil
}

func (s *Service) RegisterNewJudgeByAdmin(dataInput *RegisterNewJudgeDTO) (*entity.Judge, error) {
	err := validate.Struct(dataInput)
	if err != nil {
		return nil, err
	}
	dataToSave := &exports.RegisterNewJudgeDTO{
		FirstName: dataInput.FirstName,
		LastName:  dataInput.LastName,
		Email:     dataInput.Email,
		Gender:    dataInput.Gender,
		Password:  dataInput.Password,
		Bio:       dataInput.Bio,
		State:     dataInput.State,
	}
	created, err := s.JudgeAccountRepository.CreateJudgeAccount(dataToSave)
	if err != nil {
		return nil, err
	}

	pubData := &exports.JudgeRegisteredByAdminPublishPayload{
		Email:       dataInput.Email,
		Name:        dataInput.FirstName,
		Password:    dataInput.Password,
		InviterName: dataInput.InviterName,
	}
	by, err := json.Marshal(pubData)
	if err != nil {
		fmt.Printf("Failed to marshal: %v\n", err)
	} else {
		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.JudgesExchange,
			KeyName:      exports.JudgeRegisteredByAdminRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

		err = Publish(&PublishOpts{
			Data:         by,
			RMQConn:      s.AppResources.RabbitMQConn,
			ExchangeName: exports.ParticipantsExchange,
			KeyName:      exports.JudgeRegisteredByAdminSendWelcomeEmailRoutingKeyName,
			ExchangeKind: amqp091.ExchangeTopic,
		})
		if err != nil {
			fmt.Println(err)
		}

	}

	judge := entity.Judge{
		Id:                created.Id,
		HackathonId:       created.HackathonId,
		LastName:          created.LastName,
		FirstName:         created.FirstName,
		Email:             created.Email,
		Gender:            created.Gender,
		PhoneNumber:       created.PhoneNumber,
		ProfilePictureUrl: created.ProfilePictureUrl,
		Role:              created.Role,
		Bio:               created.Bio,
		Status:            created.Status,
		CreatedAt:         created.CreatedAt,
		UpdatedAt:         created.UpdatedAt,
	}

	return &judge, nil
}

func (s *Service) UpdateJudgeInfo(email string, dataInput *UpdateJudgeDTO) (*entity.Judge, error) {
	validate = validator.New()
	err := validate.Struct(dataInput)
	if err != nil {
		return nil, err
	}
	dataToSave := &exports.UpdateJudgeDTO{
		FirstName: dataInput.FirstName,
		LastName:  dataInput.LastName,
		Gender:    dataInput.Gender,
		Bio:       dataInput.Bio,
		State:     dataInput.State,
	}
	err = s.JudgeAccountRepository.UpdateJudgeAccount(email, dataToSave)
	if err != nil {
		return nil, err
	}

	judgeAccRepo, err := s.JudgeAccountRepository.GetJudgeAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	judge := entity.Judge{
		Id:                judgeAccRepo.Id,
		HackathonId:       judgeAccRepo.HackathonId,
		LastName:          judgeAccRepo.LastName,
		FirstName:         judgeAccRepo.FirstName,
		Email:             judgeAccRepo.Email,
		Gender:            judgeAccRepo.Gender,
		PhoneNumber:       judgeAccRepo.PhoneNumber,
		ProfilePictureUrl: judgeAccRepo.ProfilePictureUrl,
		Role:              judgeAccRepo.Role,
		Bio:               judgeAccRepo.Bio,
		Status:            judgeAccRepo.Status,
		CreatedAt:         judgeAccRepo.CreatedAt,
		UpdatedAt:         judgeAccRepo.UpdatedAt,
	}
	return &judge, nil
}

func (s *Service) GetJudgeByEmail(email string) (*entity.Judge, error) {
	if email == "" {
		return nil, errors.New("email must not be empty")
	}
	fetched, err := s.JudgeAccountRepository.GetJudgeAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	judge := entity.Judge{
		Id:                fetched.Id,
		HackathonId:       fetched.HackathonId,
		LastName:          fetched.LastName,
		FirstName:         fetched.FirstName,
		Email:             fetched.Email,
		Gender:            fetched.Gender,
		PhoneNumber:       fetched.PhoneNumber,
		ProfilePictureUrl: fetched.ProfilePictureUrl,
		Role:              fetched.Role,
		Status:            fetched.Status,
		Bio:               fetched.Bio,
		IsEmailVerified:   fetched.IsEmailVerified,
		CreatedAt:         fetched.CreatedAt,
		UpdatedAt:         fetched.UpdatedAt,
	}
	return &judge, nil
}

func (s *Service) GetJudges() ([]*entity.Judge, error) {
	repoAccounts, err := s.JudgeAccountRepository.GetJudges()
	if err != nil {
		return nil, err
	}
	var ents []*entity.Judge
	for _, re := range repoAccounts {
		ents = append(ents, &entity.Judge{
			Id:                re.Id,
			LastName:          re.LastName,
			FirstName:         re.FirstName,
			Email:             re.Email,
			Bio:               re.Bio,
			Gender:            re.Gender,
			State:             re.State,
			Status:            re.Status,
			PhoneNumber:       re.PhoneNumber,
			ProfilePictureUrl: re.ProfilePictureUrl,
			HackathonId:       re.HackathonId,
			IsEmailVerified:   re.IsEmailVerified,
			CreatedAt:         re.CreatedAt,
			UpdatedAt:         re.UpdatedAt,
		})
	}
	return ents, nil
}
