package routes_v1

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/authutils"
	"github.com/arravoco/hackathon_backend/utils/email"
	"github.com/labstack/echo/v4"
)

type BasicLoginFailureResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BasicLoginSuccessResponse struct {
	Code    int                           `json:"code"`
	Message string                        `json:"message"`
	Data    BasicLoginSuccessResponseData `json:"data"`
}

type BasicLoginSuccessResponseData struct {
	AccessToken string `json:"access_token"`
}

type BasicLoginDTO struct {
	Identifier string ` validate:"required" json:"identifier"`
	Password   string ` validate:"required" json:"password"`
}

type InitiateEmailVerificationFailureResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type InitiateEmailVerificationSuccessResponse struct {
	Code    int                                           `json:"code"`
	Message string                                        `json:"message"`
	Data    *InitiateEmailVerificationSuccessResponseData `json:"data"`
}

type InitiateEmailVerificationSuccessResponseData struct {
	Email string `json:"email"`
}

type CompleteEmailVerificationFailureResponse struct {
	ResponseData
}
type CompleteEmailVerificationSuccessResponse struct {
	ResponseData
	Data *CompleteEmailVerificationSuccessResponseData `json:"data"`
}

type CompleteEmailVerificationSuccessResponseData struct {
	Email string `json:"email"`
}

type PasswordChangeFailureResponse struct {
	ResponseData
}
type PasswordChangeSuccessResponse struct {
	ResponseData
	Data *PasswordChangeSuccessResponseData `json:"data"`
}

type PasswordChangeSuccessResponseData struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Gender        string `json:"gender"`
	ParticipantId string `json:"participant_id"`
}

type AuthUserInfoFetchFailureResponse struct {
	ResponseData
}

type AuthUserInfoFetchSuccessResponse struct {
	ResponseData
	Data *AuthUserInfoFetchSuccessResponseData `json:"data"`
}

type AuthUserInfoFetchSuccessResponseData struct {
	*entity.Participant
}

type CompletePasswordRecoveryFailureResponse struct {
	ResponseData
}
type CompletePasswordRecoverySuccessResponse struct {
	ResponseData
	Data *CompleteEmailVerificationSuccessResponseData `json:"data"`
}

type CompletePasswordRecoverySuccessResponseData struct {
	Email string `json:"email"`
}

