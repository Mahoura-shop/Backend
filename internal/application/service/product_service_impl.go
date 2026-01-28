package service

import (
	"errors"

	"github.com/Mahoura-shop/Backend/bootstrap"
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type ProductService struct {
	constants          *bootstrap.Constants
	productRepository  postgres.ProductRepository
	categoryRepository postgres.CategoryRepository
	db                 database.Database
}

type ProductServiceDeps struct {
	Constants          *bootstrap.Constants
	ProductRepository  postgres.ProductRepository
	CategoryRepository postgres.CategoryRepository
	DB                 database.Database
}

func NewProductService(deps ProductServiceDeps) *ProductService {
	return &ProductService{
		constants:          deps.Constants,
		productRepository:  deps.ProductRepository,
	    categoryRepository: deps.CategoryRepository,
		db:                 deps.DB,
	}
}

func (productService *ProductService) FindProductBySlug(slug string) (*entity.Product, error) {
	product, err := productService.productRepository.FindProductBySlug(productService.db, slug)
	if err != nil {
		return nil, err
	}
	if product == nil {
		notFoundError := exception.NotFoundError{Item: productService.constants.Field.Product}
		return nil, notFoundError
	}
	return product, nil
}

func (productService *ProductService) validateDuplicateProduct(slug string) error {
	var conflictErrors exception.ConflictErrors
	product, err := productService.productRepository.FindProductBySlug(productService.db, slug)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if product != nil {
		conflictErrors.Add(productService.constants.Field.Product, productService.constants.Tag.AlreadyExist)
		return conflictErrors
	}

	return nil
}

func (productService *ProductService) CreateProduct(productInfo productdto.CreateProductRequest) error {
	err := productService.validateDuplicateProduct(productInfo.Slug)
	if err != nil {
		return err
	}

	product := &entity.Product{
		Name:       productInfo.Name,
		Slug:       productInfo.Slug,
		Price:      productInfo.Price,
		CategoryID: productInfo.CategoryID,
	}
	if productInfo.Description != nil {
		product.Description = *productInfo.Description
	} else {
		product.Description = ""
	}

	if productInfo.IsActive != nil {
		product.IsActive = *productInfo.IsActive
	} else {
		product.IsActive = true
	}	
	
	if productInfo.IsNew != nil {
		product.IsNew = *productInfo.IsNew
	} else {
		product.IsNew = true
	}

	if productInfo.Priority != nil {
		product.Priority = *productInfo.Priority
	} else {
		product.Priority = 0
	}

	if productInfo.MinOrder != nil {
		product.MinOrder = *productInfo.MinOrder
	} else {
		product.MinOrder = 1
	}

	if productInfo.CategoryID != nil {
		product.CategoryID = productInfo.CategoryID
	}

	if productInfo.Quantity != nil {
		product.Quantity = *productInfo.Quantity
	} else {
		product.Quantity = 0
	}

	if productInfo.QuantityType != nil {
		product.QuantityType = *productInfo.QuantityType
	} else {
		product.QuantityType = "عدد"
	}

	if productInfo.CurrencyCode != nil {
		product.CurrencyCode = *productInfo.CurrencyCode
	} else {
		product.CurrencyCode = "IRR"
	}

	if productInfo.CategoryID != nil {
		category, err := productService.categoryRepository.FindCategoryByID(productService.db, *productInfo.CategoryID)
		if err != nil {
			return err
		}
		if category == nil {
			notFoundError := exception.NotFoundError{Item: productService.constants.Field.Category}
			return notFoundError
		} 
		product.Category = category
	}

	err = productService.db.WithTransaction(func(tx database.Database) error {
		err = productService.productRepository.CreateProduct(tx, product)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

// func (productService *ProductService) GetProducts() ([]productdto.ProductCredentialResponse, error) {
// 	products, err := productService.productRepository.GetProducts(productService.db)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var responses []productdto.ProductCredentialResponse
// 	for _, product := range products {
// 		response := productdto.ProductCredentialResponse{
// 			ID:       product.ID,
// 			Name:     product.Name,
// 			Slug:     product.Slug,
// 			IsActive: product.IsActive,
// 		}
		
// 		if product.Description != "" {
// 			description := product.Description
// 			response.Description = &description
// 		}
		
// 		responses = append(responses, response)
// 	}
// 	return responses, nil
// }

// func (productService *ProductService) DeleteProduct(productID uint) error {
// 	product, err := productService.productRepository.FindProductByID(productService.db, productID)
// 	if err != nil {
// 		return err
// 	}
// 	if product == nil {
// 		return exception.NotFoundError{Item: productService.constants.Field.Product}
// 	}

// 	if err := productService.productRepository.DeleteProductByID(productService.db, productID); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (productService *ProductService) applyProductUpdates(product *entity.Product, name *string, slug *string, description *string, isActive *bool) {
// 	if name != nil {
// 		product.Name = *name
// 	}

// 	if slug != nil {
// 		product.Slug = *slug
// 	}

// 	if description != nil {
// 		product.Description = *description
// 	}
	
// 	if isActive != nil {
// 		product.IsActive = *isActive
// 	}
// }

// func (productService *ProductService) UpdateProduct(productInfo productdto.UpdateProductRequest) error {
// 	product, err := productService.productRepository.FindProductByID(productService.db, productInfo.ID)
// 	if err != nil {
// 		return err
// 	}
// 	if product == nil {
// 		return exception.NotFoundError{Item: productService.constants.Field.Product}
// 	}

// 	productService.applyProductUpdates(product, productInfo.Name, productInfo.Slug, productInfo.Description, productInfo.IsActive)

// 	if (productInfo.Slug != nil) {
// 		err := productService.validateDuplicateProduct(product.Slug)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	err = productService.db.WithTransaction(func(tx database.Database) error {
// 		if err := productService.productRepository.UpdateProduct(tx, product); err != nil {
// 			return err
// 		}
// 		return nil
// 	})

// 	return err
// }