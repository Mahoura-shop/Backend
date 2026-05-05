package usecase

import coupondto "github.com/Mahoura-shop/Backend/internal/application/dto/coupon"

type CouponService interface {
	CreateCoupon(coupondto.CreateCouponRequest) error
	ValidateCoupon(code string, irrPrice uint) (*coupondto.CouponValidationResult, error)
	GetCoupons() ([]coupondto.CouponCredential, error)
	UpdateCoupon(coupondto.UpdateCouponRequest) error
	DeleteCoupon(uint) error
}
