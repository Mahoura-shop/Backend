package service

import (
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/Mahoura-shop/Backend/bootstrap"
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/domain/s3"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type ProductService struct {
	constants         *bootstrap.Constants
	productRepository postgres.ProductRepository
	categoryService   usecase.CategoryService
	brandService      usecase.BrandService
	s3Storage         s3.S3Storage
	db                database.Database
}

type ProductServiceDeps struct {
	Constants         *bootstrap.Constants
	ProductRepository postgres.ProductRepository
	CategoryService   usecase.CategoryService
	BrandService      usecase.BrandService
	S3Storage         s3.S3Storage
	DB                database.Database
}

func NewProductService(deps ProductServiceDeps) *ProductService {
	return &ProductService{
		constants:         deps.Constants,
		productRepository: deps.ProductRepository,
		categoryService:   deps.CategoryService,
		brandService:      deps.BrandService,
		s3Storage:         deps.S3Storage,
		db:                deps.DB,
	}
}

func (productService *ProductService) ParseSlug(slug string) (string, error) {
	parsedSlug := strings.ReplaceAll(slug, " ", "_")
    parsedSlug = strings.ReplaceAll(parsedSlug, "\t", "_")
    parsedSlug = strings.ReplaceAll(parsedSlug, "\n", "_")
    
    parsedSlug = strings.ToLower(parsedSlug)

	var result strings.Builder
    for _, r := range slug {
        if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
            result.WriteRune(r)
        }
    }
    
    parsedSlug = result.String()
    
    parsedSlug = strings.Trim(parsedSlug, "_")
    
    for strings.Contains(parsedSlug, "__") {
		parsedSlug = strings.ReplaceAll(parsedSlug, "__", "_")
    }

	parsedSlug = strings.Trim(parsedSlug, "_")

	var validationErrors exception.ValidationErrors
	if parsedSlug == "" {
		validationErrors.Add(productService.constants.Field.Product, productService.constants.Tag.EmptySlug)
		return "", validationErrors
	}
	
	return parsedSlug, nil
}

func (productService *ProductService) ParseProduct(product entity.Product) (productdto.ProductCredential) {
	response := productdto.ProductCredential{
		ID:           product.ID,
		Name:         product.Name,
		Slug:         product.Slug,
		Price:        product.Price,
		Description:  product.Description,
		IsActive:     product.IsActive,
		IsNew:        product.IsNew,
		Priority:     product.Priority,
		MinOrder:     product.MinOrder,
		Quantity:     product.Quantity,
		QuantityType: product.QuantityType,
		CurrencyCode: product.CurrencyCode,
		ProductPic:   product.ProductPic,
	}
	if product.Category != nil {
		category, _ := productService.categoryService.ParseCategory(*product.Category)
		response.Category = &category
	}
	if product.Brand != nil {
		brand, _ := productService.brandService.ParseBrand(*product.Brand)
		response.Brand = &brand
	}
	
	return response
}

