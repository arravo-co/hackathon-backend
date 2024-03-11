package routes

import (
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/labstack/echo/v4"
)

// @Description Register new participant
// @Summary Register new participant
// @Tags Participants
// @Produce json
// @Router /participants [post]
func RegisterParticipant(c echo.Context) error {
	data := dtos.RegisterParticipantDTO{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	newParticipant := entity.Participant{}

	newParticipant.Register(data)
	return nil
}

// @Description Update participant's information
// @Summary Update participant's information
// @Tags Participants
// @Produce json
// @Router /participants [put]
func UpdateParticipantInfo(c echo.Context) error {
	return nil
}
