package domain

import (
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

type AppState struct {
	mu           sync.Mutex
	userSessions map[pgtype.UUID]*UserSession
	CaroQueue    MmQueue
}

func NewAppState() *AppState {
	return &AppState{
		mu:           sync.Mutex{},
		userSessions: map[pgtype.UUID]*UserSession{},
		CaroQueue:    MmQueue{},
	}
}

func (s *AppState) SafeGetUserSession(userId pgtype.UUID) (*UserSession, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, existed := s.userSessions[userId]

	return session, existed
}

func (s *AppState) SafeGetOrCreateUserSession(userId pgtype.UUID) *UserSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.userSessions[userId]
	if !ok {
		session = &UserSession{
			UserId: userId,
			State:  Idle,
		}
		s.userSessions[userId] = session
	}

	return session
}
