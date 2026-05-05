package coupondto

import "time"

type CouponCredential struct {
	ID              uint       `json:"id"`
	Code            string     `json:"code"`
	DiscountPercent float64    `json:"discountPercent"`
	MinIRRPrice     uint       `json:"minIRRPrice"`
	MaxUses         uint       `json:"maxUses"`
	UsedCount       uint       `json:"usedCount"`
	ExpiresAt       *time.Time `json:"expiresAt"`
	IsActive        bool       `json:"isActive"`
}

type CouponValidationResult struct {
	Code            string  `json:"code"`
	DiscountPercent float64 `json:"discountPercent"`
	DiscountAmount  uint    `json:"discountAmount"`
	FinalPrice      uint    `json:"finalPrice"`
}
