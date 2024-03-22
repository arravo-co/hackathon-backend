package config

import (
	"errors"
	"fmt"
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

func GetPort() (int, error) {
	portString := os.Getenv("PORT")
	if portString == "" {
		return 0, errors.New("port must be set")
	}
	return strconv.Atoi(portString)
}

func GetRedisURL() (string, error) {
	port, _ := GetRedisPort()
	host, _ := GetRedisHost()
	password, _ := GetRedisPassword()
	username, _ := GetRedisUsername()
	portStr := strconv.Itoa(port)
	start := "redis://"
	cred := strings.Join([]string{username, password}, ":")
	cred_host := strings.Join([]string{cred, host}, "@")
	cred_host_port := strings.Join([]string{cred_host, portStr}, ":")
	url := strings.Join([]string{start, cred_host_port}, "")
	fmt.Printf("%s\n", url)
	return url, nil
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
