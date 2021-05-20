package command

import (
	"context"
	"github.com/ramadani/balapan/internal/app/command/model"
	"time"
)

type usageRewardsSleeperCommand struct {
	next    UsageRewardsCommander
	sleepIn time.Duration
}

func (c *usageRewardsSleeperCommand) Do(ctx context.Context, data *model.UsageRewards) error {
	time.Sleep(c.sleepIn)

	return c.next.Do(ctx, data)
}

func NewUsageRewardsSleeperCommand(next UsageRewardsCommander, sleepIn time.Duration) UsageRewardsCommander {
	return &usageRewardsSleeperCommand{
		next:    next,
		sleepIn: sleepIn,
	}
}
