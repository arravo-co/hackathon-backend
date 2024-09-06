package exports

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParticipantTeamMembersWithAccountsAggregateDocument struct {
	Id               primitive.ObjectID                             `bson:"_id"`
	ParticipantId    string                                         `bson:"participant_id"`
	HackathonId      string                                         `bson:"hackathon_id"`
	Type             string                                         `bson:"type,omitempty"`
	TeamLeadEmail    string                                         `bson:"team_lead_email,omitempty"`
	SolutionId       string                                         `bson:"solution_id,omitempty"`
	Solution         ParticipantDocumentParticipantSelectedSolution `bson:"solution"`
	TeamName         string                                         `bson:"team_name,omitempty"`
	ParticipantEmail string                                         `bson:"participant_email,omitempty"`
	GithubAddress    string                                         `bson:"github_address,omitempty"`
	InviteList       []ParticipantDocumentTeamInviteInfo            `bson:"invite_list,omitempty"`
	ReviewRanking    int                                            `bson:"review_ranking,omitempty"`
	Status           string                                         `bson:"status,omitempty"`
	CreatedAt        time.Time                                      `bson:"created_at,omitempty"`
	UpdatedAt        time.Time                                      `bson:"updated_at,omitempty"`

	TeamLeadFirstName string `bson:"team_lead_first_name,omitempty"`
	TeamLeadLastName  string `bson:"team_lead_last_name,omitempty"`
	TeamLeadGender    string `bson:"team_lead_gender,omitempty"`
	TeamLeadAccountId string `bson:"team_lead_account_id,omitempty"`

	TeamLeadInfo struct {
		HackathonId   string    `bson:"team_lead_hackathon_id"`
		AccountId     string    `bson:"id,omitempty"`
		Email         string    `bson:"email"`
		FirstName     string    `bson:"first_name"`
		LastName      string    `bson:"last_name"`
		Gender        string    `bson:"gender"`
		PhoneNumber   string    `bson:"phone_number"`
		Skillset      []string  `bson:"skillset"`
		AccountStatus string    `bson:"status"`
		AccountRole   string    `bson:"account_role"`
		TeamRole      string    `bson:"team_role"`
		PasswordHash  string    `bson:"password_hash"`
		State         string    `bson:"state"`
		CreatedAt     time.Time `bson:"created_at"`
		UpdateAt      time.Time `bson:"update_at"`
	} `bson:"team_lead_info,omitempty"`
	CoParticipants []struct {
		HackathonId   string    `bson:"team_lead_hackathon_id"`
		AccountId     string    `bson:"id,omitempty"`
		Email         string    `bson:"email"`
		FirstName     string    `bson:"first_name"`
		LastName      string    `bson:"last_name"`
		Gender        string    `bson:"gender"`
		PhoneNumber   string    `bson:"phone_number"`
		Skillset      []string  `bson:"skillset"`
		AccountStatus string    `bson:"status"`
		AccountRole   string    `bson:"account_role"`
		TeamRole      string    `bson:"team_role"`
		PasswordHash  string    `bson:"password_hash"`
		State         string    `bson:"state"`
		CreatedAt     time.Time `bson:"created_at"`
		UpdateAt      time.Time `bson:"update_at"`
	} `bson:"co_participants,omitempty"`
}

type GetParticipantsWithAccountsAggregateFilterOpts struct {
	ParticipantId            *string
	ParticipantStatus        *string `validate:"omitempty, oneof UNREVIEWED REVIEWED AI_RANKED "`
	ParticipantType          *string `validate:"omitempty, oneof TEAM "`
	ReviewRanking_Eq         *int
	ReviewRanking_Top        *int
	Solution_Like            *string
	Limit                    *int
	SortByReviewRanking_Asc  *bool
	SortByReviewRanking_Desc *bool
}
