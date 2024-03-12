package main

import (
	"strconv"
	"strings"

	"github.com/arravoco/hackathon_backend/config"
	_ "github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/routes"
	"github.com/labstack/echo/v4"
)

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
