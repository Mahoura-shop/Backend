package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CouponRepository interface {
	CreateCoupon(database.Database, entity.Coupon) (*entity.Coupon, error)
	FindCouponByCode(database.Database, string) (*entity.Coupon, error)
	FindCouponByID(database.Database, uint) (*entity.Coupon, error)
	GetCoupons(database.Database) ([]*entity.Coupon, error)
	UpdateCoupon(database.Database, entity.Coupon) error
	DeleteCouponByID(database.Database, uint) error
	IncrementUsedCount(database.Database, uint) error
}
