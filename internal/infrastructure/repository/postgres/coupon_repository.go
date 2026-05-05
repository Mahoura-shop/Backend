package postgres

import (
	"errors"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type CouponRepository struct{}

func NewCouponRepository() *CouponRepository {
	return &CouponRepository{}
}

func (r *CouponRepository) CreateCoupon(db database.Database, coupon entity.Coupon) (*entity.Coupon, error) {
	result := db.GetDB().Create(&coupon)
	if result.Error != nil {
		return nil, result.Error
	}
	return &coupon, nil
}

func (r *CouponRepository) FindCouponByCode(db database.Database, code string) (*entity.Coupon, error) {
	var coupon entity.Coupon
	result := db.GetDB().Where("code = ?", code).First(&coupon)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &coupon, nil
}

func (r *CouponRepository) FindCouponByID(db database.Database, id uint) (*entity.Coupon, error) {
	var coupon entity.Coupon
	result := db.GetDB().Where("id = ?", id).First(&coupon)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &coupon, nil
}

func (r *CouponRepository) GetCoupons(db database.Database) ([]*entity.Coupon, error) {
	var coupons []*entity.Coupon
	result := db.GetDB().Order("created_at DESC").Find(&coupons)
	if result.Error != nil {
		return nil, result.Error
	}
	return coupons, nil
}

func (r *CouponRepository) UpdateCoupon(db database.Database, coupon entity.Coupon) error {
	return db.GetDB().Save(&coupon).Error
}

func (r *CouponRepository) DeleteCouponByID(db database.Database, id uint) error {
	return db.GetDB().Where("id = ?", id).Delete(&entity.Coupon{}).Error
}

func (r *CouponRepository) IncrementUsedCount(db database.Database, couponID uint) error {
	return db.GetDB().Model(&entity.Coupon{}).Where("id = ?", couponID).UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}
