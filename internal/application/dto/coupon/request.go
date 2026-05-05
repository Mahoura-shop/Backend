package coupondto

import "time"

type CreateCouponRequest struct {
	Code            string     `json:"code" validate:"required"`
	DiscountPercent float64    `json:"discountPercent" validate:"required,min=0.01,max=100"`
	MinIRRPrice     uint       `json:"minIRRPrice"`
	MaxUses         uint       `json:"maxUses"`
	ExpiresAt       *time.Time `json:"expiresAt"`
}

type UpdateCouponRequest struct {
	ID              uint
	DiscountPercent *float64   `json:"discountPercent" validate:"omitempty,min=0.01,max=100"`
	MinIRRPrice     *uint      `json:"minIRRPrice"`
	MaxUses         *uint      `json:"maxUses"`
	ExpiresAt       *time.Time `json:"expiresAt"`
	IsActive        *bool      `json:"isActive"`
}

type ValidateCouponRequest struct {
	Code      string `json:"code" validate:"required"`
	IRRPrice  uint   `json:"irrPrice" validate:"required"`
}
