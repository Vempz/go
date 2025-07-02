package repository

import (
    "errors"
    "sync"

    "user-session-service/internal/models"
)

type UserRepository interface {
    Create(user *models.User) error
    GetByUsername(username string) (*models.User, error)
    GetByID(id string) (*models.User, error)
}

type InMemoryUserRepository struct {
    users map[string]*models.User
    mutex sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{
        users: make(map[string]*models.User),
    }
}

func (r *InMemoryUserRepository) Create(user *models.User) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.users[user.Username]; exists {
        return errors.New("username already exists")
    }

    r.users[user.Username] = user
    return nil
}

func (r *InMemoryUserRepository) GetByUsername(username string) (*models.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    user, exists := r.users[username]
    if !exists {
        return nil, errors.New("user not found")
    }

    return user, nil
}

func (r *InMemoryUserRepository) GetByID(id string) (*models.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    for _, user := range r.users {
        if user.ID == id {
            return user, nil
        }
    }

    return nil, errors.New("user not found")
}
