package domain

import (
	"backend/src/db"
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppState struct {
	mu              sync.Mutex
	userSessionsMap map[pgtype.UUID]*UserSession
	caroQueue       MmQueue
	caroMatchesMap  map[uuid.UUID]*CaroMatch
}

func NewAppState() *AppState {
	return &AppState{
		mu:              sync.Mutex{},
		userSessionsMap: map[pgtype.UUID]*UserSession{},
		caroQueue:       MmQueue{},
		caroMatchesMap:  map[uuid.UUID]*CaroMatch{},
	}
}

type DbExecDeps struct {
	c context.Context
	q *db.Queries
	p *pgxpool.Pool
}

func CreateDbExecDeps(c context.Context, q *db.Queries, p *pgxpool.Pool) *DbExecDeps {
	return &DbExecDeps{
		c: c,
		q: q,
		p: p,
	}
}

func (a *AppState) GetUserSession(userId pgtype.UUID) (*UserSession, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	session, existed := a.userSessionsMap[userId]

	return session, existed
}

func (a *AppState) GetOrCreateUserSession(userId pgtype.UUID) *UserSession {
	a.mu.Lock()
	defer a.mu.Unlock()

	session, ok := a.userSessionsMap[userId]
	if !ok {
		session = &UserSession{
			UserId:         userId,
			State:          Idle,
			CurrentMatchId: uuid.Nil,
		}
		a.userSessionsMap[userId] = session
	}

	return session
}

func (a *AppState) GetCaroMatch(matchId uuid.UUID) (*CaroMatch, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	m, existed := a.caroMatchesMap[matchId]

	return m, existed
}
