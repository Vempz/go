package models

import (
    "time"
)

type User struct {
    ID        string    `json:"id"`
    Username  string    `json:"username" binding:"required,alphanum,min=3,max=30"`
    Password  string    `json:"password" binding:"required,min=6"`
    Email     string    `json:"email" binding:"required,email"`
    CreatedAt time.Time `json:"createdAt"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required,alphanum,min=3,max=30"`
    Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,alphanum,min=3,max=30"`
    Password string `json:"password" binding:"required,min=6"`
    Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
