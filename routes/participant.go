package routes

import (
	"net/http"

	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/labstack/echo/v4"
)

type RegisterIndividualParticipantSuccessResponse struct {
	Code    int                                            `json:"code"`
	Message string                                         `json:"message"`
	Data    exports.CreateIndividualParticipantAccountData `data:"data"`
}
type RegisterIndividualParticipantFailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// @Description	Register new participant
// @Summary		Register new participant
// @Tags			Participants
// @Produce		json
// @Success		201	{object}	RegisterIndividualParticipantSuccessResponse
// @Failure		400	{object}	RegisterIndividualParticipantFailResponse
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
	Code    int                                      `json:"code"`
	Message string                                   `json:"message"`
	Data    exports.CreateTeamParticipantAccountData `data:"data"`
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
	data := dtos.RegisterNewTeamParticipantDTO{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	newParticipant := entity.Participant{}
	err = validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterTeamParticipantFailResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	resData, err := newParticipant.RegisterTeam(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterTeamParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, &RegisterTeamParticipantSuccessResponse{
		Code: http.StatusCreated,
		Data: exports.CreateTeamParticipantAccountData{
			TeamLeadEmail:       resData.TeamLeadEmail,
			CoParticipantEmails: resData.CoParticipantEmails,
			TeamName:            resData.TeamName,
			HackathonId:         resData.HackathonId,
			ParticipantId:       resData.ParticipantId,
			Type:                resData.ParticipantId,
			GithubAddress:       resData.GithubAddress,
			Role:                resData.Role,
		},
	})
}
