package authutils

import (
	"errors"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/repository"
	"github.com/arravoco/hackathon_backend/resources"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jaevor/go-nanoid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthUtils struct {
	AccountRepository            *repository.AccountRepository
	JudgeAccountRepository       *repository.JudgeAccountRepository
	TokenRepository              *repository.TokenDataRepository
	ParticipantAccountRepository *repository.ParticipantAccountRepository
	ParticipantRecordRepository  *repository.ParticipantRecordRepository
}

func GetAuthUtilsWithDefaultRepositories() *AuthUtils {
	res := resources.GetDefaultResources()
	var dataSourceInstance exports.DBInterface = data.GetDatasourceWithMongoDBInstance(res.Mongo)
	q := query.GetQueryWithConfiguredDatasource(dataSourceInstance)

	var accRepoInstance *repository.AccountRepository = repository.NewAccountRepository(q)

	var judgeRepoInstance *repository.JudgeAccountRepository = repository.NewJudgeAccountRepository(q)

	var tokenRepoInstance *repository.TokenDataRepository = repository.NewTokenDataRepository(q)

	var partAccRepoInstance *repository.ParticipantAccountRepository = repository.NewParticipantAccountRepository(q)

	var partRecordRepoInstance *repository.ParticipantRecordRepository = repository.NewParticipantRecordRepository(q)

	return &AuthUtils{
		AccountRepository:            accRepoInstance,
		ParticipantAccountRepository: partAccRepoInstance,
		TokenRepository:              tokenRepoInstance,
		JudgeAccountRepository:       judgeRepoInstance,
		ParticipantRecordRepository:  partRecordRepoInstance,
	}
}

func (auth *AuthUtils) BasicLogin(dataInput *exports.AuthUtilsBasicLoginData) (*exports.AuthUtilsBasicLoginSuccessData, error) {

	accountDoc, err := auth.AccountRepository.FindAccountIdentifier(dataInput.Identifier)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, errors.New("no email or username provided that matches record")
	}
	if accountDoc == nil {
		return nil, errors.New("no email or username provided that matches record")
	}
	_, err = exports.ComparePasswordAndHash(dataInput.Password, accountDoc.PasswordHash)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	var rr *exports.AuthUtilsPayload
	var participantAndAccDoc *exports.ParticipantTeamMembersWithAccountsAggregate
	if accountDoc.Role == "PARTICIPANT" {
		participantAndAccDoc, err = auth.ParticipantRecordRepository.GetSingleParticipantRecordAndMemberAccountsInfo(accountDoc.ParticipantId)
		if err != nil {
			return nil, err
		}
		if participantAndAccDoc == nil {
			return nil, errors.New("participant details not found")
		}
		fmt.Println(participantAndAccDoc)
		rr = &exports.AuthUtilsPayload{
			Email:           accountDoc.Email,
			LastName:        accountDoc.LastName,
			FirstName:       accountDoc.FirstName,
			Role:            accountDoc.Role,
			HackathonId:     accountDoc.HackathonId,
			IsParticipant:   true,
			ParticipantType: participantAndAccDoc.Type,
			ParticipantId:   participantAndAccDoc.ParticipantId,
		}
		accessToken, err := GenerateAccessToken(rr)
		if err != nil {
			return nil, err
		}
		team_role := ""
		if participantAndAccDoc.TeamLeadInfo.Email == accountDoc.Email {
			team_role = "TEAM_LEAD"
		} else {
			for _, v := range participantAndAccDoc.CoParticipants {
				if v.Email == accountDoc.Email {
					team_role = v.TeamRole
				}
			}
		}
		var co_parts []exports.AuthUtilsParticipantCoParticipantInfo
		for _, v := range participantAndAccDoc.CoParticipants {
			co_parts = append(co_parts, exports.AuthUtilsParticipantCoParticipantInfo{
				FirstName:           v.FirstName,
				LastName:            v.LastName,
				Gender:              v.Gender,
				YearsOfExperience:   v.YearsOfExperience,
				AccountRole:         v.AccountRole,
				ParticipantId:       v.ParticipantId,
				TeamRole:            v.TeamRole,
				AccountId:           v.AccountId,
				AccountStatus:       v.AccountStatus,
				PhoneNumber:         v.PhoneNumber,
				PreviousProjects:    v.PreviousProjects,
				DOB:                 v.DOB,
				State:               v.State,
				Skillset:            v.Skillset,
				EmploymentStatus:    v.EmploymentStatus,
				Email:               v.Email,
				ExperienceLevel:     v.ExperienceLevel,
				HackathonExperience: v.HackathonExperience,
				CreatedAt:           v.CreatedAt,
				UpdatedAt:           v.UpdateAt,
			})
		}
		sol := exports.AuthUtilsParticipantSolutionInfo{
			Title:            participantAndAccDoc.Solution.Title,
			Description:      participantAndAccDoc.Solution.Description,
			SolutionImageUrl: participantAndAccDoc.Solution.SolutionImageUrl,
			Objective:        participantAndAccDoc.Solution.Objective,
		}
		var invite_list []exports.ParticipantDocumentTeamInviteInfo
		for _, v := range participantAndAccDoc.InviteList {
			invite_list = append(invite_list, exports.ParticipantDocumentTeamInviteInfo{
				Email:     v.Email,
				InviterId: v.Email,
				Time:      v.Time,
			})
		}
		return &exports.AuthUtilsBasicLoginSuccessData{
			AccessToken:         accessToken,
			FirstName:           accountDoc.FirstName,
			LastName:            accountDoc.LastName,
			Email:               accountDoc.Email,
			ParticipantId:       accountDoc.ParticipantId,
			Gender:              accountDoc.Gender,
			TeamRole:            team_role,
			State:               accountDoc.State,
			AccountStatus:       accountDoc.Status,
			ParticipantStatus:   participantAndAccDoc.Status,
			HackathonId:         accountDoc.HackathonId,
			HackathonExperience: accountDoc.HackathonExperience,
			ExperienceLevel:     accountDoc.ExperienceLevel,
			YearsOfExperience:   accountDoc.YearsOfExperience,
			CoParticipants:      co_parts,
			AccountRole:         accountDoc.Role,
			PhoneNumber:         accountDoc.PhoneNumber,
			DOB:                 accountDoc.DOB,
			TeamLeadEmail:       participantAndAccDoc.TeamLeadEmail,
			TeamName:            participantAndAccDoc.TeamName,
			Solution:            &sol,
			Skillset:            accountDoc.Skillset,
			EmploymentStatus:    accountDoc.EmploymentStatus,
			ParticipantEmail:    participantAndAccDoc.ParticipantEmail,
			FieldOfStudy:        accountDoc.FieldOfStudy,
			InviteList:          invite_list,
			Motivation:          accountDoc.Motivation,
			PreviousProjects:    accountDoc.PreviousProjects,
			CreatedAt:           accountDoc.CreatedAt,
			UpdatedAt:           accountDoc.UpdatedAt,
		}, nil
	} else if accountDoc.Role == "JUDGE" {

		rr = &exports.AuthUtilsPayload{
			Email:           accountDoc.Email,
			LastName:        accountDoc.LastName,
			FirstName:       accountDoc.FirstName,
			Role:            accountDoc.Role,
			HackathonId:     accountDoc.HackathonId,
			IsParticipant:   true,
			ParticipantType: accountDoc.ParticipantId,
			ParticipantId:   participantAndAccDoc.Type,
		}
		accessToken, err := GenerateAccessToken(rr)
		if err != nil {
			return nil, err
		}
		return &exports.AuthUtilsBasicLoginSuccessData{
			AccessToken:       accessToken,
			FirstName:         accountDoc.FirstName,
			LastName:          accountDoc.LastName,
			Email:             accountDoc.Email,
			Bio:               accountDoc.Bio,
			AccountRole:       accountDoc.Role,
			HackathonId:       accountDoc.HackathonId,
			Gender:            accountDoc.Gender,
			State:             accountDoc.State,
			PhoneNumber:       accountDoc.PhoneNumber,
			ProfilePictureUrl: accountDoc.ProfilePictureUrl,
			CreatedAt:         accountDoc.CreatedAt,
			UpdatedAt:         accountDoc.UpdatedAt,
		}, nil
	} else if accountDoc.Role == "ADMIN" {

		rr := &exports.AuthUtilsPayload{
			Email:       accountDoc.Email,
			LastName:    accountDoc.LastName,
			FirstName:   accountDoc.FirstName,
			Role:        accountDoc.Role,
			HackathonId: accountDoc.HackathonId,
		}
		accessToken, err := GenerateAccessToken(rr)
		if err != nil {
			return nil, err
		}
		return &exports.AuthUtilsBasicLoginSuccessData{
			AccessToken:       accessToken,
			FirstName:         accountDoc.FirstName,
			LastName:          accountDoc.LastName,
			Email:             accountDoc.Email,
			Bio:               accountDoc.Bio,
			AccountRole:       accountDoc.Role,
			HackathonId:       accountDoc.HackathonId,
			Gender:            accountDoc.Gender,
			State:             accountDoc.State,
			PhoneNumber:       accountDoc.PhoneNumber,
			ProfilePictureUrl: accountDoc.ProfilePictureUrl,
			CreatedAt:         accountDoc.CreatedAt,
			UpdatedAt:         accountDoc.UpdatedAt,
		}, nil
	}

	return nil, fmt.Errorf("no role for account email %s", accountDoc.Email)
}

