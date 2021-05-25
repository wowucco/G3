package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/wowucco/G3/internal/checkout"
	"github.com/wowucco/G3/internal/checkout/strategy"
	"github.com/wowucco/G3/internal/delivery"
	"github.com/wowucco/G3/internal/entity"
	"github.com/wowucco/G3/internal/product"
	"github.com/wowucco/G3/pkg/notification"
	"log"
)

func NewOrderUseCase(
	o checkout.IOrderRepository,
	p product.ReadRepository,
	d delivery.DeliveryReadRepository,
	pr checkout.IPaymentRepository,
	n *notification.Service,
	pc *strategy.PaymentContext,
) *OrderUserCase {

	return &OrderUserCase{orderRepository: o, productRepository: p, deliveryRepository: d, notify: n, paymentContext: pc, paymentRepository: pr}
}

type OrderUserCase struct {
	orderRepository    checkout.IOrderRepository
	productRepository  product.ReadRepository
	deliveryRepository delivery.DeliveryReadRepository
	paymentRepository  checkout.IPaymentRepository

	notify         *notification.Service
	paymentContext *strategy.PaymentContext
}

func (o OrderUserCase) Create(ctx context.Context, form checkout.CreateOrderForm) (*entity.Order, error) {

	oProducts, err := o.orderProducts(ctx, form)

	if err != nil {
		return nil, err
	}

	if len(oProducts) != len(form.GetOrder().GetOrderItems()) {
		e := fmt.Sprintf("[Create order][count not equal] come: %v; query: %v", len(form.GetOrder().GetOrderItems()), len(oProducts))
		log.Printf(e)
		return nil, errors.New(e)
	}

	dMethod, err := o.deliveryRepository.GetDeliveryMethodBySlug(form.GetDelivery().GetMethod())

	if err != nil {
		return nil, err
	}

	pMethod, err := o.deliveryRepository.GetPaymentMethodBySlug(form.GetPayment().GetMethod())

	if err != nil {
		return nil, err
	}

	builder := &checkout.CreateOrderBuilder{
		DeliveryMethod: dMethod,
		PaymentMethod:  pMethod,
		Products:       oProducts,
		Warehouse:      entity.NewOrderDeliveryWarehouse(form.GetDelivery().GetCity().GetCode(), form.GetDelivery().GetCity().GetName(), form.GetDelivery().GetAddress().GetCode(), form.GetDelivery().GetAddress().GetName(), form.GetDelivery().IsCustomAddress()),
		Customer:       entity.NewOrderCustomer(form.GetClient().GetFio(), form.GetClient().GetPhone()),
		Cost:           form.GetOrder().GetCost(),
		Comment:        form.GetComment(),
		DoNotCall:      form.GetDoNotCall(),
		PayInCompany:   form.GetPayment().GetPayInCompany(),
		PayInEdrpou:    form.GetPayment().GetPayInEdrpou(),
		PayInEmail:     form.GetPayment().GetPayInEmail(),
		PayPartsPay:    form.GetPayment().GetPayPartsPay(),
	}

	order, err := o.orderRepository.Create(ctx, builder)

	if err != nil {
		return nil, err
	}

	o.notify.OrderCreated(order)

	return order, nil
}

func (o OrderUserCase) InitPayment(ctx context.Context, form checkout.InitPaymentForm) (checkout.IInitPaymentResponse, error) {

	order, err := o.orderRepository.Get(ctx, form.GetOrderId())

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[Init payment][%v]", err))
	}

	if order.CanMakePayment() != true {
		return nil, errors.New(fmt.Sprintf("[Init payment][Can't init payment for order with payment status: %v]", order.GetPayment().GetStatus()))
	}

	id, err := o.paymentRepository.NextId()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[payment init][next sequense][%v]", err))
	}

	p := entity.CreateNewPaymentByOrder(id, order)

	err = o.paymentRepository.Create(ctx, p)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[payment init][create payment][%v]", err))
	}

	pMethod := order.GetPayment().GetMethod()

	r, err := o.paymentContext.GetInitPaymentStrategy(&pMethod).Init(ctx, order, p)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[payment init][provider init]%v", err))
	}

	p.SetProvider(r.GetProviderName())

	err = o.paymentRepository.Save(ctx, p)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[payment init][update error][%v]", err))
	}

	o.notify.PaymentCreated(order, p)

	return NewIniPaymentResponse(p, order, r.GetAction(), r.GetResource()), nil
}

