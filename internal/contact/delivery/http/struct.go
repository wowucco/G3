package http

import validation "github.com/go-ozzo/ozzo-validation/v4"

type RecallForm struct {
	Message string `json:"message"`
	Phone   string `json:"phone"`
}

func (r RecallForm) GetPhone() string {

	return r.Phone
}

func (r RecallForm) GetMessage() string {

	return r.Message
}

func (r RecallForm) Validate() error {
	// todo phone regexp
	return validation.ValidateStruct(&r, validation.Field(&r.Phone, validation.Required))
}

type BuyOnClickForm struct {
	ProductId int    `json:"product_id"`
	Phone     string `json:"phone"`
}

func (b BuyOnClickForm) GetPhone() string {

	return b.Phone
}

func (b BuyOnClickForm) GetProductId() int {

	return b.ProductId
}

func (b BuyOnClickForm) Validate() error {
	// todo phone regexp
	return validation.ValidateStruct(&b,
		validation.Field(&b.Phone, validation.Required),
		validation.Field(&b.ProductId, validation.Required, validation.Min(1)),
	)
}
