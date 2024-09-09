package routes_v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/services"
	"github.com/arravoco/hackathon_backend/utils/authutils"
	"github.com/labstack/echo/v4"
)

type GetParticipantsFilterOpts struct {
	ParticipantId            string `json:"participant_id"`
	ParticipantStatus        string `validate:"omitempty,oneof=UNREVIEWED REVIEWED AI_RANKED REVIEW_DISQUALIFIED TEAM_ONBOARDING SOLUTION_SELECTION SOLUTION_IMPLEMENTATION SHORTLISTED COMPETITION_WON" json:"participant_status"`
	ParticipantType          string `validate:"omitempty,oneof=TEAM" json:"participant_type"`
	ReviewRanking_Eq         int    `json:"review_ranking_eq"`
	ReviewRanking_Top        int    `json:"review_ranking_top"`
	Solution_Like            string `json:"solution_like"`
	Limit                    int    `json:"limit"`
	SortByReviewRanking_Asc  bool   `json:"sort_by_review_ranking_asc"`
	SortByReviewRanking_Desc bool   `json:"sort_by_review_ranking_desc"`
}

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

type GetTeamIntoSuccessResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    entity.Participant `json:"data"`
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

	serv := services.GetServiceWithDefaultRepositories()
	resData, err := serv.CompleteNewTeamMemberRegistration(&services.CompleteNewTeamMemberRegistrationDTO{
		FirstName:        data.FirstName,
		Email:            data.Email,
		LastName:         data.LastName,
		Gender:           data.Gender,
		Password:         data.Password,
		Skillset:         data.Skillset,
		State:            data.State,
		HackathonId:      data.HackathonId,
		TeamLeadEmail:    data.TeamLeadEmail,
		DOB:              data.DOB,
		PhoneNumber:      data.PhoneNumber,
		ParticipantId:    participantId,
		Motivation:       data.Motivation,
		EmploymentStatus: data.EmploymentStatus,
		ExperienceLevel:  data.ExperienceLevel,
		TeamRole:         "TEAM_MEMBER",
		ConfirmPassword:  data.ConfirmPassword,
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
	mems, err := serv.GetTeamMembersInfo(participantId)

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
	fmt.Println("GetParticipants called")
	queryFilterObj := &GetParticipantsFilterOpts{}
	var limit int
	var review_ranking_top int
	limit_str := c.QueryParam("limit")
	if limit_str != "" {
		var err error
		limit, err = strconv.Atoi(limit_str)
		if err != nil {
			return c.JSON(400, &FailResponse{
				Code:    400,
				Message: fmt.Sprint("Ensure limit is an integer ", err.Error()),
			})
		}
	}

	review_ranking_top_str := c.QueryParam("review_ranking_top_str")
	if review_ranking_top_str != "" {
		var err error
		review_ranking_top, err = strconv.Atoi(review_ranking_top_str)
		if err != nil {
			return c.JSON(400, &FailResponse{
				Code:    400,
				Message: fmt.Sprint("Ensure 'review_ranking_top' is an integer ", err.Error()),
			})
		}
	}
	queryFilterObj.Limit = limit
	queryFilterObj.ReviewRanking_Top = review_ranking_top
	queryFilterObj.ParticipantId = c.QueryParam("participant_id")
	queryFilterObj.ParticipantStatus = c.QueryParam("participant_status")
	queryFilterObj.ParticipantType = c.QueryParam("participant_type")
	if queryFilterObj.Limit == 0 {
		queryFilterObj.Limit = 1000
	}
	err := validate.Struct(queryFilterObj)
	if err != nil {
		return c.JSON(400, &FailResponse{
			Code:    400,
			Message: err.Error(),
		})
	}

	fmt.Println(queryFilterObj.Limit)
	filterObj := &services.GetParticipantsWithAccountsAggregateFilterOpts{
		ParticipantId:     &queryFilterObj.ParticipantId,
		ParticipantStatus: &queryFilterObj.ParticipantStatus,
		ParticipantType:   &queryFilterObj.ParticipantType,
		/*ReviewRanking_Eq:  &queryFilterObj.ReviewRanking_Eq,
		ReviewRanking_Top: &queryFilterObj.ReviewRanking_Top,
		Solution_Like:     &queryFilterObj.Solution_Like,
		Limit:                    &queryFilterObj.Limit,
		SortByReviewRanking_Asc:  &queryFilterObj.SortByReviewRanking_Asc,
		SortByReviewRanking_Desc: &queryFilterObj.SortByReviewRanking_Desc,*/
	}
	serv := services.GetServiceWithDefaultRepositories()
	participants, err := serv.GetMultipleParticipantsWithAccounts(filterObj)
	if err != nil {
		return c.JSON(400, &FailResponse{
			Code:    400,
			Message: err.Error(),
		})
	}
	//fmt.Println(participants)
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
	participant, err := serv.GetSingleParticipantWithAccountsInfo(participantId)
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

// @Title Admin-only Auth User Info
// @Summary		Update admin-auth participant information
// @Description	Update admin-auth participant information
// @Tags			Auth
// @Param			updateMyInfoJSON	body		dtos.AuthParticipantInfoUpdateDTO	true	"the required info"
// @Produce		json
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token
// @Success		200	{object}	AuthUserInfoFetchSuccessResponse
// @Failure		400	{object}	AuthUserInfoFetchFailureResponse
// @Failure		404	{object}	AuthUserInfoFetchFailureResponse
// @Failure		500	{object}	AuthUserInfoFetchFailureResponse
// @Router			/api/v1/admin/updates/participants [post]
func UpdateAuthParticipantInfo(c echo.Context) error {
	tokenData := authutils.GetAuthPayload(c)
	hackathon_id := c.QueryParam("hackathon_id")
	participant_id := c.QueryParam("participant_id")
	if hackathon_id == "" {
		return c.JSON(http.StatusBadRequest, &PasswordChangeFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: "hackathon_id is required",
			},
		})
	}

	if participant_id == "" {
		return c.JSON(http.StatusBadRequest, &PasswordChangeFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: "participant_id is required",
			},
		})
	}
	var err error
	updateData := &dtos.AdminParticipantInfoUpdateDTO{}
	err = c.Bind(updateData)
	if err != nil {

		return c.JSON(http.StatusBadRequest, &PasswordChangeFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	err = validate.Struct(updateData)
	if err != nil {

		return c.JSON(http.StatusBadRequest, &PasswordChangeFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	serv := services.GetServiceWithDefaultRepositories()
	_, err = serv.AdminUpdateParticipantInfo(&services.UpdateSingleParticipantRecordFilter{
		ParticipantId: participant_id,
		HackathonId:   hackathon_id,
		AdminEmailId:  tokenData.Email,
	}, &services.AdminParticipantInfoUpdateDTO{
		ReviewRanking: updateData.ReviewRanking,
		Status:        updateData.Status,
		//GithubAddress:   updateData.GithubAddress,
		//LinkedInAddress: updateData.LinkedInAddress,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, &PasswordChangeFailureResponse{
			ResponseData{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	return c.JSON(200, &PasswordChangeSuccessResponse{
		ResponseData: ResponseData{
			Code:    200,
			Message: "Participant record change completed successfully",
		}, Data: nil,
	})
}
