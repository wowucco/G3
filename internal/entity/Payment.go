package entity

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

const PaymentMethodCash = "cash"
const PaymentMethodP2P = "p2p"
const PaymentMethodPayin = "pay-in"
const PaymentMethodCashOnDelivery = "cod"
const PaymentMethodToCard = "to_card"
const PaymentMethodPartsPay = "parts_pay"

const PaymentStatusNew = 1
const PaymentStatusWaitingConfirmation = 2
const PaymentStatusConfirmed = 3
const PaymentStatusPending = 4
const PaymentStatusDone = 5
const PaymentStatusRefund = 6
const PaymentStatusFailed = 7
const PaymentStatusCanceled = 8

const PaymentStatusNewLabel = "New"
const PaymentStatusWaitingConfirmationLabel = "Waiting confirmation"
const PaymentStatusConfirmedLabel = "Confirmed"
const PaymentStatusPendingLabel = "Pending"
const PaymentStatusDoneLabel = "Done"
const PaymentStatusRefundLabel = "Refund"
const PaymentStatusFailedLabel = "Failed"
const PaymentStatusCanceledLabel = "Canceled"

const PaymentInitActionForm = "form"
const PaymentInitActionRedirect = "redirect"
const PaymentInitActionNone = "none"

func CreateNewPaymentByOrder(id int, order *Order) *Payment {

	now := time.Now().Unix()
	price := order.GetPrice()

	return &Payment{
		id:            id,
		transactionId: uuid.NewString(),
		orderId:       order.GetId(),
		provider:      "",
		price:         &price,
		status:        PaymentStatusNew,
		created:       now,
		updated:       now,
	}
}

func NewPayment(id int, transactionId string, orderId int, provider string, price *Price, status int, created, updated int64) *Payment {

	return &Payment{
		id:            id,
		transactionId: transactionId,
		orderId:       orderId,
		provider:      provider,
		price:         price,
		status:        status,
		created:       created,
		updated:       updated,
	}
}

type Payment struct {
	id            int
	transactionId string
	orderId       int
	provider      string
	price         *Price
	status        int
	created       int64
	updated       int64
}

func (p *Payment) GetId() int {

	return p.id
}
func (p *Payment) GetTransactionId() string {

	return p.transactionId
}
func (p *Payment) GetOrderId() int {

	return p.orderId
}
func (p *Payment) GetDescription() string {

	return fmt.Sprintf("Оплата заказа № %d", p.id)
}
func (p *Payment) GetPrice() *Price {

	return p.price
}
func (p *Payment) GetStatus() int {

	return p.status
}
func (p *Payment) GetCreatedTimestamp() int64 {

	return p.created
}
func (p *Payment) GetCreatedTime() time.Time {

	if p.created == 0 {
		return time.Now()
	}

	return time.Unix(p.created, 0)
}
func (p *Payment) GetUpdatedTimestamp() int64 {

	return p.updated
}
func (p *Payment) GetUpdatedTime() time.Time {

	if p.created == 0 {
		return time.Now()
	}

	return time.Unix(p.updated, 0)
}
func (p *Payment) TouchUpdated() {
	p.updated = time.Now().Unix()
}
func (p *Payment) GetProvider() string {
	return p.provider
}
func (p *Payment) SetProvider(provider string) {
	p.provider = provider
}
func (p *Payment) UpdateStatus(status int) {
	p.status = status
}
func (p *Payment) HasEqualStatus(status int) bool {
	return p.status == status
}

//func (p *Payment) UpdateStatus(status int) error {
//
//}

func NewPaymentMethod(id int, name, slug string) *PaymentMethod {
	return &PaymentMethod{id, name, slug}
}

type PaymentMethod struct {
	ID   int
	Name string
	Slug string
}

func (d PaymentMethod) GetID() int {
	return d.ID
}

func (d PaymentMethod) GetName() string {
	return d.Name
}

func (d PaymentMethod) GetSlug() string {
	return d.Slug
}
func StatusLabel(status int) string {
	m := map[int]string{
		PaymentStatusNew: PaymentStatusNewLabel,
		PaymentStatusWaitingConfirmation: PaymentStatusWaitingConfirmationLabel,
		PaymentStatusConfirmed: PaymentStatusConfirmedLabel,
		PaymentStatusPending: PaymentStatusPendingLabel,
		PaymentStatusDone: PaymentStatusDoneLabel,
		PaymentStatusRefund: PaymentStatusRefundLabel,
		PaymentStatusFailed: PaymentStatusFailedLabel,
		PaymentStatusCanceled: PaymentStatusCanceledLabel,
	}

	return m[status]
}
