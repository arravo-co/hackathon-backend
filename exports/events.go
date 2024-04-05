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
	Email     string
	FirstName string
	LastName  string
	Gender    string
	Password  string
}

type ParticipantAccountCreatedEventData struct {
	EventData
	FirstName        string
	LastName         string
	ParticipantEmail string
	TeamParticipants []TeamParticipantInfo
	TeamLeadEmail    string
	TeamName         string
	TeamRole         string
	ParticipantType  string
}

type ParticipantAccountCreatedEventHandler func(data *ParticipantAccountCreatedEventData, otherParams ...interface{})

type JudgeAccountCreatedEventData struct {
	EventData
	FirstName  string
	LastName   string
	JudgeEmail string
}

type JudgeAccountCreatedEventHandler func(input *JudgeAccountCreatedEventData, otherParams ...interface{})

type AdminAccountCreatedEventHandler func(data *AdminAccountCreatedEventData, otherParams ...interface{})

type AdminAccountCreatedByAdminEventHandler func(data *AdminAccountCreatedByAdminEventData, otherParams ...interface{})
