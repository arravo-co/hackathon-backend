package routes

import (
	"net/http"

	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/labstack/echo/v4"
)

type RegisterParticipantSuccessResponse struct {
	Code    int                               `json:"code"`
	Message string                            `json:"message"`
	Data    data.CreateParticipantAccountData `data:"data"`
}
type RegisterParticipantFailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// @Description	Register new participant
// @Summary		Register new participant
// @Tags		Participants
// @Produce		json
// @Success		201							{object}	RegisterParticipantSuccessResponse
// @Failure		400	                        {object}	RegisterParticipantFailResponse
// @Router		/api/participants               [post]
func RegisterParticipant(c echo.Context) error {
	data := dtos.RegisterNewParticipantDTO{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	newParticipant := entity.Participant{}
	err = validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	responseData, err := newParticipant.Register(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, &RegisterParticipantSuccessResponse{
		Code: http.StatusCreated,
		Data: *responseData,
	})
}

// @Description	Update participant's information
// @Summary		Update participant's information
// @Tags			Participants
// @Produce		json
// @Router			/participants [put]
func UpdateParticipantInfo(c echo.Context) error {
	return nil
}
