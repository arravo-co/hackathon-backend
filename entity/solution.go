package entity

import "time"

type Solution struct {
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
