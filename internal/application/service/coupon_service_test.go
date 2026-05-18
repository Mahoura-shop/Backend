package service

import (
	"errors"
	"testing"
	"time"

	coupondto "github.com/Mahoura-shop/Backend/internal/application/dto/coupon"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type CouponServiceTestSuite struct {
	suite.Suite
	couponRepository *mocks.CouponRepositoryMock
	db               *mocks.DatabaseMock
	service          *CouponService
}

func (s *CouponServiceTestSuite) SetupTest() {
	s.couponRepository = mocks.NewCouponRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewCouponService(CouponServiceDeps{
		CouponRepository: s.couponRepository,
		DB:               s.db,
	})
}

func activeCoupon() *entity.Coupon {
	return &entity.Coupon{
		Model:           dbmodel.Model{ID: 1},
		Code:            "SAVE20",
		DiscountPercent: 20,
		MinIRRPrice:     1_000_000,
		MaxUses:         0,
		UsedCount:       0,
		IsActive:        true,
	}
}

func (s *CouponServiceTestSuite) TestValidateCoupon_Success() {
	s.couponRepository.On("FindCouponByCode", s.db, "SAVE20").Return(activeCoupon(), nil).Once()

	result, err := s.service.ValidateCoupon("SAVE20", 2_000_000)

	s.NoError(err)
	s.NotNil(result)
	s.Equal("SAVE20", result.Code)
	s.Equal(float64(20), result.DiscountPercent)
	s.Equal(uint(400_000), result.DiscountAmount)
	s.Equal(uint(1_600_000), result.FinalPrice)
	s.couponRepository.AssertExpectations(s.T())
}

func (s *CouponServiceTestSuite) TestValidateCoupon_NotFound() {
	var nilCoupon *entity.Coupon
	s.couponRepository.On("FindCouponByCode", s.db, "GHOST").Return(nilCoupon, nil).Once()

	result, err := s.service.ValidateCoupon("GHOST", 2_000_000)

	s.Nil(result)
	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.couponRepository.AssertExpectations(s.T())
}

func (s *CouponServiceTestSuite) TestValidateCoupon_Inactive() {
	coupon := activeCoupon()
	coupon.IsActive = false
	s.couponRepository.On("FindCouponByCode", s.db, "SAVE20").Return(coupon, nil).Once()

	result, err := s.service.ValidateCoupon("SAVE20", 2_000_000)

	s.Nil(result)
	s.Error(err)
	s.couponRepository.AssertExpectations(s.T())
}

func (s *CouponServiceTestSuite) TestValidateCoupon_Expired() {
	coupon := activeCoupon()
	past := time.Now().Add(-time.Hour)
	coupon.ExpiresAt = &past
	s.couponRepository.On("FindCouponByCode", s.db, "SAVE20").Return(coupon, nil).Once()

	result, err := s.service.ValidateCoupon("SAVE20", 2_000_000)

	s.Nil(result)
	s.Error(err)
	s.couponRepository.AssertExpectations(s.T())
}

func (s *CouponServiceTestSuite) TestValidateCoupon_MaxUsesReached() {
	coupon := activeCoupon()
	coupon.MaxUses = 5
	coupon.UsedCount = 5
	s.couponRepository.On("FindCouponByCode", s.db, "SAVE20").Return(coupon, nil).Once()

	result, err := s.service.ValidateCoupon("SAVE20", 2_000_000)

	s.Nil(result)
	s.Error(err)
	s.couponRepository.AssertExpectations(s.T())
}

func (s *CouponServiceTestSuite) TestValidateCoupon_BelowMinPrice() {
	s.couponRepository.On("FindCouponByCode", s.db, "SAVE20").Return(activeCoupon(), nil).Once()

	result, err := s.service.ValidateCoupon("SAVE20", 500_000)

	s.Nil(result)
	s.Error(err)
	var ve exception.ValidationErrors
	s.True(errors.As(err, &ve))
	s.couponRepository.AssertExpectations(s.T())
}

func (s *CouponServiceTestSuite) TestValidateCoupon_NotExpiredWhenFuture() {
	coupon := activeCoupon()
	future := time.Now().Add(time.Hour * 24)
	coupon.ExpiresAt = &future
	s.couponRepository.On("FindCouponByCode", s.db, "SAVE20").Return(coupon, nil).Once()

	result, err := s.service.ValidateCoupon("SAVE20", 2_000_000)

	s.NoError(err)
	s.NotNil(result)
	s.couponRepository.AssertExpectations(s.T())
}

func (s *CouponServiceTestSuite) TestCreateCoupon_Success() {
	var nilCoupon *entity.Coupon
	s.couponRepository.On("FindCouponByCode", s.db, "NEW10").Return(nilCoupon, nil).Once()
	s.couponRepository.On("CreateCoupon", s.db, entity.Coupon{
		Code:            "NEW10",
		DiscountPercent: 10,
		MinIRRPrice:     500_000,
		MaxUses:         100,
		IsActive:        true,
	}).Return(&entity.Coupon{}, nil).Once()

	err := s.service.CreateCoupon(coupondto.CreateCouponRequest{
		Code:            "NEW10",
		DiscountPercent: 10,
		MinIRRPrice:     500_000,
		MaxUses:         100,
	})

	s.NoError(err)
	s.couponRepository.AssertExpectations(s.T())
}

func (s *CouponServiceTestSuite) TestCreateCoupon_DuplicateCode() {
	s.couponRepository.On("FindCouponByCode", s.db, "SAVE20").Return(activeCoupon(), nil).Once()

	err := s.service.CreateCoupon(coupondto.CreateCouponRequest{
		Code:            "SAVE20",
		DiscountPercent: 10,
	})

	s.Error(err)
	var ce exception.ConflictErrors
	s.True(errors.As(err, &ce))
	s.couponRepository.AssertNotCalled(s.T(), "CreateCoupon")
	s.couponRepository.AssertExpectations(s.T())
}

func (s *CouponServiceTestSuite) TestDeleteCoupon_Success() {
	s.couponRepository.On("FindCouponByID", s.db, uint(1)).Return(activeCoupon(), nil).Once()
	s.couponRepository.On("DeleteCouponByID", s.db, uint(1)).Return(nil).Once()

	err := s.service.DeleteCoupon(1)

	s.NoError(err)
	s.couponRepository.AssertExpectations(s.T())
}

func (s *CouponServiceTestSuite) TestDeleteCoupon_NotFound() {
	var nilCoupon *entity.Coupon
	s.couponRepository.On("FindCouponByID", s.db, uint(99)).Return(nilCoupon, nil).Once()

	err := s.service.DeleteCoupon(99)

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.couponRepository.AssertNotCalled(s.T(), "DeleteCouponByID")
	s.couponRepository.AssertExpectations(s.T())
}

func (s *CouponServiceTestSuite) TestGetCoupons_ReturnsList() {
	coupons := []*entity.Coupon{activeCoupon()}
	s.couponRepository.On("GetCoupons", s.db).Return(coupons, nil).Once()

	result, err := s.service.GetCoupons()

	s.NoError(err)
	s.Len(result, 1)
	s.Equal("SAVE20", result[0].Code)
	s.couponRepository.AssertExpectations(s.T())
}

func TestCouponServiceSuite(t *testing.T) {
	suite.Run(t, new(CouponServiceTestSuite))
}
