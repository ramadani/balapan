package command

import (
	"context"
	"fmt"
	"github.com/ramadani/balapan/internal/app/command/model"
	"github.com/ramadani/balapan/internal/domain/rewards"
)

type UsageRewardsCommander interface {
	Do(ctx context.Context, data *model.UsageRewards) error
}

type usageRewardsCommand struct {
	rewardsRepo rewards.Repository
}

func (c *usageRewardsCommand) Do(ctx context.Context, data *model.UsageRewards) error {
	rewards, err := c.rewardsRepo.FindByID(ctx, data.ID)
	if err != nil {
		return err
	}

	rewards.TransactionUsage += 1
	rewards.RewardsUsage += data.Amount

	if rewards.TransactionUsage > rewards.TransactionLimit || rewards.RewardsUsage > rewards.RewardsLimit {
		return fmt.Errorf("quota limit exceeded")
	}

	err = c.rewardsRepo.Update(ctx, rewards)

	return err
}

func NewUsageRewardsCommand(rewardsRepo rewards.Repository) UsageRewardsCommander {
	return &usageRewardsCommand{rewardsRepo: rewardsRepo}
}
