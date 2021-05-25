package entity

import "time"

/**
 *************	Order	**************
 */
func NewOrder(id int, created int64, comment string, doNotCall bool, cost int, customer *OrderCustomer, delivery *OrderDelivery, payment *OrderPayment, items []*OrderProduct) *Order {

	return &Order{
		id:        id,
		created:   created,
		comment:   comment,
		doNotCall: doNotCall,
		totalCost: *NewPrice(cost, 0, 0, nil),
		customer:  customer,
		delivery:  delivery,
		payment:   payment,
		items:     items,
	}
}

type Order struct {
	id        int
	created   int64
	comment   string
	doNotCall bool
	totalCost Price

	customer *OrderCustomer
	delivery *OrderDelivery
	payment  *OrderPayment
	items    []*OrderProduct
}

func (o Order) GetId() int {
	return o.id
}
func (o Order) GetCreated() int64 {
	return o.created
}
func (o Order) GetComment() string {
	return o.comment
}
func (o Order) GetDoNotCall() bool {
	return o.doNotCall
}
func (o Order) GetPrice() Price {
	return o.totalCost
}
func (o Order) GetCustomer() OrderCustomer {
	return *o.customer
}
func (o Order) GetDelivery() OrderDelivery {
	return *o.delivery
}
func (o Order) GetPayment() OrderPayment {
	return *o.payment
}
func (o Order) GetItems() []OrderProduct {
	p := make([]OrderProduct, len(o.items))

	for k, v := range o.items {
		p[k] = *v
	}

	return p
}

func (o Order) CanMakePayment() bool {
	return o.payment.status == PaymentStatusNew || o.payment.status == PaymentStatusFailed
}

func (o *Order) UpdatePaymentStatus(status int, comment string) {
	o.payment.status = status
	h := NewOrderPaymentStatusHistory(status, time.Now().Unix(), comment)
	o.payment.statusHistory = append(o.payment.statusHistory, h)
}

/**
****************	Customer	*************
 */

func NewOrderCustomer(name, phone string) *OrderCustomer {

	return &OrderCustomer{
		name:  name,
		phone: phone,
	}
}

type OrderCustomer struct {
	name  string
	phone string
}

func (c OrderCustomer) GetName() string {

	return c.name
}
func (c OrderCustomer) GetPhone() string {

	return c.phone
}

/**
****************	Delivery	*****************
 */
func NewOrderDelivery(status int, method *DeliveryMethod, warehouse *OrderDeliveryWarehouse, statusHistory []*OrderDeliveryStatusHistory) *OrderDelivery {
	return &OrderDelivery{
		status:        status,
		method:        method,
		warehouse:     warehouse,
		statusHistory: statusHistory,
	}
}

type OrderDelivery struct {
	status        int
	method        *DeliveryMethod
	warehouse     *OrderDeliveryWarehouse
	statusHistory []*OrderDeliveryStatusHistory
}

func (o OrderDelivery) GetStatus() int {
	return o.status
}
func (o OrderDelivery) GetMethod() DeliveryMethod {
	return *o.method
}
func (o OrderDelivery) GetWarehouse() OrderDeliveryWarehouse {
	return *o.warehouse
}
func (o OrderDelivery) GetStatusHistory() []OrderDeliveryStatusHistory {
	s := make([]OrderDeliveryStatusHistory, len(o.statusHistory))

	for k, v := range o.statusHistory {
		s[k] = *v
	}

	return s
}

func NewOrderDeliveryStatusHistory(status int, created int64, comment string) *OrderDeliveryStatusHistory {
	return &OrderDeliveryStatusHistory{status, created, comment}
}

type OrderDeliveryStatusHistory struct {
	status  int
	created int64
	comment string
}

func (s OrderDeliveryStatusHistory) GetStatus() int {
	return s.status
}
func (s OrderDeliveryStatusHistory) GetCreated() int64 {
	return s.created
}
func (s OrderDeliveryStatusHistory) GetComment() string {
	return s.comment
}

type OrderDeliveryWarehouseCity struct {
	id   string
	name string
}

