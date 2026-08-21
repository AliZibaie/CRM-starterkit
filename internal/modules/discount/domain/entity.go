package domain

import "time"

const (
	TypePercent = "percent"
	TypeFixed   = "fixed"
)

type Discount struct {
	ID             string
	Code           string
	Type           string
	Value          int64
	MinOrderAmount int64
	MaxUses        int64
	UsedCount      int64
	StartsAt       *time.Time
	EndsAt         *time.Time
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ApplyResult struct {
	Code           string
	DiscountAmount int64
	FinalAmount    int64
}
