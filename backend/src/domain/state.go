package domain

type UserState string

const (
	Idle UserState = "IDLE"
	Queuing UserState = "QUEUING"
	Playing UserState = "PLAYING"
	Watching UserState = "WATCHING"
)