func GenerateAccessToken(payload *exports.AuthUtilsPayload) (string, error) {
	claims := &exports.MyJWTCustomClaims{
		Email:           payload.Email,
		FirstName:       payload.FirstName,
		LastName:        payload.LastName,
		Role:            payload.Role,
		ParticipantId:   payload.ParticipantId,
		ParticipantType: payload.ParticipantType,
		IsParticipant:   payload.IsParticipant,
		HackathonId:     payload.HackathonId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.GetSecretKey()))
	return t, err
}

func GetAuthPayload(c echo.Context) *exports.Payload {

	jwtData := c.Get("user").(*jwt.Token)
	claims := jwtData.Claims.(*exports.MyJWTCustomClaims)
	tokenData := exports.Payload{
		Email:           claims.Email,
		LastName:        claims.LastName,
		FirstName:       claims.FirstName,
		Role:            claims.Role,
		ParticipantType: claims.ParticipantType,
		IsParticipant:   claims.IsParticipant,
		ParticipantId:   claims.ParticipantId,
		HackathonId:     claims.HackathonId,
	}
	return &tokenData
}

func GetJWTConfig() echojwt.Config {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(exports.MyJWTCustomClaims)
		},
		SigningKey: []byte(config.GetSecretKey()),
	}
	return config
}

