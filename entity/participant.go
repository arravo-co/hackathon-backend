package entity

import (
	"time"
)

// AddMemberToParticipatingTeam
type Participant struct {
	//FirstName           string                               `json:"first_name"`
	//LastName            string                               `json:"last_name"`
	//Email               string                               `json:"email"`
	//Gender              string                               `json:"gender"`
	//State               string                               `json:"state"`
	//Age                 int                                  `json:"age"`
	//DOB                 time.Time                            `json:"dob"`
	//AccountRole         string                               `json:"role"`
	ParticipantId       string                               `json:"participant_id"`
	TeamLeadEmail       string                               `json:"team_lead_email"`
	TeamName            string                               `json:"team_name"`
	TeamRole            string                               `json:"team_role"`
	HackathonId         string                               `json:"hackathon_id"`
	ParticipantType     string                               `json:"type"`
	CoParticipants      []ParticipantEntityCoParticipantInfo `json:"co_participants"`
	ParticipantEmail    string                               `json:"participant_email"`
	InviteList          []InviteInfo                         `json:"invite_list"`
	AccountStatus       string                               `json:"account_status"`
	ParticipationStatus string                               `json:"participation_status"`
	ReviewRanking       int                                  `json:"review_ranking"`
	TeamLeadInfo        ParticipantEntityTeamLeadInfo
	Solution            *ParticipantEntitySelectedSolution `json:"solution"`
	CreatedAt           time.Time                          `json:"created_at"`
	UpdatedAt           time.Time                          `json:"updated_at"`
}

type ParticipantEntityCoParticipantInfo struct {
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Email               string    `json:"email"`
	Gender              string    `json:"gender"`
	State               string    `json:"state"`
	Age                 int       `json:"age"`
	DOB                 time.Time `json:"dob"`
	AccountStatus       string    `json:"account_status"`
	AccountRole         string    `json:"account_role"`
	TeamRole            string    `json:"team_role"`
	ParticipantId       string    `json:"participant_id"`
	HackathonId         string    `json:"hackathon_id"`
	Skillset            []string  `json:"skillset"`
	PhoneNumber         string    `json:"phone_number"`
	EmploymentStatus    string    `json:"employment_status"`
	ExperienceLevel     string    `json:"experience_level"`
	Motivation          string    `json:"motivation"`
	HackathonExperience string    `json:"hackathon_experience"`
	YearsOfExperience   int       `json:"years_of_experience"`
	FieldOfStudy        string    `json:"field_of_study"`
	PreviousProjects    []string  `json:"previous_projects"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type InviteInfo struct {
	Email     string    `bson:"email,omitempty" json:"email,omitempty"`
	Time      time.Time `bson:"time,omitempty" json:"time,omitempty"`
	InviterId string    `bson:"inviter_id,omitempty" json:"inviter_id,omitempty"`
}
type ParticipantEntityTeamLeadInfo struct {
	HackathonId   string
	AccountId     string
	Email         string
	FirstName     string
	LastName      string
	Gender        string
	PhoneNumber   string
	Skillset      []string
	AccountStatus string
	AccountRole   string
	TeamRole      string
	PasswordHash  string
	State         string
	CreatedAt     string
	UpdateAt      string
}

type ParticipantEntitySelectedSolution struct {
	Id               string    `json:"id"`
	HackathonId      string    `json:"hackathon_id"`
	Title            string    `json:"name"`
	Description      string    `json:"description"`
	Objective        string    `json:"objective"`
	CreatorId        string    `json:"creator_id"`
	SolutionImageUrl string    `json:"solution_image_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
