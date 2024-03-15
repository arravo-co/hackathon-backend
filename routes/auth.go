package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arravoco/hackathon_backend/dtos"
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
// @Tags		Auth
// @Produce		json
// @Success		200							{object}	BasicLoginSuccessResponse
// @Failure		400	                        {object}	BasicLoginFailureResponse
// @Router		/api/auth/login             [post]
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

// @Summary      Verify user email address
// @Description  Verify user email address
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        email     query     string  false  "Email to verify"  Format(email)
// @Success      200  	   {object}  InitiateEmailVerificationSuccessResponse
// @Failure      400       {object}  InitiateEmailVerificationFailureResponse
// @Failure      404       {object}  InitiateEmailVerificationFailureResponse
// @Failure      500       {object}  InitiateEmailVerificationFailureResponse
// @Router       /api/auth/verification/email/initiation [get]
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
