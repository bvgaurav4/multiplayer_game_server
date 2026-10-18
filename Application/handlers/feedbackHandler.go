package handlers

import (
	Provider "GameApplication/entity/entityProvider"
	feedback "GameApplication/entity/feedback"
	reqres "GameApplication/reqres"
	"context"
)

func CreateFeedback(ctx context.Context, p Provider.Provider, req reqres.CreateFeedbackRequest) (*reqres.CreateFeedbackResponse, error) {
	f := &feedback.Feedback{
		UserID:  req.UserID,
		Message: req.Message,
	}
	if err := feedback.CreateFeedback(ctx, p, f); err != nil {
		return nil, err
	}
	return &reqres.CreateFeedbackResponse{
		ID:        f.ID,
		UserID:    f.UserID,
		Message:   f.Message,
		CreatedAt: f.CreatedAt,
	}, nil
}
