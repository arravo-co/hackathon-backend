package entity

type Judge struct {
	FirstName       string
	LastName        string
	Email           string
	PasswordHash    string
	Gender          string
	State           string
	GithubAddress   string
	LinkedInAddress string
	Role            string
}

func (judge *Judge) Register() {

}
