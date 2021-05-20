package command

import (
	"context"
	"fmt"
	"github.com/ramadani/balapan/internal/app/command/model"
	"github.com/ramadani/balapan/internal/domain/rewards"
)

type claimRewardsUsageCommand struct {
	rewardsRepo rewards.Repository
}

func (c *claimRewardsUsageCommand) Do(ctx context.Context, data *model.ClaimRewards) error {
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

func NewClaimRewardsUsageCommand(rewardsRepo rewards.Repository) ClaimRewardsCommander {
	return &claimRewardsUsageCommand{rewardsRepo: rewardsRepo}
}
