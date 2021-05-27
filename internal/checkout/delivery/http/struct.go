package http

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/wowucco/G3/internal/checkout"
	"github.com/wowucco/G3/internal/entity"
)

type Client struct {
	Fio   string `json:"fio"`
	Phone string `json:"phone"`
}

func (c Client) GetFio() string {

	return c.Fio
}

func (c Client) GetPhone() string {

	return c.Phone
}

func (c Client) Validate() error {

	return validation.ValidateStruct(&c,
		validation.Field(&c.Fio, validation.Required),
		// todo phone regexp
		validation.Field(&c.Phone, validation.Required),
	)
}

type City struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (c City) GetCode() string {

	return c.Code
}

func (c City) GetName() string {

	return c.Name
}

func (c City) Validate() error {

	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Code, validation.Required),
	)
}

type Address struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (a Address) GetCode() string {

	return a.Code
}

func (a Address) GetName() string {

	return a.Name
}

func (a Address) Validate() error {

	return validation.ValidateStruct(&a, validation.Field(&a.Name, validation.Required))
}

type Delivery struct {
	Method        string  `json:"method"`
	City          City    `json:"city"`
	CustomAddress bool    `json:"is_custom_address"`
	Address       Address `json:"address"`
}

func (d Delivery) GetMethod() string {

	return d.Method
}

func (d Delivery) GetCity() checkout.DeliveryCityForm {

	return d.City
}

func (d Delivery) IsCustomAddress() bool {

	return d.CustomAddress
}

func (d Delivery) GetAddress() checkout.DeliveryAddressForm {

	return d.Address
}

func (d Delivery) Validate() error {

	return validation.ValidateStruct(&d,
		validation.Field(&d.Method, validation.Required, validation.In(entity.DeliveryMethodYourself, entity.DeliveryMethodNovaposhta, entity.DeliveryMethodCourier)),
		validation.Field(&d.City),
		validation.Field(&d.Address),
	)
}

type Payment struct {
	Method       string `json:"method"`
	PayInEdrpou  string `json:"pay_in_edrpou"`
	PayInEmail   string `json:"pay_in_email"`
	PayInCompany string `json:"pay_in_company"`
	PayPartsPay  int    `json:"pay_parts_pay"`
}

func (p Payment) GetMethod() string {

	return p.Method
}
func (p Payment) GetPayInEdrpou() string {

	return p.PayInEdrpou
}
func (p Payment) GetPayInEmail() string {

	return p.PayInEmail
}
func (p Payment) GetPayInCompany() string {

	return p.PayInCompany
}
func (p Payment) GetPayPartsPay() int {

	return p.PayPartsPay
}

func (p Payment) Validate() error {

	return validation.ValidateStruct(&p,
		validation.Field(&p.Method, validation.Required, validation.In(entity.PaymentMethodCash, entity.PaymentMethodP2P, entity.PaymentMethodPayin, entity.PaymentMethodCashOnDelivery, entity.PaymentMethodToCard, entity.PaymentMethodPartsPay)),
		validation.Field(&p.PayInCompany, validation.When(p.Method == entity.PaymentMethodPayin, validation.Required)),
		validation.Field(&p.PayInEdrpou, validation.When(p.Method == entity.PaymentMethodPayin, validation.Required)),
		validation.Field(&p.PayInEmail, validation.When(p.Method == entity.PaymentMethodPayin, validation.Required, is.Email)),
		validation.Field(&p.PayPartsPay, validation.When(p.Method == entity.PaymentMethodPartsPay, validation.Required)),
	)
}

type OrderItem struct {
	ProductId int `json:"product_id"`
	Count     int `json:"count"`
	Price     int `json:"price"`
}

func (o OrderItem) GetProductId() int {

	return o.ProductId
}
func (o OrderItem) GetCount() int {

	return o.Count
}
func (o OrderItem) GetPrice() int {

	return o.Price
}
func (o OrderItem) Validate() error {

	return validation.ValidateStruct(&o,
		validation.Field(&o.ProductId, validation.Required, validation.Min(1)),
		validation.Field(&o.Count, validation.Required, validation.Min(1)),
		validation.Field(&o.Price, validation.Required, validation.Min(1)),
	)
}

type Order struct {
	Cost  int         `json:"cost"`
	Items []OrderItem `json:"items"`
}

func (o Order) GetCost() int {

	return o.Cost
}
func (o Order) GetOrderItems() []checkout.OrderItemForm {

	i := make([]checkout.OrderItemForm, len(o.Items))
	for k, v := range o.Items {
		i[k] = v
	}

	return i
}
func (o Order) Validate() error {

	err := validation.ValidateStruct(&o,
		validation.Field(&o.Cost, validation.Required, validation.Min(1)),
		validation.Field(&o.Items, validation.Required),
	)

	if err != nil {
		return err
	}

	return validation.Validate(o.Items)
}

type CreateOrder struct {
	Client   Client   `json:"client"`
	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`
	Order    Order    `json:"order"`

	Comment   string `json:"comment"`
	DoNotCall bool `json:"do_not_call"`
}

func (c CreateOrder) GetDelivery() checkout.DeliveryForm {

	return c.Delivery
}
func (c CreateOrder) GetClient() checkout.ClientForm {

	return c.Client
}
func (c CreateOrder) GetPayment() checkout.PaymentForm {

	return c.Payment
}
func (c CreateOrder) GetOrder() checkout.OrderForm {

	return c.Order
}
func (c CreateOrder) GetComment() string {

	return c.Comment
}
func (c CreateOrder) GetDoNotCall() bool {

	return c.DoNotCall
}
func (c CreateOrder) Validate() error {

	return validation.ValidateStruct(&c,
		validation.Field(&c.Client),
		validation.Field(&c.Delivery),
		validation.Field(&c.Order),
	)
}
func NewCreateOrderResponse(order *entity.Order) *OrderResponse {

	return &OrderResponse{
		OrderId: order.GetId(),
		Created: order.GetCreated(),
	}
}
type OrderResponse struct {
	OrderId int	`json:"order_id"`
	Created int64 `json:"created"`
}

type InitPaymentForm struct {
	OrderId int `json:"order_id"`
}

func (f InitPaymentForm) GetOrderId() int {
	return f.OrderId
}
func (f InitPaymentForm) Validate() error {
	return validation.ValidateStruct(&f, validation.Field(&f.OrderId, validation.Required, validation.Min(1)))
}

type AccentPaymentForm struct {
	TransactionId string `json:"transaction_id"`
}

func (f AccentPaymentForm) GetTransactionId() string {
	return f.TransactionId
}
func (f AccentPaymentForm) Validate() error {
	return validation.ValidateStruct(&f, validation.Field(&f.TransactionId, validation.Required))
}

type ProviderCallbackPaymentForm struct {
	Provider string
}

func (f ProviderCallbackPaymentForm) GetProvider() string {
	return f.Provider
}

