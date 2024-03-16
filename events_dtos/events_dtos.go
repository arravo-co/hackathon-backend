package eventsdtos

type EventData struct {
	EventName string
}
type ParticipantAccountCreatedEventData struct {
	EventData
	FirstName        string
	LastName         string
	ParticipantEmail string
}

type ParticipantAccountCreatedEventHandler func(data *ParticipantAccountCreatedEventData, otherParams ...interface{})

type JudgeAccountCreatedEventData struct {
	EventData
	FirstName  string
	LastName   string
	JudgeEmail string
}

type JudgeAccountCreatedEventHandler func(*JudgeAccountCreatedEventData)
