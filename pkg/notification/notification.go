package notification

import (
	"fmt"
	"github.com/wowucco/G3/internal/entity"
	"github.com/wowucco/G3/pkg/sms"
	"github.com/wowucco/G3/pkg/telegram"
)

const TelegramOrderChat = "t_order_char"
const TelegramRecallChat = "t_recall_char"

func NewNotificationService(smsChan chan sms.Message, telegramChan chan telegram.Message, telegramChats map[string]string, cardNumber, boOrderLinkMask, webProductLinkMask string) *Service {

	return &Service{
		smsChan,
		telegramChan,
		telegramChats,
		cardNumber,
		boOrderLinkMask,
		webProductLinkMask,
	}
}

type Service struct {
	smsChan            chan sms.Message
	telegramChan       chan telegram.Message
	telegramChats      map[string]string
	cartNumber         string
	boOrderLinkMask    string
	webProductLinkMask string
}

// Send sms to client "Vashe zamovlenya # %d vzyato na opracyvaya."
// Send messege to telegram "*New order created:*..."
func (s *Service) OrderCreated(order *entity.Order) {
	var (
		smsMessage, telegramMessage string
	)

	smsMessage = fmt.Sprintf("Vashe zamovlenya # %d vzyato na opracyvaya.", order.GetId())

	if order.NeedToCall() == true {
		smsMessage = fmt.Sprintf("%s Mi zatelefonyemo naiblishim chasom", smsMessage)
	}

	s.smsSend([]string{order.GetCustomer().GetPhone()}, smsMessage)

	telegramMessage = fmt.Sprintf("*New order created:* [%d](%s)", order.GetId(), s.makeLinkToOrder(order))

	if order.NeedToCall() == false {
		telegramMessage += "_DO NOT CALL_"
	}

	s.telegramMessageCustomerBlock(&telegramMessage, order)
	s.telegramMessageSummaryBlock(&telegramMessage, order)
	s.telegramMessageItemsBlock(&telegramMessage, order)

	if order.GetComment() != "" {
		telegramMessage += fmt.Sprintf("\n_Comment_\n%s", order.GetComment())
	}

	s.telegramSend(s.telegramChats[TelegramOrderChat], telegramMessage)
}

// send sms to client with card number if payment method to_card
func (s *Service) PaymentCreated(order *entity.Order, payment *entity.Payment) {

	if order.GetPayment().HasToCardPayment() == true {
		totalCost := order.GetPrice()
		smsMessage := fmt.Sprintf("Dlya oplati zamovlenya # %d zdiysnit perekaz na kartu %s. Summa do splatu %s UAH", order.GetId(), s.cartNumber, (&totalCost).CentToCurrency())

		s.smsSend([]string{order.GetCustomer().GetPhone()}, smsMessage)
	}
}

func (s *Service) PaymentStatusUpdated(order *entity.Order, payment *entity.Payment) {

	switch order.GetPayment().GetStatus() {
	case entity.PaymentStatusWaitingConfirmation:
		message := fmt.Sprintf("*Customer paid order and waiting confirmation: * [%d](%s)", order.GetId(), s.makeLinkToOrder(order))
		s.telegramMessageCustomerBlock(&message, order)
		s.telegramMessageOrderInfoBlock(&message, order)
		s.telegramSend(s.telegramChats[TelegramOrderChat], message)
	case entity.PaymentStatusDone:
		tmessage := fmt.Sprintf("*Payment accepted and successful for order: * [%d](%s)", order.GetId(), s.makeLinkToOrder(order))
		s.telegramMessageCustomerBlock(&tmessage, order)
		s.telegramMessageOrderInfoBlock(&tmessage, order)
		s.telegramSend(s.telegramChats[TelegramOrderChat], tmessage)

		s.smsSend([]string{order.GetCustomer().GetPhone()}, fmt.Sprintf("Vashe zamovlenya # %d oplacheno", order.GetId()))
	case entity.PaymentStatusFailed:
		tmessage := fmt.Sprintf("*Payment was failed: * [%d](%s)", order.GetId(), s.makeLinkToOrder(order))
		s.telegramMessageCustomerBlock(&tmessage, order)
		s.telegramMessageOrderInfoBlock(&tmessage, order)
		s.telegramSend(s.telegramChats[TelegramOrderChat], tmessage)

		smessage := fmt.Sprintf("Vashe zamovlennya # %d ne bulo oplacheno. My zatelefonuemo Vam najblizhchim chasom", order.GetId())
		s.smsSend([]string{order.GetCustomer().GetPhone()}, smessage)
	}
}

