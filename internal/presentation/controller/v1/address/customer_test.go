package address

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	addressdto "github.com/Mahoura-shop/Backend/internal/application/dto/address"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func init() { gin.SetMode(gin.TestMode) }

const addrUserID = uint(42)

type CustomerAddressTestSuite struct {
	suite.Suite
	constants      *bootstrap.Constants
	addressService *mocks.AddressServiceMock
	controller     *CustomerAddressController
}

func (s *CustomerAddressTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.addressService = mocks.NewAddressServiceMock()
	s.controller = NewCustomerAddressController(s.constants, s.addressService)
}

func (s *CustomerAddressTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Set(s.constants.Context.ID, addrUserID)
		c.Next()
	})
	if withRecovery {
		r.Use(addrPanicToHTTP())
	}
	r.GET("/address", s.controller.GetCustomerAddresses)
	r.POST("/address", s.controller.CreateUserAddress)
	return r
}

func addrPanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case exception.ValidationErrors:
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation"})
				case exception.NotFoundError:
					c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

// C-26: POST /address → 200, address saved
func (s *CustomerAddressTestSuite) TestC26_CreateAddress_Returns200() {
	req := addressdto.CreateAddressRequest{
		OwnerID:       addrUserID,
		ProvinceID:    1,
		CityID:        2,
		StreetAddress: "Valiasr St",
		PostalCode:    "1234567890",
		HouseNumber:   "12",
		Unit:          3,
	}
	s.addressService.On("CreateAddress", req).Return(addressdto.AddressResponse{
		ID:            1,
		ProvinceID:    1,
		CityID:        2,
		StreetAddress: "Valiasr St",
		PostalCode:    "1234567890",
		HouseNumber:   "12",
		Unit:          3,
	}, nil)

	body, _ := json.Marshal(map[string]any{
		"provinceID":    1,
		"cityID":        2,
		"streetAddress": "Valiasr St",
		"postalCode":    "1234567890",
		"houseNumber":   "12",
		"unit":          3,
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/address", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.addressService.AssertExpectations(s.T())
}

// C-26b: POST /address missing required field → 422
func (s *CustomerAddressTestSuite) TestC26b_CreateAddress_MissingField_Returns422() {
	// Missing postalCode
	body, _ := json.Marshal(map[string]any{
		"provinceID":    1,
		"cityID":        2,
		"streetAddress": "Valiasr St",
		"houseNumber":   "12",
		"unit":          3,
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/address", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(true).ServeHTTP(w, r)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
	s.addressService.AssertNotCalled(s.T(), "CreateAddress")
}

// C-27: GET /address → 200 with address list
func (s *CustomerAddressTestSuite) TestC27_GetAddresses_Returns200() {
	s.addressService.On("GetAddresses", addressdto.GetOwnerAddressesRequest{
		OwnerID: addrUserID,
	}).Return([]addressdto.AddressResponse{
		{ID: 1, StreetAddress: "Valiasr St"},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/address", nil))

	s.Equal(http.StatusOK, w.Code)
	var resp map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	addresses := resp["data"].([]any)
	s.Len(addresses, 1)
	s.addressService.AssertExpectations(s.T())
}

func TestCustomerAddressSuite(t *testing.T) {
	suite.Run(t, new(CustomerAddressTestSuite))
}
