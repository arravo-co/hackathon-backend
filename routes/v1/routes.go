package routes_v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	//_ "github.com/arravoco/hackathon_backend/docs"
	othermiddleware "github.com/arravoco/hackathon_backend/other_middleware"
)

type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var participantsRoutes *echo.Group
var authRoutes *echo.Group
var judgesRoutes *echo.Group
var validate *validator.Validate

// @Version 1.0.0
func StartAllRoutes(e *echo.Echo) {
	validate = validator.New()
	e.Renderer = t
	e.GET("/hello", Hello)
	e.Static("/api/v1/docs", "static/docs")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	api := e.Group("/api/v1")
	setupAuthRoutes(api)
	setupJudgesRoutes(api)
	setupParticipantsRoutes(api)
	setupScoreRoutes(api)
	setupAdminsRoutes(api)
}

func setupAdminsRoutes(api *echo.Group) {
	adminsRoutes := api.Group("/admin")
	adminsRoutes.POST("/register_admin", RegisterAnotherAdmin, othermiddleware.AuthRole([]string{"ADMIN", "SUPER_ADMIN"}))
	adminsRoutes.POST("/register_judge", RegisterJudgeByAdmin, othermiddleware.AuthRole([]string{"ADMIN", "SUPER_ADMIN"}))
}

func setupParticipantsRoutes(api *echo.Group) {
	participantsRoutes = api.Group("/participants")
	participantsRoutes.POST("", RegisterParticipant)
	participantsRoutes.POST("/participants/:participantId/members ", RegisterNewTeamMember)
	participantsRoutes.POST("/:participantId/invite", InviteMemberToTeam, othermiddleware.AuthRole([]string{"PARTICIPANT"}))
}

func setupJudgesRoutes(api *echo.Group) {
	judgesRoutes = api.Group("/judges")
	judgesRoutes.POST("", RegisterJudge)
}

func setupAuthRoutes(api *echo.Group) {
	authRoutes = api.Group("/auth")
	authRoutes.POST("/login", BasicLogin)
	authRoutes.GET("/verification/email/initiation", InitiateEmailVerification)
	authRoutes.GET("/verification/email/completion", CompleteEmailVerificationViaGet)
	authRoutes.POST("/verification/email/completion", CompleteEmailVerification)
	authRoutes.POST("/password/change", ChangePassword, othermiddleware.Auth())
	authRoutes.GET("/password/recovery/initiation", InitiatePasswordRecovery)
	authRoutes.POST("/password/recovery/completion", CompletePasswordRecovery)
	authRoutes.GET("/me", GetAuthUserInfo, othermiddleware.Auth())
	authRoutes.PUT("/me", UpdateAuthUserInfo, othermiddleware.Auth())
	authRoutes.GET("/team/invite", ValidateTeamInviteLink)
	authRoutes.GET("/password/recovery/link/verification", ValidatePasswordRecoveryLink)
}

func setupScoreRoutes(api *echo.Group) {
	authRoutes = api.Group("/score_management")
	authRoutes.GET("/create_score_record", ScoreParticipant, othermiddleware.CheckIfIsRole("JUDGE"))
	authRoutes.GET("/score_participant", ScoreParticipant, othermiddleware.CheckIfIsRole("JUDGE"))
}
