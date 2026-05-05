package service

import (
	wishlistdto "github.com/Mahoura-shop/Backend/internal/application/dto/wishlist"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	domainPostgres "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type WishlistService struct {
	wishlistRepository domainPostgres.WishlistRepository
	productService     usecase.ProductService
	db                 database.Database
}

type WishlistServiceDeps struct {
	WishlistRepository domainPostgres.WishlistRepository
	ProductService     usecase.ProductService
	DB                 database.Database
}

func NewWishlistService(deps WishlistServiceDeps) *WishlistService {
	return &WishlistService{
		wishlistRepository: deps.WishlistRepository,
		productService:     deps.ProductService,
		db:                 deps.DB,
	}
}

func (s *WishlistService) AddToWishlist(userID, productID uint) error {
	existing, err := s.wishlistRepository.FindWishlistItem(s.db, userID, productID)
	if err != nil {
		return err
	}
	if existing != nil {
		var ce exception.ConflictErrors
		ce.Add("wishlist", "alreadyExist")
		return ce
	}

	return s.wishlistRepository.AddToWishlist(s.db, entity.Wishlist{
		UserID:    userID,
		ProductID: productID,
	})
}

func (s *WishlistService) RemoveFromWishlist(userID, productID uint) error {
	return s.wishlistRepository.RemoveFromWishlist(s.db, userID, productID)
}

func (s *WishlistService) GetMyWishlist(userID uint) ([]wishlistdto.WishlistItemCredential, error) {
	items, err := s.wishlistRepository.GetWishlistByUserID(s.db, userID)
	if err != nil {
		return nil, err
	}
	var result []wishlistdto.WishlistItemCredential
	for _, item := range items {
		parsed := s.productService.ParseProduct(item.Product)
		result = append(result, wishlistdto.WishlistItemCredential{
			ID:        item.ID,
			ProductID: item.ProductID,
			Product:   parsed,
		})
	}
	return result, nil
}
