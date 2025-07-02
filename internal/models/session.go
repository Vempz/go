package models

import (
    "time"
)

type Session struct {
    SessionID    string    `json:"sessionId"`
    UserID       string    `json:"userId"`
    Username     string    `json:"username"`
    CreatedAt    time.Time `json:"createdAt"`
    LastAccessed time.Time `json:"lastAccessed"`
    IPAddress    string    `json:"ipAddress"`
    UserAgent    string    `json:"userAgent"`
}

type SessionResponse struct {
    SessionID    string    `json:"sessionId"`
    UserID       string    `json:"userId"`
    Username     string    `json:"username"`
    CreatedAt    time.Time `json:"createdAt"`
    LastAccessed time.Time `json:"lastAccessed"`
}

type LoginResponse struct {
    Message   string       `json:"message"`
    SessionID string       `json:"sessionId"`
    Token     string       `json:"token"`
    User      UserResponse `json:"user"`
}
