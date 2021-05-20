package model

type ClaimRewardsRequest struct {
	UserID string `json:"userId"`
	Amount int64  `json:"amount"`
}
