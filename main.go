package main

import (
	"strconv"
	"strings"

	"github.com/arravoco/hackathon_backend/config"
	_ "github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/routes"
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
// @Server localhost:5000 Localhost
// @Server http://www.fake2.com Main
func main() {
	e := echo.New()
	port, err := config.GetPort()
	routes.StartAllRoutes(e)
	e.Logger.Info(port)
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Fatal(e.Start(getURL(port)))
}

func getURL(port int) string {
	return strings.Join([]string{"", strconv.Itoa(port)}, ":")
}
