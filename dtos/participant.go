package dtos

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type GenderEnum int

const (
	Male GenderEnum = iota
	Female
	InvalidGender
)

func (g GenderEnum) String() string {
	return [...]string{"MALE", "FEMALE"}[g]
}

func ValidateGender(fl validator.FieldLevel) bool {
	fmt.Printf("%#v", fl)
	return fl.Field().String() != ""
}

type RegisterNewParticipantDTO struct {
	FirstName       string `validate:"min=2" json:"first_name"`
	LastName        string `validate:"min=2" json:"last_name"`
	Email           string `validate:"email" json:"email"`
	Password        string `validate:"min=7" json:"password"`
	ConfirmPassword string `validate:"eqfield=Password" json:"confirm_password"`
	Gender          string `validate:"oneof=MALE FEMALE" json:"gender"`
	GithubAddress   string `validate:"url" json:"github_address"`
	LinkedInAddress string `json:"linkedIn_address"`
	State           string `json:"state"`
}

type ParticipantAddResponseDTO struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Gender          string `json:"gender"`
	State           string `json:"state"`
	GithubAddress   string `json:"github_address"`
	LinkedInAddress string `json:"LinkedIn_address"`
	Role            string `json:"role"`
}
