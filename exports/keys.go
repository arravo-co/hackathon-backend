package exports

type LinkPayload struct {
	InviteeEmail string `json:"invitee_email"`
	Token        string `json:"token"`
	InviterEmail string `json:"inviter_email"`
	TTL          int64  `json:"ttl"`
}
