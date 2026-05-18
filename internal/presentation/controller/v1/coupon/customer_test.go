package coupon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	coupondto "github.com/Mahoura-shop/Backend/internal/application/dto/coupon"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type CustomerCouponTestSuite struct {
	suite.Suite
	constants     *bootstrap.Constants
	couponService *mocks.CouponServiceMock
	controller    *CustomerCouponController
}

func (s *CustomerCouponTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.couponService = mocks.NewCouponServiceMock()
	s.controller = NewCustomerCouponController(s.constants, s.couponService)
}

func (s *CustomerCouponTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Next()
	})
	if withRecovery {
		r.Use(couponPanicToHTTP())
	}
	r.POST("/coupons/validate", s.controller.ValidateCoupon)
	return r
}

func (s *CustomerCouponTestSuite) TestValidateCoupon_Returns200() {
	result := &coupondto.CouponValidationResult{
		Code:            "SAVE10",
		DiscountPercent: 10,
		DiscountAmount:  100000,
		FinalPrice:      900000,
	}
	s.couponService.On("ValidateCoupon", "SAVE10", uint(1000000)).Return(result, nil).Once()

	body, _ := json.Marshal(map[string]any{"code": "SAVE10", "irrPrice": 1000000})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/coupons/validate", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.couponService.AssertExpectations(s.T())
}

func (s *CustomerCouponTestSuite) TestValidateCoupon_NotFound_NonOK() {
	s.couponService.On("ValidateCoupon", "GHOST", uint(1000000)).Return((*coupondto.CouponValidationResult)(nil), exception.NotFoundError{Item: "coupon"}).Once()

	body, _ := json.Marshal(map[string]any{"code": "GHOST", "irrPrice": 1000000})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/coupons/validate", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(true).ServeHTTP(w, r)

	s.NotEqual(http.StatusOK, w.Code)
	s.couponService.AssertExpectations(s.T())
}

func TestCustomerCouponSuite(t *testing.T) {
	suite.Run(t, new(CustomerCouponTestSuite))
}
