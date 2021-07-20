package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wowucco/G3/internal/contact"
	"io/ioutil"
	"log"
	"net/http"
)

func NewHandler(contactManage contact.IContactUseCase) *Handler {

	return &Handler{contactManage}
}

type Handler struct {
	contactManage contact.IContactUseCase
}

func (h *Handler) recall(c *gin.Context) {

	b, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		log.Printf("[error][contact recall request][read body][%v]", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var form RecallForm

	if err := json.Unmarshal(b, &form); err != nil {
		log.Printf("[error][contact recall request][decode body][%v]", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = form.Validate()

	if err != nil {
		log.Printf("[error][contact recall request][validate][%v]", err)
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	err = h.contactManage.Recall(c, form)

	if err != nil {
		log.Printf("[error][contact recall request][recall][%v]", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
