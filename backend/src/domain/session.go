package domain

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserSession struct {
	UserId         pgtype.UUID      `json:"userId"`
	State          UserSessionState `json:"state"`
	CurrentMatchId uuid.UUID        `json:"currentMatchId"` // Nil = not playing/watching
}

type UserSessionState string

const (
	Idle        UserSessionState = "IDLE"
	QueuingCaro UserSessionState = "QUEUING_CARO"
	PlayingCaro UserSessionState = "PLAYING_CARO"
)
