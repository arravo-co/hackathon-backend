package routes_v1

import (
	"net/http"

	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/labstack/echo/v4"
)

type RegisterJudgeSuccessResponse struct {
	Code    int                            `json:"code"`
	Message string                         `json:"message"`
	Data    exports.CreateJudgeAccountData `data:"data"`
}
type RegisterJudgeFailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetJudgeSuccessResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *entity.Judge `data:"data"`
}

type GetJudgesSuccessResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    []*entity.Judge `data:"data"`
}

// @Title Register New Judge
// @Description	Register new judge
// @Summary		Register New Judge
// @Tags			Judges
// @Produce		json
// @Param registerJudgeJSON body object dtos.RegisterNewJudgeDTO true "Create Judge profile"
// @Success		201	{object}	RegisterJudgeSuccessResponse
// @Failure		400	{object}	RegisterJudgeFailResponse
// @Router			/api/v1/judges              [post]
func RegisterJudge(c echo.Context) error {
	data := dtos.RegisterNewJudgeDTO{}
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterJudgeFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
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

func UpdateJudge(c echo.Context) error {
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

// @Title Register New Judge
// @Description	Register new judge
// @Summary		Register New Judge
// @Tags			Judges
// @Produce		json
// @Success		200	{object}	GetJudgesSuccessResponse
// @Failure		400	{object}	RegisterJudgeFailResponse
// @Router			/api/v1/judges             [get]
func GetJudges(c echo.Context) error {
	judges, err := entity.GetJudges()
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterJudgeFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, &GetJudgesSuccessResponse{
		Code: http.StatusCreated,
		Data: judges,
	})
}

// @Title Register New Judge
// @Description	Register new judge
// @Summary		Register New Judge
// @Tags			Judges
// @Produce		json
// @Param email path string true "Create Judge profile"
// @Success		200	{object}	GetJudgeSuccessResponse
// @Failure		400	{object}	RegisterJudgeFailResponse
// @Router			/api/v1/judges/{email}               [get]
func GetJudgeByEmailAddress(c echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, &RegisterJudgeFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: "Email address cannot be empty",
		})
	}
	newJudge := &entity.Judge{}
	err := newJudge.FillJudgeEntity(email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterJudgeFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, &GetJudgeSuccessResponse{
		Code: http.StatusCreated,
		Data: newJudge,
	})
}
