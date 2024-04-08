package main

import (
	"strconv"
	"strings"

	"github.com/arravoco/hackathon_backend/config"
	_ "github.com/arravoco/hackathon_backend/db"
	_ "github.com/arravoco/hackathon_backend/jobs"
	routes_v1 "github.com/arravoco/hackathon_backend/routes/v1"
	"github.com/arravoco/hackathon_backend/security"
	"github.com/labstack/echo/v4"
)

// @Version 1.0.0
// @Title Hackathon Backend API
// @Description API usually works as expected. But sometimes its not true.
// @ContactName David Alabi
// @ContactEmail appdev@arravo.co
// @ContactURL http://arravo.co/contact
// @TermsOfServiceUrl http://arravo.co/contact
// @LicenseName MIT
// @LicenseURL https://en.wikipedia.org/wiki/MIT_License
// @Server http://localhost:5000 Localhost
// @Server https://hackathon-backend-2cvk.onrender.com Development
func main() {
	security.GenerateKeys()
	e := echo.New()
	port := config.GetPort()
	routes_v1.StartAllRoutes(e)
	e.Logger.Info(port)
	e.Logger.Fatal(e.Start(getURL(port)))
}

func getURL(port int) string {
	return strings.Join([]string{"", strconv.Itoa(port)}, ":")
}
