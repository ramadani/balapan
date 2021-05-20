package command

import (
	"context"
	"fmt"
	"github.com/ramadani/balapan/internal/app/command/model"
	"time"
)

type usageRewardsLoggerCommand struct {
	next UsageRewardsCommander
}

func (c *usageRewardsLoggerCommand) Do(ctx context.Context, data *model.UsageRewards) error {
	startAt := time.Now()

	err := c.next.Do(ctx, data)

	rt := time.Now().Sub(startAt).Milliseconds()

	fmt.Println(rt)
	return err
}

func NewUsageRewardsLoggerCommand(next UsageRewardsCommander) UsageRewardsCommander {
	return &usageRewardsLoggerCommand{next: next}
}
