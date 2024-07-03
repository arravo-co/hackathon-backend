package routes_v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/publish"
	"github.com/arravoco/hackathon_backend/services"
	taskmgt "github.com/arravoco/hackathon_backend/task_mgt"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/authutils"
	"github.com/labstack/echo/v4"
)

type CreateSolutionResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *entity.Solution `json:"data"`
}

type GetSolutionsDataResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    []entity.Solution `json:"data"`
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
	dataDto.Title = c.FormValue("title")
	dataDto.Description = c.FormValue("description")
	dataDto.Objective = c.FormValue("objective")
	err := validate.Struct(dataDto)
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
		Objective:   dataDto.Objective,
		HackathonId: tokenData.HackathonId,
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, CreateSolutionResponse{
			Code:    500,
			Message: "Error creating solution: ",
		})
	}

	profPic, err := c.FormFile("solution_picture")
	ch := make(chan interface{})
	if profPic != nil {
		go func(solutionId string) {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println(err.Error())
			}
			opt := utils.UploadOpts{
				Folder: filepath.Join(dir, "uploads"),
			}
			filePath, err := utils.GetUploadedPic(profPic, []utils.UploadOpts{opt}...)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			tsk := taskmgt.GenerateTask(&exports.AddTaskDTO{
				Label: "upload solution pic",
			})
			err = taskmgt.SaveTaskById(tsk)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			payload := exports.UploadSolutionPicQueuePayload{
				SolutionId: solutionId,
				FilePath:   filePath,
				QueuePayload: exports.QueuePayload{
					TaskId: tsk.Id,
				},
			}
			byt, err := json.Marshal(payload)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = publish.Publish(&exports.PublisherConfig{
				RabbitMQExchange: "",
				RabbitMQKey:      "upload.solution_picture.cloudinary",
			}, byt)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Queue payload published")
			ch <- struct{}{}
			//judge :=result.SecureURL
		}(solution.Id)
		<-ch
	}
	return c.JSON(201, CreateSolutionResponse{
		Code:    201,
		Message: "Successfully created solution",
		Data:    solution,
	})
}

func UpdateSolutionDataById(c echo.Context) error {
	solId := c.Param("id")
	tokenData := authutils.GetAuthPayload(c)
	fmt.Println(tokenData)
	dataDto := dtos.UpdateSolutionData{}
	dataDto.Title = c.FormValue("title")
	dataDto.Description = c.FormValue("description")
	dataDto.Objective = c.FormValue("objective")

	solutionService := services.NewSolutionService()
	solution, err := solutionService.UpdateSolutionDataById(solId, &exports.UpdateSolutionData{
		CreatorId:   tokenData.Email,
		Title:       dataDto.Title,
		Description: dataDto.Description,
		Objective:   dataDto.Objective,
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, CreateSolutionResponse{
			Code:    500,
			Message: "Error updating solution: ",
		})
	}

	c.JSON(200, CreateSolutionResponse{
		Code:    200,
		Message: "Successfully updated solution",
		Data:    solution,
	})
	profPic, err := c.FormFile("solution_picture")
	ch := make(chan interface{})
	if err == nil && profPic != nil {
		go func(solutionId string) {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println(err.Error())
			}
			opt := utils.UploadOpts{
				Folder: filepath.Join(dir, "uploads"),
			}
			filePath, err := utils.GetUploadedPic(profPic, []utils.UploadOpts{opt}...)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			tsk := taskmgt.GenerateTask(&exports.AddTaskDTO{
				Label: "upload solution pic",
			})
			err = taskmgt.SaveTaskById(tsk)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			payload := exports.UploadSolutionPicQueuePayload{
				SolutionId: solutionId,
				FilePath:   filePath,
				QueuePayload: exports.QueuePayload{
					TaskId: tsk.Id,
				},
			}
			byt, err := json.Marshal(payload)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = publish.Publish(&exports.PublisherConfig{
				RabbitMQExchange: "",
				RabbitMQKey:      "upload.solution_picture.cloudinary",
			}, byt)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Queue payload published")
			ch <- struct{}{}
			//judge :=result.SecureURL
		}(solution.Id)
	}
	<-ch
	return nil
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

func GetSolutionsData(c echo.Context) error {

	hackathonId := c.QueryParam("hackathon_id")

	solutionService := services.NewSolutionService()
	solutions, err := solutionService.GetSolutionData(&exports.GetSolutionsQueryData{
		HackathonId: hackathonId,
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, CreateSolutionResponse{
			Code:    500,
			Message: "Error fetch solution: ",
		})
	}

	return c.JSON(200, GetSolutionsDataResponse{
		Code:    200,
		Message: "Successfully fetch solutions",
		Data:    solutions,
	})
}
