package mapper

import (
	"address-book-server/dto"
	"address-book-server/model"
)

func ToListAddressResponse(address model.Address) dto.ListAddressResponse {
	return dto.ListAddressResponse{
		Id:           address.ID,
		UserId:       address.UserID,
		FirstName:    address.FirstName,
		LastName:     address.LastName,
		Email:        address.Email,
		Phone:        address.Phone,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		State:        address.State,
		Country:      address.Country,
		Pincode:      address.Pincode,
	}
}
