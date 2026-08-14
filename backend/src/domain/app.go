package domain

import (
	"backend/src/db"
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

type AppState struct {
	mu           sync.Mutex
	userSessions map[pgtype.UUID]*UserSession
	caroQueue    MmQueue
}

type DbExecDeps struct {
	c context.Context
	q *db.Queries
}

func NewAppState() *AppState {
	return &AppState{
		mu:           sync.Mutex{},
		userSessions: map[pgtype.UUID]*UserSession{},
		caroQueue:    MmQueue{},
	}
}

func CreateDbExecDeps(c context.Context, q *db.Queries) *DbExecDeps {
	return &DbExecDeps{
		c: c,
		q: q,
	}
}

func (a *AppState) GetUserSession(userId pgtype.UUID) (*UserSession, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	session, existed := a.userSessions[userId]

	return session, existed
}

func (a *AppState) GetOrCreateUserSession(userId pgtype.UUID) *UserSession {
	a.mu.Lock()
	defer a.mu.Unlock()

	session, ok := a.userSessions[userId]
	if !ok {
		session = &UserSession{
			UserId: userId,
			State:  Idle,
		}
		a.userSessions[userId] = session
	}

	return session
}