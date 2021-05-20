package sqlx

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/ramadani/balapan/internal/domain/rewards"
)

type rewardsRepository struct {
	db *sqlx.DB
}

func (r *rewardsRepository) FindAll(ctx context.Context) ([]*rewards.Rewards, error) {
	res := make([]*rewards.Rewards, 0)
	query := "SELECT * FROM rewards"
	err := r.db.SelectContext(ctx, &res, query)

	return res, err
}

func (r *rewardsRepository) FindByID(ctx context.Context, id string) (*rewards.Rewards, error) {
	res := &rewards.Rewards{}
	query := "SELECT * FROM rewards WHERE id = $1"
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(res)

	return res, err
}

func (r *rewardsRepository) Update(ctx context.Context, data *rewards.Rewards) error {
	query := "UPDATE rewards SET transaction_limit = $1, transaction_usage = $2, rewards_limit = $3, rewards_usage = $4 WHERE id = $5"
	_, err := r.db.ExecContext(ctx, query, data.TransactionLimit, data.TransactionUsage, data.RewardsLimit, data.RewardsUsage, data.ID)

	return err
}

func NewRewardsRepository(db *sqlx.DB) rewards.Repository {
	return &rewardsRepository{db: db}
}
