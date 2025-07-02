package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "user-session-service/internal/config"
    "user-session-service/internal/handlers"
    "user-session-service/internal/middleware"
    "user-session-service/internal/repository"
    "user-session-service/internal/service"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/redis"
)

func main() {
    // Load configuration
    cfg := config.Load()

    // Setup Redis store for sessions
    store, err := redis.NewStore(10, "tcp", cfg.RedisAddr, cfg.RedisPassword, []byte(cfg.SessionSecret))
    if err != nil {
        log.Fatal("Failed to create Redis store:", err)
    }

    // Initialize repositories
    userRepo := repository.NewInMemoryUserRepository()
    sessionRepo := repository.NewInMemorySessionRepository()

    // Initialize services
    authService := service.NewAuthService(userRepo, sessionRepo, cfg)
    sessionService := service.NewSessionService(sessionRepo, cfg)

    // Initialize handlers
    authHandler := handlers.NewAuthHandler(authService)
    sessionHandler := handlers.NewSessionHandler(sessionService)

    // Setup Gin router
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    // CORS middleware
    corsConfig := cors.DefaultConfig()
    corsConfig.AllowOrigins = cfg.AllowedOrigins
    corsConfig.AllowCredentials = true
    corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
    r.Use(cors.New(corsConfig))

    // Session middleware
    r.Use(sessions.Sessions("session", store))

    // Custom middleware
    r.Use(middleware.RequestLogger())

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":    "healthy",
            "timestamp": time.Now().Format(time.RFC3339),
        })
    })

    // API routes
    api := r.Group("/api")
    {
        auth := api.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.POST("/logout", middleware.AuthRequired(), authHandler.Logout)
        }

        sessions := api.Group("/sessions")
        sessions.Use(middleware.AuthRequired())
        {
            sessions.GET("/current", sessionHandler.GetCurrentSession)
            sessions.GET("/user/:userId", sessionHandler.GetUserSessions)
            sessions.DELETE("/:sessionId", sessionHandler.RevokeSession)
        }
    }

    // 404 handler
    r.NoRoute(func(c *gin.Context) {
        c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
    })

    // Start server with graceful shutdown
    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: r,
    }

    go func() {
        log.Printf("User Session Service starting on port %s", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the server
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exited")
}