func (auth *AuthUtils) VerifyToken(dataInput *exports.AuthUtilsVerifyTokenData) error {
	_, err := auth.AccountRepository.GetAccountByEmail(dataInput.TokenTypeValue)
	if err != nil {
		return errors.New("email not found in record")
	}
	err = auth.TokenRepository.VerifyToken(&exports.VerifyTokenData{
		Token:          dataInput.Token,
		TokenType:      dataInput.TokenType,
		TokenTypeValue: dataInput.TokenTypeValue,
		Scope:          dataInput.Scope,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	return nil
}

func (auth *AuthUtils) GenerateToken(dataInput *GenerateTokenArgs) (*TokenData, error) {
	data, err := auth.TokenRepository.UpsertToken(&exports.UpsertTokenData{
		Token:          dataInput.Token,
		TokenType:      dataInput.TokenType,
		TokenTypeValue: dataInput.TokenTypeValue,
		TTL:            dataInput.TTL,
		Scope:          dataInput.Scope,
		Status:         dataInput.Status,
	})
	if err != nil {
		return nil, err
	}
	return &TokenData{
		Id:             data.Id,
		Token:          data.Token,
		TokenType:      data.TokenType,
		TokenTypeValue: data.TokenType,
		TTL:            data.TTL,
		Scope:          data.Scope,
		Status:         data.Status,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}, nil
}

func (auth *AuthUtils) InitiateEmailVerification(dataInput *exports.AuthUtilsConfigTokenData) (*exports.TokenDataRepository, error) {
	_, err := auth.AccountRepository.GetAccountByEmail(dataInput.Email)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return nil, errors.New("email not found in record")
	}
	tokenFunc, _ := nanoid.Custom("1234567890", 6)
	token := tokenFunc()
	tokenData, err := auth.TokenRepository.UpsertToken(&exports.UpsertTokenData{
		Token:          token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		TTL:            dataInput.TTL,
		Status:         "PENDING",
		Scope:          "EMAIL_VERIFICATION",
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return nil, errors.New("failed to generate token: ")
	}
	return tokenData, nil
}

func (auth *AuthUtils) CompleteEmailVerification(dataInput *exports.AuthUtilsCompleteEmailVerificationData) error {
	err := auth.VerifyToken(&exports.AuthUtilsVerifyTokenData{
		Email:          dataInput.Email,
		Token:          dataInput.Token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		Scope:          "EMAIL_VERIFICATION",
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	is_email_verified := true
	err = auth.AccountRepository.UpdateAccount(&exports.UpdateAccountFilter{
		Email: dataInput.Email,
	}, &exports.UpdateAccountDTO{
		IsEmailVerified:   is_email_verified,
		IsEmailVerifiedAt: time.Now(),
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return errors.New("unable to complete token verification")
	}

	return nil
}

func (auth *AuthUtils) ChangePassword(dataInput *exports.AuthUtilsChangePasswordData) error {
	accountDoc, err := auth.AccountRepository.GetAccountByEmail(dataInput.Email)
	if err != nil {
		return errors.New("user info not found in record")
	}
	_, err = exports.ComparePasswordAndHash(dataInput.OldPassword, accountDoc.PasswordHash)
	if err != nil {
		return err
	}
	hash, err := exports.GenerateHashPassword(dataInput.NewPassword)
	if err != nil {
		return err
	}
	accountDoc, err = auth.AccountRepository.UpdatePasswordByEmail(&exports.UpdateAccountFilter{Email: dataInput.Email}, hash)

	// emit an emit here
	return nil
}

func (auth *AuthUtils) InitiatePasswordRecovery(dataInput *exports.AuthUtilsConfigTokenData) (*exports.TokenData, error) {
	_, err := auth.AccountRepository.GetAccountByEmail(dataInput.Email)
	if err != nil {
		return nil, errors.New("email not found in record")
	}
	tokenFunc, _ := nanoid.Custom("1234567890", 6)
	token := tokenFunc()
	tokenData, err := data.UpsertToken(&exports.UpsertTokenData{
		Token:          token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		TTL:            dataInput.TTL,
		Status:         "PENDING",
		Scope:          "PASSWORD_RECOVERY",
	})
	return tokenData, err
}

func (auth *AuthUtils) CompletePasswordRecovery(dataInput *exports.AuthUtilsCompletePasswordRecoveryData) (interface{}, error) {
	err := auth.VerifyToken(&exports.AuthUtilsVerifyTokenData{
		Token:          dataInput.Token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		Scope:          "PASSWORD_RECOVERY",
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return nil, err
	}
	newPasswordHash, err := exports.GenerateHashPassword(dataInput.NewPassword)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return nil, err
	}
	_, err = data.UpdatePasswordByEmail(&exports.UpdateAccountFilter{Email: dataInput.Email}, newPasswordHash)

	if err != nil {
		exports.MySugarLogger.Error(err)
		return nil, errors.New("unable to complete password recovery")
	}

	return nil, nil
}
