package dto

type ListAddressQuery struct {
	Page int `form:"page"`
	Limit int `form:"limit"`
	Search string `form:"search"`
	City string `form:"city"`
	Country string `form:"country"`
}