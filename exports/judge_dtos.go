package exports

type RegisterNewJudgeDTO struct {
	FirstName string `validate:"min=2,required"`
	LastName  string `validate:"min=2,required"`
	Email     string `validate:"email,required"`
	Password  string `validate:"min=7"`
	Gender    string `validate:"oneof=MALE FEMALE"`
	State     string
	Bio       string `validate:"min=1"`
}

type UpdateJudgeDTO struct {
	FirstName         string `validate:"omitempty,min=2"`
	LastName          string `validate:"omitempty,min=2"`
	Gender            string `validate:"omitempty, oneof=MALE FEMALE"`
	State             string
	Bio               string
	ProfilePictureUrl string
}
