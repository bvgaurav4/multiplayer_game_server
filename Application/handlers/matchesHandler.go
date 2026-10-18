package handlers

import (
	Provider "GameApplication/entity/entityProvider"
	match "GameApplication/entity/matches"
	reqres "GameApplication/reqres"
	"context"
)

func GetMatches(ctx context.Context, p Provider.Provider, req reqres.GetMatchesRequest) (*reqres.GetMatchesResponse, error) {
	list, err := match.GetMatchesByUserID(ctx, p, req.UserID)
	if err != nil {
		return nil, err
	}
	return &reqres.GetMatchesResponse{
		UserID:  req.UserID,
		Matches: list,
		Count:   len(list),
	}, nil
}

func GetMatchDetails(ctx context.Context, p Provider.Provider, req reqres.GetMatchDetailsRequest) (*reqres.GetMatchDetailsResponse, error) {
	return match.GetMatchByID(ctx, p, req.MatchID)
}
