package checkout

import "github.com/wowucco/G3/internal/entity"

type CreateOrderBuilder struct {
	DeliveryMethod *entity.DeliveryMethod
	PaymentMethod  *entity.PaymentMethod
	Products       []*entity.OrderProduct
	Warehouse      *entity.OrderDeliveryWarehouse
	Customer       *entity.OrderCustomer
	PayInEdrpou    string
	PayInEmail     string
	PayInCompany   string
	PayPartsPay    int
	Comment        string
	DoNotCall      bool
	Cost           int
}
