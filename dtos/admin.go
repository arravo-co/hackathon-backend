package dtos

type CreateNewAdminDTO struct {
	Email           string `validate:"email" json:"email"`
	FirstName       string `validate:"min=2" json:"first_name"`
	LastName        string `validate:"min=2" json:"last_name"`
	PhoneNumber     string ` validate:"e164" json:"phone_number"`
	Password        string `validate:"min=2" json:"password"`
	ComfirmPassword string `validate:"min=2" json:"confirm_password"`
}

type CreateNewAdminByAuthAdminDTO struct {
	Email       string `validate:"email" json:"email"`
	FirstName   string `validate:"required" json:"first_name"`
	LastName    string `validate:"required" json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
}