func (s *Service) Recall(phone, message string) {

	if message == "" {
		message = "message is empty"
	}
	
	msg := fmt.Sprintf("*Phone*\n%s\n*Message*\n%s\n", phone, message)
	s.telegramSend(s.telegramChats[TelegramRecallChat], msg)
}

func (s *Service) BuyOnClick(phone string, product entity.Product) {

	telegramMessage := "*Buy on click request*"
	telegramMessage += fmt.Sprintf("\n_Customer_\n%s\n", phone)
	telegramMessage += "_Item_\n"
	telegramMessage += fmt.Sprintf("[%s](%s)\n", product.Name, s.makeLinkToProduct(product.ID))

	cost := product.Price
	telegramMessage += fmt.Sprintf("_cost:_ _%s_\n", (&cost).CentToCurrency())

	s.telegramSend(s.telegramChats[TelegramRecallChat], telegramMessage)
}

func (s *Service) makeLinkToOrder(o *entity.Order) string {
	return fmt.Sprintf(s.boOrderLinkMask, o.GetId())
}

func (s *Service) makeLinkToProduct(id int) string {
	return fmt.Sprintf(s.webProductLinkMask,  id)
}

func (s *Service) telegramSend(chat, message string) {
	s.telegramChan <- telegram.NewMsg(chat, message, "Markdown")
}

func (s *Service) smsSend(numbers []string, message string) {
	s.smsChan <- sms.NewMsg(numbers, message)
}

func (s *Service) telegramMessageCustomerBlock(message *string, order *entity.Order) {
	*message += fmt.Sprintf("\n_Customer_\n%s\n%s\n", order.GetCustomer().GetPhone(), order.GetCustomer().GetName())
}
func (s *Service) telegramMessageSummaryBlock(message *string, order *entity.Order) {
	totalCost := order.GetPrice()

	*message += "\n_Summary_\n"
	*message += fmt.Sprintf("total items: _%d_\n", len(order.GetItems()))
	*message += fmt.Sprintf("total cost: _%s_\n", (&totalCost).CentToCurrency())
	*message += fmt.Sprintf("delivery: _%s_\n", order.GetDelivery().GetMethod().GetName())
	*message += fmt.Sprintf("payment: _%s_\n", order.GetPayment().GetMethod().GetName())
}
func (s *Service) telegramMessageItemsBlock(message *string, order *entity.Order) {
	*message += "\n_Items_\n"

	for _, v := range order.GetItems() {
		cost := v.GetPrice()
		*message += fmt.Sprintf("[%s](%s)", v.GetProduct().Name, s.makeLinkToProduct(v.GetProduct().ID))
		*message += fmt.Sprintf("_%d_\tcost: _%s_\n", v.GetQuantity(), (&cost).CentToCurrency())
	}
}
func (s *Service) telegramMessageOrderInfoBlock(message *string, order *entity.Order) {
	totalCost := order.GetPrice()

	*message += "\nOrder info\n"
	*message += fmt.Sprintf("total items: _%d_\n", len(order.GetItems()))
	*message += fmt.Sprintf("total cost: _%s_\n", (&totalCost).CentToCurrency())
	*message += fmt.Sprintf("delivery: _%s_\n", order.GetDelivery().GetMethod().GetName())
	*message += fmt.Sprintf("payment: _%s_\n", order.GetPayment().GetMethod().GetName())
	*message += fmt.Sprintf("status: _%s_\n", order.GetPayment().GetStatusLabel())
}