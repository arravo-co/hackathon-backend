package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}
func GetPort() (int, error) {
	portString := os.Getenv("PORT")
	if portString == "" {
		return 0, errors.New("port must be set")
	}
	return strconv.Atoi(portString)
}
