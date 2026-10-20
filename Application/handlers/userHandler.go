package handlers

import (
	"context"

	Provider "GameApplication/entity/entityProvider"
	user "GameApplication/entity/user"
)

func GetUserByID(ctx context.Context, p Provider.Provider, id string) (*user.User, error) {
	return user.GetUserByID(ctx, p, id)
}

func UpdateUser(ctx context.Context, p Provider.Provider, id string, u *user.User) (*user.User, error) {
	u.ID = id
	if u.Password != "" {
		u.Password = encription(ctx, p, u.Password, "")
	}
	if err := user.UpdateUser(ctx, p, u); err != nil {
		return nil, err
	}
	return user.GetUserByID(ctx, p, id)
}

func UpdateUserStatus(ctx context.Context, p Provider.Provider, id string, status string) (*user.User, error) {
	if err := user.UpdateUser(ctx, p, &user.User{ID: id, Status: status}); err != nil {
		return nil, err
	}
	return user.GetUserByID(ctx, p, id)
}
