package schedule

import (
	"backend/src/domain"
	"backend/src/router"
	"context"
	"time"
)

func EverySecond(a *domain.AppState, hub *router.WsConnHub) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		checkAllCaroMatchesTimeout(a, hub)
	}
}

func checkAllCaroMatchesTimeout(a *domain.AppState, hub *router.WsConnHub) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	matches := a.GetCaroMatches()
	for _, m := range matches {
		lastAction := m.StartedAt
		if len(m.Moves) > 0 {
			lastAction = m.Moves[len(m.Moves)-1].PlayedAt
		}
		if time.Since(lastAction) > domain.CARO_MAX_MOVE_TIME {
			m.OutOfTime()
			endedMatch, err := a.ProcessCaroMatchEnded(ctx, *m)
			if err != nil {
				continue
			}

			xSes, xExisted := a.GetUserSession(endedMatch.XPlayerId)
			oSes, oExisted := a.GetUserSession(endedMatch.OPlayerId)

			if xExisted {
				xMsg := router.BuildServerMessage(router.CaroMatchEndedOutOfTime, *xSes, a)
				xMsg.Data = map[string]any{
					"match": endedMatch,
				}
				hub.Send(xSes.UserId, xMsg)
			}

			if oExisted {
				oMsg := router.BuildServerMessage(router.CaroMatchEndedOutOfTime, *oSes, a)
				oMsg.Data = map[string]any{
					"match": endedMatch,
				}
				hub.Send(oSes.UserId, oMsg)
			}
		}
	}
}
