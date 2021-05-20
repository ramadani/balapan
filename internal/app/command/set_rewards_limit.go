package command

import (
	"context"
)

type SetRewardsLimitCommander interface {
	Do(ctx context.Context, id string, value int64) error
}
