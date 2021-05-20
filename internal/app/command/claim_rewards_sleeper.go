package command

import (
	"context"
	"github.com/ramadani/balapan/internal/app/command/model"
	"time"
)

type claimRewardsSleeperCommand struct {
	next    ClaimRewardsCommander
	sleepIn time.Duration
}

func (c *claimRewardsSleeperCommand) Do(ctx context.Context, data *model.ClaimRewards) error {
	time.Sleep(c.sleepIn)

	return c.next.Do(ctx, data)
}

func NewClaimRewardsSleeperCommand(next ClaimRewardsCommander, sleepIn time.Duration) ClaimRewardsCommander {
	return &claimRewardsSleeperCommand{
		next:    next,
		sleepIn: sleepIn,
	}
}
