package sqlx

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/ramadani/balapan/internal/domain/history"
)

type historyRepository struct {
	db *sqlx.DB
}

func (r *historyRepository) Create(ctx context.Context, data *history.History) error {
	query := "INSERT INTO histories (id, rewards_id, user_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.ExecContext(ctx, query, data.ID, data.RewardsID, data.UserID, data.Amount, data.CreatedAt)

	return err
}

func NewHistoryRepository(db *sqlx.DB) history.Repository {
	return &historyRepository{db: db}
}