// @Title Basic Log in
// @Description	Log a user in
// @Summary		Log a user in
// @Tags		Auth
// @Produce		json
// @Param       loginJSON   body BasicLoginDTO    true                   "login Request JSON"
// @Success		200	  object 	BasicLoginSuccessResponse "Users Responses JSON"
// @Failure		400	object	BasicLoginFailureResponse "Login failed"
// @Router			/api/v1/auth/login             [post]
func BasicLogin(c echo.Context) error {
	data := dtos.BasicLoginDTO{}
	c.Bind(&data)
	err := validate.Struct(&data)
	if err != nil {
		utils.MySugarLogger.Error(err)
		return c.JSON(http.StatusBadRequest, &BasicLoginFailureResponse{

			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	dataResponse, err := authutils.BasicLogin(&exports.AuthUtilsBasicLoginData{
		Identifier: data.Identifier,
		Password:   data.Password,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &BasicLoginFailureResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	fmt.Printf("%#v", dataResponse)
	return c.JSON(200, &BasicLoginSuccessResponse{
		Code:    200,
		Message: "Successfully logged in",
		Data: BasicLoginSuccessResponseData{
			AccessToken: dataResponse.AccessToken,
		},
	})
}

// @Title Generate Email Verification Link
// @Summary		Generate token to verify user email address
// @Description	Generate token to verify user email address
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			email	query		string	true	"Email to verify"	Format(email)
// @Success		200		object	InitiateEmailVerificationSuccessResponse "Verification succeeded"
// @Failure		400		object	InitiateEmailVerificationFailureResponse "Verification failed"
// @Failure		404		object	InitiateEmailVerificationFailureResponse "Verification failed"
// @Failure		500		object	InitiateEmailVerificationFailureResponse "Verification failed"
// @Router			/api/v1/auth/verification/email/initiation [get]
func InitiateEmailVerification(c echo.Context) error {
	emailToVerify := c.QueryParam("email")
	if emailToVerify == "" {
		return c.JSON(http.StatusBadRequest, &InitiateEmailVerificationFailureResponse{
			Code:    http.StatusBadRequest,
			Message: "'email' query parameter is required",
		})
	}

	ttl := time.Now().Add(time.Minute * 15)
	tokenData, err := authutils.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
		TTL:   ttl,
		Email: emailToVerify,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &InitiateEmailVerificationFailureResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	link, err := utils.GenerateEmailVerificationLink(&exports.EmailVerificationLinkPayload{
		Token: tokenData.Token,
		TTL:   tokenData.TTL,
		Email: tokenData.TokenTypeValue,
	})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		email.SendEmailVerificationEmail(&email.SendEmailVerificationEmailData{
			Email:    tokenData.TokenTypeValue,
			Token:    tokenData.Token,
			TokenTTL: tokenData.TTL,
			Subject:  "Email Verification",
			Link:     link,
		})
	}
	return c.JSON(200, &InitiateEmailVerificationSuccessResponse{
		Code:    200,
		Message: "Verification email sent successfully",
		Data:    &InitiateEmailVerificationSuccessResponseData{},
	})
}

// @Title Verify Email Via Link
// @Summary		Verify user email address
// @Description	Verify user email address
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Success		200				object	CompleteEmailVerificationSuccessResponse "Email verification successful"
// @Failure		400				object	CompleteEmailVerificationFailureResponse "Email verification failed"
// @Failure		404				object	CompleteEmailVerificationFailureResponse "Email verification failed"
// @Failure		500				object	CompleteEmailVerificationFailureResponse "Email verification failed"
// @Router			/api/v1/auth/verification/email/completion [get]
func CompleteEmailVerificationViaGet(c echo.Context) error {
	queryToken := c.QueryParam("token")
	payload, err := utils.ProcessEmailVerificationLink(queryToken)
	if err != nil {
		return c.Redirect(302, strings.Join([]string{config.GetFrontendURL(),
			strings.Join([]string{"verify_fail", strings.Join([]string{"err", err.Error()}, "=")}, "?")}, "/"))
	}
	fmt.Println(payload)
	if payload.TTL.Before(time.Now()) {
		return c.JSON(http.StatusBadRequest, &CompleteEmailVerificationFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: "Link has expired",
			},
		})
	}

	err = authutils.CompleteEmailVerification(&exports.AuthUtilsCompleteEmailVerificationData{
		Token: payload.Token,
		Email: payload.Email,
	})
	if err != nil {
		return c.Redirect(302, strings.Join([]string{config.GetFrontendURL(),
			strings.Join([]string{"verify_fail", strings.Join([]string{"err", err.Error()}, "=")}, "?")}, "/"))
	}
	email.SendEmailVerificationCompleteEmail(&email.SendEmailVerificationCompleteEmailData{
		Email:   payload.Email,
		Subject: "Email Verification Success",
	})
	return c.Redirect(302, strings.Join([]string{config.GetFrontendURL(), "verify"}, "/"))
}

// @Title Verify Email Via Token
// @Summary		Verify user email address
// @Description	Verify user email address
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			completeToken	body dtos.CompleteEmailVerificationDTO	true	"the required info"
// @Success		200				object	CompleteEmailVerificationSuccessResponse "Email verification successful"
// @Failure		400				object	CompleteEmailVerificationFailureResponse "Email verification failed"
// @Failure		404				object	CompleteEmailVerificationFailureResponse "Email verification failed"
// @Failure		500				object	CompleteEmailVerificationFailureResponse "Email verification failed"
// @Router			/api/v1/auth/verification/email/completion [post]
func CompleteEmailVerification(c echo.Context) error {
	dataDto := dtos.CompleteEmailVerificationDTO{}
	err := c.Bind(&dataDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &CompleteEmailVerificationFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	err = validate.Struct(dataDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &CompleteEmailVerificationFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}

	err = authutils.CompleteEmailVerification(&exports.AuthUtilsCompleteEmailVerificationData{
		Token: dataDto.Token,
		Email: dataDto.Email,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &CompleteEmailVerificationFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	email.SendEmailVerificationCompleteEmail(&email.SendEmailVerificationCompleteEmailData{
		Email:   dataDto.Email,
		Subject: "Email Verification Success",
	})
	return c.JSON(200, &CompleteEmailVerificationSuccessResponse{
		ResponseData{
			Code:    200,
			Message: "Verification email completed successfully",
		}, &CompleteEmailVerificationSuccessResponseData{
			Email: dataDto.Email,
		},
	})
}

// @Title Change Password
// @Summary		Change User Password
// @Description	Change User Password
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			changePasswordJSON	body		dtos.ChangePasswordDTO	true	"the required info"
// @Success		200				object	PasswordChangeSuccessResponse ""
// @Failure		400					object	PasswordChangeFailureResponse ""
// @Failure		404					object	PasswordChangeFailureResponse ""
// @Failure		500					object	PasswordChangeFailureResponse ""
// @Router			/api/v1/auth/password/change [post]
func ChangePassword(c echo.Context) error {
	tokenData := authutils.GetAuthPayload(c)
	dataDto := dtos.ChangePasswordDTO{}
	err := c.Bind(&dataDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &PasswordChangeFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	err = validate.Struct(dataDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &PasswordChangeFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	fmt.Println(tokenData)
	acc, err := entity.ChangePassword(&entity.PasswordChangeData{
		Email:       tokenData.Email,
		OldPassword: dataDto.OldPassword,
		NewPassword: dataDto.NewPassword,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &PasswordChangeFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	return c.JSON(200, &PasswordChangeSuccessResponse{
		ResponseData{
			Code:    200,
			Message: "Password change completed successfully",
		}, &PasswordChangeSuccessResponseData{FirstName: acc.FirstName},
	})
}

// @Title  Get Auth User Information
// @Summary		Get auth user information
// @Description	Get auth user information
// @Tags			Auth
// @in				header
// @name			Authorization
// @description	"Type 'Bearer TOKEN' to correctly set the Bearer token"
// @Produce		json
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token
// @Success		200	object	AuthUserInfoFetchSuccessResponse "UserInfo fetched successfully"
// @Failure		400	object	AuthUserInfoFetchFailureResponse "UserInfo fetch failed"
// @Failure		404	object	AuthUserInfoFetchFailureResponse "UserInfo fetch failed"
// @Failure		500	object	AuthUserInfoFetchFailureResponse "UserInfo fetch failed"
// @Router			/api/v1/auth/me [get]
func GetAuthUserInfo(c echo.Context) error {
	tokenData := authutils.GetAuthPayload(c)
	user := AuthUserInfoFetchSuccessResponseData{}
	fmt.Println(tokenData)
	var err error
	if tokenData.Role == "PARTICIPANT" {
		participant := entity.Participant{}
		err = participant.FillParticipantInfo(tokenData.Email)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &AuthUserInfoFetchFailureResponse{
				ResponseData{
					Code:    http.StatusBadRequest,
					Message: "Error getting user info",
				},
			})
		}
		fmt.Println(participant)
		user.Participant = &participant
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, &AuthUserInfoFetchFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	return c.JSON(200, &AuthUserInfoFetchSuccessResponse{
		ResponseData: ResponseData{
			Code:    200,
			Message: "Auth user info fetched successfully",
		}, Data: &user,
	})
}

// @Title Update Auth User Info
// @Summary		Update auth user information
// @Description	Update auth user information
// @Tags			Auth
// @Param			updateMyInfoJSON	body		dtos.AuthParticipantInfoUpdateDTO	true	"the required info"
// @Produce		json
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token
// @Success		200	{object}	AuthUserInfoFetchSuccessResponse
// @Failure		400	{object}	AuthUserInfoFetchFailureResponse
// @Failure		404	{object}	AuthUserInfoFetchFailureResponse
// @Failure		500	{object}	AuthUserInfoFetchFailureResponse
// @Router			/api/v1/auth/me [put]
func UpdateAuthUserInfo(c echo.Context) error {
	tokenData := c.Get("user").(exports.Payload)
	var err error
	if tokenData.Role == "PARTICIPANT" {
		updateData := dtos.AuthParticipantInfoUpdateDTO{}
		participant := entity.Participant{
			Email: tokenData.Email,
		}
		err = participant.UpdateParticipantInfo(&dtos.AuthParticipantInfoUpdateDTO{
			AuthUserInfoUpdateDTO: dtos.AuthUserInfoUpdateDTO{
				FirstName: updateData.FirstName,
				LastName:  updateData.LastName,
				Email:     updateData.Email,
				Gender:    updateData.Gender,
			},
			GithubAddress:   updateData.GithubAddress,
			LinkedInAddress: updateData.LinkedInAddress,
		})

	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, &PasswordChangeFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	return c.JSON(200, &PasswordChangeSuccessResponse{
		ResponseData{
			Code:    200,
			Message: "Password change completed successfully",
		}, &PasswordChangeSuccessResponseData{},
	})
}

// @Title Initiate Password Recovery
// @Summary			Initiate Password Recovery
// @Description	Initiate Password Recovery
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			email	query		string	false	"Email to verify"	Format(email)
// @Success		200		object	InitiateEmailVerificationSuccessResponse
// @Failure		400		object	InitiateEmailVerificationFailureResponse
// @Failure		404		object	InitiateEmailVerificationFailureResponse
// @Failure		500		object	InitiateEmailVerificationFailureResponse
// @Router		/api/v1/auth/password/recovery/initiation [get]
func InitiatePasswordRecovery(c echo.Context) error {
	emailToVerify := c.QueryParam("email")
	if emailToVerify == "" {
		return c.JSON(http.StatusBadRequest, &InitiateEmailVerificationFailureResponse{

			Code:    http.StatusBadRequest,
			Message: "'email' query parameter is required",
		})
	}

	ttl := time.Now().Add(time.Minute * 15)
	dataResult, err := authutils.InitiatePasswordRecovery(&exports.AuthUtilsConfigTokenData{
		TTL:   ttl,
		Email: emailToVerify,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &InitiateEmailVerificationFailureResponse{

			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	email.SendPasswordRecoveryEmail(&email.SendPasswordRecoveryEmailData{
		Email:    dataResult.TokenTypeValue,
		Token:    dataResult.Token,
		TokenTTL: dataResult.TTL,
		Subject:  "Password Recovery",
	})
	return c.JSON(200, &InitiateEmailVerificationSuccessResponse{

		Code:    200,
		Message: "Verification of email sent successfully",
		Data:    &InitiateEmailVerificationSuccessResponseData{},
	})
}

// @Title Complete Password Recovery
// @Summary		Complete  password recovery
// @Description	Complete password recovery
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			completeToken	body		dtos.CompletePasswordRecoveryDTO	true	"the required info"
// @Success		200				{object}	CompletePasswordRecoverySuccessResponse
// @Failure		400				{object}	CompletePasswordRecoveryFailureResponse
// @Failure		404				{object}	CompletePasswordRecoveryFailureResponse
// @Failure		500				{object}	CompletePasswordRecoveryFailureResponse
// @Router			/api/v1/auth/password/recovery/completion [post]
func CompletePasswordRecovery(c echo.Context) error {
	dataDto := dtos.CompletePasswordRecoveryDTO{}
	err := c.Bind(&dataDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &CompleteEmailVerificationFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	err = validate.Struct(dataDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &CompleteEmailVerificationFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}

	_, err = authutils.CompletePasswordRecovery(&exports.AuthUtilsCompletePasswordRecoveryData{
		Token:       dataDto.Token,
		Email:       dataDto.Email,
		NewPassword: dataDto.NewPassword,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &CompletePasswordRecoveryFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}

	email.SendPasswordRecoveryCompleteEmail(&email.SendPasswordRecoveryCompleteEmailData{
		Email:   "",
		Subject: "Password Recovery Success",
	})
	return c.JSON(200, &CompletePasswordRecoverySuccessResponse{
		ResponseData{
			Code:    200,
			Message: "Password recovery completed successfully",
		}, &CompleteEmailVerificationSuccessResponseData{
			Email: dataDto.Email,
		},
	})
}

// @Title Validate Team Invite Link
// @Summary		Validate team invite link
// @Description	Validate team invite link
// @Tags			Auth
// @Produce		json
// @Param			completeToken	path		dt	true	"the required info"
// @Success		200				{object}	string
// @Failure		400				{object}	string
// @Failure		404				{object}	string
// @Failure		500				{object}	string
// @Router			/api/v1/auth/team/invite [get]
func ValidateTeamInviteLink(c echo.Context) error {
	tokenStr := c.QueryParam("token")
	if tokenStr == "" {
		return c.HTML(400, "<p>Invalid link</p>")
	}
	t, err := utils.ProcessTeamInviteLink(tokenStr)
	if err != nil {
		fmt.Printf("")
		c.HTML(400, err.Error())
	}
	fmt.Println(t)
	return c.Redirect(301, "https://hackathon-dev.onrender.com/")
}
