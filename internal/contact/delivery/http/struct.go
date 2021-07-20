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
