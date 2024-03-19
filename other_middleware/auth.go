package othermiddleware

import (
	"fmt"

	"github.com/arravoco/hackathon_backend/config"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Auth() echo.MiddlewareFunc {

	y := echojwt.WithConfig(echojwt.Config{
		// ...
		SigningKey: []byte(config.GetSecretKey()),
		SuccessHandler: func(c echo.Context) {
			fmt.Printf("%s", "Successful authentication")
		},
	})
	var v echo.HandlerFunc = func(c echo.Context) error {
		fmt.Printf("in v middleware")
		return nil
	}
	var g echo.MiddlewareFunc = func(next echo.HandlerFunc) echo.HandlerFunc {
		return y(v)
	}
	return g
}
