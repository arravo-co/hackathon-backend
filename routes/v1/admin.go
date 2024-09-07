package routes_v1

import (
	"fmt"

	"github.com/aidarkhanov/nanoid"
	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/services"
	"github.com/labstack/echo/v4"
)

type RegisterAdminResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    entity.Admin
}
type RegisterAnotherAdminResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RegisterJudgeByAdminResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func RegisterAdmin(c echo.Context) error {
	dataInput := dtos.CreateNewAdminDTO{}
	c.Bind(&dataInput)
	err := validate.Struct(dataInput)
	if err != nil {
		return c.JSON(400, &RegisterAnotherAdminResponseData{
			Code:    400,
			Message: err.Error(),
		})
	}
	serv := services.GetServiceWithDefaultRepositories()
	admin, err := serv.RegisterAdmin(&services.CreateNewAdminDTO{
		Email:       dataInput.Email,
		LastName:    dataInput.LastName,
		FirstName:   dataInput.FirstName,
		PhoneNumber: dataInput.PhoneNumber,
		Password:    dataInput.Password,
		HackathonId: config.GetHackathonId(),
	})
	if err != nil {
		return c.JSON(400, &RegisterAnotherAdminResponseData{
			Code:    400,
			Message: err.Error(),
		})
	}
	return c.JSON(200, &RegisterAdminResponseData{
		Code:    200,
		Message: "Admin created successfully",
		Data:    *admin,
	})
}

// @Title Register A New Admin By Another Admin
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token
// @Summary Register new admin
// @Description Register new admin
// @Param bodyJSON body dtos.CreateNewAdminDTO true "data of new admin to register"
// @Success 201 object RegisterAnotherAdminResponseData "Admin successfully registered"
// @Failure 400 object RegisterAnotherAdminResponseData "Failed to register new Admin"
// @Router /api/v1/admin/register_admin [post]
func RegisterAnotherAdmin(c echo.Context) error {
	authPayload := exports.GetPayload(c)
	dataInput := dtos.CreateNewAdminByAuthAdminDTO{}
	c.Bind(&dataInput)
	err := validate.Struct(dataInput)
	if err != nil {
		return c.JSON(400, &RegisterAnotherAdminResponseData{
			Code:    400,
			Message: err.Error()})
	}
	serv := services.GetServiceWithDefaultRepositories()
	_, err = serv.AdminCreateNewAdminProfile(&services.CreateNewAdminByAuthAdminDTO{
		Email:        dataInput.Email,
		LastName:     dataInput.LastName,
		FirstName:    dataInput.FirstName,
		Gender:       dataInput.Gender,
		PhoneNumber:  dataInput.PhoneNumber,
		InviterEmail: authPayload.Email,
		InviterName:  authPayload.FirstName,
		HackathonId:  config.GetHackathonId(),
	})
	if err != nil {
		return c.JSON(400, &RegisterAnotherAdminResponseData{
			Code:    400,
			Message: err.Error(),
		})
	}
	return c.JSON(200, &RegisterAnotherAdminResponseData{
		Code:    200,
		Message: "Invite sent!!!",
	})
}

// @Title Register A New Judge
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token
// @Summary Register new judge
// @Description Register new judge
// @Param bodyJSON body dtos.CreateNewJudgeByAdminDTO true "data of new judge to register"
// @Success 201 object RegisterJudgeByAdminResponseData "Judge successfully registered"
// @Failure 400 object RegisterJudgeByAdminResponseData "Failed to register new Judge"
// @Router /api/v1/admin/register_judge [post]
func RegisterJudgeByAdmin(c echo.Context) error {
	authPayload := exports.GetPayload(c)
	dataInput := dtos.CreateNewJudgeByAdminDTO{
		Email:       c.FormValue("email"),
		LastName:    c.FormValue("last_name"),
		FirstName:   c.FormValue("first_name"),
		Gender:      c.FormValue("gender"),
		PhoneNumber: c.FormValue("phone_number"),
		Bio:         c.FormValue("bio"),
	}
	/*	err := c.Bind(&dataInput)
		if err != nil {
			return c.JSON(400, &RegisterJudgeByAdminResponseData{
				Code:    400,
				Message: "Failed to process request"})
		}
	*/
	err := validate.Struct(dataInput)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, &RegisterJudgeByAdminResponseData{
			Code:    400,
			Message: err.Error()})
	}

	profPic, err := c.FormFile("profile_picture")
	if err != nil {
		fmt.Println("Failed to load judge image")
		return c.JSON(400, &RegisterJudgeByAdminResponseData{
			Code:    400,
			Message: err.Error()})
	}
	password := nanoid.Must(nanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456456789@#$%^&*()+_", 10))
	serv := services.GetServiceWithDefaultRepositories()
	_, err = serv.RegisterNewJudge(&services.RegisterNewJudgeDTO{
		Email:           dataInput.Email,
		LastName:        dataInput.LastName,
		FirstName:       dataInput.FirstName,
		Gender:          dataInput.Gender,
		PhoneNumber:     dataInput.PhoneNumber,
		Password:        password,
		ConfirmPassword: password,
		Bio:             dataInput.Bio,
		InviterEmail:    authPayload.Email,
		InviterName:     authPayload.FirstName,
	})
	if err != nil {
		return c.JSON(400, &RegisterAnotherAdminResponseData{
			Code:    400,
			Message: err.Error(),
		})
	}
	err = c.JSON(201, &RegisterAnotherAdminResponseData{
		Code:    201,
		Message: "Judge Profile Created. Invite sent!!!",
	})

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, &RegisterAnotherAdminResponseData{
			Code:    400,
			Message: err.Error(),
		})
	}
	//ch := make(chan interface{})
	if profPic != nil {
		serv.UploadJudgeProfile(services.UploadJudgeProfilePictureOpt{
			Email:       dataInput.Email,
			PictureFile: profPic,
		})
		/*go func() {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println(err.Error())
			}
			opt := utils.UploadOpts{
				Folder: filepath.Join(dir, "uploads"),
				FileNamePrefix: strings.Join([]string{email},""),
			}
			filePath, err := utils.SaveFile(profPic, []utils.UploadOpts{opt}...)
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
				Email:    dataInput.Email,
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
		<-ch*/
	}
	return nil
}
