package service

import (
	"time"

	coupondto "github.com/Mahoura-shop/Backend/internal/application/dto/coupon"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	domainPostgres "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CouponService struct {
	couponRepository domainPostgres.CouponRepository
	db               database.Database
}

type CouponServiceDeps struct {
	CouponRepository domainPostgres.CouponRepository
	DB               database.Database
}

func NewCouponService(deps CouponServiceDeps) *CouponService {
	return &CouponService{
		couponRepository: deps.CouponRepository,
		db:               deps.DB,
	}
}

func (s *CouponService) parseCoupon(c entity.Coupon) coupondto.CouponCredential {
	return coupondto.CouponCredential{
		ID:              c.ID,
		Code:            c.Code,
		DiscountPercent: c.DiscountPercent,
		MinIRRPrice:     c.MinIRRPrice,
		MaxUses:         c.MaxUses,
		UsedCount:       c.UsedCount,
		ExpiresAt:       c.ExpiresAt,
		IsActive:        c.IsActive,
	}
}

func (s *CouponService) CreateCoupon(req coupondto.CreateCouponRequest) error {
	existing, err := s.couponRepository.FindCouponByCode(s.db, req.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		var ce exception.ConflictErrors
		ce.Add("code", "alreadyExist")
		return ce
	}
	coupon := entity.Coupon{
		Code:            req.Code,
		DiscountPercent: req.DiscountPercent,
		MinIRRPrice:     req.MinIRRPrice,
		MaxUses:         req.MaxUses,
		ExpiresAt:       req.ExpiresAt,
		IsActive:        true,
	}
	_, err = s.couponRepository.CreateCoupon(s.db, coupon)
	return err
}

func (s *CouponService) ValidateCoupon(code string, irrPrice uint) (*coupondto.CouponValidationResult, error) {
	coupon, err := s.couponRepository.FindCouponByCode(s.db, code)
	if err != nil {
		return nil, err
	}
	if coupon == nil || !coupon.IsActive {
		return nil, exception.NotFoundError{Item: "coupon"}
	}
	if coupon.ExpiresAt != nil && time.Now().After(*coupon.ExpiresAt) {
		return nil, exception.NotFoundError{Item: "coupon"}
	}
	if coupon.MaxUses > 0 && coupon.UsedCount >= coupon.MaxUses {
		return nil, exception.NotFoundError{Item: "coupon"}
	}
	if irrPrice < coupon.MinIRRPrice {
		var ve exception.ValidationErrors
		ve.Add("irrPrice", "belowMinimum")
		return nil, ve
	}

	discount := uint(float64(irrPrice) * coupon.DiscountPercent / 100)
	return &coupondto.CouponValidationResult{
		Code:            coupon.Code,
		DiscountPercent: coupon.DiscountPercent,
		DiscountAmount:  discount,
		FinalPrice:      irrPrice - discount,
	}, nil
}

func (s *CouponService) GetCoupons() ([]coupondto.CouponCredential, error) {
	coupons, err := s.couponRepository.GetCoupons(s.db)
	if err != nil {
		return nil, err
	}
	var result []coupondto.CouponCredential
	for _, c := range coupons {
		result = append(result, s.parseCoupon(*c))
	}
	return result, nil
}

func (s *CouponService) UpdateCoupon(req coupondto.UpdateCouponRequest) error {
	coupon, err := s.couponRepository.FindCouponByID(s.db, req.ID)
	if err != nil {
		return err
	}
	if coupon == nil {
		return exception.NotFoundError{Item: "coupon"}
	}
	if req.DiscountPercent != nil {
		coupon.DiscountPercent = *req.DiscountPercent
	}
	if req.MinIRRPrice != nil {
		coupon.MinIRRPrice = *req.MinIRRPrice
	}
	if req.MaxUses != nil {
		coupon.MaxUses = *req.MaxUses
	}
	if req.ExpiresAt != nil {
		coupon.ExpiresAt = req.ExpiresAt
	}
	if req.IsActive != nil {
		coupon.IsActive = *req.IsActive
	}
	return s.couponRepository.UpdateCoupon(s.db, *coupon)
}

func (s *CouponService) DeleteCoupon(couponID uint) error {
	coupon, err := s.couponRepository.FindCouponByID(s.db, couponID)
	if err != nil {
		return err
	}
	if coupon == nil {
		return exception.NotFoundError{Item: "coupon"}
	}
	return s.couponRepository.DeleteCouponByID(s.db, couponID)
}
