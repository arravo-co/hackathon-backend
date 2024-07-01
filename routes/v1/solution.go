package routes_v1

import (
	"fmt"
	"net/http"

	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/services"
	"github.com/arravoco/hackathon_backend/utils/authutils"
	"github.com/labstack/echo/v4"
)

type CreateSolutionResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *entity.Solution `json:"data"`
}

// @Title Basic Log in
// @Description	Log a user in
// @Summary		Log a user in
// @Tags		Auth
// @Produce		json
// @Param       createSolutionJSON   body dtos.CreateSolutionData    true                   "Create solution Request JSON"
// @Success		201	  object 	CreateSolutionResponse "Solution created JSON"
// @Failure		400	object	CreateSolutionResponse "Creation of solution failed"
// @Router			/api/v1/solutions             [post]

func CreateSolution(c echo.Context) error {

	tokenData := authutils.GetAuthPayload(c)
	fmt.Println(tokenData)
	dataDto := dtos.CreateSolutionData{}
	err := c.Bind(&dataDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &CreateSolutionResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	err = validate.Struct(dataDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &CreateSolutionResponse{

			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	solutionService := services.NewSolutionService()
	solution, err := solutionService.CreateSolution(&exports.CreateSolutionData{
		CreatorId:   tokenData.Email,
		Title:       dataDto.Title,
		Description: dataDto.Description,
		HackathonId: tokenData.HackathonId,
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, CreateSolutionResponse{
			Code:    500,
			Message: "Error creating solution: ",
		})
	}

	return c.JSON(201, CreateSolutionResponse{
		Code:    201,
		Message: "Successfully created solution",
		Data:    solution,
	})
}

func GetSolutionDataById(c echo.Context) error {

	solId := c.Param("id")
	if solId == "" {
		return c.JSON(http.StatusBadRequest, &CreateSolutionResponse{
			Code:    http.StatusBadRequest,
			Message: "Solution ID is required",
		})
	}

	solutionService := services.NewSolutionService()
	solution, err := solutionService.GetSolutionDataById(solId)

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, CreateSolutionResponse{
			Code:    500,
			Message: "Error fetch solution: ",
		})
	}

	return c.JSON(200, CreateSolutionResponse{
		Code:    200,
		Message: "Successfully fetch solution",
		Data:    solution,
	})
}
