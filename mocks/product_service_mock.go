package mocks

import (
	"mime/multipart"

	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type ProductServiceMock struct {
	mock.Mock
}

func NewProductServiceMock() *ProductServiceMock {
	return &ProductServiceMock{}
}

func (m *ProductServiceMock) ParseSlug(s string) (string, error) {
	args := m.Called(s)
	return args.String(0), args.Error(1)
}

func (m *ProductServiceMock) ParseProduct(p entity.Product) productdto.ProductCredential {
	args := m.Called(p)
	return args.Get(0).(productdto.ProductCredential)
}

func (m *ProductServiceMock) FindProductBySlug(slug string) (*productdto.ProductCredential, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productdto.ProductCredential), args.Error(1)
}

func (m *ProductServiceMock) GetProduct(id uint) (*productdto.ProductCredential, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productdto.ProductCredential), args.Error(1)
}

func (m *ProductServiceMock) GetProductBySlug(slug string) (*productdto.ProductCredential, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productdto.ProductCredential), args.Error(1)
}

func (m *ProductServiceMock) GetProducts() ([]productdto.ProductCredential, error) {
	args := m.Called()
	return args.Get(0).([]productdto.ProductCredential), args.Error(1)
}

func (m *ProductServiceMock) SearchProducts(filter productdto.ProductFilterRequest) ([]productdto.ProductCredential, error) {
	args := m.Called(filter)
	return args.Get(0).([]productdto.ProductCredential), args.Error(1)
}

func (m *ProductServiceMock) SearchProductsWithPagination(filter productdto.ProductFilterRequest) (*productdto.ProductSearchResponse, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productdto.ProductSearchResponse), args.Error(1)
}

func (m *ProductServiceMock) GetRelatedProducts(id uint, limit int) ([]productdto.ProductCredential, error) {
	args := m.Called(id, limit)
	return args.Get(0).([]productdto.ProductCredential), args.Error(1)
}

func (m *ProductServiceMock) CreateProduct(req productdto.CreateProductRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *ProductServiceMock) UpdateProduct(req productdto.UpdateProductRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *ProductServiceMock) DeleteProduct(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ProductServiceMock) GetCategoryProducts(categoryID uint) ([]productdto.ProductCredential, error) {
	args := m.Called(categoryID)
	return args.Get(0).([]productdto.ProductCredential), args.Error(1)
}

func (m *ProductServiceMock) GetBrandProducts(brandID uint) ([]productdto.ProductCredential, error) {
	args := m.Called(brandID)
	return args.Get(0).([]productdto.ProductCredential), args.Error(1)
}

func (m *ProductServiceMock) UpdateProductsPrice(prices []productdto.ProductPriceUpdateCredentials) error {
	args := m.Called(prices)
	return args.Error(0)
}

func (m *ProductServiceMock) UpdateProductsStock(items []productdto.ProductStockUpdateCredentials, updateType string) error {
	args := m.Called(items, updateType)
	return args.Error(0)
}

func (m *ProductServiceMock) GetProductPrices() ([]productdto.ProductPrices, error) {
	args := m.Called()
	return args.Get(0).([]productdto.ProductPrices), args.Error(1)
}

func (m *ProductServiceMock) AddProductImage(id uint, file *multipart.FileHeader) error {
	args := m.Called(id, file)
	return args.Error(0)
}

func (m *ProductServiceMock) DeleteProductImage(productID uint, imageID uint) error {
	args := m.Called(productID, imageID)
	return args.Error(0)
}

func (m *ProductServiceMock) UpdateInventoryFromExcel(rows []productdto.ExcelInventoryRow) (int, error) {
	args := m.Called(rows)
	return args.Int(0), args.Error(1)
}