func (c OrderDeliveryWarehouseCity) GetId() string {
	return c.id
}

func (c OrderDeliveryWarehouseCity) GetName() string {
	return c.name
}

type OrderDeliveryWarehouseAddress struct {
	id     string
	name   string
	custom bool
}

func (a OrderDeliveryWarehouseAddress) GetId() string {
	return a.id
}

func (a OrderDeliveryWarehouseAddress) GetName() string {
	return a.name
}

func (a OrderDeliveryWarehouseAddress) IsCustom() bool {
	return a.custom
}

func NewOrderDeliveryWarehouse(cityId, cityName, addressId, addressName string, isCustom bool) *OrderDeliveryWarehouse {

	return &OrderDeliveryWarehouse{
		city:    OrderDeliveryWarehouseCity{cityId, cityName},
		address: OrderDeliveryWarehouseAddress{addressId, addressName, isCustom},
	}
}

type OrderDeliveryWarehouse struct {
	city    OrderDeliveryWarehouseCity
	address OrderDeliveryWarehouseAddress
}

func (w OrderDeliveryWarehouse) GetCity() OrderDeliveryWarehouseCity {

	return w.city
}

func (w OrderDeliveryWarehouse) GetAddress() OrderDeliveryWarehouseAddress {

	return w.address
}

/**
******************	Payment ************
 */
func NewOrderPayment(status int, method *PaymentMethod, statusHistory []*OrderPaymentStatusHistory, edrpou, company, email string, partsPay int) *OrderPayment {
	return &OrderPayment{
		method: method,
		status: status,
		statusHistory: statusHistory,
		extra: OrderPaymentExtra{
			edrpou:   edrpou,
			company:  company,
			email:    email,
			partsPay: partsPay,
		},
	}
}

type OrderPayment struct {
	method        *PaymentMethod
	status        int
	statusHistory []*OrderPaymentStatusHistory
	extra         OrderPaymentExtra
}

func (o OrderPayment) GetStatus() int {
	return o.status
}
func (o OrderPayment) GetMethod() PaymentMethod {
	return *o.method
}
func (o OrderPayment) GetStatusHistory() []OrderPaymentStatusHistory {
	p := make([]OrderPaymentStatusHistory, len(o.statusHistory))

	for k, v := range o.statusHistory {
		p[k] = *v
	}

	return p
}
func (o OrderPayment) GetExtra() OrderPaymentExtra {
	return o.extra
}

func NewOrderPaymentStatusHistory(status int, created int64, comment string) *OrderPaymentStatusHistory {
	return &OrderPaymentStatusHistory{status, created, comment}
}

type OrderPaymentStatusHistory struct {
	status  int
	created int64
	comment string
}

func (s OrderPaymentStatusHistory) GetStatus() int {
	return s.status
}
func (s OrderPaymentStatusHistory) GetCreated() int64 {
	return s.created
}
func (s OrderPaymentStatusHistory) GetComment() string {
	return s.comment
}
type OrderPaymentExtra struct {
	edrpou   string
	company  string
	email    string
	partsPay int
}

func (e OrderPaymentExtra) GetEdrpou() string {
	return e.edrpou
}
func (e OrderPaymentExtra) GetCompany() string {
	return e.company
}
func (e OrderPaymentExtra) GetEmail() string {
	return e.email
}
func (e OrderPaymentExtra) GetPartsPay() int {
	return e.partsPay
}

/**
***********	Products ***************
 */
func NewOrderProduct(quantity, price int, product SimpleProduct) *OrderProduct {

	return &OrderProduct{
		quantity: quantity,
		price:    *NewPrice(price, 0, 0, nil),
		product:  product,
	}
}

type OrderProduct struct {
	quantity int
	price    Price
	product  SimpleProduct
}

func (o OrderProduct) GetQuantity() int {
	return o.quantity
}
func (o OrderProduct) GetPrice() Price {
	return o.price
}
func (o OrderProduct) GetProduct() SimpleProduct {
	return o.product
}
