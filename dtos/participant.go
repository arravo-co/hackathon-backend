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

type RegisterNewIndividualParticipantDTO struct {
	FirstName       string `validate:"min=2,omitempty" json:"first_name"`
	LastName        string `validate:"min=2,omitempty" json:"last_name"`
	Email           string `validate:"email,omitempty" json:"email"`
	Password        string `validate:"min=7,omitempty" json:"password"`
	ConfirmPassword string `validate:"eqfield=Password,omitempty" json:"confirm_password"`
	Gender          string `validate:"oneof=MALE FEMALE,omitempty" json:"gender"`
	GithubAddress   string `validate:"url,omitempty" json:"github_address"`
	LinkedInAddress string `json:"linkedIn_address,omitempty"`
	State           string `json:"state,omitempty"`
}

type RegisterNewTeamParticipantDTO struct {
	FirstName           string   `validate:"min=2,omitempty" json:"first_name"`
	TeamName            string   `validate:"min=2" json:"team_name"`
	TeamLeadEmail       string   `validate:"email,omitempty" json:"team_lead_email"`
	CoParticipantEmails []string `validate:"max=8,unique,notblank,dive,email," json:"co_participant_emails"`
	Password            string   `validate:"min=7,omitempty" json:"password"`
	ConfirmPassword     string   `validate:"eqfield=Password,omitempty" json:"confirm_password"`
	GithubAddress       string   `validate:"url,omitempty" json:"github_address"`
}

type IndividualParticipantCreatedResponseDTO struct {
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	Email           string `json:"email,omitempty"`
	Gender          string `json:"gender,omitempty"`
	State           string `json:"state,omitempty"`
	GithubAddress   string `json:"github_address,omitempty"`
	LinkedInAddress string `json:"LinkedIn_address,omitempty"`
	Role            string `json:"role,omitempty"`
}
