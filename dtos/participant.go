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

type InviteToTeamData struct {
	Email string ` validate:"email" json:"email"`
	Role  string `validate:"oneof= TEAM_MEMBER" json:"role"`
}

type CompleteNewTeamMemberRegistrationDTO struct {
	FirstName       string   `validate:"min=2" json:"first_name"`
	LastName        string   `validate:"min=2" json:"last_name"`
	Email           string   `validate:"email" json:"email"`
	Password        string   `validate:"min=7" json:"password"`
	PhoneNumber     string   `validate:"e164" json:"phone_number"`
	ConfirmPassword string   `validate:"eqfield=Password" json:"confirm_password"`
	Gender          string   `validate:"oneof=MALE FEMALE" json:"gender"`
	Skillset        []string `validate:"min=1" json:"skillset"`
	State           string   `validate:"min=3" json:"state"`
	DOB             string   ` json:"dob"`
	TeamLeadEmail   string   `json:"team_lead_email"`
	HackathonId     string   `validate:"min=2" json:"hackathon_id"`
}

type RegisterNewParticipantDTO struct {
	FirstName       string   `validate:"min=2" json:"first_name"`
	LastName        string   `validate:"min=2" json:"last_name"`
	Email           string   `validate:"email" json:"email"`
	Password        string   `validate:"min=7" json:"password"`
	PhoneNumber     string   `validate:"e164" json:"phone_number"`
	ConfirmPassword string   `validate:"eqfield=Password" json:"confirm_password"`
	Gender          string   `validate:"oneof=MALE FEMALE" json:"gender"`
	Skillset        []string `validate:"min=1" json:"skillset"`
	State           string   `validate:"min=3" json:"state"`
	Type            string   `validate:"oneof=INDIVIDUAL TEAM" json:"type"`
	TeamSize        int      `json:"team_size"`
	DOB             string   ` json:"dob"`
	TeamName        string   `validate:"omitempty" json:"team_name"`
}

type ParticipantCreatedResponseDTO struct {
	FirstName       string   `validate:"min=2,omitempty" json:"first_name"`
	LastName        string   `validate:"min=2,omitempty" json:"last_name"`
	Email           string   `validate:"email,omitempty" json:"email"`
	Password        string   `validate:"min=7,omitempty" json:"password"`
	PhoneNumber     string   `json:"phone_number"`
	ConfirmPassword string   `validate:"eqfield=Password,omitempty" json:"confirm_password"`
	Gender          string   `validate:"oneof=MALE FEMALE,omitempty" json:"gender"`
	Skillset        []string `json:"skillset,omitempty"`
	State           string   `json:"state,omitempty"`
	Type            string   `json:"type"`
	TeamSize        int      `json:"team_size,omitempty"`
	Age             int      `json:"age,omitempty"`
}
