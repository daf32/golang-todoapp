package domain

import "time"

type UserOAuthIdentity struct {
	ID          int
	UserID      int
	Provider    string
	ProviderSub string
	Email       string
	CreateAt    time.Time
}
