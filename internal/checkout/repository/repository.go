package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/wowucco/G3/internal/checkout"
	"github.com/wowucco/G3/internal/entity"
	"log"
	"strconv"
	"time"
)

const tableNameOrder = "shop_order"
const tableNameProducts = "shop_products"
const tableNameCurrency = "shop_currency"
const tableOrderSeqNextValID = "shop_order_id_seq"
const tablePaymentSeqNextValID = "payment_id_seq"
const tableNameOrderItems = "shop_order_items"
const tableNamePayments = "payment"
const tableNameDeliveryMethods = "shop_delivery_method"
const tableNamePaymentMethods = "shop_payment_method"

func NewOrderRepository(db *dbx.DB) *OrderRepository {

	return &OrderRepository{db: db}
}

type OrderRepository struct {
	db *dbx.DB
}

func (r OrderRepository) Get(ctx context.Context, orderId int) (*entity.Order, error) {

	var (
		deliveryInfo DeliveryInfo
		row          Order
		rowsOI       []OrderItem
		ds           []DeliveryStatus
		ps           []PaymentStatus

		order *entity.Order
	)

	err := r.db.Select(
		"dm.id delivery_id", "dm.slug delivery_slug", "dm.name delivery_name",
		"pm.id payment_id", "pm.slug payment_slug", "pm.name payment_name",
		"o.*").
		From(tableWithAlias(tableNameOrder, "o")).
		InnerJoin(tableWithAlias(tableNameDeliveryMethods, "dm"), dbx.NewExp("dm.id = o.delivery_method_id")).
		InnerJoin(tableWithAlias(tableNamePaymentMethods, "pm"), dbx.NewExp("pm.id = o.payment_method_id")).
		Where(dbx.NewExp("o.id={:id}", dbx.Params{"id": orderId})).
		One(&row)

	if err != nil {
		log.Printf("[Get order %v][query error] err: %v", orderId, err)
		return nil, err
	}

	if row.Id == 0 {
		log.Printf("[Get order %v][order not fount]", orderId)
		return nil, fmt.Errorf("order %d not found", orderId)
	}

	err = r.db.Select(
		"p.id product_id", "p.name product_name", "p.code product_code", "p.exist product_exist", "p.status product_status",
		"p.price product_price", "p.sale_price product_sale_price", "p.sale_count product_sale_count",
		"c.id currency_id", "c.name currency_name", "c.rate currency_rate", "c.iso currency_iso",
		"oi.id", "oi.price", "oi.quantity").
		From(tableWithAlias(tableNameOrderItems, "oi")).
		InnerJoin(tableWithAlias(tableNameProducts, "p"), dbx.NewExp("oi.product_id=p.id")).
		InnerJoin(tableWithAlias(tableNameCurrency, "c"), dbx.NewExp("p.currency_id=c.id")).
		Where(dbx.NewExp("oi.order_id={:order_id}", dbx.Params{"order_id": orderId})).
		All(&rowsOI)

	if err != nil {
		log.Printf("[Get order %v][items query error] err: %v", orderId, err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(row.DeliveryInfo), &deliveryInfo); err != nil {
		log.Printf("[Get order %v][decode delivery info] err: %v", orderId, err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(row.DeliveryStatusesJson), &ds); err != nil {
		log.Printf("[Get order %v][decode delivery statuses] err: %v", orderId, err)
		return nil, err
	}

	deliveryStatuses := make([]*entity.OrderDeliveryStatusHistory, len(ds))

	for k, v := range ds {
		deliveryStatuses[k] = entity.NewOrderDeliveryStatusHistory(v.Status, v.Time, v.Comment)
	}

	if err = json.Unmarshal([]byte(row.PaymentStatusesJson), &ps); err != nil {
		log.Printf("[Get order %v][decode payment statuses] err: %v", orderId, err)
		return nil, err
	}

	paymentStatuses := make([]*entity.OrderPaymentStatusHistory, len(ps))

	for k, v := range ps {
		paymentStatuses[k] = entity.NewOrderPaymentStatusHistory(v.Status, v.Time, v.Comment)
	}

	oi := make([]*entity.OrderProduct, len(rowsOI))
	for k, v := range rowsOI {
		sp := 0
		sc := 0

		if v.ProductSalePrice.Valid == true && v.ProductSaleCount.Valid == true {
			sp, _ = strconv.Atoi(v.ProductSalePrice.String)
			sc, _ = strconv.Atoi(v.ProductSaleCount.String)
		}

		price := entity.NewPrice(v.ProductPrice, sp, sc, entity.NewCurrency(v.CurrencyID, v.CurrencyName, v.CurrencyRate, v.CurrencyISO))
		pr := entity.NewSimpleProduct(
			v.ProductId,
			v.ProductName,
			v.ProductCode,
			v.ProductExist,
			v.ProductStatus,
			*price,
		)
		oi[k] = entity.NewOrderProduct(v.Quantity, v.Cost, *pr)
	}

	order = entity.NewOrder(
		row.Id,
		row.Created,
		row.Comment,
		row.DoNotCall,
		row.Cost,
		entity.NewOrderCustomer(row.Customer.Name, row.Customer.Phone),
		entity.NewOrderDelivery(
			row.DeliveryStatus,
			entity.NewDeliveryMethod(row.DeliveryMethod.DeliveryId, row.DeliveryMethod.DeliveryName, row.DeliveryMethod.DeliverySlug),
			entity.NewOrderDeliveryWarehouse(deliveryInfo.City.Code, deliveryInfo.City.Name, deliveryInfo.Address.Code, deliveryInfo.Address.Address, deliveryInfo.Address.IsCustom),
			deliveryStatuses,
		),
		entity.NewOrderPayment(
			row.PaymentStatus,
			entity.NewPaymentMethod(row.PaymentMethod.PaymentId, row.PaymentMethod.PaymentName, row.PaymentMethod.PaymentSlug),
			paymentStatuses,
			deliveryInfo.Payment.Edrpou,
			deliveryInfo.Payment.Company,
			deliveryInfo.Payment.Email,
			deliveryInfo.Payment.PartsPay,
		),
		oi,
	)

	return order, nil
}

func (r OrderRepository) Save(ctx context.Context, order *entity.Order) error {

	price := order.GetPrice()

	dInfo, err := json.Marshal(DeliveryInfo{
		City: City{
			Name: order.GetDelivery().GetWarehouse().GetCity().GetName(),
			Code: order.GetDelivery().GetWarehouse().GetCity().GetId(),
		},
		Address: Address{
			IsCustom: order.GetDelivery().GetWarehouse().GetAddress().IsCustom(),
			Address:  order.GetDelivery().GetWarehouse().GetAddress().GetName(),
			Code:     order.GetDelivery().GetWarehouse().GetAddress().GetId(),
		},
		Payment: PaymentExtra{
			Edrpou:   order.GetPayment().GetExtra().GetEdrpou(),
			Company:  order.GetPayment().GetExtra().GetCompany(),
			Email:    order.GetPayment().GetExtra().GetEmail(),
			PartsPay: order.GetPayment().GetExtra().GetPartsPay(),
		},
	})

	if err != nil {
		return err
	}

	ds := make([]DeliveryStatus, len(order.GetDelivery().GetStatusHistory()))
	for k, v := range order.GetDelivery().GetStatusHistory() {
		ds[k] = DeliveryStatus{
			Status:  v.GetStatus(),
			Time:    v.GetCreated(),
			Comment: v.GetComment(),
		}
	}
	dStatuses, err := json.Marshal(ds)

	if err != nil {
		return err
	}

	ps := make([]PaymentStatus, len(order.GetPayment().GetStatusHistory()))
	for k, v := range order.GetPayment().GetStatusHistory() {
		ps[k] = PaymentStatus{
			Status:  v.GetStatus(),
			Time:    v.GetCreated(),
			Comment: v.GetComment(),
		}
	}
	pStatuses, err := json.Marshal(ps)

	if err != nil {
		return err
	}

	_, err = r.db.Update(tableNameOrder, dbx.Params{
		"customer_phone": order.GetCustomer().GetPhone(),
		"customer_name":  order.GetCustomer().GetName(),

		"delivery_method_id": order.GetDelivery().GetMethod().GetID(),
		"delivery_status":    order.GetDelivery().GetStatus(),

		"payment_method_id": order.GetPayment().GetMethod().GetID(),
		"payment_status":    order.GetPayment().GetStatus(),

		"cost":        price.GetInCent(),
		"comment":     order.GetComment(),
		"do_not_call": order.GetDoNotCall(),

		"delivery_info":          string(dInfo),
		"delivery_statuses_json": string(dStatuses),
		"payment_statuses_json":  string(pStatuses),
	}, dbx.NewExp("id={:id}", dbx.Params{"id": order.GetId()})).
		Execute()

	return err
}

func (r OrderRepository) NextId() (int, error) {

	var seq NextId

	err := r.db.NewQuery(fmt.Sprintf("SELECT nextval('%s') as id", tableOrderSeqNextValID)).One(&seq)

	if err != nil {
		return 0, errors.New(fmt.Sprintf("[get order next sequence][%v]", err))
	}

	return seq.Id, nil
}

func (r OrderRepository) Create(ctx context.Context, builder *checkout.CreateOrderBuilder) (*entity.Order, error) {

	var seq NextId

	err := r.db.NewQuery(fmt.Sprintf("SELECT nextval('%s') as id", tableOrderSeqNextValID)).One(&seq)

	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	dInfo, err := json.Marshal(DeliveryInfo{
		City: City{
			Name: builder.Warehouse.GetCity().GetName(),
			Code: builder.Warehouse.GetCity().GetId(),
		},
		Address: Address{
			IsCustom: builder.Warehouse.GetAddress().IsCustom(),
			Address:  builder.Warehouse.GetAddress().GetName(),
			Code:     builder.Warehouse.GetAddress().GetId(),
		},
		Payment: PaymentExtra{
			Edrpou:   builder.PayInEdrpou,
			Company:  builder.PayInCompany,
			Email:    builder.PayInEmail,
			PartsPay: builder.PayPartsPay,
		},
	})

	if err != nil {
		return nil, err
	}

	dStatuses, err := json.Marshal([]DeliveryStatus{{
		Status:  entity.DeliveryStatusNew,
		Time:    now,
		Comment: "new order",
	}})

	if err != nil {
		return nil, err
	}

	pStatuses, err := json.Marshal([]PaymentStatus{{
		Status:  entity.PaymentStatusNew,
		Time:    now,
		Comment: "new order",
	}})

	if err != nil {
		return nil, err
	}

	_, err = r.db.Insert(tableNameOrder, dbx.Params{
		"id":             seq.Id,
		"customer_phone": builder.Customer.GetPhone(),
		"customer_name":  builder.Customer.GetName(),

		"delivery_method_id": builder.DeliveryMethod.GetID(),
		"delivery_status":    entity.DeliveryStatusNew,

		"payment_method_id": builder.PaymentMethod.GetID(),
		"payment_status":    entity.PaymentStatusNew,

		"cost":        builder.Cost,
		"comment":     builder.Comment,
		"do_not_call": builder.DoNotCall,

		"delivery_info":          string(dInfo),
		"delivery_statuses_json": string(dStatuses),
		"payment_statuses_json":  string(pStatuses),

		"created_at": now,
	}).Execute()

	if err != nil {
		return nil, err
	}

	for _, p := range builder.Products {
		pr := p.GetProduct()
		price := (&pr).Price.GetPriceByQuantity(p.GetQuantity())
		_, e := r.db.Insert(tableNameOrderItems, dbx.Params{
			"order_id":   seq.Id,
			"product_id": p.GetProduct().ID,
			"price":      price,
			"quantity":   p.GetQuantity(),
		}).Execute()

		if e != nil {
			log.Printf("order item insert err: %v", err)
		}
	}

	odsh := make([]*entity.OrderDeliveryStatusHistory, 1)
	odsh[0] = entity.NewOrderDeliveryStatusHistory(entity.DeliveryStatusNew, now, "new order")

	oDelivery := entity.NewOrderDelivery(entity.DeliveryStatusNew, builder.DeliveryMethod, builder.Warehouse, odsh)

	opsh := make([]*entity.OrderPaymentStatusHistory, 1)
	opsh[0] = entity.NewOrderPaymentStatusHistory(entity.DeliveryStatusNew, now, "new order")

	oPayment := entity.NewOrderPayment(
		entity.PaymentStatusNew,
		builder.PaymentMethod,
		opsh,
		builder.PayInEdrpou,
		builder.PayInCompany,
		builder.PayInEmail,
		builder.PayPartsPay,
	)

	order := entity.NewOrder(seq.Id, now, builder.Comment, builder.DoNotCall, builder.Cost, builder.Customer, oDelivery, oPayment, builder.Products)

	return order, nil
}

func NewPaymentRepository(db *dbx.DB) *PaymentRepository {

	return &PaymentRepository{db: db}
}

type PaymentRepository struct {
	db *dbx.DB
}

func (r PaymentRepository) Get(ctx context.Context, transactionId string) (*entity.Payment, error) {

	var (
		row              Payment
		created, updated time.Time
	)

	err := r.db.Select("*").
		From(tableNamePayments).
		Where(dbx.NewExp("transaction_id={:id}", dbx.Params{"id": transactionId})).
		One(&row)

	if err != nil {
		return nil, err
	}

	if row.Created.Valid == true {
		created, _ = time.Parse(time.RFC3339, row.Created.String)
	} else {
		created = time.Now()
	}

	if row.Updated.Valid == true {
		updated, _ = time.Parse(time.RFC3339, row.Updated.String)
	} else {
		updated = time.Now()
	}

	return entity.NewPayment(
		row.ID,
		row.TransactionID,
		row.OrderId,
		row.Provider,
		entity.NewPrice(row.Amount, 0, 0, nil),
		row.Status,
		created.Unix(),
		updated.Unix(),
	), nil
}

func (r PaymentRepository) NextId() (int, error) {

	var seq NextId

	err := r.db.NewQuery(fmt.Sprintf("SELECT nextval('%s') as id", tablePaymentSeqNextValID)).One(&seq)

	if err != nil {
		return 0, errors.New(fmt.Sprintf("[get payment next sequence][%v]", err))
	}

	return seq.Id, nil
}

func (r PaymentRepository) Save(ctx context.Context, p *entity.Payment) error {

	p.TouchUpdated()

	_, err := r.db.Update(tableNamePayments, dbx.Params{
		"order_id":       p.GetOrderId(),
		"amount":         p.GetPrice().GetInCent(),
		"transaction_id": p.GetTransactionId(),
		"provider":       p.GetProvider(),
		"status":         p.GetStatus(),
		"created_at":     p.GetCreatedTime(),
		"updated_at":     p.GetUpdatedTime(),
	}, dbx.NewExp("id={:id}", dbx.Params{"id": p.GetId()})).
		Execute()

	return err
}

func (r PaymentRepository) Create(ctx context.Context, p *entity.Payment) error {

	_, err := r.db.Insert(tableNamePayments, dbx.Params{
		"id":             p.GetId(),
		"order_id":       p.GetOrderId(),
		"amount":         p.GetPrice().GetInCent(),
		"transaction_id": p.GetTransactionId(),
		"provider":       "",
		"status":         p.GetStatus(),
		"created_at":     p.GetCreatedTime(),
		"updated_at":     p.GetUpdatedTime(),
	}).Execute()

	if err != nil {
		return errors.New(fmt.Sprintf("[create payment][%s]", err.Error()))
	}

	return nil
}

func tableWithAlias(tableName, alias string) string {
	return tableName + " " + alias
}
