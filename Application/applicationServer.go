package main

import (
	"errors"
	"net/http"
	"strings"

	// played "GameApplication/entity/played"
	Provider "GameApplication/entity/entityProvider"
	match "GameApplication/entity/matches"
	user "GameApplication/entity/user"
	"GameApplication/handlers"
	reqres "GameApplication/reqres"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Blank import registers the driver
)

type Config struct {
	Database DatabaseConfig `toml:"database"`
}

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
	SSLMode  string `toml:"sslmode"`

	Pool PoolConfig `toml:"pool"`
}

type PoolConfig struct {
	MaxOpenConns           int `toml:"max_open_conns"`
	MaxIdleConns           int `toml:"max_idle_conns"`
	ConnMaxIdleTimeMinutes int `toml:"conn_max_idle_time_minutes"`
	ConnMaxLifetimeMinutes int `toml:"conn_max_lifetime_minutes"`
}

func respondSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func respondFailure(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   message,
	})
}

func main() {
	r := gin.Default()
	provider := Provider.Provider{}
	// db := pq.
	Provider.Provider.Conn(provider)

	r.POST("/login", func(ctx *gin.Context) {
		var loginRequest reqres.LoginRequest
		if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
			respondFailure(ctx, http.StatusBadRequest, "invalid request body")
			return
		}

		user, session, err := handlers.SignIn(ctx, provider, loginRequest)
		if err != nil {
			respondFailure(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		respondSuccess(ctx, reqres.LoginResponse{
			UserID:                user.ID,
			Token:                 session.Token,
			RefreshToken:          session.RefreshToken,
			TokenExpiresAt:        session.TokenExpiresAt,
			RefreshTokenExpiresAt: session.RefreshTokenExpiresAt,
		})
	})

	r.POST("/create-account", func(ctx *gin.Context) {
		var signUpRequest reqres.SignUpRequest
		handlers.SignUp(ctx, provider, signUpRequest)
	})

	// POST /users - create a user, returns a success or failure object
	r.POST("/users", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			respondFailure(c, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.Username == "" || req.Email == "" || req.Password == "" {
			respondFailure(c, http.StatusBadRequest, "username, email and password are required")
			return
		}

		created, err := handlers.SignUp(c, provider, reqres.SignUpRequest{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			if errors.Is(err, user.ErrUserAlreadyExists) {
				respondFailure(c, http.StatusConflict, "email or username already registered")
				return
			}
			respondFailure(c, http.StatusInternalServerError, "could not create user")
			return
		}

		respondSuccess(c, gin.H{
			"id":       created.ID,
			"username": created.Username,
			"email":    created.Email,
			"status":   created.Status,
		})
	})

	// GET /users/:id - get the basic details of a user
	r.GET("/users/:id", func(c *gin.Context) {
		found, err := handlers.GetUserByID(c, provider, c.Param("id"))
		if err != nil {
			respondFailure(c, http.StatusNotFound, "user not found")
			return
		}

		respondSuccess(c, gin.H{
			"id":         found.ID,
			"username":   found.Username,
			"email":      found.Email,
			"status":     found.Status,
			"created_at": found.CreatedAt,
			"updated_at": found.UpdatedAt,
		})
	})

	// POST /users/:id - update a user, returns a success or failure object
	r.POST("/users/:id", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			respondFailure(c, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.Username == "" && req.Email == "" && req.Password == "" {
			respondFailure(c, http.StatusBadRequest, "no fields to update")
			return
		}

		updated, err := handlers.UpdateUser(c, provider, c.Param("id"), &user.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			switch {
			case errors.Is(err, user.ErrUserNotFound):
				respondFailure(c, http.StatusNotFound, "user not found")
			case errors.Is(err, user.ErrUserAlreadyExists):
				respondFailure(c, http.StatusConflict, "email or username already registered")
			default:
				respondFailure(c, http.StatusInternalServerError, "could not update user")
			}
			return
		}

		respondSuccess(c, gin.H{
			"id":         updated.ID,
			"username":   updated.Username,
			"email":      updated.Email,
			"updated_at": updated.UpdatedAt,
		})
	})

	// PUT /users/:id/status - update the status of a user, returns the updated object
	r.PUT("/users/:id/status", func(c *gin.Context) {
		var req struct {
			Status string `json:"status"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			respondFailure(c, http.StatusBadRequest, "invalid request body")
			return
		}

		switch req.Status {
		case "online", "offline", "in_match", "away":
		default:
			respondFailure(c, http.StatusBadRequest, "status must be one of online, offline, in_match, away")
			return
		}

		updated, err := handlers.UpdateUserStatus(c, provider, c.Param("id"), req.Status)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				respondFailure(c, http.StatusNotFound, "user not found")
				return
			}
			respondFailure(c, http.StatusInternalServerError, "could not update user status")
			return
		}

		respondSuccess(c, gin.H{
			"id":         updated.ID,
			"username":   updated.Username,
			"status":     updated.Status,
			"updated_at": updated.UpdatedAt,
		})
	})

	// GET /users/:id/matches - get all the matches of a user, returns an object
	r.GET("/users/:id/matches", func(c *gin.Context) {
		var req reqres.GetMatchesRequest
		if err := c.ShouldBindUri(&req); err != nil || req.UserID == "" {
			respondFailure(c, http.StatusBadRequest, "user id is required")
			return
		}

		res, err := handlers.GetMatches(c, provider, req)
		if err != nil {
			respondFailure(c, http.StatusInternalServerError, "could not fetch matches")
			return
		}

		respondSuccess(c, res)
	})

	// GET /matches/:id - get the details of a match, returns an object
	r.GET("/matches/:id", func(c *gin.Context) {
		var req reqres.GetMatchDetailsRequest
		if err := c.ShouldBindUri(&req); err != nil || req.MatchID == "" {
			respondFailure(c, http.StatusBadRequest, "match id is required")
			return
		}

		res, err := handlers.GetMatchDetails(c, provider, req)
		if err != nil {
			if errors.Is(err, match.ErrMatchNotFound) {
				respondFailure(c, http.StatusNotFound, "match not found")
				return
			}
			respondFailure(c, http.StatusInternalServerError, "could not fetch match")
			return
		}

		respondSuccess(c, *res)
	})

	// POST /feedback - submit feedback, returns the created object
	r.POST("/feedback", func(c *gin.Context) {
		var req reqres.CreateFeedbackRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondFailure(c, http.StatusBadRequest, "invalid request body")
			return
		}
		if strings.TrimSpace(req.Message) == "" {
			respondFailure(c, http.StatusBadRequest, "message is required")
			return
		}

		res, err := handlers.CreateFeedback(c, provider, req)
		if err != nil {
			respondFailure(c, http.StatusInternalServerError, "could not save feedback")
			return
		}

		respondSuccess(c, res)
	})

	// GET /friends - get the friends, returns objects
	r.GET("/friends", func(c *gin.Context) {
		list, err := handlers.GetFriends(c, provider, c.Query("user_id"))
		if err != nil {
			respondFailure(c, http.StatusInternalServerError, "could not fetch friends")
			return
		}

		respondSuccess(c, gin.H{
			"friends": list,
			"count":   len(list),
		})
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
