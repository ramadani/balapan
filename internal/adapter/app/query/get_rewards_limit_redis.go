package query

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ramadani/balapan/internal/app/query"
	"strconv"
)

type getRewardsLimitRedisQueryer struct {
	redisClient *redis.Client
	keyFormat   string
}

func (q *getRewardsLimitRedisQueryer) Do(ctx context.Context, id string) (int64, error) {
	key := fmt.Sprintf(q.keyFormat, id)

	limit, err := q.redisClient.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	res, _ := strconv.Atoi(limit)

	return int64(res), nil
}

func NewGetRewardsLimitRedisQueryer(redisClient *redis.Client, keyFormat string) query.GetRewardsLimitQueryer {
	return &getRewardsLimitRedisQueryer{
		redisClient: redisClient,
		keyFormat:   keyFormat,
	}
}
