package notification

import (
	"github.com/wowucco/G3/internal/entity"
	"github.com/wowucco/G3/pkg/sms"
	"github.com/wowucco/G3/pkg/telegram"
)

const TelegramOrderChat  = "t_order_char"
const TelegramRecallChat  = "t_recall_char"

func NewNotificationService(smsChan chan sms.Message, telegramChan chan telegram.Message, telegramChats map[string]string) *Service {

	return &Service{smsChan, telegramChan, telegramChats}
}

type Service struct {
	smsChan       chan sms.Message
	telegramChan  chan telegram.Message
	telegramChats map[string]string
}

func (s *Service) OrderCreated(order *entity.Order) {

	s.smsChan <- sms.NewMsg([]string{"0123456789"}, "test message")
	s.telegramChan <- telegram.NewMsg(s.telegramChats[TelegramOrderChat], "Test chan", "")
}

func (s *Service) PaymentCreated(order *entity.Order, payment *entity.Payment) {

	s.smsChan <- sms.NewMsg([]string{"0123456789"}, "test message")
	s.telegramChan <- telegram.NewMsg(s.telegramChats[TelegramOrderChat], "Test chan", "")
}

func (s *Service) AcceptPayment(order *entity.Order, payment *entity.Payment) {

	s.smsChan <- sms.NewMsg([]string{"0123456789"}, "test message")
	s.telegramChan <- telegram.NewMsg(s.telegramChats[TelegramOrderChat], "Test chan", "")
}

func (s *Service) UpdatePayment(order *entity.Order, payment *entity.Payment) {

	s.smsChan <- sms.NewMsg([]string{"0123456789"}, "test message")
	s.telegramChan <- telegram.NewMsg(s.telegramChats[TelegramOrderChat], "Test chan", "")
}
