package service

import (
    "errors"
    "time"

    "user-session-service/internal/config"
    "user-session-service/internal/models"
    "user-session-service/internal/repository"
    "user-session-service/internal/utils"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo    repository.UserRepository
    sessionRepo repository.SessionRepository
    config      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, config *config.Config) *AuthService {
    return &AuthService{
        userRepo:    userRepo,
        sessionRepo: sessionRepo,
        config:      config,
    }
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
    // Check if user exists
    if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
        return nil, errors.New("username already exists")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.config.BCryptCost)
    if err != nil {
        return nil, err
    }

    user := &models.User{
        ID:        uuid.New().String(),
        Username:  req.Username,
        Password:  string(hashedPassword),
        Email:     req.Email,
        CreatedAt: time.Now(),
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *AuthService) Login(req *models.LoginRequest, ipAddress, userAgent string) (*models.LoginResponse, error) {
    // Get user
    user, err := s.userRepo.GetByUsername(req.Username)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, errors.New("invalid credentials")
    
