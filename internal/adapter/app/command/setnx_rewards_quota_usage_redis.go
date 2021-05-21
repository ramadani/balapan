package command

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ramadani/balapan/internal/app/command"
	"github.com/ramadani/balapan/internal/domain/rewards"
	"time"
)

type GetQuotaUsageFunc func(rewards2 *rewards.Rewards) int64

type setNXRewardsQuotaUsageRedisCommand struct {
	rewardsRepo       rewards.Repository
	redisClient       *redis.Client
	keyFormat         string
	usageExpIn        time.Duration
	lockExpIn         time.Duration
	getQuotaUsageFunc GetQuotaUsageFunc
}

func (c *setNXRewardsQuotaUsageRedisCommand) Do(ctx context.Context, id string) error {
	key := fmt.Sprintf(c.keyFormat, id)

	exists, err := c.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return err
	} else if exists == 1 {
		return nil
	}

	lockKey := fmt.Sprintf("%s-lock", key)
	lockVal, err := c.redisClient.GetSet(ctx, lockKey, 1).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	if lockVal != "" {
		return fmt.Errorf("lock: key has been locked")
	}

	if _, err = c.redisClient.Expire(ctx, lockKey, c.lockExpIn).Result(); err != nil {
		return err
	}

	rewards, err := c.rewardsRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if _, err = c.redisClient.SetNX(ctx, key, c.getQuotaUsageFunc(rewards), c.usageExpIn).Result(); err != nil {
		return err
	}

	if _, er := c.redisClient.Del(ctx, lockKey).Result(); er != nil {
		err = er
	}
	return err
}

func NewSetNXRewardsQuotaUsageRedisCommand(
	rewardsRepo rewards.Repository,
	redisClient *redis.Client,
	keyFormat string,
	usageExpIn time.Duration,
	lockExpIn time.Duration,
	getQuotaUsageFunc GetQuotaUsageFunc,
) command.SetNXRewardsQuotaCommander {
	return &setNXRewardsQuotaUsageRedisCommand{
		rewardsRepo:       rewardsRepo,
		redisClient:       redisClient,
		keyFormat:         keyFormat,
		usageExpIn:        usageExpIn,
		lockExpIn:         lockExpIn,
		getQuotaUsageFunc: getQuotaUsageFunc,
	}
}
