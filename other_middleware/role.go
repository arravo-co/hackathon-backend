package othermiddleware

import (
	"net/http"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/labstack/echo/v4"
)

func CheckIfIsRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenData := c.Get("user").(exports.Payload)
			if tokenData.Role != role {
				return c.JSON(http.StatusUnauthorized, struct {
					Message string `json:"message"`
					Code    int    `json:"code"`
				}{Message: "Unauthorized access",
					Code: http.StatusUnauthorized})
			}
			return next(c)
		}
	}
}
