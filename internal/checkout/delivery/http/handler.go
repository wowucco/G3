package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wowucco/G3/internal/checkout"
	"io/ioutil"
	"log"
	"net/http"
)

func NewHandler(ordUC checkout.IOrderUseCase) *Handler {

	return &Handler{
		orderManage: ordUC,
	}
}

type Handler struct {
	orderManage checkout.IOrderUseCase
}

func (h *Handler) create(c *gin.Context) {

	b, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		log.Printf("[Checkout create request][read body][%v]", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var form CreateOrder

	if err := json.Unmarshal(b, &form); err != nil {
		log.Printf("[Checkout create request][decode body][%v]", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = form.Validate()

	if err != nil {
		log.Printf("[Checkout create request][validate][%v]", err)
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	order, err := h.orderManage.Create(c, form)

	if err != nil {
		log.Printf("[Checkout create request][create][%v]", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, NewCreateOrderResponse(order))
}

func (h *Handler) initPayment(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		log.Printf("[error][Init payment request][read body][%v]", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var form InitPaymentForm

	if err := json.Unmarshal(b, &form); err != nil {
		log.Printf("[error][Init payment request][decode body][%v]", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = form.Validate()

	if err != nil {
		log.Printf("[error][Init payment request][validate][%v]", err)
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	resp, err := h.orderManage.InitPayment(c, form)

	if err != nil {
		log.Printf("[error][Init payment request][init][%v]", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"action":         resp.GetAction(),
		"resource":       resp.GetResource(),
		"payment_id":     resp.GetPaymentTransactionID(),
		"payment_method": resp.GetPaymentMethod(),
		"order_id":       resp.GetOrderId(),
		"do_not_call":    resp.GetDoNotCall(),
	})
}

func (h *Handler) acceptHolden(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		log.Printf("[error][accept holden payment request][read body][%v]", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var form AccentPaymentForm

	if err := json.Unmarshal(b, &form); err != nil {
		log.Printf("[error][accept holden payment request][decode body][%v]", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = form.Validate()

	if err != nil {
		log.Printf("[error][accept holden payment request][validate][%v]", err)
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	err = h.orderManage.AcceptHoldenPayment(c, form)

	if err != nil {
		log.Printf("[error][accept holden payment request]%v", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) orderInfo(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		log.Printf("[error][order info request][read body][%v]", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var form OrderIdForm

	if err := json.Unmarshal(b, &form); err != nil {
		log.Printf("[error][order info request][decode body][%v]", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = form.Validate()

	if err != nil {
		log.Printf("[error][order info request][validate][%v]", err)
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	order, err := h.orderManage.OrderInfo(c, form)

	if err != nil {
		log.Printf("[error][order info request][get order][%v]", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, NewOrderInfoResponse(order))
}

func (h *Handler) callback(c *gin.Context) {

	provider := c.Param("provider")

	form := ProviderCallbackPaymentForm{
		Provider: provider,
	}
	_, err := h.orderManage.ProviderCallback(c, form)

	if err != nil {
		log.Printf("[error][provider callback][handle][%v]", err)
		c.JSON(http.StatusFailedDependency, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
