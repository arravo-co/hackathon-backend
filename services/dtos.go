package services

type RegisterNewJudgeDTO struct {
	FirstName       string `validate:"min=2,required" json:"first_name"`
	LastName        string `validate:"min=2,required" json:"last_name"`
	Email           string `validate:"email,required" json:"email"`
	Password        string `validate:"min=7" json:"password"`
	ConfirmPassword string `validate:"eqfield=Password" json:"confirm_password"`
	Gender          string `validate:"oneof=MALE FEMALE" json:"gender"`
	State           string `json:"state"`
	Bio             string `json:"bio"`
}

type UpdateJudgeDTO struct {
	FirstName         string `validate:"omitempty,min=2"`
	LastName          string `validate:"omitempty,min=2"`
	Gender            string `validate:"omitempty,oneof=MALE FEMALE"`
	State             string
	Bio               string
	ProfilePictureUrl string
}
