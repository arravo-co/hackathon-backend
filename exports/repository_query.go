package exports

import "time"

type CreateAdminAccountRepositoryDTO struct {
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Gender       string
	HackathonId  string
	Role         string
	PhoneNumber  string
	Status       string
}

type FilterGetManyAccountRepositories struct {
	Email_eq string
}

type UpdateAccountRepositoryFilter struct {
	Email       string `bson:"email"`
	PhoneNumber string `bson:"phone_number"`
}

type UpdateAdminAccountRepository struct {
	FirstName         string    `bson:"first_name,omitempty"`
	LastName          string    `bson:"last_name,omitempty"`
	Gender            string    `bson:"gender,omitempty"`
	State             string    `bson:"state,omitempty"`
	Bio               string    `bson:"bio,omitempty"`
	IsEmailVerified   bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt time.Time `bson:"is_email_verified_at,omitempty"`
	ProfilePictureUrl string    `bson:"profile_picture_url"`
}
