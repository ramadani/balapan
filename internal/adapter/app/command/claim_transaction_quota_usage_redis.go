package command

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ramadani/balapan/internal/app/command"
	"github.com/ramadani/balapan/internal/app/command/model"
	"github.com/ramadani/balapan/internal/app/query"
	"github.com/ramadani/balapan/internal/domain/rewards"
)

type claimTransactionQuotaUsageRedisCommand struct {
	next                 command.ClaimRewardsCommander
	getRewardsQuotaLimit query.GetRewardsQuotaQueryer
	redisClient          *redis.Client
	keyFormat            string
}

func (c *claimTransactionQuotaUsageRedisCommand) Do(ctx context.Context, data *model.ClaimRewards) (err error) {
	key := fmt.Sprintf(c.keyFormat, data.ID)

	limit, err := c.getRewardsQuotaLimit.Do(ctx, data.ID)
	if err != nil {
		return
	}

	usage, err := c.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			if _, er := c.redisClient.Decr(ctx, key).Result(); er != nil {
				err = er
			}
		}
	}()

	if usage > limit {
		err = rewards.ErrTransactionQuotaLimitExceeded
		return
	}

	err = c.next.Do(ctx, data)
	return
}

func NewClaimTransactionQuotaUsageRedisCommand(
	next command.ClaimRewardsCommander,
	getRewardsLimitQuery query.GetRewardsQuotaQueryer,
	redisClient *redis.Client,
	keyFormat string,
) command.ClaimRewardsCommander {
	return &claimTransactionQuotaUsageRedisCommand{
		next:                 next,
		getRewardsQuotaLimit: getRewardsLimitQuery,
		redisClient:          redisClient,
		keyFormat:            keyFormat,
	}
}
