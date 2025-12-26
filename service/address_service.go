package service

import (
	"address-book-server/dto"
	appError "address-book-server/error"
	"address-book-server/logger"
	"address-book-server/mapper"
	"address-book-server/model"
	"address-book-server/repository"
	"address-book-server/utils"

	"errors"
	"strconv"

	"go.uber.org/zap"
)

type AddressService interface {
	Create(userId uint64, address *model.Address) error
	List(userId uint64) ([]dto.ListAddressResponse, error)
	Update(id, userId uint64, req *dto.UpdateAddressRequest) error
	Delete(id, userId uint64) error
	ExportCSV(userId uint64, fields []string) ([]byte, error)
	ListWithFilters(userId uint64, query dto.ListAddressQuery) ([]dto.ListAddressResponse, int64, error)
}

type addressService struct {
	repo repository.AddressRepository
}

func NewAddressService(repo repository.AddressRepository) AddressService {
	return &addressService{repo: repo}
}

func (s *addressService) Create(userId uint64, address *model.Address) error {
	address.UserID = userId

	logger.Log.Info(
		"Adding new Address",
		zap.Uint64("user_id", userId),
	)

	if err := s.repo.Create(address); err != nil {

		logger.Log.Error(
			"Failed to create address",
			zap.String("error", err.Error()),
		)

		return appError.Internal(
			"Failed to create address",
			err,
		)
	}

	logger.Log.Info(
		"New Address added",
		zap.Uint64("id", address.ID),
		zap.String("email", address.Email),
	)

	return nil
}

func (s *addressService) List(userId uint64) ([]dto.ListAddressResponse, error) {

	logger.Log.Info(
		"Finding Addresses",
		zap.Uint64("user_id", userId),
	)

	addresses, err := s.repo.FindByUser(userId)

	if err != nil {

		logger.Log.Error(
			"Failed to fetch addresses",
			zap.String("error", err.Error()),
		)

		return nil, appError.Internal(
			"Failed to fetch addresses",
			err,
		)
	}

	response := make([]dto.ListAddressResponse, 0, len(addresses))

	for _, a := range addresses {
		response = append(response, mapper.ToListAddressResponse(a))
	}

	return response, nil
}

func (s *addressService) ListWithFilters(userId uint64, query dto.ListAddressQuery) ([]dto.ListAddressResponse, int64, error) {

	logger.Log.Info(
		"Finding Addresses",
		zap.Uint64("user_id", userId),
	)

	addresses, total, err := s.repo.FindUserWithFilters(userId, query)

	if err != nil {

		logger.Log.Error(
			"Failed to fetch addresses",
			zap.String("error", err.Error()),
		)

		return nil, 0, appError.Internal(
			"Failed to fetch addresses",
			err,
		)
	}

	resp := make([]dto.ListAddressResponse, 0, len(addresses))
	for _, a := range addresses {
		resp = append(resp, mapper.ToListAddressResponse(a))
	}

	logger.Log.Info(
		"Addresses found",
		zap.Int64("total", total),
	)

	return resp, total, nil
}

func (s *addressService) Update(id, userId uint64, req *dto.UpdateAddressRequest) error {

	logger.Log.Info(
		"Updating Address",
		zap.Uint64("address_id", id),
		zap.Uint64("user_id", userId),
	)

	address, err := s.repo.FindByIDAndUser(id, userId)
	if err != nil {

		logger.Log.Error(
			"Address not found",
			zap.String("error", err.Error()),
		)

		return appError.NotFound(
			"Address not found",
			err,
		)
	}

	if req.FirstName != nil {
		address.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		address.LastName = *req.LastName
	}
	if req.Email != nil {
		address.Email = *req.Email
	}
	if req.Phone != nil {
		address.Phone = *req.Phone
	}
	if req.AddressLine1 != nil {
		address.AddressLine1 = *req.AddressLine1
	}
	if req.AddressLine2 != nil {
		address.AddressLine2 = *req.AddressLine2
	}
	if req.City != nil {
		address.City = *req.City
	}
	if req.State != nil {
		address.State = *req.State
	}
	if req.Country != nil {
		address.Country = *req.Country
	}
	if req.Pincode != nil {
		address.Pincode = *req.Pincode
	}

	if err := s.repo.Update(address); err != nil {

		logger.Log.Error(
			"Failed to update address",
			zap.String("error", err.Error()),
		)

		return appError.Internal(
			"Failed to update address",
			err,
		)
	}

	logger.Log.Info(
		"Update Successfull",
		zap.Uint64("address_id", id),
		zap.Uint64("user_id", userId),
	)

	return nil
}

func (s *addressService) Delete(id, userId uint64) error {

	logger.Log.Info(
		"Deleting Address",
		zap.Uint64("address_id", id),
		zap.Uint64("user_id", userId),
	)

	_, err := s.repo.FindByIDAndUser(id, userId)
	if err != nil {

		logger.Log.Error(
			"Address not found",
			zap.String("error", err.Error()),
		)

		return appError.NotFound(
			"Address not found",
			err,
		)
	}

	if err := s.repo.SoftDelete(id, userId); err != nil {

		logger.Log.Error(
			"Failed to delete address",
			zap.String("error", err.Error()),
		)

		return appError.Internal(
			"Failed to delete address",
			err,
		)
	}

	logger.Log.Info(
		"Soft delete Successfull",
		zap.Uint64("address_id", id),
		zap.Uint64("user_id", userId),
	)

	return nil
}

func (s *addressService) ExportCSV(userId uint64, fields []string) ([]byte, error) {

	logger.Log.Info(
		"Exporting started in CSV format",
		zap.Strings("fields", fields),
	)

	for _, f := range fields {
		if _, ok := utils.AllowedAddressExportFields[f]; !ok {
			return nil, appError.BadRequest(
				"Invalid export field: "+f,
				errors.New(f),
			)
		}
	}

	addresses, err := s.repo.FindByUser(userId)
	if err != nil {
		return nil, appError.NotFound(
			"Address not found",
			err,
		)
	}

	records := make([]map[string]string, 0, len(addresses))
	for _, a := range addresses {
		records = append(records, map[string]string{
			"id":            strconv.FormatUint(a.ID, 10),
			"user_id":       strconv.FormatUint(a.UserID, 10),
			"first_name":    a.FirstName,
			"last_name":     a.LastName,
			"email":         a.Email,
			"phone":         a.Phone,
			"address_line1": a.AddressLine1,
			"address_line2": a.AddressLine2,
			"city":          a.City,
			"state":         a.State,
			"country":       a.Country,
			"pincode":       a.Pincode,
		})
	}

	return utils.GenerateAddressCSV(fields, records)
}
