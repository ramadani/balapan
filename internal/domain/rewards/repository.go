package rewards

import "context"

type Repository interface {
	FindAll(ctx context.Context) ([]*Rewards, error)
	FindByID(ctx context.Context, id string) (*Rewards, error)
	Update(ctx context.Context, data *Rewards) error
}
