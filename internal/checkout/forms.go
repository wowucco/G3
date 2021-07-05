package checkout

type ClientForm interface {
	GetFio() string
	GetPhone() string
}

type DeliveryCityForm interface {
	GetCode() string
	GetName() string
}

type DeliveryAddressForm interface {
	GetCode() string
	GetName() string
}

type DeliveryForm interface {
	GetMethod() string
	GetCity() DeliveryCityForm
	IsCustomAddress() bool
	GetAddress() DeliveryAddressForm
}

type PaymentForm interface {
	GetMethod() string
	GetPayInEdrpou() string
	GetPayInEmail() string
	GetPayInCompany() string
	GetPayPartsPay() int
}

type OrderItemForm interface {
	GetProductId() int
	GetCount() int
	GetPrice() int
}

type OrderForm interface {
	GetCost() int
	GetOrderItems() []OrderItemForm
}

type CreateOrderForm interface {
	GetClient() ClientForm
	GetDelivery() DeliveryForm
	GetPayment() PaymentForm
	GetOrder() OrderForm

	GetComment() string
	GetDoNotCall() bool
}

type InitPaymentForm interface {
	GetOrderId() int
}

type OrderIdForm interface {
	GetOrderId() int
}

type IAcceptHoldenPaymentForm interface {
	GetTransactionId() string
}

type IProviderCallbackPaymentForm interface {
	GetProvider() string
}
