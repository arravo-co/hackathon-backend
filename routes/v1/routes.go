package routes_v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	//_ "github.com/arravoco/hackathon_backend/docs"
	othermiddleware "github.com/arravoco/hackathon_backend/other_middleware"
)

type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var participantsRoutes *echo.Group
var solutionRoutes *echo.Group
var scoreRoutes *echo.Group
var judgesRoutes *echo.Group
var authRoutes *echo.Group
var validate *validator.Validate

// @Version 1.0.0
func StartAllRoutes(e *echo.Echo) {
	validate = validator.New()
	e.Renderer = t
	e.GET("/hello", Hello)
	e.Static("/api/v1/docs", "docs/v1/html")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	rrr := otelecho.Middleware("")
	api := e.Group("/api/v1" /**/, rrr)
	setupSolutionRoutes(api)
	setupAuthRoutes(api)
	setupJudgesRoutes(api)
	setupParticipantsRoutes(api)
	setupScoreRoutes(api)
	setupAdminsRoutes(api)
}

func setupAdminsRoutes(api *echo.Group) {
	adminsRoutes := api.Group("/admin")
	adminsRoutes.POST("/updates/participants", UpdateAuthParticipantInfo, othermiddleware.AuthRole([]string{"ADMIN", "SUPER_ADMIN"}))
	adminsRoutes.POST("/register_admin", RegisterAnotherAdmin, othermiddleware.AuthRole([]string{"ADMIN", "SUPER_ADMIN"}))
	adminsRoutes.POST("/register_judge", RegisterJudgeByAdmin, othermiddleware.AuthRole([]string{"ADMIN", "SUPER_ADMIN"}))
	adminsRoutes.POST("", RegisterAdmin, middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"localhost"},
		AllowMethods: []string{"POST"},
	}))
}

func setupParticipantsRoutes(api *echo.Group) {
	participantsRoutes = api.Group("/participants")
	participantsRoutes.POST("/:participantId/members", CompleteNewTeamMemberRegistration)
	participantsRoutes.GET("/:participantId", GetParticipant)
	participantsRoutes.POST("", RegisterParticipant)
	participantsRoutes.GET("", GetParticipants)
	//participantsRoutes.POST("/invite", InviteMemberToTeam, othermiddleware.AuthRole([]string{"PARTICIPANT"}))
}

func setupJudgesRoutes(api *echo.Group) {
	judgesRoutes = api.Group("/judges")
	judgesRoutes.GET("/:email", GetJudgeByEmailAddress)
	judgesRoutes.GET("", GetJudges)
	judgesRoutes.POST("", RegisterJudge, othermiddleware.AuthRole([]string{"ADMIN", "JUDGE"}))
}

func setupAuthRoutes(api *echo.Group) {
	authRoutes = api.Group("/auth")
	authRoutes.POST("/login", BasicLogin)
	authRoutes.GET("/verification/email/initiation", InitiateEmailVerification)
	authRoutes.GET("/verification/email/completion", CompleteEmailVerificationViaGet)
	authRoutes.POST("/verification/email/completion", CompleteEmailVerification)
	authRoutes.GET("/password/recovery/initiation", InitiatePasswordRecovery)
	authRoutes.POST("/password/recovery/completion", CompletePasswordRecovery)
	authRoutes.POST("/password/change", ChangePassword, othermiddleware.Auth())
	authRoutes.GET("/me", GetAuthUserInfo, othermiddleware.Auth())
	authRoutes.PUT("/me", UpdateAuthParticipantInfo, othermiddleware.Auth())
	authRoutes.POST("/me/team/invite", InviteMemberToTeam, othermiddleware.AuthRole([]string{"PARTICIPANT"}))
	authRoutes.POST("/me/team/solution", ChooseSolutionForMyTeam, othermiddleware.AuthRole([]string{"PARTICIPANT"}))
	authRoutes.GET("/team/invite", ValidateTeamInviteLink)
	authRoutes.GET("/me/team", GetMyTeamMembersInfo, othermiddleware.AuthRole([]string{"PARTICIPANT"}))

	authRoutes.DELETE("/me/team/:team_member_email", RemoveMemberFromMyTeam, othermiddleware.AuthRole([]string{"PARTICIPANT"}))
	authRoutes.GET("/password/recovery/link/verification", ValidatePasswordRecoveryLink)
}

func setupScoreRoutes(api *echo.Group) {
	scoreRoutes = api.Group("/score_management")
	scoreRoutes.GET("/create_score_record", ScoreParticipant, othermiddleware.CheckIfIsRole("JUDGE"))
	scoreRoutes.GET("/score_participant", ScoreParticipant, othermiddleware.CheckIfIsRole("JUDGE"))
}

func setupSolutionRoutes(api *echo.Group) {
	solutionRoutes = api.Group("/solutions")
	solutionRoutes.GET("/:id", GetSolutionDataById)
	solutionRoutes.PUT("/:id", UpdateSolutionDataById, othermiddleware.Auth(), othermiddleware.CheckIfIsRole("ADMIN"))
	solutionRoutes.GET("", GetSolutionsData)
	solutionRoutes.POST("", CreateSolution, othermiddleware.Auth(), othermiddleware.CheckIfIsRole("ADMIN"))
}
