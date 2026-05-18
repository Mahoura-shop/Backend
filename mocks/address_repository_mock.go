package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type AddressRepositoryMock struct {
	mock.Mock
}

func NewAddressRepositoryMock() *AddressRepositoryMock {
	return &AddressRepositoryMock{}
}

func (m *AddressRepositoryMock) CreateAddress(db database.Database, address *entity.Address) error {
	args := m.Called(db, address)
	return args.Error(0)
}

func (m *AddressRepositoryMock) CreateCity(db database.Database, city *entity.City) error {
	args := m.Called(db, city)
	return args.Error(0)
}

func (m *AddressRepositoryMock) CreateProvince(db database.Database, province *entity.Province) error {
	args := m.Called(db, province)
	return args.Error(0)
}

func (m *AddressRepositoryMock) DeleteAddress(db database.Database, address *entity.Address) error {
	args := m.Called(db, address)
	return args.Error(0)
}

func (m *AddressRepositoryMock) GetAddressByID(db database.Database, id uint) (*entity.Address, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Address), args.Error(1)
}

func (m *AddressRepositoryMock) GetCityByID(db database.Database, id uint) (*entity.City, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.City), args.Error(1)
}

func (m *AddressRepositoryMock) GetCityByName(db database.Database, name string) (*entity.City, error) {
	args := m.Called(db, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.City), args.Error(1)
}

func (m *AddressRepositoryMock) GetOwnerAddress(db database.Database, ownerID uint, ownerType string) (*entity.Address, error) {
	args := m.Called(db, ownerID, ownerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Address), args.Error(1)
}

func (m *AddressRepositoryMock) GetOwnerAddresses(db database.Database, ownerID uint, ownerType string) ([]*entity.Address, error) {
	args := m.Called(db, ownerID, ownerType)
	return args.Get(0).([]*entity.Address), args.Error(1)
}

func (m *AddressRepositoryMock) GetProvinceByID(db database.Database, id uint) (*entity.Province, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Province), args.Error(1)
}

func (m *AddressRepositoryMock) GetProvinceByName(db database.Database, name string) (*entity.Province, error) {
	args := m.Called(db, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Province), args.Error(1)
}

func (m *AddressRepositoryMock) GetProvinceCities(db database.Database, provinceID uint) ([]*entity.City, error) {
	args := m.Called(db, provinceID)
	return args.Get(0).([]*entity.City), args.Error(1)
}

func (m *AddressRepositoryMock) GetProvinceList(db database.Database) ([]*entity.Province, error) {
	args := m.Called(db)
	return args.Get(0).([]*entity.Province), args.Error(1)
}
