package routes

import (
	"fmt"
	"net/http"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/utils/authutils"
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

// @Description	Register new participant
// @Summary		Register new participant
// @Tags			Participants
// @Param registerIndividualJSON body dtos.RegisterNewParticipantDTO true "register individual participant"
// @Produce		json
// @Success		201	{object}	RegisterParticipantSuccessResponse
// @Failure		400	{object}	RegisterParticipantFailResponse
// @Router			/api/participants               [post]
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
	if data.Type == "PARTICIPANT" {
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
	Code    int                                      `json:"code"`
	Message string                                   `json:"message"`
	Data    exports.CreateTeamParticipantAccountData `json:"data"`
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

// @Description	Fully register new member to the participating team
// @Summary		Fully Register New Member To The Participating Team
// @Tags		Participants
// @Produce		json
// @Param  participantId  path  string  true  "participant id of the participating team"
// @Param registerTeam body dtos.RegisterNewTeamMemberDTO true "register new member to participating team"
// @Success		201	{object}	RegisterTeamParticipantSuccessResponse
// @Failure		400	{object}	RegisterTeamParticipantFailResponse
// @Router		/api/participants/{participantId}/members     [post]
func RegisterNewTeamMember(c echo.Context) error {
	participantId := c.Param("participantId")
	data := dtos.RegisterNewTeamMemberDTO{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	data.ParticipantId = participantId
	err = validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterTeamParticipantFailResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	newParticipant := entity.Participant{}
	resData, err := newParticipant.RegisterNewTeamMember(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RegisterTeamParticipantFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, &RegisterTeamParticipantSuccessResponse{
		Code: http.StatusCreated,
		Data: exports.CreateTeamParticipantAccountData{
			TeamLeadEmail:       resData.TeamLeadEmail,
			CoParticipantEmails: resData.CoParticipantEmails,
			TeamName:            resData.TeamName,
			HackathonId:         resData.HackathonId,
			ParticipantId:       resData.ParticipantId,
			Type:                resData.Type,
			Role:                resData.Role,
		},
	})
}

// @Description	Invite new member
// @Summary		Invite new member
// @Tags			Participants
// @Param  participantId  path  string  true  "participant id of the participating team"
// @Param registerIndividualJSON body dtos.InviteToTeamData true "invite member to team"
// @Produce		json
// @Success		201	{object}	InviteTeamMemberSuccessResponse
// @Failure		400	{object}	InviteTeamMemberFailResponse
// @Router			/api/participants/{participantId}/invite               [post]
func InviteMemberToTeam(c echo.Context) error {
	participantId := c.Param("participantId")
	tokenData := authutils.GetAuthPayload(c)
	data := dtos.InviteToTeamData{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	data.ParticipantId = participantId
	err = validate.Struct(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &InviteTeamMemberFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	participant := entity.Participant{}
	err = participant.FillParticipantInfo(tokenData.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &InviteTeamMemberFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	if participant.Type == "INDIVIDUAL" {
		return c.JSON(http.StatusBadRequest, &InviteTeamMemberFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: "Only a participating team can invite new members.",
		})
	}
	responseData, err := participant.InviteToTeam(&exports.AddToTeamInviteListData{
		HackathonId:   config.GetHackathonId(),
		ParticipantId: data.ParticipantId,
		Email:         data.Email,
		Role:          data.Role,
		InviterEmail:  participant.Email,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &InviteTeamMemberFailResponse{
			Code:    echo.ErrBadRequest.Code,
			Message: err.Error(),
		})
	}
	fmt.Println(responseData)
	return c.JSON(http.StatusCreated, &InviteTeamMemberSuccessResponse{
		Code:    http.StatusCreated,
		Message: "Member will be invited!!!",
	})
}
