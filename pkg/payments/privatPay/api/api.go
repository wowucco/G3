package api

import (
	"net/http"
)

const (
	StateCreated    = "CREATED"     //Платеж создан
	StateCanceled   = "CANCELED"    //Платеж отменен (клиентом)
	StateSuccess    = "SUCCESS"     //Платеж успешно совершен
	StateFail       = "FAIL"        //Ошибка при создании платежа
	StateClientWait = "CLIENT_WAIT" //Ожидание оплаты клиента
	StateOtpWaiting = "OTP_WAITING" //Подтверждения клиентом ОТП пароля
	StatePpCreation = "PP_CREATION" //создание контракта для платежа
	StateLocked     = "LOCKED"      //Платеж подтвержден клиентом и ожидает подтверждение магазином.
)

type Transport interface {
	Perform(*http.Request) (*http.Response, error)
}

func NewConfig(storeId, password, responseUrl, redirectUrl string) Config {
	return Config{
		storeId:     storeId,
		passport:    password,
		responseUrl: responseUrl,
		redirectUrl: redirectUrl,
	}
}

type Config struct {
	storeId     string
	passport    string
	responseUrl string
	redirectUrl string
}

func New(t Transport, cfg Config) *API {
	return &API{
		Pay: &Pay{
			Hold:   newPaymentHoldFunc(t, cfg),
			Accept: newAcceptHoldenFunc(t, cfg),
		},
		Sign: &Sign{
			CallbackCheck: newSignCallbackCheckFunc(cfg),
		},
	}
}

type API struct {
	Pay  *Pay
	Sign *Sign
}