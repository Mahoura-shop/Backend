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

func init() { gin.SetMode(gin.TestMode) }

const testCouponID = uint(7)

type AdminCouponTestSuite struct {
	suite.Suite
	constants     *bootstrap.Constants
	couponService *mocks.CouponServiceMock
	controller    *AdminCouponController
}

func (s *AdminCouponTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.couponService = mocks.NewCouponServiceMock()
	s.controller = NewAdminCouponController(s.constants, s.couponService)
}

func (s *AdminCouponTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Next()
	})
	if withRecovery {
		r.Use(couponPanicToHTTP())
	}
	r.GET("/admin/coupons", s.controller.GetCoupons)
	r.POST("/admin/coupons", s.controller.CreateCoupon)
	r.PATCH("/admin/coupons/:couponID", s.controller.UpdateCoupon)
	r.DELETE("/admin/coupons/:couponID", s.controller.DeleteCoupon)
	return r
}

func couponPanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case exception.NotFoundError:
					c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
				case exception.ConflictErrors:
					c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

func (s *AdminCouponTestSuite) TestGetCoupons_Returns200() {
	s.couponService.On("GetCoupons").Return([]coupondto.CouponCredential{}, nil).Once()

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/coupons", nil))

	s.Equal(http.StatusOK, w.Code)
	s.couponService.AssertExpectations(s.T())
}

func (s *AdminCouponTestSuite) TestCreateCoupon_Returns201() {
	req := coupondto.CreateCouponRequest{
		Code:            "SAVE10",
		DiscountPercent: 10,
		MinIRRPrice:     500000,
	}
	s.couponService.On("CreateCoupon", req).Return(nil).Once()

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/coupons", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusCreated, w.Code)
	s.couponService.AssertExpectations(s.T())
}

func (s *AdminCouponTestSuite) TestCreateCoupon_DuplicateCode_NonOK() {
	req := coupondto.CreateCouponRequest{
		Code:            "SAVE10",
		DiscountPercent: 10,
		MinIRRPrice:     500000,
	}
	conflict := exception.ConflictErrors{}
	conflict.Add("code", "alreadyExist")
	s.couponService.On("CreateCoupon", req).Return(conflict).Once()

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/coupons", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(true).ServeHTTP(w, r)

	s.NotEqual(http.StatusOK, w.Code)
	s.couponService.AssertExpectations(s.T())
}

func (s *AdminCouponTestSuite) TestDeleteCoupon_Returns200() {
	s.couponService.On("DeleteCoupon", testCouponID).Return(nil).Once()

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/admin/coupons/7", nil))

	s.Equal(http.StatusOK, w.Code)
	s.couponService.AssertExpectations(s.T())
}

func (s *AdminCouponTestSuite) TestDeleteCoupon_NotFound_NonOK() {
	s.couponService.On("DeleteCoupon", testCouponID).Return(exception.NotFoundError{Item: "coupon"}).Once()

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/admin/coupons/7", nil))

	s.Equal(http.StatusNotFound, w.Code)
	s.couponService.AssertExpectations(s.T())
}

func TestAdminCouponSuite(t *testing.T) {
	suite.Run(t, new(AdminCouponTestSuite))
}
