package controller

import (
	"address-book-server/dto"
	"address-book-server/model"
	"address-book-server/service"
	"address-book-server/utils"
	"address-book-server/validator"
	appError "address-book-server/error"
	"log"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AddressController interface {
	List(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Export(ctx *gin.Context)
	runExportJob(userId uint64, req dto.ExportAddressRequest)
}

type addressController struct {
	addressService service.AddressService
}

func NewAddressController(addressService service.AddressService) AddressController {
	return &addressController{addressService: addressService}
}

func (c *addressController) List(ctx *gin.Context) {

	userId := ctx.GetUint64("user_id")

	var query dto.ListAddressQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {

		ctx.Error(
			appError.BadRequest(
				"Invalid request body",
				err,
			),
		)
		return
	}

	response, total, err := c.addressService.ListWithFilters(userId, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	page := query.Page
	limit := query.Limit

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"addresses": response,
		},
		"meta": gin.H{
			"page": page,
			"limit": limit,
			"total": total,
			"total_pages": (total + int64(limit) - 1) / int64(limit), 
		},
	})
}

func (c *addressController) Create(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")

	var req dto.CreateAddressRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(
			appError.BadRequest(
				"Invalid request Body",
				err,
			),
		)
		return 
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.Error(
			appError.NewValidationError(
				utils.FormatValidationErrors(err),
			),
		)
		return
	}

	address := model.Address{
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		Phone: req.Phone,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City: req.City,
		State: req.State,
		Country: req.Country,
		Pincode: req.Pincode,
	}

	if err := c.addressService.Create(userId, &address); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "Contact added to the Dictionary",
		},
	})
}

func (c *addressController) Update(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")
	
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	
	if err != nil {
		ctx.Error(
			appError.BadRequest(
				"Invalid ID",
				err,
			),
		)
		return
	}

	var req dto.UpdateAddressRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(
			appError.BadRequest(
				"Invalid request Body",
				err,
			),
		)
		return
	}

	if err := c.addressService.Update(id, userId, &req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "Contact updated",
		},
	})
}

func (c *addressController) Delete(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.Error(
			appError.BadRequest(
				"Invalid address ID",
				err,
			),
		)
		return
	}

	if err := c.addressService.Delete(id, userId); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "Contact deleted from the Dictionary",
		},
	})
}

func (c *addressController) Export(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")

	var req dto.ExportAddressRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(
			appError.BadRequest(
				"Invalid request Body",
				err,
			),
		)
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.Error(
			appError.NewValidationError(
				utils.FormatValidationErrors(err),
			),
		)
		return
	}

	go func(userId uint64, req dto.ExportAddressRequest) {
		c.runExportJob(userId, req)
	}(userId, req)
	
	ctx.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "Export started. CSV will be sent to your email shortly.",
		},
	})

}

func (c *addressController) runExportJob(userId uint64, req dto.ExportAddressRequest) {
	csvData, err := c.addressService.ExportCSV(userId, req.Fields)
	
	if err != nil {
		log.Printf("Export failed for user %d: %v", userId, err)
		return
	}

	err = utils.SendEmailWithCSV(
		req.Email,
		"Address Book Export",
		"Please find your address book export attached.",
		"address_book.csv",
		csvData,
	)

	if err != nil {
		log.Printf("Email sending failed for user %d: %v", userId, err)
		return
	}
}