func (productService *ProductService) FindProductBySlug(slug string) (*productdto.ProductCredential, error) {
	product, err := productService.productRepository.FindProductBySlug(productService.db, slug)
	if err != nil {
		return nil, err
	}
	if product == nil {
		notFoundError := exception.NotFoundError{Item: productService.constants.Field.Product}
		return nil, notFoundError
	}
	
	if product.ProductPic != "" {
		_, err := productService.s3Storage.GetPresignedURL(enum.ProductPic, product.ProductPic, 8*time.Hour)
		if err != nil {
			return nil, err
		}
	}

	parsedProduct := productService.ParseProduct(*product)
	return &parsedProduct, nil
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

func (productService *ProductService) GetProduct(productID uint) (*productdto.ProductCredential, error) {
	product, err := productService.productRepository.FindProductByID(productService.db, productID)
	if err != nil {
		return nil, err
	}
	
	if product == nil {
		return nil, exception.NotFoundError{Item: productService.constants.Field.Product}
	}

	if product.ProductPic != "" {
		productPic, err := productService.s3Storage.GetPresignedURL(enum.ProductPic, product.ProductPic, 8*time.Hour)
		if err != nil {
			return nil, err
		}
		product.ProductPic = productPic
	}
	
	parsedProduct := productService.ParseProduct(*product)
	return &parsedProduct, nil
}

func (productService *ProductService) GetProducts() ([]productdto.ProductCredential, error) {
	products, err := productService.productRepository.GetProducts(productService.db)
	if err != nil {
		return nil, err
	}

	var responses []productdto.ProductCredential
	for _, product := range products {
		response := productService.ParseProduct(*product)
		if product.Brand != nil {
			brand, _ := productService.brandService.ParseBrand(*product.Brand)
			response.Brand = &brand
		}
		if product.Category != nil {
			category, _ := productService.categoryService.ParseCategory(*product.Category)
			response.Category = &category
		}
		if product.ProductPic != "" {
			productPic, err := productService.s3Storage.GetPresignedURL(enum.ProductPic, product.ProductPic, 8*time.Hour)
			if err != nil {
				return nil, err
			}
			response.ProductPic = productPic
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (productService *ProductService) applyProductInitial(product entity.Product, productInfo productdto.CreateProductRequest) {
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

	if productInfo.BrandID != nil {
		product.BrandID = productInfo.BrandID
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
}

func (productService *ProductService) CreateProduct(productInfo productdto.CreateProductRequest) error {
	parsedSlug, err := productService.ParseSlug(productInfo.Slug)
	if err != nil {
		return err
	}

	err = productService.validateDuplicateProduct(parsedSlug)
	if err != nil {
		return err
	}

	product := entity.Product{
		Name:       productInfo.Name,
		Slug:       parsedSlug,
		Price:      productInfo.Price,
		CategoryID: productInfo.CategoryID,
		BrandID:    productInfo.BrandID,
	}

	productService.applyProductInitial(product, productInfo)

	if productInfo.CategoryID != nil {
		category, err := productService.categoryService.FindCategoryByID(*productInfo.CategoryID)
		if err != nil {
			return err
		}
		if category == nil {
			notFoundError := exception.NotFoundError{Item: productService.constants.Field.Category}
			return notFoundError
		}
	}

	if productInfo.BrandID != nil {
		brand, err := productService.brandService.FindBrandByID(*productInfo.BrandID)
		if err != nil {
			return err
		}
		if brand == nil {
			notFoundError := exception.NotFoundError{Item: productService.constants.Field.Brand}
			return notFoundError
		}
	}

	err = productService.db.WithTransaction(func(tx database.Database) error {
		createdProduct, err := productService.productRepository.CreateProduct(tx, &product)
		if err != nil {
			return err
		}
		if productInfo.ProductPic != nil {
			product.ProductPic = productService.constants.S3BucketPath.GetProductPicPath(createdProduct.ID, productInfo.ProductPic.Filename)
			if err := productService.s3Storage.UploadObject(enum.ProductPic, product.ProductPic, productInfo.ProductPic); err != nil {
				return  err
			}
		
			if err := productService.productRepository.UpdateProduct(tx, &product); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (productService *ProductService) applyProductUpdates(product *entity.Product, productInfo productdto.UpdateProductRequest) error {
	if productInfo.Name != nil {
		product.Name = *productInfo.Name
	}
	if productInfo.Slug != nil {
		product.Slug = *productInfo.Slug
		parsedSlug, err := productService.ParseSlug(product.Slug)
		if err != nil {
			return err
		}
		product.Slug = parsedSlug
	}
	if productInfo.Price != nil {
		product.Price = *productInfo.Price
	}
	if productInfo.Description != nil {
		product.Description = *productInfo.Description
	}
	if productInfo.IsActive != nil {
		product.IsActive = *productInfo.IsActive
	}
	if productInfo.IsNew != nil {
		product.IsNew = *productInfo.IsNew
	}
	if productInfo.Priority != nil {
		product.Priority = *productInfo.Priority
	}
	if productInfo.MinOrder != nil {
		product.MinOrder = *productInfo.MinOrder
	}
	if productInfo.CategoryID != nil {
		product.CategoryID = productInfo.CategoryID
	}
	if productInfo.BrandID != nil {
		product.BrandID = productInfo.BrandID
	}
	if productInfo.Quantity != nil {
		product.Quantity = *productInfo.Quantity
	}
	if productInfo.QuantityType != nil {
		product.QuantityType = *productInfo.QuantityType
	}
	if productInfo.CurrencyCode != nil {
		product.CurrencyCode = *productInfo.CurrencyCode
	}
	return nil
}

func (productService *ProductService) newSlugAvailable(slug string, productID uint) error {
	var conflictErrors exception.ConflictErrors
	product, err := productService.productRepository.FindProductBySlug(productService.db, slug)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if product != nil {
		if product.ID != productID {
			conflictErrors.Add(productService.constants.Field.Product, productService.constants.Tag.AlreadyExist)
			return conflictErrors
		}
	}

	return nil
}

func (productService *ProductService) UpdateProduct(productInfo productdto.UpdateProductRequest) error {
	product, err := productService.productRepository.FindProductByID(productService.db, productInfo.ID)
	if err != nil {
		return err
	}
	if product == nil {
		return exception.NotFoundError{Item: productService.constants.Field.Product}
	}

	err = productService.applyProductUpdates(product, productInfo)
	if err != nil {
		return err
	}

	if productInfo.CategoryID != nil {
		category, err := productService.categoryService.FindCategoryByID(*productInfo.CategoryID)
		if err != nil {
			return err
		}
		if category == nil {
			notFoundError := exception.NotFoundError{Item: productService.constants.Field.Category}
			return notFoundError
		} 
	}

	if productInfo.BrandID != nil {
		brand, err := productService.brandService.FindBrandByID(*productInfo.BrandID)
		if err != nil {
			return err
		}
		if brand == nil {
			notFoundError := exception.NotFoundError{Item: productService.constants.Field.Brand}
			return notFoundError
		} 
	}

	if (productInfo.Slug != nil) {
		err := productService.newSlugAvailable(product.Slug, productInfo.ID)
		if err != nil {
			return err
		}
	}
	err = productService.db.WithTransaction(func(tx database.Database) error {
		if productInfo.ProductPic != nil {
			productPicPath := productService.constants.S3BucketPath.GetProductPicPath(productInfo.ID, productInfo.ProductPic.Filename)
			productService.s3Storage.UploadObject(enum.ProductPic, productPicPath, productInfo.ProductPic)
			product.ProductPic = productPicPath
		}
		if err := productService.productRepository.UpdateProduct(tx, product); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (productService *ProductService) DeleteProduct(productID uint) error {
	product, err := productService.productRepository.FindProductByID(productService.db, productID)
	if err != nil {
		return err
	}
	if product == nil {
		return exception.NotFoundError{Item: productService.constants.Field.Product}
	}

	if err := productService.productRepository.DeleteProductByID(productService.db, productID); err != nil {
		return err
	}
	return nil
}