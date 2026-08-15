package domain

import (
	"backend/src/shared"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type MmQueue []pgtype.UUID

// CaroMatch == nil && error == nil: join queue, not found match yet
func (a *AppState) UserJoinCaroQueue(us *UserSession) (*CaroMatch, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if us.State != Idle {
		return nil, fmt.Errorf("invalid state")
	}

	if len(a.caroQueue) > 0 {
		xPlayerId := a.caroQueue[0]
		a.caroQueue = removedUserFromMmQueue(xPlayerId, a.caroQueue)

		oPlayerId := us.UserId

		matchId := uuid.New()
		a.caroMatchesMap[matchId] = NewCaroMatch(true, xPlayerId, oPlayerId)
		
		xSession := a.userSessionsMap[xPlayerId]
		oSession := a.userSessionsMap[oPlayerId]

		xSession.State = Playing
		xSession.CurrentGame = shared.Caro
		xSession.CurrentMatchId = matchId

		oSession.State = Playing
		oSession.CurrentGame = shared.Caro
		oSession.CurrentMatchId = matchId

		return a.caroMatchesMap[matchId], nil
	}

	a.caroQueue = append(a.caroQueue, us.UserId)
	us.State = Queuing
	us.CurrentGame = shared.Caro

	return nil, nil
}

func (a *AppState) UserLeaveCaroQueue(us *UserSession) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !(us.State == Queuing && us.CurrentGame == shared.Caro) {
		return fmt.Errorf("invalid state")
	}

	a.caroQueue = removedUserFromMmQueue(us.UserId, a.caroQueue)
	us.State = Idle
	us.CurrentGame = shared.None

	//fmt.Println(a.caroQueue)
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