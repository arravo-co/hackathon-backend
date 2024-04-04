package othermiddleware

import (
	"net/http"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Auth() echo.MiddlewareFunc {

	y := echojwt.WithConfig(echojwt.Config{
		// ...
		SigningKey: []byte(config.GetSecretKey()),
		SuccessHandler: func(c echo.Context) {
			//fmt.Printf("%s\n", "Successful authentication")
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(exports.MyJWTCustomClaims)
		},
	})
	return y
}

func AuthRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return Auth()(func(c echo.Context) error {
			tokenData := exports.GetPayload(c)
			if tokenData.Role != role {
				return c.JSON(http.StatusUnauthorized, struct {
					Message string `json:"message"`
					Code    int    `json:"code"`
				}{Message: "Unauthorized access",
					Code: http.StatusUnauthorized})
			}
			return next(c)
		})
	}
}
