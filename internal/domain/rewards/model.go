package rewards

type Rewards struct {
	ID               string `db:"id"`
	TransactionLimit int64  `db:"transaction_limit"`
	TransactionUsage int64  `db:"transaction_usage"`
	RewardsLimit     int64  `db:"rewards_limit"`
	RewardsUsage     int64  `db:"rewards_usage"`
}
