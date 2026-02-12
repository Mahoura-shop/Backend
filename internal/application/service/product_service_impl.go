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
		validationErrors.Add(productService.constants.Field.Slug, productService.constants.Tag.EmptySlug)
		return "", validationErrors
	}
	
	return parsedSlug, nil
}

func (productService *ProductService) ParseProduct(product entity.Product) (productdto.ProductCredential) {
	response := productdto.ProductCredential{
		ID:            product.ID,
		Name:          product.Name,
		Slug:          product.Slug,
		Price:         product.Price,
		CurrencyCode:  product.CurrencyCode,
		IRRPrice:      product.IRRPrice,
		ConsumerPrice: product.ConsumerPrice,
		Step1Percent:  product.Step1Percent,
		Step2Percent:  product.Step2Percent,
		Step3Percent:  product.Step3Percent,
		Step1Price:    product.Step1Price,
		Step2Price:    product.Step2Price,
		Step3Price:    product.Step3Price,
		Quantity:      product.Quantity,
		QuantityType:  product.QuantityType,
		Priority:      product.Priority,
		MinOrder:      product.MinOrder,
		Description:   product.Description,
		IsActive:      product.IsActive,
		IsNew:         product.IsNew,
		ProductPic:    product.ProductPic,
	}
	if product.Category != nil {
		category, _ := productService.categoryService.ParseCategory(*product.Category)
		response.Category = &category
		response.CategoryID = *product.CategoryID
	}
	if product.Brand != nil {
		brand, _ := productService.brandService.ParseBrand(*product.Brand)
		response.Brand = &brand
		response.BrandID = *product.BrandID
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

func (productService *ProductService) validateDuplicateProduct(slug string, name string) error {
	var conflictErrors exception.ConflictErrors

	product, err := productService.productRepository.FindProductBySlug(productService.db, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if product != nil {
		conflictErrors.Add(productService.constants.Field.Slug, productService.constants.Tag.AlreadyExist)
		return conflictErrors
	}

	product, err = productService.productRepository.FindProductByName(productService.db, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if product != nil {
		conflictErrors.Add(productService.constants.Field.Name, productService.constants.Tag.AlreadyExist)
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

func (productService *ProductService) applyProductInitial(product *entity.Product, productInfo productdto.CreateProductRequest) {
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

	if productInfo.IRRPrice != nil {
		product.IRRPrice = *productInfo.IRRPrice
	} else {
		product.IRRPrice = 0
	}
	
	if productInfo.ConsumerPrice != nil {
		product.ConsumerPrice = *productInfo.ConsumerPrice
	} else {
		product.ConsumerPrice = 0
	}

	if productInfo.Step1Percent != nil {
		product.Step1Percent = *productInfo.Step1Percent
	} else {
		product.Step1Percent = 0
	}

	if productInfo.Step2Percent != nil {
		product.Step2Percent = *productInfo.Step2Percent
	} else {
		product.Step2Percent = 0
	}

	if productInfo.Step3Percent != nil {
		product.Step3Percent = *productInfo.Step3Percent
	} else {
		product.Step3Percent = 0
	}

	if productInfo.Step1Price != nil {
		product.Step1Price = *productInfo.Step1Price
	} else {
		product.Step1Price = 0
	}

	if productInfo.Step2Price != nil {
		product.Step2Price = *productInfo.Step2Price
	} else {
		product.Step2Price = 0
	}

	if productInfo.Step3Price != nil {
		product.Step3Price = *productInfo.Step3Price
	} else {
		product.Step3Price = 0
	}
}

func (productService *ProductService) CreateProduct(productInfo productdto.CreateProductRequest) error {
	parsedSlug, err := productService.ParseSlug(productInfo.Slug)
	if err != nil {
		return err
	}

	err = productService.validateDuplicateProduct(parsedSlug, productInfo.Name)
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

	productService.applyProductInitial(&product, productInfo)

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
				networkErr := exception.ClassifyNetworkError(err, "ArvanStorage", "UploadObject")
				_ = productService.productRepository.DeleteProductByID(tx, createdProduct.ID);
				return networkErr
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
	if productInfo.IRRPrice != nil {
		product.IRRPrice = *productInfo.IRRPrice
	}
	if productInfo.ConsumerPrice != nil {
		product.ConsumerPrice = *productInfo.ConsumerPrice
	}
	if productInfo.Step1Percent != nil {
		product.Step1Percent = *productInfo.Step1Percent
	}
	if productInfo.Step2Percent != nil {
		product.Step2Percent = *productInfo.Step2Percent
	}
	if productInfo.Step3Percent != nil {
		product.Step3Percent = *productInfo.Step3Percent
	}
	if productInfo.Step1Price != nil {
		product.Step1Price = *productInfo.Step1Price
	}
	if productInfo.Step2Price != nil {
		product.Step2Price = *productInfo.Step2Price
	}
	if productInfo.Step3Price != nil {
		product.Step3Price = *productInfo.Step3Price
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
			conflictErrors.Add(productService.constants.Field.Slug, productService.constants.Tag.AlreadyExist)
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
			if err := productService.s3Storage.UploadObject(enum.ProductPic, productPicPath, productInfo.ProductPic); err != nil {
				networkErr := exception.ClassifyNetworkError(err, "ArvanStorage", "UploadObject")
				_ = productService.productRepository.DeleteProductByID(tx, productInfo.ID);
				return networkErr
			}
			product.ProductPic = productPicPath
		} else {
			if product.ProductPic != "" {
				if err := productService.s3Storage.DeleteObject(enum.ProductPic, product.ProductPic); err != nil {
					networkErr := exception.ClassifyNetworkError(err, "ArvanStorage", "DeleteObject")
					return networkErr
				}
				product.ProductPic = ""
			}
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