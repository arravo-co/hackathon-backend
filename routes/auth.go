package routes

import (
	"fmt"
	"net/http"

	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/auth"
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

	dataResponse, err := auth.BasicLogin(&auth.BasicLoginData{
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
	Data BasicLoginSuccessResponseData `json:"data"`
}

type InitiateEmailVerificationSuccessResponseData struct {
	AccessToken string `json:"access_token"`
}

func InitiateEmailVerification(c echo.Context) error {
	emailToVerify := c.QueryParam("email")
	if emailToVerify == "" {
		c.JSON(http.StatusBadRequest, &InitiateEmailVerificationFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: "'email' query parameter is required",
			},
		})
	}
}
