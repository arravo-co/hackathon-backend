package exports

import "time"

type QueuePayload struct {
	RecipientQueueName    string `json:"recipient_queue_name"`
	TaskId                string `json:"task_id,omitempty"`
	TaskName              string `json:"task_name,omitempty"`
	ParentTaskID          string `json:"parent_task_id,omitempty"`
	RequestingJobTaskID   string `json:"requesting_job_task_id,omitempty"`
	RequestingJobTaskName string `json:"requesting_job_task_name,omitempty"`
	RespondingJobTaskID   string `json:"responding_job_task_id,omitempty"`
	RespondingJobTaskName string `json:"responding_job_task_name,omitempty"`
	NextJob               string `json:"next_job,omitempty"`
	ParentJob             string `json:"parent_job,omitempty"`
	TriggerNextJob        string `json:"trigger_next_job,omitempty"`
	Direction             string `json:"direction,omitempty"` // REQUEST RESPONSE COMMAND
}

type UploadPicQueuePayload struct {
	QueuePayload
	Email    string `json:"account_email"`
	FilePath string `json:"file_path"`
}

type CoordinateSendParticipantEmailPayload struct {
	QueuePayload
	Email   string `json:"email"`
	TTL     int    `json:"ttl"`
	TTLType string `json:"ttl_type"`
}

type SendWelcomeAndEmailVerificationTokenJobQueuePayload struct {
	QueuePayload
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateEmailTokenJobResponsePayload struct {
	QueuePayload
	TokenId        interface{} `json:"token_id"`
	Token          string      `json:"token"`
	TokenType      string      `json:"token_type"`
	TokenTypeValue string      `json:"token_type_value"`
	Scope          string      `json:"scope"`
	TTL            time.Time   `json:"ttl"`
	Status         string      `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type SendWelcomeAndEmailVerificationQueueRequestPayload struct {
	QueuePayload
	Email        string `json:"email"`
	TeamName     string `json:"team_name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	TeamRole     string `json:"team_role"` // TEAM_LEAD TEAM_MEMBER
	TeamLeadName string `json:"team_leader_name"`
}
