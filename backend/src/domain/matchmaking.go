package domain

import (
	"backend/src/shared"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type MmQueue []pgtype.UUID

func (a *AppState) UserJoinCaroQueue(us *UserSession) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if us.State != Idle {
		return fmt.Errorf("invalid state")
	}

	a.caroQueue = append(a.caroQueue, us.UserId)
	us.State = Queuing
	us.CurrentGame = shared.Caro

	//fmt.Println(a.caroQueue)
	return nil
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
