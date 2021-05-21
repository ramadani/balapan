package command

import (
	"context"
	"github.com/ramadani/balapan/internal/app/command/model"
)

type claimRewardsSetNxQuotaCommand struct {
	setNXRewardsQuotaCommand SetNXRewardsQuotaCommander
}

func (c *claimRewardsSetNxQuotaCommand) Do(ctx context.Context, data *model.ClaimRewards) error {
	return c.setNXRewardsQuotaCommand.Do(ctx, data.ID)
}

func NewClaimRewardsSetNXQuotaCommand(setNXRewardsQuotaCommand SetNXRewardsQuotaCommander) ClaimRewardsCommander {
	return &claimRewardsSetNxQuotaCommand{setNXRewardsQuotaCommand: setNXRewardsQuotaCommand}
}
