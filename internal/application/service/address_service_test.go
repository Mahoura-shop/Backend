package service

import (
	"errors"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	addressdto "github.com/Mahoura-shop/Backend/internal/application/dto/address"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AddressServiceTestSuite struct {
	suite.Suite
	constants         *bootstrap.Constants
	addressRepository *mocks.AddressRepositoryMock
	db                *mocks.DatabaseMock
	service           *AddressService
}

func (s *AddressServiceTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.addressRepository = mocks.NewAddressRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewAddressService(s.constants, s.addressRepository, s.db)
}

func (s *AddressServiceTestSuite) TestCreateAddress_Success() {
	province := &entity.Province{Model: dbmodel.Model{ID: 1}, Name: "Tehran"}
	city := &entity.City{Model: dbmodel.Model{ID: 2}, Name: "Tehran City"}

	s.addressRepository.On("GetProvinceByID", s.db, uint(1)).Return(province, nil).Once()
	s.addressRepository.On("GetCityByID", s.db, uint(2)).Return(city, nil).Once()
	s.addressRepository.On("CreateAddress", s.db, mock.AnythingOfType("*entity.Address")).Return(nil).Once()

	req := addressdto.CreateAddressRequest{
		ProvinceID:    1,
		CityID:        2,
		StreetAddress: "123 Main St",
		PostalCode:    "1234567890",
		HouseNumber:   "5",
		Unit:          2,
		OwnerID:       10,
		OwnerType:     "user",
	}
	result, err := s.service.CreateAddress(req)

	s.NoError(err)
	s.Equal("Tehran", result.Province)
	s.Equal("Tehran City", result.City)
	s.addressRepository.AssertExpectations(s.T())
}

func (s *AddressServiceTestSuite) TestCreateAddress_ProvinceNotFound() {
	var nilProvince *entity.Province
	s.addressRepository.On("GetProvinceByID", s.db, uint(99)).Return(nilProvince, nil).Once()

	req := addressdto.CreateAddressRequest{ProvinceID: 99, CityID: 1}
	_, err := s.service.CreateAddress(req)

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.addressRepository.AssertExpectations(s.T())
}

func (s *AddressServiceTestSuite) TestCreateAddress_CityNotFound() {
	province := &entity.Province{Model: dbmodel.Model{ID: 1}, Name: "Tehran"}
	var nilCity *entity.City
	s.addressRepository.On("GetProvinceByID", s.db, uint(1)).Return(province, nil).Once()
	s.addressRepository.On("GetCityByID", s.db, uint(99)).Return(nilCity, nil).Once()

	req := addressdto.CreateAddressRequest{ProvinceID: 1, CityID: 99}
	_, err := s.service.CreateAddress(req)

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.addressRepository.AssertExpectations(s.T())
}

func (s *AddressServiceTestSuite) TestDeleteAddress_Success() {
	address := &entity.Address{Model: dbmodel.Model{ID: 5}}
	s.addressRepository.On("GetAddressByID", s.db, uint(5)).Return(address, nil).Once()
	s.addressRepository.On("DeleteAddress", s.db, address).Return(nil).Once()

	err := s.service.DeleteAddress(5)

	s.NoError(err)
	s.addressRepository.AssertExpectations(s.T())
}

func (s *AddressServiceTestSuite) TestDeleteAddress_NotFound() {
	var nilAddress *entity.Address
	s.addressRepository.On("GetAddressByID", s.db, uint(99)).Return(nilAddress, nil).Once()

	err := s.service.DeleteAddress(99)

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.addressRepository.AssertNotCalled(s.T(), "DeleteAddress")
	s.addressRepository.AssertExpectations(s.T())
}

func (s *AddressServiceTestSuite) TestGetProvinceList_ReturnsMapped() {
	provinces := []*entity.Province{
		{Model: dbmodel.Model{ID: 1}, Name: "Tehran"},
		{Model: dbmodel.Model{ID: 2}, Name: "Isfahan"},
	}
	s.addressRepository.On("GetProvinceList", s.db).Return(provinces, nil).Once()

	result, err := s.service.GetProvinceList()

	s.NoError(err)
	s.Len(result, 2)
	s.Equal("Tehran", result[0].Name)
	s.addressRepository.AssertExpectations(s.T())
}

func (s *AddressServiceTestSuite) TestGetAddress_NotFound() {
	var nilAddress *entity.Address
	s.addressRepository.On("GetOwnerAddress", s.db, uint(1), "user").Return(nilAddress, nil).Once()

	_, err := s.service.GetAddress(1, "user")

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.addressRepository.AssertExpectations(s.T())
}

func (s *AddressServiceTestSuite) TestGetAddress_RepoError() {
	repoErr := errors.New("db error")
	s.addressRepository.On("GetOwnerAddress", s.db, uint(1), "user").Return((*entity.Address)(nil), repoErr).Once()

	_, err := s.service.GetAddress(1, "user")

	s.ErrorIs(err, repoErr)
	s.addressRepository.AssertExpectations(s.T())
}

func TestAddressServiceSuite(t *testing.T) {
	suite.Run(t, new(AddressServiceTestSuite))
}
