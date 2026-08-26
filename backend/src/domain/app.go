package domain

import (
	"backend/src/db"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppState struct {
	// states
	mu              sync.Mutex
	userSessionsMap map[uuid.UUID]*UserSession
	caroQueue       MmQueue
	caroMatchesMap  map[uuid.UUID]*CaroMatch

	// db exec
	q *db.Queries
	p *pgxpool.Pool
}

func InitAppState(q *db.Queries, p *pgxpool.Pool) *AppState {
	return &AppState{
		mu:              sync.Mutex{},
		userSessionsMap: map[uuid.UUID]*UserSession{},
		caroQueue:       MmQueue{},
		caroMatchesMap:  map[uuid.UUID]*CaroMatch{},
		q: q,
		p: p,
	}
}

func (a *AppState) GetUserSession(userId uuid.UUID) (*UserSession, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	session, existed := a.userSessionsMap[userId]

	return session, existed
}

func (a *AppState) GetOrCreateUserSession(userId uuid.UUID) *UserSession {
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

func (a *AppState) GetCaroMatches() []*CaroMatch {
	a.mu.Lock()
	defer a.mu.Unlock()

	matches := []*CaroMatch{}
	for _, m := range a.caroMatchesMap {
		matches = append(matches, m)
	}
	return matches
}
