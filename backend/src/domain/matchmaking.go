package domain

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type MmQueue []User

// CaroMatch == nil && error == nil: join queue, not found match yet
func (a *AppState) UserJoinCaroQueue(dbe *DbExecDeps, uSes *UserSession) (*CaroMatch, error) {
	newU, err := GetUserById(dbe, uSes.UserId)
	if err != nil {
		return nil, err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if uSes.State != Idle {
		return nil, fmt.Errorf("invalid state")
	}

	if len(a.caroQueue) > 0 {
		xPlayer := a.caroQueue[0]
		a.caroQueue = removedUserFromMmQueue(xPlayer.ID, a.caroQueue)

		oPlayer := newU

		matchId := uuid.Must(uuid.NewV7())
		a.caroMatchesMap[matchId] = NewCaroMatch(matchId, true, xPlayer.ID, int(xPlayer.CaroRating), oPlayer.ID, int(oPlayer.CaroRating))

		xSession := a.userSessionsMap[xPlayer.ID]
		oSession := a.userSessionsMap[oPlayer.ID]

		xSession.State = PlayingCaro
		xSession.CurrentMatchId = matchId

		oSession.State = PlayingCaro
		oSession.CurrentMatchId = matchId

		return a.caroMatchesMap[matchId], nil
	}

	a.caroQueue = append(a.caroQueue, newU)
	uSes.State = QueuingCaro

	return nil, nil
}

func (a *AppState) UserLeaveCaroQueue(us *UserSession) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if us.State != QueuingCaro {
		return fmt.Errorf("invalid state")
	}

	a.caroQueue = removedUserFromMmQueue(us.UserId, a.caroQueue)
	us.State = Idle

	//fmt.Println(a.caroQueue)
	return nil
}

func removedUserFromMmQueue(uId pgtype.UUID, q MmQueue) MmQueue {
	for i, u := range q {
		if u.ID == uId {
			return append(q[:i], q[i+1:]...)
		}
	}
	return q
}
