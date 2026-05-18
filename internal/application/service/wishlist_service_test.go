package service

import (
	"errors"
	"testing"

	wishlistdto "github.com/Mahoura-shop/Backend/internal/application/dto/wishlist"
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type WishlistServiceTestSuite struct {
	suite.Suite
	wishlistRepository *mocks.WishlistRepositoryMock
	productService     *mocks.ProductServiceMock
	db                 *mocks.DatabaseMock
	service            *WishlistService
}

func (s *WishlistServiceTestSuite) SetupTest() {
	s.wishlistRepository = mocks.NewWishlistRepositoryMock()
	s.productService = mocks.NewProductServiceMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewWishlistService(WishlistServiceDeps{
		WishlistRepository: s.wishlistRepository,
		ProductService:     s.productService,
		DB:                 s.db,
	})
}

func (s *WishlistServiceTestSuite) TestAddToWishlist_Success() {
	var nilItem *entity.Wishlist
	s.wishlistRepository.On("FindWishlistItem", s.db, uint(1), uint(7)).Return(nilItem, nil).Once()
	s.wishlistRepository.On("AddToWishlist", s.db, entity.Wishlist{UserID: 1, ProductID: 7}).Return(nil).Once()

	err := s.service.AddToWishlist(1, 7)

	s.NoError(err)
	s.wishlistRepository.AssertExpectations(s.T())
}

func (s *WishlistServiceTestSuite) TestAddToWishlist_AlreadyExists() {
	existing := &entity.Wishlist{Model: dbmodel.Model{ID: 3}, UserID: 1, ProductID: 7}
	s.wishlistRepository.On("FindWishlistItem", s.db, uint(1), uint(7)).Return(existing, nil).Once()

	err := s.service.AddToWishlist(1, 7)

	s.Error(err)
	var ce exception.ConflictErrors
	s.True(errors.As(err, &ce))
	s.wishlistRepository.AssertNotCalled(s.T(), "AddToWishlist")
	s.wishlistRepository.AssertExpectations(s.T())
}

func (s *WishlistServiceTestSuite) TestRemoveFromWishlist_CallsRepo() {
	s.wishlistRepository.On("RemoveFromWishlist", s.db, uint(1), uint(7)).Return(nil).Once()

	err := s.service.RemoveFromWishlist(1, 7)

	s.NoError(err)
	s.wishlistRepository.AssertExpectations(s.T())
}

func (s *WishlistServiceTestSuite) TestGetMyWishlist_ReturnsMapped() {
	items := []*entity.Wishlist{
		{Model: dbmodel.Model{ID: 1}, UserID: 10, ProductID: 5, Product: entity.Product{}},
	}
	s.wishlistRepository.On("GetWishlistByUserID", s.db, uint(10)).Return(items, nil).Once()
	s.productService.On("ParseProduct", entity.Product{}).Return(productdto.ProductCredential{ID: 5}).Once()

	result, err := s.service.GetMyWishlist(10)

	s.NoError(err)
	s.Len(result, 1)
	s.Equal(wishlistdto.WishlistItemCredential{
		ID:        1,
		ProductID: 5,
		Product:   productdto.ProductCredential{ID: 5},
	}, result[0])
	s.wishlistRepository.AssertExpectations(s.T())
	s.productService.AssertExpectations(s.T())
}

func (s *WishlistServiceTestSuite) TestGetMyWishlist_Empty() {
	s.wishlistRepository.On("GetWishlistByUserID", s.db, uint(10)).Return([]*entity.Wishlist{}, nil).Once()

	result, err := s.service.GetMyWishlist(10)

	s.NoError(err)
	s.Nil(result)
	s.wishlistRepository.AssertExpectations(s.T())
}

func (s *WishlistServiceTestSuite) TestGetMyWishlist_RepoError() {
	repoErr := errors.New("db error")
	s.wishlistRepository.On("GetWishlistByUserID", s.db, uint(10)).Return([]*entity.Wishlist{}, repoErr).Once()

	result, err := s.service.GetMyWishlist(10)

	s.Nil(result)
	s.ErrorIs(err, repoErr)
	s.wishlistRepository.AssertExpectations(s.T())
}

func TestWishlistServiceSuite(t *testing.T) {
	suite.Run(t, new(WishlistServiceTestSuite))
}
