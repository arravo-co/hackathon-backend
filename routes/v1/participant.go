package routes_v1

import (
	"fmt"
	"net/http"

	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
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
	newParticipant := entity.Participant{}
	err = validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	var responseData *entity.Participant
	if data.Type == "INDIVIDUAL" {
		responseData, err = newParticipant.RegisterIndividual(data)
	} else if data.Type == "TEAM" {
		responseData, err = newParticipant.RegisterTeamLead(data)
	}
	fmt.Println(err)
	if err != nil {
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
	Code    int                        `json:"code"`
	Message string                     `json:"message"`
	Data    []entity.TeamMemberAccount `json:"data"`
}

// @Title Fully Register New Team Member
// @Description	Fully register new member to the participating team
// @Summary		Fully Register New Member To The Participating Team
// @Tags		Participants
// @Produce		json
// @Param  participantId  path  string  true  "participant id of the participating team"
// @Param registerTeam body dtos.RegisterNewTeamMemberDTO true "register new member to participating team"
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
	resData, err := entity.CompleteNewTeamMemberRegistration(&entity.CompleteNewTeamMemberRegistrationEntityData{
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
		ParticipantId: participantId,
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
	participant := &entity.Participant{}
	err := participant.FillParticipantInfo(participantId)
	if err != nil {
		return err
	}
	participants, err := participant.GetTeamMembersInfo()
	if err != nil {
		return ctx.JSON(400, GetTeamMembersSuccessResponse{
			Message: "",
			Data:    participants,
		})
	}
	return ctx.JSON(200, GetTeamMembersSuccessResponse{
		Message: "",
		Data:    participants,
		Code:    200,
	})
}
