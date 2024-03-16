package routes

import (
	"net/http"

	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/labstack/echo/v4"
)

type RegisterJudgeSuccessResponse struct {
	Code    int                         `json:"code"`
	Message string                      `json:"message"`
	Data    data.CreateJudgeAccountData `data:"data"`
}
type RegisterJudgeFailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// @Description	Register new participant
// @Summary		Register new participant
// @Tags			Judges
// @Produce		json
// @Success		201	{object}	RegisterParticipantSuccessResponse
// @Failure		400	{object}	RegisterParticipantFailResponse
// @Router			/api/judges               [post]
func RegisterJudge(c echo.Context) error {
	data := dtos.RegisterNewJudgeDTO{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	newJudge := entity.Judge{}
	err = validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterJudgeFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	responseData, err := newJudge.Register(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterJudgeFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, &RegisterJudgeSuccessResponse{
		Code: http.StatusCreated,
		Data: *responseData,
	})
}
