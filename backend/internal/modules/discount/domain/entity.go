package domain

import "time"

const (
	TypePercent = "percent"
	TypeFixed   = "fixed"
)

type Discount struct {
	ID             string     `json:"id"`
	Code           string     `json:"code"`
	Type           string     `json:"type"`
	Value          int64      `json:"value"`
	MinOrderAmount int64      `json:"min_order_amount"`
	MaxUses        int64      `json:"max_uses"`
	UsedCount      int64      `json:"used_count"`
	StartsAt       *time.Time `json:"starts_at,omitempty"`
	EndsAt         *time.Time `json:"ends_at,omitempty"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ApplyResult struct {
	Code           string `json:"code"`
	DiscountAmount int64  `json:"discount_amount"`
	FinalAmount    int64  `json:"final_amount"`
}
