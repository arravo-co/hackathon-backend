package routes_v1

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/rmqUtils"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/labstack/echo/v4"
)

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
	authAdmin := entity.Admin{}
	err = authAdmin.RegisterNewAdmin(&dtos.CreateNewAdminDTO{
		Email:       dataInput.Email,
		LastName:    dataInput.LastName,
		FirstName:   dataInput.FirstName,
		PhoneNumber: dataInput.PhoneNumber,
		Password:    dataInput.Password,
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

// @Title Register A New Admin
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
	validate.Struct(dataInput)
	authAdmin := entity.Admin{}
	err := authAdmin.FillAdminEntity(authPayload.Email)
	if err != nil {
		return c.JSON(400, &RegisterAnotherAdminResponseData{
			Code:    400,
			Message: "Failed at fully authenticating admin"})
	}
	err = authAdmin.AdminCreateNewAdminProlife(&dtos.CreateNewAdminByAuthAdminDTO{
		Email:       dataInput.Email,
		LastName:    dataInput.LastName,
		FirstName:   dataInput.FirstName,
		Gender:      dataInput.Gender,
		PhoneNumber: dataInput.PhoneNumber,
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
		fmt.Println("Failed to fully authenticate admin")
		return c.JSON(400, &RegisterJudgeByAdminResponseData{
			Code:    400,
			Message: err.Error()})
	}

	profPic, err := c.FormFile("profile_picture")
	if err != nil {
		fmt.Println("Failed to fully authenticate admin")
		return c.JSON(400, &RegisterJudgeByAdminResponseData{
			Code:    400,
			Message: err.Error()})
	}
	authAdmin := entity.Admin{}
	err = authAdmin.FillAdminEntity(authPayload.Email)
	if err != nil {
		return c.JSON(400, &RegisterJudgeByAdminResponseData{
			Code:    400,
			Message: "Failed to fully authenticate admin"})
	}
	err = authAdmin.AdminCreateNewJudgeProlife(&dtos.CreateNewJudgeByAdminDTO{
		Email:       dataInput.Email,
		LastName:    dataInput.LastName,
		FirstName:   dataInput.FirstName,
		Gender:      dataInput.Gender,
		PhoneNumber: dataInput.PhoneNumber,
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
	ch := make(chan interface{})
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
		fmt.Println(filePath)
		queue, err := rmqUtils.GetQueue("upload_pic_cloudinary")

		if err != nil {
			fmt.Printf("Error getting queue: %s\n", err.Error())
			return
		}
		fmt.Println("Queue created")
		payload := exports.UploadPicQueuePayload{
			Email:    dataInput.Email,
			FilePath: filePath,
		}
		byt, err := json.Marshal(payload)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = queue.PublishBytes(byt)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Queue payload published")
		ch <- struct{}{}
		//judge :=result.SecureURL
	}()
	<-ch
	return nil
}
