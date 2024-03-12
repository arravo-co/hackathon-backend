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

type ParticipantAccountCreatedEventHandler func(*ParticipantAccountCreatedEventData)
