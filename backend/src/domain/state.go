package domain

import (
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

type AppState struct {
	Mu           sync.Mutex
	UserSessions map[pgtype.UUID]*UserSession
}

func NewAppState() *AppState {
	return &AppState{
		Mu:           sync.Mutex{},
		UserSessions: map[pgtype.UUID]*UserSession{},
	}
}

func (s *AppState) GetUserSession(userId pgtype.UUID) (*UserSession, bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	session, existed := s.UserSessions[userId]

	return session, existed
}

func (s *AppState) GetOrCreateUserSession(userId pgtype.UUID) *UserSession {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	session, ok := s.UserSessions[userId]
	if !ok {
		session = &UserSession{
			UserId: userId,
			State:  Idle,
		}
		s.UserSessions[userId] = session
	}

	return session
}
