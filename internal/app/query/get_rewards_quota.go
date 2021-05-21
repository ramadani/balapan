package query

import "context"

type GetRewardsQuotaQueryer interface {
	Do(ctx context.Context, id string) (int64, error)
}
