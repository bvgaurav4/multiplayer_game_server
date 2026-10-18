package main

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	// played "GameApplication/entity/played"
	Provider "GameApplication/entity/entityProvider"
	friend "GameApplication/entity/friends"
	match "GameApplication/entity/matches"
	user "GameApplication/entity/user"
	"GameApplication/handlers"
	reqres "GameApplication/reqres"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // Blank import registers the driver
)

var (
	storeMu sync.RWMutex
	users   = map[string]*user.User{}
	friends = map[string]*friend.Friend{}
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
	provider := Provider.Provider{}
	// db := pq.
	Provider.Provider.Conn(provider)

	r.POST("/login", func(ctx *gin.Context) {
		var loginRequest reqres.LoginRequest
		handlers.SignIn(ctx, provider, loginRequest)

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

		storeMu.Lock()
		defer storeMu.Unlock()

		for _, u := range users {
			if u.Email == req.Email {
				respondFailure(c, http.StatusConflict, "email already registered")
				return
			}
		}

		now := time.Now().UTC()
		user := &user.User{
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

		// user, ok := users[c.Param("id")]
		user, err := handlers.GetUser(c, provider, c.Param("id"))
		if err != nil {
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
		userID := c.Query("user_id")

		storeMu.RLock()
		defer storeMu.RUnlock()

		list := make([]*friend.Friend, 0, len(friends))
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
