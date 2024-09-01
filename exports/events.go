package exports

type EventData struct {
	EventName string
}

type TeamParticipantInfo struct {
	Email    string
	Password string
}
type AdminAccountCreatedEventData struct {
	EventData
	Email     string
	FirstName string
	LastName  string
	Gender    string
}

type AdminAccountCreatedByAdminEventData struct {
	EventData
	Email       string
	AdminName   string
	Gender      string
	Password    string
	InviterName string
}

type ParticipantAccountCreatedEventData struct {
	EventData
	FirstName        string
	LastName         string
	ParticipantEmail string
	ParticipantId    string
	TeamParticipants []TeamParticipantInfo
	TeamLeadEmail    string
	TeamName         string
	TeamRole         string
	ParticipantType  string
}

type ParticipantAccountCreatedEventHandler func(data *ParticipantAccountCreatedEventData, otherParams ...interface{})

type JudgeAccountCreatedByAdminEventData struct {
	EventData
	InviteeEmail string
	InviterName  string
	JudgeName    string
	JudgeEmail   string
	Password     string
}

type JudgeAccountCreatedByAdminEventHandler func(input *JudgeAccountCreatedByAdminEventData, otherParams ...interface{})

type AdminAccountCreatedEventHandler func(data *AdminAccountCreatedEventData, otherParams ...interface{})

type AdminAccountCreatedByAdminEventHandler func(data *AdminAccountCreatedByAdminEventData, otherParams ...interface{})
