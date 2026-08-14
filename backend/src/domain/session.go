package domain

import (
	"backend/src/shared"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserSession struct {
	UserId         pgtype.UUID      `json:"userId"`
	WsConn         *websocket.Conn  `json:"-"`
	State          UserSessionState `json:"state"`
	CurrentGame    shared.GameCode  `json:"currentGame"`    // None = not playing/watching
	CurrentMatchId string           `json:"currentMatchId"` // "" = not playing/watching
}

type UserSessionState string

const (
	Idle     UserSessionState = "IDLE"
	Queuing  UserSessionState = "QUEUING"
	Playing  UserSessionState = "PLAYING"
	Watching UserSessionState = "WATCHING"
)

