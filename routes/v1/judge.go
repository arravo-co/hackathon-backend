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
	"github.com/arravoco/hackathon_backend/repository"
	taskmgt "github.com/arravoco/hackathon_backend/task_mgt"
	"github.com/arravoco/hackathon_backend/utils"
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

type UpdateJudgeSuccessResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *entity.Judge `data:"data"`
}
type UpdateJudgeFailResponse struct {
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
	newJudge := repository.JudgeAccountRepository{}
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

// @Title Update Judge info
// @Description	Register new judge
// @Summary		Register New Judge
// @Tags			Judges
// @Produce		json
// @Param updateJudgeFORM form dtos.UpdateJudgeDTO true "Update Judge profile"
// @Success		200	{object}	UpdateJudgeSuccessResponse
// @Failure		400	{object}	UpdateJudgeFailResponse
// @Router			/api/v1/judges/{email}              [put]
func UpdateJudge(c echo.Context) error {
	email := c.Param("email")
	data := dtos.UpdateJudgeDTO{
		FirstName: c.FormValue("first_name"),
		LastName:  c.FormValue("last_name"),
		Gender:    c.FormValue("gender"),
		State:     c.FormValue("state"),
		Bio:       c.FormValue("bio"),
	}
	err := validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterJudgeFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}

	err = repository.UpdateJudgeProfile(email, dtos.UpdateJudgeDTO{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Gender:    data.Gender,
		Bio:       data.Bio,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterJudgeFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	judge, err := repository.GetJudgeByEmail(email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterJudgeFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}

	profPic, err := c.FormFile("profile_picture")
	ch := make(chan interface{})
	if profPic != nil {
		go func() {
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
				Label: "upload profile pic",
			})
			err = taskmgt.SaveTaskById(tsk)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			payload := exports.UploadJudgeProfilePicQueuePayload{
				Email:    email,
				FilePath: filePath,
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
				RabbitMQKey:      "upload.profile_picture.cloudinary",
			}, byt)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Queue payload published")
			ch <- struct{}{}
			//judge :=result.SecureURL
		}()
		<-ch
	}
	return c.JSON(http.StatusCreated, &UpdateJudgeSuccessResponse{
		Code: http.StatusCreated,
		Data: judge,
	})

}

// @Title Get Judges
// @Description	Register new judge
// @Summary		Register New Judge
// @Tags			Judges
// @Produce		json
// @Success		200	{object}	GetJudgesSuccessResponse
// @Failure		400	{object}	RegisterJudgeFailResponse
// @Router			/api/v1/judges             [get]
func GetJudges(c echo.Context) error {
	judges, err := repository.GetJudges()
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

// @Title Get Judge by Email Address
// @Description	Get Judge by Email Address
// @Summary		Get Judge by Email Address
// @Tags			Judges
// @Produce		json
// @Param email path string true "Get Judge by Email Address"
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
	judge, err := repository.GetJudgeByEmail(email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterJudgeFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, &GetJudgeSuccessResponse{
		Code: http.StatusCreated,
		Data: judge,
	})
}
