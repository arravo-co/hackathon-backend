package dtos

type RegisterNewJudgeDTO struct {
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
