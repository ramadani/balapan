package command

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ramadani/balapan/internal/app/command"
	"time"
)

type setRewardsLimitRedisCommand struct {
	redisClient *redis.Client
	keyFormat   string
	ttl         time.Duration
}

func (c *setRewardsLimitRedisCommand) Do(ctx context.Context, id string, value int64) error {
	key := fmt.Sprintf(c.keyFormat, id)

	_, err := c.redisClient.Set(ctx, key, value, c.ttl).Result()

	return err
}

func NewSetRewardsLimitRedisCommand(redisClient *redis.Client, keyFormat string, ttl time.Duration) command.SetRewardsLimitCommander {
	return &setRewardsLimitRedisCommand{
		redisClient: redisClient,
		keyFormat:   keyFormat,
		ttl:         ttl,
	}
}
