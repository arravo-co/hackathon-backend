package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

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

func GetSecretKey() string {
	return os.Getenv("SECRET_KEY")
}

func GetPort() int {
	portString := os.Getenv("PORT")
	if portString == "" {
		return 8080
	}
	port, _ := strconv.Atoi(portString)
	return port
}

func GetMongoDBURL() string {
	url := os.Getenv("MONGODB_URL")
	return url
}

func GetServerURL() string {
	url, notFound := os.LookupEnv("SERVER_URL")
	if !notFound {
		return strings.Join([]string{"http://localhost", strconv.Itoa(GetPort())}, ":")
	}
	return url
}

func GetFrontendURL() string {
	url := os.Getenv("FRONTEND_URL")
	return url
}

func GetRedisURL() string {
	url := os.Getenv("REDIS_URL")
	return url
}

func GetHackathonId() string {
	portString := os.Getenv("HACKATHON_ID")
	return portString
}

func GetRedisPort() (int, error) {
	portString := os.Getenv("PORT")
	if portString == "" {
		return 0, errors.New("port must be set")
	}
	return strconv.Atoi(portString)
}

func GetRedisHost() (string, error) {
	hostString := os.Getenv("REDIS_HOST")
	if hostString == "" {
		return "", errors.New("host must be set")
	}
	return hostString, nil
}

func GetRedisUsername() (string, error) {
	username := os.Getenv("REDIS_USERNAME")
	if username == "" {
		return "", errors.New("host must be set")
	}
	return username, nil
}

func GetRedisPassword() (string, error) {
	username := os.Getenv("REDIS_PASSWORD")
	if username == "" {
		return "", errors.New("password must be set")
	}
	return username, nil
}
