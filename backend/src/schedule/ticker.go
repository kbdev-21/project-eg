package schedule

import (
	"backend/src/db"
	"backend/src/domain"
	"backend/src/router"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EverySecond(a *domain.AppState, q *db.Queries, p *pgxpool.Pool) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		checkAllCaroMatchesTimeout(a, domain.CreateDbExec(ctx, q, p))
		cancel()
	}
}

func checkAllCaroMatchesTimeout(a *domain.AppState, dbe *domain.DbExec) {
	fmt.Println("Check time out")
	matches := a.GetCaroMatches()
	for _, m := range matches {
		lastAction := m.StartedAt
		if len(m.Moves) > 0 {
			lastAction = m.Moves[len(m.Moves)-1].PlayedAt
		}
		if time.Since(lastAction) > domain.CARO_MAX_MOVE_TIME {
			m.OutOfTime()
			matchResult, err := a.ProcessCaroMatchEnded(dbe, *m)
			if err != nil {
				continue
			}

			xSes, _ := a.GetUserSession(matchResult.XPlayerID)
			oSes, _ := a.GetUserSession(matchResult.OPlayerID)

			xMsg := router.BuildServerMessage(router.CaroMatchEndedOutOfTime, *xSes, a)
			xMsg.Data = map[string]any{
				"endedMatch": matchResult,
			}
			xSes.WsConn.WriteJSON(xMsg)

			oMsg := router.BuildServerMessage(router.CaroMatchEndedOutOfTime, *oSes, a)
			oMsg.Data = map[string]any{
				"endedMatch": matchResult,
			}
			oSes.WsConn.WriteJSON(oMsg)
		}
	}
}
