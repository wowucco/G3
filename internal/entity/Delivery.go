package entity

const DeliveryMethodYourself = "yourself"
const DeliveryMethodNovaposhta = "novaposhta"
const DeliveryMethodCourier = "courier"

const DeliveryStatusNew = 1
const DeliveryStatusCheck = 2
const DeliveryStatusWaitingDelivery = 3
const DeliveryStatusDelivery = 4
const DeliveryStatusReadyToReceive = 5
const DeliveryStatusCanceled = 6

type City struct {
	ID   string
	Name string
}

type DeliveryInfo struct {
	DeliveryMethod DeliveryMethod
	PaymentMethods []PaymentMethod
	Warehouses     []Warehouse
}

func NewDeliveryMethod(id int, name, slug string, ) *DeliveryMethod {
	return &DeliveryMethod{id, name, slug}
}

type DeliveryMethod struct {
	ID   int
	Name string
	Slug string
}

func (d DeliveryMethod) GetID() int {
	return d.ID
}
func (d DeliveryMethod) GetName() string {
	return d.Name
}
func (d DeliveryMethod) GetSlug() string {
	return d.Slug
}

type Warehouse struct {
	ID        string
	Name      string
	Address   string
	Phone     string
	Number    int
	MaxWeight int
}
