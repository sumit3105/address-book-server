package repository

import (
	"address-book-server/dto"
	"address-book-server/model"

	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(address *model.Address) error
	FindByUser(userID uint64) ([]model.Address, error)
	FindByIDAndUser(id, userID uint64) (*model.Address, error)
	Update(address *model.Address) error
	SoftDelete(id, userID uint64) error
	FindUserWithFilters(userId uint64, query dto.ListAddressQuery) ([]model.Address, int64, error)
}

type addressRepository struct {
	db *gorm.DB
}


func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (repository *addressRepository) Create(address *model.Address) error {
	return repository.db.Create(address).Error
}

func (repository *addressRepository) FindByIDAndUser(id uint64, userID uint64) (*model.Address, error) {
	var address model.Address

	err := repository.db.Where("id = ? AND user_id = ? AND is_deleted = false", id, userID).First(&address).Error

	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (repository *addressRepository) FindByUser(userID uint64) ([]model.Address, error) {
	var addresses []model.Address

	err := repository.db.Where("user_id = ? AND is_deleted = false", userID).Find(&addresses).Error

	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (repository *addressRepository) SoftDelete(id uint64, userID uint64) error {
	return repository.db.Model(&model.Address{}).Where("id = ? AND user_id = ?", id, userID).Update("is_deleted", true).Error
}

func (repository *addressRepository) Update(address *model.Address) error {
	return repository.db.Save(address).Error
}

func (repository *addressRepository) FindUserWithFilters(userId uint64, query dto.ListAddressQuery) ([]model.Address, int64, error) {
	var addresses []model.Address
	var total int64

	db := repository.db.Model(&model.Address{}).Where("user_id = ? AND is_deleted = false", userId)

	if query.Search != "" {
		like := "%" + query.Search + "%"
		db = db.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", like, like, like, like)
	}

	if query.City != "" {
		db = db.Where("city ILIKE ?", query.City)
	}

	if query.Country != "" {
		db = db.Where("country ILIKE ?", query.Country)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	page := query.Page
	limit := query.Limit

	if page <= 0 {
		page = 1
	}

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page-1) * limit

	err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&addresses).Error

	return addresses, total, err
}