func (o OrderUserCase) AcceptHoldenPayment(ctx context.Context, form checkout.IAcceptHoldenPaymentForm) error {

	p, err := o.paymentRepository.Get(ctx, form.GetTransactionId())

	if err != nil {
		return errors.New(fmt.Sprintf("[error][accept holden][payment not found][%s][%v]", form.GetTransactionId(), err))
	}

	order, err := o.orderRepository.Get(ctx, p.GetOrderId())

	if err != nil {
		return errors.New(fmt.Sprintf("[error][accept holden][order not found][%d][%v]", p.GetOrderId(), err))
	}

	s, err := o.paymentContext.GetAcceptHoldenPaymentStrategy(p.GetProvider())

	if err != nil {
		return errors.New(fmt.Sprintf("[error][accept holden][unresolved strategy][%s][%v]", p.GetTransactionId(), err))
	}

	r, err := s.Accept(ctx, order, p)

	if err != nil {
		return errors.New(fmt.Sprintf("[error][accept holden][unresolved strategy][%s]%v", p.GetTransactionId(), err))
	}

	p.UpdateStatus(r.GetStatus())
	err = o.paymentRepository.Save(ctx, p)

	if err != nil {
		return errors.New(fmt.Sprintf("[error][accept holden][payment save][%s][%v]", p.GetTransactionId(), err))
	}

	order.UpdatePaymentStatus(r.GetStatus(), r.GetDescription())
	err = o.orderRepository.Save(ctx, order)
	if err != nil {
		return errors.New(fmt.Sprintf("[error][accept holden][order save][%d][%v]", order.GetId(), err))
	}

	o.notify.AcceptPayment(order, p)

	return nil
}

func (o OrderUserCase) ProviderCallback(ctx context.Context, form checkout.IProviderCallbackPaymentForm) (checkout.IProviderCallbackPaymentResponse, error) {

	s, err := o.paymentContext.GetProviderCallbackPaymentStrategy(form.GetProvider())

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[error][provider callback][unresolved strategy][%v][%v]", form.GetParams(), err))
	}

	if s.IsValidSignature(form.GetParams()) == false {
		return nil, errors.New(fmt.Sprintf("[error][provider callback][invalid signature][%v][%v]", form.GetParams(), err))
	}

	transactionId := s.GetTransactionId(form.GetParams())

	payment, err := o.paymentRepository.Get(ctx, transactionId)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[error][provider callback][payment not found][%s][%v]", transactionId, err))
	}

	order, err := o.orderRepository.Get(ctx, payment.GetOrderId())

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[error][provider callback][order not found][%d][%v]", payment.GetOrderId(), err))
	}

	resp, err := s.ProcessingCallback(form.GetParams())

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[error][provider callback][processing error][%d][%v]", payment.GetOrderId(), err))
	}

	payment.UpdateStatus(resp.GetStatus())
	err = o.paymentRepository.Save(ctx, payment)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[error][provider callback][payment save][%s][%v]", payment.GetTransactionId(), err))
	}

	order.UpdatePaymentStatus(resp.GetStatus(), resp.GetDescription())
	err = o.orderRepository.Save(ctx, order)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[error][provider callback][order save][%d][%v]", order.GetId(), err))
	}

	o.notify.AcceptPayment(order, payment)

	return nil, nil
}

func (o OrderUserCase) orderProducts(ctx context.Context, form checkout.CreateOrderForm) ([]*entity.OrderProduct, error) {

	pIds := make([]int, len(form.GetOrder().GetOrderItems()))
	mTemp := make(map[int]checkout.OrderItemForm, len(form.GetOrder().GetOrderItems()))

	for k, v := range form.GetOrder().GetOrderItems() {
		pIds[k] = v.GetProductId()
		mTemp[v.GetProductId()] = v
	}

	products, err := o.productRepository.GetByIdsWithSequence(ctx, pIds)

	if err != nil {
		return nil, err
	}

	simple := make([]entity.SimpleProduct, len(products))

	for k, v := range products {
		simple[k] = *entity.NewSimpleProduct(v.ID, v.Name, v.Code, v.Exist, v.Status, v.Price)
	}

	oProducts := make([]*entity.OrderProduct, len(simple))

	for k, v := range simple {
		oProducts[k] = entity.NewOrderProduct(mTemp[v.ID].GetCount(), mTemp[v.ID].GetPrice(), v)
	}

	return oProducts, nil
}

func (o OrderUserCase) Callback(provider string, ) {

}

func NewIniPaymentResponse(payment *entity.Payment, order *entity.Order, action, resource string) checkout.IInitPaymentResponse {

	return &InitPaymentResponse{
		payment: payment,
		order: order,
		action:   action,
		resource: resource,
	}
}

type InitPaymentResponse struct {
	payment *entity.Payment
	order *entity.Order
	action   string
	resource string
	provider string
}

func (r *InitPaymentResponse) GetAction() string {
	return r.action
}
func (r *InitPaymentResponse) GetResource() string {
	return r.resource
}
func (r *InitPaymentResponse) GetPaymentTransactionID() string {
	return r.payment.GetTransactionId()
}
func (r *InitPaymentResponse) GetPaymentMethod() string {
	return r.order.GetPayment().GetMethod().GetSlug()
}
func (r *InitPaymentResponse) GetOrderId() int {
	return r.order.GetId()
}
func (r *InitPaymentResponse) GetDoNotCall() bool {
	return r.order.GetDoNotCall()
}