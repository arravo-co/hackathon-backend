package dtos

type BasicLoginDTO struct {
	Identifier string ` validate:"required" json:"identifier"`
	Password   string ` validate:"required" json:"password"`
}

type CompleteEmailVerificationDTO struct {
	Email string ` validate:"required,email" json:"email"`
	Token string ` validate:"required" json:"token"`
}

type ChangePasswordDTO struct {
	OldPassword string ` validate:"required" json:"old_password"`
	NewPassword string ` validate:"required" json:"new_password"`
}

type AuthUserInfoUpdateDTO struct {
	FirstName string `validate:"min=2,omitempty" json:"first_name"`
	LastName  string `validate:"min=2,omitempty" json:"last_name"`
	Email     string `validate:"email,omitempty" json:"email"`
	Gender    string `validate:"oneof=MALE FEMALE,omitempty" json:"gender"`
	State     string `json:"state,omitempty"`
}

type AuthParticipantInfoUpdateDTO struct {
	AuthUserInfoUpdateDTO
	GithubAddress   string `validate:"url,omitempty" json:"github_address"`
	LinkedInAddress string `validate:"url,omitempty" json:"linkedIn_address"`
	SolutionId      string `validate:"url,omitempty" json:"solution_id"`
}

type CompletePasswordRecoveryDTO struct {
	Email       string ` validate:"required,email" json:"email"`
	Token       string ` validate:"required" json:"token"`
	NewPassword string ` validate:"required" json:"new_password"`
}

type CreateSolutionData struct {
	Title       string `validate:"required" json:"title"`
	Description string `validate:"required" json:"description"`
	Objective   string `validate:"required" json:"objective"`
}

type UpdateSolutionData struct {
	Title       string ` json:"title"`
	Description string ` json:"description"`
	Objective   string ` json:"objective"`
}
