package rewards

import "errors"

var (
	ErrTransactionQuotaLimitExceeded = errors.New("rewards: transaction quota limit exceeded")
	ErrRewardsQuotaLimitExceeded     = errors.New("rewards: rewards quota limit exceeded")
)
