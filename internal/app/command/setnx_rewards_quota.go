package command

import (
	"context"
	"fmt"
	"time"
)

type SetNXRewardsQuotaCommander interface {
	Do(ctx context.Context, id string) error
}

type setNXRewardsQuotaRetryableCommand struct {
	next     SetNXRewardsQuotaCommander
	maxRetry int
	sleepIn  time.Duration
}

func (c *setNXRewardsQuotaRetryableCommand) Do(ctx context.Context, id string) error {
	for i := 0; i < c.maxRetry; i++ {
		if err := c.next.Do(ctx, id); err == nil {
			return nil
		}

		time.Sleep(c.sleepIn)
	}

	return fmt.Errorf("max retry")
}

func NewSetNXRewardsQuotaRetryableCommand(next SetNXRewardsQuotaCommander, maxRetry int, sleepIn time.Duration) SetNXRewardsQuotaCommander {
	return &setNXRewardsQuotaRetryableCommand{
		next:     next,
		maxRetry: maxRetry,
		sleepIn:  sleepIn,
	}
}
