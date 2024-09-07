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
	"github.com/arravoco/hackathon_backend/repository"
	"github.com/arravoco/hackathon_backend/services"
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
	Code    int                                     `json:"code"`
	Message string                                  `json:"message"`
	Data    *exports.AuthUtilsBasicLoginSuccessData `json:"data"`
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
	Data interface{} `json:"data"`
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

type DeleteTeamMemberSuccessResponse struct {
	Code    int                                     `json:"code"`
	Message string                                  `json:"message"`
	Data    *entity.TeamMemberWithParticipantRecord `json:"data"`
}

type DeleteTeamMemberFailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AddSolutionSuccessResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    *entity.Participant `json:"data"`
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
	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	dataResponse, err := auth.BasicLogin(&exports.AuthUtilsBasicLoginData{
		Identifier: data.Identifier,
		Password:   data.Password,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &BasicLoginFailureResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return c.JSON(200, &BasicLoginSuccessResponse{
		Code:    200,
		Message: "Successfully logged in",
		Data:    dataResponse,
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
	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	tokenData, err := auth.InitiateEmailVerification(&exports.AuthUtilsConfigTokenData{
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
	if payload.TTL.Before(time.Now()) {
		return c.JSON(http.StatusBadRequest, &CompleteEmailVerificationFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: "Link has expired",
			},
		})
	}

	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	err = auth.CompleteEmailVerification(&exports.AuthUtilsCompleteEmailVerificationData{
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
	redirectUrl := payload.RedirectUrl
	fmt.Println(payload)
	if redirectUrl == "" {
		redirectUrl = strings.Join([]string{config.GetFrontendURL(), "verify"}, "/")
	}
	return c.Redirect(302, redirectUrl)
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

	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	err = auth.CompleteEmailVerification(&exports.AuthUtilsCompleteEmailVerificationData{
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
	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	err = auth.ChangePassword(&exports.AuthUtilsChangePasswordData{
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
		}, &PasswordChangeSuccessResponseData{},
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
	var user interface{}
	var err error
	if tokenData.Role == "PARTICIPANT" {
		serv := services.GetServiceWithDefaultRepositories()
		participant, err := serv.GetSingleParticipantWithAccountsInfo(tokenData.Email)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &AuthUserInfoFetchFailureResponse{
				ResponseData{
					Code:    http.StatusBadRequest,
					Message: "Error getting user info",
				},
			})
		}
		user = participant
	}

	if tokenData.Role == "ADMIN" {
		serv := services.GetServiceWithDefaultRepositories()
		admin, err := serv.GetAdminInfo(tokenData.Email)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &AuthUserInfoFetchFailureResponse{
				ResponseData{
					Code:    http.StatusBadRequest,
					Message: "Error getting user info",
				},
			})
		}
		user = &admin
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

/*
// @Title Admin-only Auth User Info
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
func UpdateAuthParticipantInfo(c echo.Context) error {
	tokenData := c.Get("user").(exports.Payload)
	var err error
	if tokenData.Role == "PARTICIPANT" {
		updateData := dtos.AdminParticipantInfoUpdateDTO{}
		serv := services.GetServiceWithDefaultRepositories()
		_, err = serv.AdminUpdateParticipantInfo(&services.UpdateSingleParticipantRecordFilter {

		}, &services.AdminParticipantInfoUpdateDTO{
			ReviewRanking: updateData.ReviewRanking,
			Status:        updateData.Status,
			//GithubAddress:   updateData.GithubAddress,
			//LinkedInAddress: updateData.LinkedInAddress,
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
}*/

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
	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	dataResult, err := auth.InitiatePasswordRecovery(&exports.AuthUtilsConfigTokenData{
		TTL:   ttl,
		Email: emailToVerify,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &InitiateEmailVerificationFailureResponse{

			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	payload, err := utils.GeneratePasswordRecoveryLinkPayload(&exports.PaswordRecoveryPayload{
		Email: dataResult.TokenTypeValue,
		Token: dataResult.Token,
		TTL:   dataResult.TTL,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		})
	}
	err = email.SendPasswordRecoveryEmail(&email.SendPasswordRecoveryEmailData{
		Email: dataResult.TokenTypeValue,
		Token: dataResult.Token,
		TTL:   uint32(dataResult.TTL.Sub((time.Now())).Minutes()),
		//Link:strings.Join([]string{config.GetFrontendURL(),"/password_reset_complete"},""),
		Link: strings.Join([]string{
			config.GetServerURL(),
			strings.Join(
				[]string{"api/v1/auth/password/recovery/link/verification",
					strings.Join([]string{"?token", payload}, "=")}, ""),
		}, "/"),
		Subject: "Password Recovery for Arravo Hackathon Account",
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	return c.JSON(200, &InitiateEmailVerificationSuccessResponse{

		Code:    200,
		Message: "Verification of email sent successfully",
		Data:    &InitiateEmailVerificationSuccessResponseData{Email: emailToVerify},
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

	auth := authutils.GetAuthUtilsWithDefaultRepositories()
	_, err = auth.CompletePasswordRecovery(&exports.AuthUtilsCompletePasswordRecoveryData{
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
		Email:   dataDto.Email,
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

// @Title Invite New Member
// @Description	Invite new member
// @Summary		Invite new member
// @Tags			Participants
// @Param registerIndividualJSON body dtos.InviteToTeamData true "invite member to team"
// @Produce		json
// @Success		201	{object}	InviteTeamMemberSuccessResponse
// @Failure		400	{object}	InviteTeamMemberFailResponse
// @Router			/api/v1/auth/me/team/invite               [post]
func InviteMemberToTeam(c echo.Context) error {
	tokenData := authutils.GetAuthPayload(c)
	participantId := tokenData.ParticipantId
	hackathonId := tokenData.HackathonId
	data := dtos.InviteToTeamData{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	err = validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &InviteTeamMemberFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	serv := services.GetServiceWithDefaultRepositories()
	participant, err := serv.GetSingleParticipantWithAccountsInfo(tokenData.ParticipantId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &InviteTeamMemberFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	if participantId != participant.ParticipantId {
		return c.JSON(http.StatusBadRequest, &InviteTeamMemberFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: "Wrong authentication.",
		})
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, &InviteTeamMemberFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	if participant.ParticipantType == "INDIVIDUAL" {
		return c.JSON(http.StatusBadRequest, &InviteTeamMemberFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: "Only a participating team can invite new members.",
		})
	}
	responseData, err := serv.InviteToTeam(&services.AddToTeamInviteListData{
		HackathonId:      hackathonId,
		ParticipantId:    participant.ParticipantId,
		Email:            data.Email,
		Role:             data.Role,
		InviterEmail:     participant.TeamLeadInfo.Email,
		InviterFirstName: participant.TeamLeadInfo.FirstName,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &InviteTeamMemberFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	fmt.Println(responseData)
	return c.JSON(http.StatusCreated, &InviteTeamMemberSuccessResponse{
		Code:    http.StatusCreated,
		Message: "Member will be invited!!!",
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
	data, err := utils.ProcessTeamInviteLink(tokenStr)
	if err != nil {
		fmt.Printf(err.Error())
		c.HTML(400, err.Error())
	}
	return c.Redirect(301,
		strings.Join([]string{
			strings.Join([]string{config.GetFrontendURL(), "team_member_complete"}, "/"),
			strings.Join([]string{
				strings.Join([]string{"participant_id", data.ParticipantId}, "="),
				strings.Join([]string{"hackathon_id", data.HackathonId}, "="),
				strings.Join([]string{"invitee_email", data.InviteeEmail}, "="),
				strings.Join([]string{"teamlead_email", data.TeamLeadEmailEmail}, "="),
			}, "&")}, "?"))
}

// @Title Validate Password Recovery Link
// @Summary		Validate Password Recovery link
// @Description	Validate Password Recovery link
// @Tags			Auth
// @Produce		json
// @Param			completeToken	path		dt	true	"the required info"
// @Success		200				{object}	string
// @Router			/api/v1/auth/password/recovery/link/verification [get]
func ValidatePasswordRecoveryLink(c echo.Context) error {
	tokenStr := c.QueryParam("token")
	if tokenStr == "" {
		return c.HTML(400, strings.Join([]string{config.GetFrontendURL()}, "/"))
	}
	t, err := utils.ProcessPasswordRecoveryLink(tokenStr)
	if err != nil {
		fmt.Printf(err.Error())
		c.HTML(400, err.Error())
	}
	return c.Redirect(301, strings.Join([]string{
		strings.Join([]string{config.GetFrontendURL(), "password_reset_complete"}, "/"),
		strings.Join([]string{
			strings.Join([]string{"email", t.Email}, "="),
			strings.Join([]string{"token", t.Token}, "="),
		}, "&"),
	}, "?"))
}

// @Title Get Team Members Info
// @Description	 Get Team Members Info
// @Summary		 Get Team Members Info
// @Tags			Participants
// @Produce		json
// @Success		200	{object}	GetTeamMembersSuccessResponse
// @Failure		400	{object}	FailResponse
// @Router			/api/v1/auth/me/team              [get]
func GetMyTeamMembersInfo(ctx echo.Context) error {
	payload := authutils.GetAuthPayload(ctx)
	serv := services.GetServiceWithDefaultRepositories()
	participants, err := serv.GetTeamMembersInfo(payload.ParticipantId)
	fmt.Println(participants)
	if err != nil {
		return ctx.JSON(400, GetTeamMembersSuccessResponse{
			Message: "",
			Data:    participants,
		})
	}
	return ctx.JSON(200, GetTeamMembersSuccessResponse{
		Message: "",
		Data:    participants,
		Code:    200,
	})
}

// @Title Get Team Members Info
// @Description	 Delete Team Member from Team
// @Summary		 Remove Team Member from Team
// @Tags			Participants
// @Param  memberId  path  string  true  "email of team member"
// @Produce		json
// @Success		200	{object}	DeleteTeamMemberSuccessResponse
// @Failure		400	{object}	FailResponse
// @Router			/api/v1/auth/me/team/{memberId}              [delete]
func RemoveMemberFromMyTeam(ctx echo.Context) error {
	payload := authutils.GetAuthPayload(ctx)

	memberId := ctx.Param("team_member_email")

	serv := services.GetServiceWithDefaultRepositories()
	participant, err := serv.GetSingleParticipantWithAccountsInfo(payload.Email)

	if err != nil {
		return err
	}

	if participant.TeamLeadEmail != payload.Email {
		return ctx.JSON(401, DeleteTeamMemberFailResponse{
			Message: "Unauthorized",
			Code:    401,
		})
	}

	member, err := serv.RemoveMemberFromTeam(&repository.RemoveMemberFromTeamData{
		MemberEmail:   memberId,
		HackathonId:   payload.HackathonId,
		ParticipantId: payload.ParticipantId,
	})
	if err != nil {
		return ctx.JSON(400, FailResponse{
			Message: err.Error(),
			Code:    400,
		})
	}
	return ctx.JSON(200, DeleteTeamMemberSuccessResponse{
		Message: "",
		Data:    member,
		Code:    200,
	})
}

// @Title Get Team Members Info
// @Description	 Delete Team Member from Team
// @Summary		 Remove Team Member from Team
// @Tags			Participants
// @Param  memberId  path  string  true  "email of team member"
// @Produce		json
// @Success		200	{object}	DeleteTeamMemberSuccessResponse
// @Failure		400	{object}	FailResponse
// @Router			/api/v1/auth/me/team/solution              [put]
func ChooseSolutionForMyTeam(ctx echo.Context) error {
	authPayload := authutils.GetAuthPayload(ctx)
	solDataBodyPayload := dtos.SelectTeamSolutionData{}
	err := ctx.Bind(&solDataBodyPayload)

	if err != nil {
		return ctx.JSON(400, &ResponseData{
			Code:    400,
			Message: err.Error(),
		})
	}
	err = validate.Struct(&solDataBodyPayload)

	if err != nil {
		return ctx.JSON(400, &ResponseData{
			Code:    400,
			Message: err.Error(),
		})
	}
	serv := services.GetServiceWithDefaultRepositories()
	participant, err := serv.GetSingleParticipantWithAccountsInfo(authPayload.Email)

	if err != nil {
		return ctx.JSON(400, &ResponseData{
			Code:    400,
			Message: err.Error(),
		})
	}

	if participant.TeamLeadEmail != authPayload.Email {
		return ctx.JSON(401, DeleteTeamMemberFailResponse{
			Message: "Unauthorized",
			Code:    401,
		})
	}

	err = serv.SelectTeamSolution(&services.SelectTeamSolutionData{
		HackathonId:   authPayload.HackathonId,
		ParticipantId: authPayload.ParticipantId,
		SolutionId:    solDataBodyPayload.SolutionId,
	})
	if err != nil {
		return ctx.JSON(400, FailResponse{
			Message: err.Error(),
			Code:    400,
		})
	}
	return ctx.JSON(200, AddSolutionSuccessResponse{
		Message: "",
		Data:    participant,
		Code:    200,
	})
}
