package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type CouponRepositoryMock struct {
	mock.Mock
}

func NewCouponRepositoryMock() *CouponRepositoryMock {
	return &CouponRepositoryMock{}
}

func (m *CouponRepositoryMock) CreateCoupon(db database.Database, coupon entity.Coupon) (*entity.Coupon, error) {
	args := m.Called(db, coupon)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Coupon), args.Error(1)
}

func (m *CouponRepositoryMock) FindCouponByCode(db database.Database, code string) (*entity.Coupon, error) {
	args := m.Called(db, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Coupon), args.Error(1)
}

func (m *CouponRepositoryMock) FindCouponByID(db database.Database, id uint) (*entity.Coupon, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Coupon), args.Error(1)
}

func (m *CouponRepositoryMock) GetCoupons(db database.Database) ([]*entity.Coupon, error) {
	args := m.Called(db)
	return args.Get(0).([]*entity.Coupon), args.Error(1)
}

func (m *CouponRepositoryMock) UpdateCoupon(db database.Database, coupon entity.Coupon) error {
	args := m.Called(db, coupon)
	return args.Error(0)
}

func (m *CouponRepositoryMock) DeleteCouponByID(db database.Database, id uint) error {
	args := m.Called(db, id)
	return args.Error(0)
}

func (m *CouponRepositoryMock) IncrementUsedCount(db database.Database, id uint) error {
	args := m.Called(db, id)
	return args.Error(0)
}
