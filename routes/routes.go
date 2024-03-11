package routes

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/arravoco/hackathon_backend/docs"
)

var participantsRoutes *echo.Group
var validate *validator.Validate

// @title			Hackathons API
// @version		1.0
// @description	This is the documentation website for all Arravo hackathons
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.email	appdev@arravo.co
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:5000
// @BasePath		/api
func StartAllRoutes(e *echo.Echo) {
	validate = validator.New()
	e.GET("/api/docs/*", echoSwagger.WrapHandler) // default
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	participantsRoutes = e.Group("/participants")
	participantsRoutes.POST("/", RegisterParticipant)
}
