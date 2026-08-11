package ws

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UserStateStatus string

const (
	Idle    UserStateStatus = "IDLE"
	InQueue UserStateStatus = "IN_QUEUE"
)

type UserState struct {
	Id     pgtype.UUID     `json:"id"`
	Status UserStateStatus `json:"status"`
}
