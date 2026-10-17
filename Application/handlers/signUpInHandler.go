package handlers

import (
	// "sync"
	"context"
	"fmt"
	"time"

	// friend "GameApplication/entity/friends"
	// match "GameApplication/entity/matches"
	Provider "GameApplication/entity/entityProvider"
	User "GameApplication/entity/user"
	reqres "GameApplication/reqres"

	"github.com/google/uuid"
)

func SignUp(ctx context.Context, p Provider.Provider, loginRequest reqres.SignUpRequest) (*User.User, error) {
	// var req struct {
	// 	Username string `json:"username"`
	// 	Email    string `json:"email"`
	// 	Password string `json:"password"`
	// }
	now := time.Now().UTC()
	user := &User.User{
		ID:        uuid.New().String(),
		Username:  loginRequest.Username,
		Email:     loginRequest.Email,
		Password:  encription(ctx, p, loginRequest.Password),
		Status:    "offline",
		CreatedAt: now,
		UpdatedAt: now,
	}
	error := User.CreateUser(ctx, p, user)
	if error != nil {
		return nil, error
	}
	return user, nil
}

func SignIn(ctx context.Context, p Provider.Provider, login reqres.LoginRequest) (*User.User, error) {
	// var req struct {
	// 	identifier string `json:"identifier"`
	// 	Password   string `json:"password"`
	// }
	var user *User.User
	user, err := User.GetUserWithUsernameOrEmail(ctx, p, login.Identifier)
	if err != nil {
		return nil, err
	}
	if user.Password != encription(ctx, p, login.Password) {
		return nil, fmt.Errorf("%s password is incorrect for the user ", login.Identifier)
	} else {
		return user, nil
	}
}
