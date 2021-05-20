package history

import "time"

type History struct {
	ID        string    `db:"id"`
	RewardsID string    `db:"rewards_id"`
	UserID    string    `db:"user_id"`
	Amount    int64     `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}
