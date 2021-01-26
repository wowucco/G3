package entity

const DeliveryMethodYourself = "yourself"
const DeliveryMethodNovaposhta = "novaposhta"
const DeliveryMethodCourier = "courier"

const PaymentMethodCash = "cash"
const PaymentMethodP2P = "p2p"
const PaymentMethodPayin = "pay-in"
const PaymentMethodCashOnDelivery = "cod"
const PaymentMethodToCard = "to_card"

type City struct {
	ID   string
	Name string
}

type DeliveryInfo struct {
	DeliveryMethod DeliveryMethod
	PaymentMethods []PaymentMethod
	Warehouses     []Warehouse
}

type DeliveryMethod struct {
	ID   int
	Name string
	Slug string
}

type PaymentMethod struct {
	ID   int
	Name string
	Slug string
}

type Warehouse struct {
	ID      string
	Name    string
	Address string
	Phone   string
}
