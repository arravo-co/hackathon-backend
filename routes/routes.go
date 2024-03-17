package routes

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

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

func StartAllRoutes(e *echo.Echo) {
	validate = validator.New()
	e.Renderer = t
	e.GET("/hello", Hello)
	e.GET("/api/docs/*", echoSwagger.WrapHandler, func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			var re = regexp.MustCompile(`^(.*/)([^?].*)?[?|.]*$`)
			matches := re.FindStringSubmatch(c.Request().RequestURI)
			path := matches[2]
			fmt.Printf("%+v\n", path)
			return next(c)
		}
	}) // default
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	api := e.Group("/api")
	setupAuthRoutes(api)
	setupJudgesRoutes(api)
	setupParticipantsRoutes(api)
}

func setupParticipantsRoutes(api *echo.Group) {
	participantsRoutes = api.Group("/participants")
	participantsRoutes.POST("", RegisterParticipant)
}

func setupJudgesRoutes(api *echo.Group) {
	judgesRoutes = api.Group("/judges")
	judgesRoutes.POST("", RegisterJudge)
}

func setupAuthRoutes(api *echo.Group) {
	authRoutes = api.Group("/auth")
	authRoutes.POST("/login", BasicLogin)
	authRoutes.GET("/verification/email/initiation", InitiateEmailVerification)
	authRoutes.POST("/verification/email/completion", CompleteEmailVerification)
	authRoutes.POST("/password/change", ChangePassword, othermiddleware.Auth())
	authRoutes.POST("/password/recovery/initiation", InitiatePasswordRecovery)
	authRoutes.POST("/password/recovery/completion", ChangePassword)
	authRoutes.GET("/me", GetAuthUserInfo, othermiddleware.Auth())
	authRoutes.PUT("/me", UpdateAuthUserInfo, othermiddleware.Auth())
}
