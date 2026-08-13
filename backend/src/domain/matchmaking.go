package domain

import (
	"backend/src/shared"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type MmQueue []pgtype.UUID

func UserJoinCaroQueue(session *UserSession, app *AppState) error {
	if session.State != Idle {
		return fmt.Errorf("invalid state")
	}

	app.CaroQueue = append(app.CaroQueue, session.UserId)
	session.State = Queuing
	session.CurrentGame = shared.Caro

	return nil
}

func UserLeaveCaroQueue(session *UserSession, app *AppState) error {
	if !(session.State == Queuing && session.CurrentGame == shared.Caro) {
		return fmt.Errorf("invalid state")
	}

	app.CaroQueue = removedUserFromMmQueue(session.UserId, app.CaroQueue)
	session.State = Idle
	session.CurrentGame = shared.None

	return nil
}

func removedUserFromMmQueue(uId pgtype.UUID, q MmQueue) MmQueue {
	for i, id := range q {
		if id == uId {
			return append(q[:i], q[i+1:]...)
		}
	}
	return q
}
