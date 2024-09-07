package entity

import (
	"time"
)

// AddMemberToParticipatingTeam
type Judge struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	State     string `json:"state"`
	Age       int    `json:"age"`
	//DOB                 time.Time                        `json:"dob"`
	Role              string    `json:"role"`
	HackathonId       string    `json:"hackathon_id"`
	Status            string    `json:"account_status"`
	PhoneNumber       string    `json:"phone_number"`
	Bio               string    `json:"bio"`
	ProfilePictureUrl string    `json:"profile_picture_url"`
	IsEmailVerified   bool      `json:"is_email_verified"`
	EmailVerifiedAt   time.Time `json:"email_verified_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func GetNewJudgeByEmail(email string) (*Judge, error) {
	judge := Judge{}
	return &judge, nil
}
