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

func GetResendAPIKey() string {
	return os.Getenv("RESEND_API_KEY")
}

func GetResendFromEmail() string {
	return os.Getenv("RESEND_FROM_EMAIL")
}

func GetPort() (int, error) {
	portString := os.Getenv("PORT")
	if portString == "" {
		return 0, errors.New("port must be set")
	}
	return strconv.Atoi(portString)
}
