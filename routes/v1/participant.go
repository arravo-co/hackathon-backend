package routes_v1

import (
	"fmt"
	"net/http"

	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/services"
	"github.com/labstack/echo/v4"
)

type RegisterParticipantSuccessResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    *entity.Participant `data:"data"`
}

type RegisterParticipantFailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetParticipantsResponseData struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    []entity.Participant `json:"data"`
}

type GetParticipantResponseData struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    *entity.Participant `json:"data"`
}

type RegisterTeamParticipantSuccessResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    entity.Participant `json:"data"`
}
type RegisterTeamParticipantFailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type InviteTeamMemberFailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type InviteTeamMemberSuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type FailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetTeamMembersSuccessResponse struct {
	Code    int                                      `json:"code"`
	Message string                                   `json:"message"`
	Data    []entity.TeamMemberWithParticipantRecord `json:"data"`
}

// @Title Register New Participant
// @Description	Register new participant
// @Summary		Register new participant
// @Tags			Participants
// @Param registerIndividualJSON body dtos.RegisterNewParticipantDTO true "register individual participant"
// @Produce		json
// @Success		201	{object}	RegisterParticipantSuccessResponse
// @Failure		400	{object}	RegisterParticipantFailResponse
// @Router			/api/v1/participants               [post]
func RegisterParticipant(c echo.Context) error {
	data := dtos.RegisterNewParticipantDTO{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	err = validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	var responseData *entity.Participant
	if data.Type == "INDIVIDUAL" {
		return c.JSON(http.StatusBadRequest, &RegisterParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: "Only teams are allowed to register",
		})
	} else if data.Type == "TEAM" {
		serv := services.GetServiceWithDefaultRepositories()
		responseData, err = serv.RegisterTeamLead(&services.RegisterNewParticipantDTO{
			FirstName:           data.FirstName,
			LastName:            data.LastName,
			Gender:              data.Gender,
			Password:            data.Password,
			ConfirmPassword:     data.ConfirmPassword,
			PhoneNumber:         data.PhoneNumber,
			PreviousProjects:    data.PreviousProjects,
			Skillset:            data.Skillset,
			State:               data.State,
			TeamSize:            data.TeamSize,
			DOB:                 data.DOB,
			EmploymentStatus:    data.EmploymentStatus,
			Email:               data.Email,
			ExperienceLevel:     data.ExperienceLevel,
			HackathonExperience: data.HackathonExperience,
			YearsOfExperience:   data.YearsOfExperience,
			FieldOfStudy:        data.FieldOfStudy,
			Type:                data.Type,
			TeamName:            data.TeamName,
			Motivation:          data.Motivation,
		})
		if err != nil {
			return c.JSON(http.StatusBadRequest, &RegisterParticipantFailResponse{
				Code:    echo.ErrBadRequest.Code,
				Message: err.Error(),
			})
		}
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, &RegisterParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &RegisterParticipantSuccessResponse{
		Code: http.StatusCreated,
		Data: responseData,
	})
}

// @Title Fully Register New Team Member
// @Description	Fully register new member to the participating team
// @Summary		Fully Register New Member To The Participating Team
// @Tags		Participants
// @Produce		json
// @Param  participantId  path  string  true  "participant id of the participating team"
// @Param registerTeam body dtos.CompleteNewTeamMemberRegistrationDTO true "register new member to participating team"
// @Success		201	{object}	RegisterTeamParticipantSuccessResponse
// @Failure		400	{object}	RegisterTeamParticipantFailResponse
// @Router		/api/v1/participants/{participantId}/members     [post]
func CompleteNewTeamMemberRegistration(c echo.Context) error {
	participantId := c.Param("participantId")
	data := dtos.CompleteNewTeamMemberRegistrationDTO{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	err = validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterTeamParticipantFailResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	q := query.GetDefaultQuery()
	if q != nil {
		return c.JSON(http.StatusBadRequest, &RegisterTeamParticipantFailResponse{
			Code:    http.StatusBadRequest,
			Message: "Failed to complete registration. Try again later",
		})
	}

	serv := services.GetServiceWithDefaultRepositories()
	resData, err := serv.CompleteNewTeamMemberRegistration(&services.CompleteNewTeamMemberRegistrationDTO{
		FirstName:     data.FirstName,
		Email:         data.Email,
		LastName:      data.LastName,
		Gender:        data.Gender,
		Password:      data.Password,
		Skillset:      data.Skillset,
		State:         data.State,
		HackathonId:   data.HackathonId,
		TeamLeadEmail: data.TeamLeadEmail,
		DOB:           data.DOB,
		PhoneNumber:   data.PhoneNumber,
		ParticipantId: participantId, Motivation: data.Motivation,
		EmploymentStatus: data.EmploymentStatus,
		ExperienceLevel:  data.ExperienceLevel,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterTeamParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, &RegisterTeamParticipantSuccessResponse{
		Code: http.StatusCreated,
		Data: *resData,
	})
}

// @Title Get Team Members Info
// @Description	 Get Team Members Info
// @Summary		 Get Team Members Info
// @Tags			Participants
// @Param  participantId  path  string  true  "participant id of the participating team"
// @Produce		json
// @Success		200	{object}	GetTeamMembersSuccessResponse
// @Failure		400	{object}	FailResponse
// @ Router			/api/v1/participants/{participantId}/team              [get]
func GetTeamMembersInfo(ctx echo.Context) error {
	participantId := ctx.Param("participantId")
	serv := services.GetServiceWithDefaultRepositories()
	participant, err := serv.GetParticipantInfo(participantId)
	if err != nil {
		return ctx.JSON(400, GetTeamMembersSuccessResponse{
			Message: "Failed to fetch team members information",
		})
	}
	mems, err := serv.GetTeamMembersInfo(participant)

	if err != nil {
		return ctx.JSON(400, GetTeamMembersSuccessResponse{
			Message: "Failed to fetch team members information",
		})
	}
	return ctx.JSON(200, GetTeamMembersSuccessResponse{
		Message: "",
		Data:    mems,
		Code:    200,
	})
}

// @Title Get participant info
// @Summary Get participant info
// @Description Get participant info
// @Success 200 object GetParticipantsResponseData "Get participant info"
// @Failure 400 object RegisterAnotherAdminResponseData "Failed to get participant info"
// @Router /api/v1/participants [get]
func GetParticipants(c echo.Context) error {
	serv := services.GetServiceWithDefaultRepositories()
	participants, err := serv.GetParticipantsInfo()
	if err != nil {
		return c.JSON(400, &FailResponse{
			Code:    400,
			Message: err.Error(),
		})
	}
	return c.JSON(200, &GetParticipantsResponseData{
		Code:    200,
		Message: "List of participants",
		Data:    participants,
	})
}

// @Title Get participant info
// @Summary Get participant info
// @Description Register new admin
// @Param  participantId  path  string  true  "participant id of the participant"
// @Success 200 object GetParticipantResponseData "Get participant info"
// @Failure 400 object RegisterAnotherAdminResponseData "Failed to get participant info"
// @Router /api/v1/participants/{participantId} [get]
func GetParticipant(c echo.Context) error {
	participantId := c.Param("participantId")
	serv := services.GetServiceWithDefaultRepositories()
	participant, err := serv.GetParticipantInfo(participantId)
	if err != nil {
		return c.JSON(400, &FailResponse{
			Code:    400,
			Message: err.Error(),
		})
	}
	return c.JSON(200, &GetParticipantResponseData{
		Code:    200,
		Message: "Participant",
		Data:    participant,
	})
}
