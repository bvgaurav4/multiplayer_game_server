package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Match struct {
	ID       string    `json:"id"`
	Players  []string  `json:"players"`
	Mode     string    `json:"mode"`
	Result   string    `json:"result"`
	PlayedAt time.Time `json:"played_at"`
}

type Friend struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	FriendID  string    `json:"friend_id"`
	Alias     string    `json:"alias"`
	Status    string    `json:"status"`
	AddedAt   time.Time `json:"added_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	storeMu sync.RWMutex
	users   = map[string]*User{}
	matches = map[string]*Match{}
	friends = map[string]*Friend{}
)

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

		storeMu.Lock()
		defer storeMu.Unlock()

		for _, u := range users {
			if u.Email == req.Email {
				respondFailure(c, http.StatusConflict, "email already registered")
				return
			}
		}

		now := time.Now().UTC()
		user := &User{
			ID:        uuid.New().String(),
			Username:  req.Username,
			Email:     req.Email,
			Password:  req.Password,
			Status:    "offline",
			CreatedAt: now,
			UpdatedAt: now,
		}
		users[user.ID] = user

		respondSuccess(c, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"status":   user.Status,
		})
	})

	// GET /users/:id - get the basic details of a user
	r.GET("/users/:id", func(c *gin.Context) {
		storeMu.RLock()
		defer storeMu.RUnlock()

		user, ok := users[c.Param("id")]
		if !ok {
			respondFailure(c, http.StatusNotFound, "user not found")
			return
		}

		respondSuccess(c, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"status":     user.Status,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
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

		storeMu.Lock()
		defer storeMu.Unlock()

		user, ok := users[c.Param("id")]
		if !ok {
			respondFailure(c, http.StatusNotFound, "user not found")
			return
		}
		if req.Email != "" {
			for id, u := range users {
				if id != user.ID && u.Email == req.Email {
					respondFailure(c, http.StatusConflict, "email already registered")
					return
				}
			}
		}

		if req.Username != "" {
			user.Username = req.Username
		}
		if req.Email != "" {
			user.Email = req.Email
		}
		if req.Password != "" {
			user.Password = req.Password
		}
		user.UpdatedAt = time.Now().UTC()

		respondSuccess(c, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"updated_at": user.UpdatedAt,
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

		storeMu.Lock()
		defer storeMu.Unlock()

		user, ok := users[c.Param("id")]
		if !ok {
			respondFailure(c, http.StatusNotFound, "user not found")
			return
		}

		user.Status = req.Status
		user.UpdatedAt = time.Now().UTC()

		respondSuccess(c, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"status":     user.Status,
			"updated_at": user.UpdatedAt,
		})
	})

	// GET /matches - get the matches, returns an object
	r.GET("/matches", func(c *gin.Context) {
		storeMu.RLock()
		defer storeMu.RUnlock()

		list := make([]*Match, 0, len(matches))
		for _, m := range matches {
			list = append(list, m)
		}

		respondSuccess(c, gin.H{
			"matches": list,
			"count":   len(list),
		})
	})

	// GET /friends - get the friends, returns objects
	r.GET("/friends", func(c *gin.Context) {
		userID := c.Query("user_id")

		storeMu.RLock()
		defer storeMu.RUnlock()

		list := make([]*Friend, 0, len(friends))
		for _, f := range friends {
			if userID != "" && f.UserID != userID {
				continue
			}
			list = append(list, f)
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
