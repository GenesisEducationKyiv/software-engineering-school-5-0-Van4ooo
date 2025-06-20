package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services/email"
)

type EmailSender interface {
	Send(template email.Template) error
}

type SubscriptionService interface {
	RegisterNew(sub *models.SubscriptionRequest) (*models.Subscription, error)
	ConfirmByToken(token string) error
	CancelByToken(token string) error
}

type SubscriptionHandler struct {
	service     SubscriptionService
	emailSender EmailSender
}

func NewSubscriptionHandler(
	service SubscriptionService,
	sender EmailSender) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, emailSender: sender}
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	req := h.parseAndValidate(c)
	if req == nil {
		return
	}

	sub := h.registerSubscription(c, req)
	if sub == nil {
		return
	}

	mail := h.prepareConfirmationMail(c, sub)
	h.sendConfirmation(c, mail)
}

func (h *SubscriptionHandler) parseAndValidate(
	c *gin.Context) *models.SubscriptionRequest {
	req, err := h.parseRequest(c)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err.Error())
		return nil
	}
	return req
}

func (h *SubscriptionHandler) registerSubscription(c *gin.Context,
	req *models.SubscriptionRequest) *models.Subscription {
	sub, err := h.service.RegisterNew(req)
	if err != nil {
		h.respondError(c, http.StatusConflict, "email already subscribed")
		return nil
	}
	return sub
}

func (h *SubscriptionHandler) sendConfirmation(c *gin.Context, tmpl email.Template) {
	if err := h.emailSender.Send(tmpl); err != nil {
		h.respondError(c, http.StatusInternalServerError, "failed to send confirmation email")
		return
	}
	h.respondJSON(c, http.StatusOK, gin.H{"message": "Confirmation email sent"})
}

func (h *SubscriptionHandler) prepareConfirmationMail(c *gin.Context,
	sub *models.Subscription) email.Template {
	link := email.GenerateConfirmationLink(h.getBaseURL(c), sub.Token)
	return email.NewConfirmationMail(sub.Email, link)
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
