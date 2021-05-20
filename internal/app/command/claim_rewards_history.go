package command

import (
	"context"
	"github.com/ramadani/balapan/internal/app/command/model"
	"github.com/ramadani/balapan/internal/domain/history"
	uuid "github.com/satori/go.uuid"
	"time"
)

type claimRewardsHistoryCommand struct {
	historyRepo history.Repository
}

func (c *claimRewardsHistoryCommand) Do(ctx context.Context, data *model.ClaimRewards) error {
	history := &history.History{
		ID:        uuid.NewV4().String(),
		RewardsID: data.ID,
		UserID:    data.UserID,
		Amount:    data.Amount,
		CreatedAt: time.Now(),
	}

	return c.historyRepo.Create(ctx, history)
}

func NewClaimRewardsHistoryCommand(historyRepo history.Repository) ClaimRewardsCommander {
	return &claimRewardsHistoryCommand{historyRepo: historyRepo}
}
