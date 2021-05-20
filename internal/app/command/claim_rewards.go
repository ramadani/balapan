package command

import (
	"context"
	"fmt"
	"github.com/ramadani/balapan/internal/app/command/model"
	"github.com/ramadani/balapan/internal/domain/rewards"
)

type ClaimRewardsCommander interface {
	Do(ctx context.Context, data *model.ClaimRewards) error
}

type claimRewardsCommand struct {
	rewardsRepo rewards.Repository
}

func (c *claimRewardsCommand) Do(ctx context.Context, data *model.ClaimRewards) error {
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

func NewClaimRewardsCommand(rewardsRepo rewards.Repository) ClaimRewardsCommander {
	return &claimRewardsCommand{rewardsRepo: rewardsRepo}
}
