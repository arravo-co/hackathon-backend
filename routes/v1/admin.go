package routes_v1

import (
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/labstack/echo/v4"
)

type RegisterAnotherAdminResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
	dataInput := dtos.CreateNewAdminDTO{}
	c.Bind(&dataInput)
	validate.Struct(dataInput)
	authAdmin := entity.Admin{}
	err := authAdmin.FillEntity(authPayload.Email)
	if err != nil {
		return c.JSON(400, &RegisterAnotherAdminResponseData{
			Code:    400,
			Message: "Failed at fully authenticating admin"})
	}
	err = authAdmin.AdminCreateNewAdminProlife(&dtos.CreateNewAdminDTO{
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
	return c.JSON(201, &RegisterAnotherAdminResponseData{
		Code:    201,
		Message: "Invite sent!!!",
	})
}
