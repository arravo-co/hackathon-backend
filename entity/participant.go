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
	ParticipantId string `json:"participant_id"`
	TeamLeadEmail string `json:"team_lead_email"`
	TeamName      string `json:"team_name"`
	//TeamRole            string                               `json:"team_role"`
	HackathonId         string                               `json:"hackathon_id"`
	ParticipantType     string                               `json:"type"`
	CoParticipants      []ParticipantEntityCoParticipantInfo `json:"co_participants"`
	ParticipantEmail    string                               `json:"participant_email"`
	InviteList          []InviteInfo                         `json:"invite_list"`
	AccountStatus       string                               `json:"account_status"`
	ParticipationStatus string                               `json:"participation_status"`
	ReviewRanking       int                                  `json:"review_ranking"`
	TeamLeadInfo        ParticipantEntityTeamLeadInfo        `json:"team_leader_info"`
	Solution            *ParticipantEntitySelectedSolution   `json:"solution"`
	CreatedAt           time.Time                            `json:"created_at"`
	UpdatedAt           time.Time                            `json:"updated_at"`
}

type ParticipantEntityCoParticipantInfo struct {
	AccountId           string    `json:"account_id"`
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
	Email     string    `json:"email,omitempty" json:"email,omitempty"`
	Time      time.Time `json:"time,omitempty" json:"time,omitempty"`
	InviterId string    `json:"inviter_id,omitempty" json:"inviter_id,omitempty"`
}
type ParticipantEntityTeamLeadInfo struct {
	HackathonId         string   `json:"hackathon_id,omitempty"`
	AccountId           string   `json:"account_id,omitempty"`
	Email               string   `json:"email,omitempty"`
	FirstName           string   `json:"first_name,omitempty"`
	LastName            string   `json:"last_name,omitempty"`
	Gender              string   `json:"gender,omitempty"`
	PhoneNumber         string   `json:"phone_number,omitempty"`
	Skillset            []string `json:"skillset,omitempty"`
	AccountStatus       string   `json:"account_status,omitempty"`
	AccountRole         string   `json:"account_role,omitempty"`
	TeamRole            string   `json:"team_role,omitempty"`
	passwordHash        string
	State               string    `json:"state,omitempty"`
	LinkedInAddress     string    `json:"linkedin_address,omitempty"`
	ParticipantId       string    `json:"participant_id,omitempty"`
	DOB                 time.Time `json:"dob,omitempty"`
	EmploymentStatus    string    `json:"employment_status,omitempty"`
	ExperienceLevel     string    `json:"experience_level,omitempty"`
	Motivation          string    `json:"motivation,omitempty"`
	HackathonExperience string    `json:"hackathon_experience,omitempty"`
	YearsOfExperience   int       `json:"years_of_experience,omitempty"`
	FieldOfStudy        string    `json:"field_of_study,omitempty"`
	PreviousProjects    []string  `json:"previous_projects,omitempty"`
	IsEmailVerified     bool      `json:"is_email_verified,omitempty"`
	IsEmailVerifiedAt   time.Time `json:"is_email_verified_at,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdateAt            time.Time `json:"updated_at,omitempty"`
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
