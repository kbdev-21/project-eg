package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type MmQueue []User

// CaroMatch == nil && error == nil: join queue, not found match yet
func (a *AppState) UserJoinCaroQueue(ctx context.Context, uSes *UserSession) (*CaroMatch, error) {
	newU, err := a.GetUserById(ctx, uSes.UserId)
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
		a.caroQueue = removedUserFromMmQueue(xPlayer.Id, a.caroQueue)

		oPlayer := newU

		matchId := uuid.Must(uuid.NewV7())
		a.caroMatchesMap[matchId] = NewCaroMatch(matchId, true, xPlayer.Id, xPlayer.CaroRating, oPlayer.Id, oPlayer.CaroRating)

		xSession := a.userSessionsMap[xPlayer.Id]
		oSession := a.userSessionsMap[oPlayer.Id]

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

func removedUserFromMmQueue(uId uuid.UUID, q MmQueue) MmQueue {
	for i, u := range q {
		if u.Id == uId {
			return append(q[:i], q[i+1:]...)
		}
	}
	return q
}
