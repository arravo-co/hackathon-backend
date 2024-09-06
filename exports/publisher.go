package exports

import "time"

type PublisherConfig struct {
	RabbitMQExchange string
	RabbitMQKey      string
}

type ParticipantInvitedEmailPayload struct {
	Email       string    `json:"email"`
	TTL         time.Time `json:"ttl"`
	Link        string    `json:"link"`
	Token       string    `json:"token"`
	InviterName string    `json:"inviter_name"`
}
