package repository

import (
    "errors"
    "sync"

    "user-session-service/internal/models"
)

type SessionRepository interface {
    Create(session *models.Session) error
    GetByID(sessionID string) (*models.Session, error)
    GetByUserID(userID string) ([]*models.Session, error)
    Update(session *models.Session) error
    Delete(sessionID string) error
}

type InMemorySessionRepository struct {
    sessions map[string]*models.Session
    mutex    sync.RWMutex
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
    return &InMemorySessionRepository{
        sessions: make(map[string]*models.Session),
    }
}

func (r *InMemorySessionRepository) Create(session *models.Session) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    r.sessions[session.SessionID] = session
    return nil
}

func (r *InMemorySessionRepository) GetByID(sessionID string) (*models.Session, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    session, exists := r.sessions[sessionID]
    if !exists {
        return nil, errors.New("session not found")
    }

    return session, nil
}

func (r *InMemorySessionRepository) GetByUserID(userID string) ([]*models.Session, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    var userSessions []*models.Session
    for _, session := range r.sessions {
        if session.UserID == userID {
            userSessions = append(userSessions, session)
        }
    }

    return userSessions, nil
}

func (r *InMemorySessionRepository) Update(session *models.Session) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.sessions[session.SessionID]; !exists {
        return errors.New("session not found")
    }

    r.sessions[session.SessionID] = session
    return nil
}

func (r *InMemorySessionRepository) Delete(sessionID string) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.sessions[sessionID]; !exists {
        return errors.New("session not found")
    }

    delete(r.sessions, sessionID)
    return nil
}
