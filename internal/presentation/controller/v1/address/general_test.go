package address

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	addressdto "github.com/Mahoura-shop/Backend/internal/application/dto/address"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

const testProvinceID = uint(8)

type GeneralAddressTestSuite struct {
	suite.Suite
	constants      *bootstrap.Constants
	addressService *mocks.AddressServiceMock
	controller     *GeneralAddressController
}

func (s *GeneralAddressTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.addressService = mocks.NewAddressServiceMock()
	s.controller = NewGeneralAddressController(s.constants, s.addressService)
}

func (s *GeneralAddressTestSuite) newRouter() *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Next()
	})
	r.GET("/provinces", s.controller.GetProvince)
	r.GET("/provinces/:provinceID/cities", s.controller.GetProvinceCities)
	return r
}

// LOC-01: GET /provinces → 200 with all provinces
func (s *GeneralAddressTestSuite) TestLOC01_GetProvinces_Returns200() {
	s.addressService.On("GetProvinceList").Return([]addressdto.ProvinceResponse{
		{ID: 1, Name: "Tehran"},
		{ID: testProvinceID, Name: "Isfahan"},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/provinces", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	provinces := body["data"].([]any)
	s.Len(provinces, 2)
	s.addressService.AssertExpectations(s.T())
}

// LOC-01b: Province list is non-empty
func (s *GeneralAddressTestSuite) TestLOC01b_GetProvinces_ReturnsNonEmpty() {
	s.addressService.On("GetProvinceList").Return([]addressdto.ProvinceResponse{
		{ID: 1, Name: "Tehran"},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/provinces", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	provinces := body["data"].([]any)
	s.Greater(len(provinces), 0)
	s.addressService.AssertExpectations(s.T())
}

// LOC-02: GET /provinces/:provinceID/cities → 200 with cities for that province
func (s *GeneralAddressTestSuite) TestLOC02_GetProvinceCities_Returns200() {
	s.addressService.On("GetCityProvinceCities", addressdto.GetProvinceCitiesRequest{
		ProvinceID: testProvinceID,
	}).Return([]addressdto.CityResponse{
		{ID: 101, Name: "Isfahan City"},
		{ID: 102, Name: "Kashan"},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter().ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/provinces/8/cities", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	cities := body["data"].([]any)
	s.Len(cities, 2)
	s.addressService.AssertExpectations(s.T())
}

// LOC-02b: Cities belong to the requested province (filter correctness)
func (s *GeneralAddressTestSuite) TestLOC02b_GetProvinceCities_WrongProvince_ReturnsDifferentCities() {
	otherProvinceID := uint(1)
	s.addressService.On("GetCityProvinceCities", addressdto.GetProvinceCitiesRequest{
		ProvinceID: otherProvinceID,
	}).Return([]addressdto.CityResponse{
		{ID: 1, Name: "Tehran City"},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter().ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/provinces/1/cities", nil))

	s.Equal(http.StatusOK, w.Code)
	s.addressService.AssertExpectations(s.T())
}

func TestGeneralAddressSuite(t *testing.T) {
	suite.Run(t, new(GeneralAddressTestSuite))
}
