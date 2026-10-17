package handlers

import (
	Provider "GameApplication/entity/entityProvider"
	user "GameApplication/entity/user"
	"context"
)

func GetUser(ctx context.Context, p Provider.Provider, username string) (*user.User, error) {
	return user.GetUserWithUsernameOrEmail(ctx, p, username)
}
