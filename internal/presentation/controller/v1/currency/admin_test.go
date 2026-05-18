package currency

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	currencydto "github.com/Mahoura-shop/Backend/internal/application/dto/currency"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func init() { gin.SetMode(gin.TestMode) }

const testCurrencyID = uint(1)

type AdminCurrencyTestSuite struct {
	suite.Suite
	constants       *bootstrap.Constants
	currencyService *mocks.CurrencyServiceMock
	controller      *AdminCurrencyController
}

func (s *AdminCurrencyTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.currencyService = mocks.NewCurrencyServiceMock()
	s.controller = NewAdminCurrencyController(s.constants, &bootstrap.Pagination{}, s.currencyService)
}

func (s *AdminCurrencyTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Next()
	})
	if withRecovery {
		r.Use(currencyPanicToHTTP())
	}
	r.GET("/admin/currencies", s.controller.GetCurrencies)
	r.PATCH("/admin/currencies/:currencyID", s.controller.UpdateCurrency)
	return r
}

func currencyPanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case exception.NotFoundError:
					c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
				case exception.ValidationErrors:
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation"})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

// ADM-60: GET /admin/currencies → 200 with all currencies
func (s *AdminCurrencyTestSuite) TestADM60_GetCurrencies_Returns200() {
	s.currencyService.On("GetCurrencies").Return([]currencydto.CurrencyCredential{
		{ID: testCurrencyID, Name: "US Dollar", Code: "USD", ConvertRate: 60000},
		{ID: 2, Name: "Euro", Code: "EUR", ConvertRate: 65000},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/currencies", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	currencies := body["data"].([]any)
	s.Len(currencies, 2)
	s.currencyService.AssertExpectations(s.T())
}

// ADM-61: PATCH /admin/currencies/:currencyID → 200, convert rate updated
func (s *AdminCurrencyTestSuite) TestADM61_UpdateCurrency_Returns200() {
	newRate := uint(62000)
	req := currencydto.UpdateCurrencyRequest{
		ID:          testCurrencyID,
		ConvertRate: &newRate,
	}
	s.currencyService.On("UpdateCurrency", req).Return(nil)

	payload, _ := json.Marshal(map[string]any{"convertRate": newRate})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/admin/currencies/1", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.currencyService.AssertExpectations(s.T())
}

// ADM-61b: Update non-existent currency → service returns NotFoundError → 404
func (s *AdminCurrencyTestSuite) TestADM61b_UpdateCurrency_NotFound_Returns404() {
	newRate := uint(62000)
	req := currencydto.UpdateCurrencyRequest{
		ID:          testCurrencyID,
		ConvertRate: &newRate,
	}
	s.currencyService.On("UpdateCurrency", req).Return(exception.NotFoundError{Item: "currency"})

	payload, _ := json.Marshal(map[string]any{"convertRate": newRate})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/admin/currencies/1", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(true).ServeHTTP(w, r)

	s.Equal(http.StatusNotFound, w.Code)
	s.currencyService.AssertExpectations(s.T())
}

func TestAdminCurrencySuite(t *testing.T) {
	suite.Run(t, new(AdminCurrencyTestSuite))
}
