package handlers

import (
	// "sync"
	"context"
	"fmt"
	"strings"
	"time"

	// friend "GameApplication/entity/friends"
	// match "GameApplication/entity/matches"
	Auth "GameApplication/entity/auth"
	Provider "GameApplication/entity/entityProvider"
	User "GameApplication/entity/user"
	reqres "GameApplication/reqres"

	"github.com/google/uuid"
)

func SignUp(ctx context.Context, p Provider.Provider, loginRequest reqres.SignUpRequest) (*User.User, error) {
	now := time.Now().UTC()
	user := &User.User{
		ID:        uuid.New().String(),
		Username:  loginRequest.Username,
		Email:     loginRequest.Email,
		Password:  encription(ctx, p, loginRequest.Password, ""),
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

func SignIn(ctx context.Context, p Provider.Provider, login reqres.LoginRequest) (*User.User, *Auth.Auth, error) {
	if strings.TrimSpace(login.RefreshToken) != "" {
		return signInWithRefreshToken(ctx, p, login.RefreshToken)
	}

	if strings.TrimSpace(login.Identifier) == "" || login.Password == "" {
		return nil, nil, fmt.Errorf("identifier and password are required")
	}

	user, err := User.GetUserWithUsernameOrEmail(ctx, p, login.Identifier)
	if err != nil {
		return nil, nil, err
	}
	if user.Password != encription(ctx, p, login.Password, user.Password) {
		return nil, nil, fmt.Errorf("%s password is incorrect for the user ", login.Identifier)
	}

	session, err := Auth.CreateAuth(ctx, p, user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func signInWithRefreshToken(ctx context.Context, p Provider.Provider, refreshToken string) (*User.User, *Auth.Auth, error) {
	session, err := Auth.GetByRefreshToken(ctx, p, refreshToken)
	if err != nil {
		return nil, nil, err
	}

	user, err := User.GetUserByID(ctx, p, session.UserID)
	if err != nil {
		return nil, nil, err
	}

	session, err = Auth.RenewAccessToken(ctx, p, session)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}
