package handlers

import (
	"context"

	Provider "GameApplication/entity/entityProvider"
	friend "GameApplication/entity/friends"
)

func GetFriends(ctx context.Context, p Provider.Provider, userID string) ([]*friend.Friend, error) {
	return friend.GetFriends(ctx, p, userID)
}
