package mocks

import (
	coupondto "github.com/Mahoura-shop/Backend/internal/application/dto/coupon"
	"github.com/stretchr/testify/mock"
)

type CouponServiceMock struct {
	mock.Mock
}

func NewCouponServiceMock() *CouponServiceMock {
	return &CouponServiceMock{}
}

func (m *CouponServiceMock) CreateCoupon(req coupondto.CreateCouponRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *CouponServiceMock) ValidateCoupon(code string, irrPrice uint) (*coupondto.CouponValidationResult, error) {
	args := m.Called(code, irrPrice)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*coupondto.CouponValidationResult), args.Error(1)
}

func (m *CouponServiceMock) GetCoupons() ([]coupondto.CouponCredential, error) {
	args := m.Called()
	return args.Get(0).([]coupondto.CouponCredential), args.Error(1)
}

func (m *CouponServiceMock) UpdateCoupon(req coupondto.UpdateCouponRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *CouponServiceMock) DeleteCoupon(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
