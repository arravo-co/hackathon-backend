package routes

import (
	"net/http"

	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/labstack/echo/v4"
)

type RegisterIndividualParticipantSuccessResponse struct {
	Code    int                                         `json:"code"`
	Message string                                      `json:"message"`
	Data    data.CreateIndividualParticipantAccountData `data:"data"`
}
type RegisterIndividualParticipantFailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// @Description	Register new participant
// @Summary		Register new participant
// @Tags			Participants
// @Produce		json
// @Success		201	{object}	RegisterParticipantSuccessResponse
// @Failure		400	{object}	RegisterParticipantFailResponse
// @Router			/api/participants/individual               [post]
func RegisterIndividualParticipant(c echo.Context) error {
	data := dtos.RegisterNewIndividualParticipantDTO{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	newParticipant := entity.Participant{}
	err = validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterIndividualParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	responseData, err := newParticipant.RegisterIndividual(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterIndividualParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, &RegisterIndividualParticipantSuccessResponse{
		Code: http.StatusCreated,
		Data: *responseData,
	})
}

type RegisterTeamParticipantSuccessResponse struct {
	Code    int                                   `json:"code"`
	Message string                                `json:"message"`
	Data    data.CreateTeamParticipantAccountData `data:"data"`
}
type RegisterTeamParticipantFailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// @Description	Register new participant
// @Summary		Register new participant
// @Tags			Participants
// @Produce		json
// @Success		201	{object}	RegisterTeamParticipantSuccessResponse
// @Failure		400	{object}	RegisterTeamParticipantFailResponse
// @Router			/api/participants/team               [post]
func RegisterTeamParticipant(c echo.Context) error {
	data := dtos.RegisterNewIndividualParticipantDTO{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	newParticipant := entity.Participant{}
	err = validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterTeamParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	_, err = newParticipant.RegisterIndividual(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterTeamParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, &RegisterTeamParticipantSuccessResponse{
		Code: http.StatusCreated,
		//Data: *responseData,
	})
}
