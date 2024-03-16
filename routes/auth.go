package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/authutils"
	"github.com/labstack/echo/v4"
)

type BasicLoginFailureResponse struct {
	ResponseData
}

type BasicLoginSuccessResponse struct {
	ResponseData
	Data BasicLoginSuccessResponseData `json:"data"`
}

type BasicLoginSuccessResponseData struct {
	AccessToken string `json:"access_token"`
}

// @Description	Log a user in
// @Summary		Log a user in
// @Tags			Auth
// @Produce		json
// @Success		200	{object}	BasicLoginSuccessResponse
// @Failure		400	{object}	BasicLoginFailureResponse
// @Router			/api/auth/login             [post]
func BasicLogin(c echo.Context) error {
	data := dtos.BasicLoginDTO{}
	c.Bind(&data)
	err := validate.Struct(&data)
	if err != nil {
		utils.MySugarLogger.Error(err)
		return c.JSON(http.StatusBadRequest, &BasicLoginFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}

	dataResponse, err := authutils.BasicLogin(&authutils.BasicLoginData{
		Identifier: data.Identifier,
		Password:   data.Password,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &BasicLoginFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	fmt.Printf("%#v", dataResponse)
	return c.JSON(200, &BasicLoginSuccessResponse{
		ResponseData{
			Code:    200,
			Message: "Successfully logged in",
		},
		BasicLoginSuccessResponseData{
			AccessToken: dataResponse.AccessToken,
		},
	})
}

type InitiateEmailVerificationFailureResponse struct {
	ResponseData
}
type InitiateEmailVerificationSuccessResponse struct {
	ResponseData
	Data *InitiateEmailVerificationSuccessResponseData `json:"data"`
}

type InitiateEmailVerificationSuccessResponseData struct {
	Email string `json:"email"`
}

// @Summary		Verify user email address
// @Description	Verify user email address
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			email	query		string	false	"Email to verify"	Format(email)
// @Success		200		{object}	InitiateEmailVerificationSuccessResponse
// @Failure		400		{object}	InitiateEmailVerificationFailureResponse
// @Failure		404		{object}	InitiateEmailVerificationFailureResponse
// @Failure		500		{object}	InitiateEmailVerificationFailureResponse
// @Router			/api/auth/verification/email/initiation [get]
func InitiateEmailVerification(c echo.Context) error {
	emailToVerify := c.QueryParam("email")
	if emailToVerify == "" {
		return c.JSON(http.StatusBadRequest, &InitiateEmailVerificationFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: "'email' query parameter is required",
			},
		})
	}

	ttl := time.Now().Add(time.Minute * 15)
	err := authutils.InitiateEmailVerification(&authutils.ConfigTokenData{
		TTL:   ttl,
		Email: emailToVerify,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &InitiateEmailVerificationFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	return c.JSON(200, &InitiateEmailVerificationSuccessResponse{
		ResponseData{
			Code:    200,
			Message: "Verification email sent successfully",
		}, &InitiateEmailVerificationSuccessResponseData{},
	})
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

// @Summary		Verify user email address
// @Description	Verify user email address
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			completeToken	body		dtos.CompleteEmailVerificationDTO	true	"the required info"
// @Success		200				{object}	CompleteEmailVerificationSuccessResponse
// @Failure		400				{object}	CompleteEmailVerificationFailureResponse
// @Failure		404				{object}	CompleteEmailVerificationFailureResponse
// @Failure		500				{object}	CompleteEmailVerificationFailureResponse
// @Router			/api/auth/verification/email/completion [post]
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

	err = authutils.CompleteEmailVerification(&authutils.VerifyTokenData{
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
	return c.JSON(200, &CompleteEmailVerificationSuccessResponse{
		ResponseData{
			Code:    200,
			Message: "Verification email completed successfully",
		}, &CompleteEmailVerificationSuccessResponseData{
			Email: dataDto.Email,
		},
	})
}

type PasswordChangeFailureResponse struct {
	ResponseData
}
type PasswordChangeSuccessResponse struct {
	ResponseData
	Data *PasswordChangeSuccessResponseData `json:"data"`
}

type PasswordChangeSuccessResponseData struct {
}

// @Summary		Verify user email address
// @Description	Verify user email address
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			changePasswordJSON	body		dtos.ChangePasswordDTO	true	"the required info"
// @Success		200					{object}	PasswordChangeSuccessResponse
// @Failure		400					{object}	PasswordChangeFailureResponse
// @Failure		404					{object}	PasswordChangeFailureResponse
// @Failure		500					{object}	PasswordChangeFailureResponse
// @Router			/api/auth/password/change [post]
func ChangePassword(c echo.Context) error {
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
	return c.JSON(200, &PasswordChangeSuccessResponse{
		ResponseData{
			Code:    200,
			Message: "Password change completed successfully",
		}, &PasswordChangeSuccessResponseData{},
	})
}

type AuthUserInfoFetchFailureResponse struct {
	ResponseData
}

type AuthUserInfoFetchSuccessResponse struct {
	ResponseData
	Data *AuthUserInfoFetchSuccessResponseData `json:"data"`
}

type AuthUserInfoFetchSuccessResponseData struct {
	entity.Participant
	entity.Judge
}

// @Summary		Get auth user information
// @Description	Get auth user information
// @Tags			Auth
// @in				header
// @name			Authorization
// @description	"Type 'Bearer TOKEN' to correctly set the Bearer token"
// @Produce		json
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token
// @Success		200	{object}	AuthUserInfoFetchSuccessResponse
// @Failure		400	{object}	AuthUserInfoFetchFailureResponse
// @Failure		404	{object}	AuthUserInfoFetchFailureResponse
// @Failure		500	{object}	AuthUserInfoFetchFailureResponse
// @Router			/api/auth/me [get]
func GetAuthUserInfo(c echo.Context) error {
	tokenData := c.Get("user").(authutils.Payload)
	var err error
	if tokenData.Role == "PARTICIPANT" {
		participant := entity.Participant{}
		err = participant.GetParticipant(tokenData.Email)

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

// @Summary		Get auth user information
// @Description	Get auth user information
// @Tags			Auth
// @in				header
// @name			Authorization
// @description	"Type 'Bearer TOKEN' to correctly set the Bearer token"
// @Produce		json
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token
// @Success		200	{object}	AuthUserInfoFetchSuccessResponse
// @Failure		400	{object}	AuthUserInfoFetchFailureResponse
// @Failure		404	{object}	AuthUserInfoFetchFailureResponse
// @Failure		500	{object}	AuthUserInfoFetchFailureResponse
// @Router			/api/auth/me [put]
func UpdateAuthUserInfo(c echo.Context) error {
	tokenData := c.Get("user").(authutils.Payload)
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
