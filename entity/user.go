package entity

type User interface {
	Register()
	InitiatePasswordChange()
	CompletePasswordChange()
	InitiatePasswordRecovery()
	CompletePasswordRecovery()
}
