package othermiddleware

import (
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
