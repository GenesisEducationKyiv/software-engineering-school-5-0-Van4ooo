package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"

	"github.com/gin-gonic/gin"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services"
)

type SubscriptionHandler struct {
	service     services.SubscriptionService
	emailSender services.EmailSender
}

func NewSubscriptionHandler(
	service services.SubscriptionService,
	sender services.EmailSender) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, emailSender: sender}
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	req, err := h.parseRequest(c)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.service.RegisterNew(req)
	if err != nil {
		h.respondError(c, http.StatusConflict, "email already subscribed")
		return
	}
	if err := h.emailSender.SendConfirmation(req.Email, h.getBaseURL(c), token); err != nil {
		h.respondError(c, http.StatusInternalServerError, "failed to send email")
		return
	}
	h.respondJSON(c, http.StatusOK, gin.H{"message": "Confirmation email sent"})
}

func (h *SubscriptionHandler) Confirm(c *gin.Context) {
	token := c.Param("token")
	if err := h.service.ConfirmByToken(token); err != nil {
		h.respondError(c, http.StatusNotFound, "token not found")
		return
	}
	h.respondJSON(c, http.StatusOK, gin.H{"message": "Subscription confirmed"})
}

func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	token := c.Param("token")
	if err := h.service.CancelByToken(token); err != nil {
		h.respondError(c, http.StatusNotFound, "token not found")
		return
	}
	h.respondJSON(c, http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

func (h *SubscriptionHandler) RenderSubscribePage(c *gin.Context) {
	c.File("static/subscribe.html")
}

func (h *SubscriptionHandler) parseRequest(
	c *gin.Context,
) (*models.SubscriptionRequest, error) {
	var req models.SubscriptionRequest
	if strings.Contains(c.GetHeader("Content-Type"), "application/json") {
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, err
		}
	} else {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read body: %v", err)
		}
		vals, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, fmt.Errorf("invalid form data: %v", err)
		}
		req.Email = vals.Get("email")
		req.City = vals.Get("city")
		req.Frequency = vals.Get("frequency")
	}
	return &req, nil
}

func (h *SubscriptionHandler) respondError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}

func (h *SubscriptionHandler) respondJSON(
	c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

func (h *SubscriptionHandler) getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}
