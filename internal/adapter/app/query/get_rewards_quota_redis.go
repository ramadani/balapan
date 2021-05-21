package query

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ramadani/balapan/internal/app/query"
	"strconv"
)

type getRewardsQuotaRedisQueryer struct {
	redisClient *redis.Client
	keyFormat   string
}

func (q *getRewardsQuotaRedisQueryer) Do(ctx context.Context, id string) (int64, error) {
	key := fmt.Sprintf(q.keyFormat, id)

	limit, err := q.redisClient.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	res, _ := strconv.Atoi(limit)

	return int64(res), nil
}

func NewGetRewardsQuotaRedisQueryer(redisClient *redis.Client, keyFormat string) query.GetRewardsQuotaQueryer {
	return &getRewardsQuotaRedisQueryer{
		redisClient: redisClient,
		keyFormat:   keyFormat,
	}
}
