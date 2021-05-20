package query

import "context"

type GetRewardsLimitQueryer interface {
	Do(ctx context.Context, id string) (int64, error)
}
