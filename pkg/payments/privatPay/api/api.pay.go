package api

import "fmt"

const MerchantTypeII = "II"
const MerchantTypePP = "PP"
const paymentRedirectUrl = "https://payparts2.privatbank.ua/ipp/v2/payment?token="

func NewProduct(name string, count int, price float64) Product {
	return Product{name, count, price}
}

func PaymentRedirectUrl(token string) string {
	return fmt.Sprintf("%s%s", paymentRedirectUrl, token)
}

type Product struct {
	name  string
	count int
	price float64
}

type Pay struct {
	Hold   PaymentHold
	Accept PaymentAcceptHolden
}