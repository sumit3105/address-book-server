package dto

type CreateAddressRequest struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"address_line1" validate:"required"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	Pincode      string `json:"pincode"`
}

type UpdateAddressRequest struct {
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	AddressLine1 *string `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	City         *string `json:"city"`
	State        *string `json:"state"`
	Country      *string `json:"country"`
	Pincode      *string `json:"pincode"`
}

type ExportAddressRequest struct {
	Fields []string `json:"fields" validate:"required,min=1"`
	Email  string   `json:"email" validate:"required,email"`
}
