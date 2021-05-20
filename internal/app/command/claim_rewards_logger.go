package command

import (
	"context"
	"fmt"
	"github.com/ramadani/balapan/internal/app/command/model"
	"time"
)

type claimRewardsLoggerCommand struct {
	next ClaimRewardsCommander
}

func (c *claimRewardsLoggerCommand) Do(ctx context.Context, data *model.ClaimRewards) error {
	startAt := time.Now()

	err := c.next.Do(ctx, data)

	rt := time.Now().Sub(startAt).Milliseconds()

	fmt.Println(rt)
	return err
}

func NewClaimRewardsLoggerCommand(next ClaimRewardsCommander) ClaimRewardsCommander {
	return &claimRewardsLoggerCommand{next: next}
}
