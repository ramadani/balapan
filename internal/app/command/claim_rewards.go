package command

import (
	"context"
	"github.com/ramadani/balapan/internal/app/command/model"
)

type ClaimRewardsCommander interface {
	Do(ctx context.Context, data *model.ClaimRewards) error
}

type claimRewardsMiddlewareCommand struct {
	prev, next ClaimRewardsCommander
}

func (c *claimRewardsMiddlewareCommand) Do(ctx context.Context, data *model.ClaimRewards) error {
	if err := c.prev.Do(ctx, data); err != nil {
		return err
	}

	return c.next.Do(ctx, data)
}

func NewClaimRewardsMiddlewareCommand(prev, next ClaimRewardsCommander) ClaimRewardsCommander {
	return &claimRewardsMiddlewareCommand{
		prev: prev,
		next: next,
	}
}
