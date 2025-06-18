package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/repositories"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services"
)

type SubscriptionHandler struct {
	repo        *repositories.SubscriptionRepository
	emailSender services.EmailSender
}

func NewSubscriptionHandler(
	cfg *config.AppConfig,
	db *gorm.DB,
) *SubscriptionHandler {
	repo := repositories.NewSubscriptionRepository(db)
	sender := services.NewSMTPEmailSender(
		cfg.SMTP.Host,
		cfg.SMTP.Addr,
		cfg.SMTP.Name,
		cfg.SMTP.Pass,
	)
	return &SubscriptionHandler{repo: repo, emailSender: sender}
}

type subscriptionRequest struct {
	Email string `form:"email" json:"email" binding:"required,email"`
	City  string `form:"city" json:"city" binding:"required"`
	// nolint: lll
	Frequency string `form:"frequency" json:"frequency" binding:"required,oneof=hourly daily"`
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	req, err := h.parseRequest(c)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.createSubscription(req)
	if err != nil {
		h.respondError(c, http.StatusConflict,
			"email already subscribed")
		return
	}

	if err := h.sendConfirmation(
		req.Email, token, c,
	); err != nil {
		h.respondError(c, http.StatusInternalServerError,
			"failed to send email")
		return
	}

	h.respondJSON(c, http.StatusOK,
		gin.H{"message": "Confirmation email sent"})
}

func (h *SubscriptionHandler) parseRequest(
	c *gin.Context,
) (*subscriptionRequest, error) {
	var req subscriptionRequest
	if strings.Contains(
		c.GetHeader("Content-Type"),
		"application/json",
	) {
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, err
		}
	} else {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return nil,
				fmt.Errorf("cannot read body: %v", err)
		}
		vals, err := url.ParseQuery(string(body))
		if err != nil {
			return nil,
				fmt.Errorf("invalid form data: %v", err)
		}
		req.Email = vals.Get("email")
		req.City = vals.Get("city")
		req.Frequency = vals.Get("frequency")
	}
	return &req, nil
}

func (h *SubscriptionHandler) createSubscription(
	req *subscriptionRequest,
) (string, error) {
	token := uuid.NewString()
	sub := &models.Subscription{
		Email:     req.Email,
		City:      req.City,
		Frequency: req.Frequency,
		Token:     token,
	}
	if err := h.repo.Create(sub); err != nil {
		return "", err
	}
	return token, nil
}

func (h *SubscriptionHandler) sendConfirmation(
	email, token string,
	c *gin.Context,
) error {
	base := h.getBaseURL(c)
	link := fmt.Sprintf("%s/api/confirm/%s", base, token)
	subj := "Confirm your subscription"
	body := fmt.Sprintf(
		"Click to confirm your subscription: %s", link,
	)
	return h.emailSender.Send(email, subj, body)
}

func (h *SubscriptionHandler) respondError(
	c *gin.Context,
	code int,
	msg string,
) {
	c.JSON(code, gin.H{"error": msg})
}

func (h *SubscriptionHandler) respondJSON(
	c *gin.Context,
	code int,
	payload interface{},
) {
	c.JSON(code, payload)
}

func (h *SubscriptionHandler) Confirm(
	c *gin.Context,
) {
	token := c.Param("token")
	if err := h.repo.Confirm(token); err != nil {
		h.respondError(c, http.StatusNotFound,
			"token not found")
		return
	}
	h.respondJSON(c, http.StatusOK,
		gin.H{"message": "Subscription confirmed"})
}

func (h *SubscriptionHandler) Unsubscribe(
	c *gin.Context,
) {
	token := c.Param("token")
	if err := h.repo.Delete(token); err != nil {
		h.respondError(c, http.StatusNotFound,
			"token not found")
		return
	}
	h.respondJSON(c, http.StatusOK,
		gin.H{"message": "Unsubscribed successfully"})
}

func (h *SubscriptionHandler) RenderSubscribePage(
	c *gin.Context,
) {
	c.File("static/subscribe.html")
}

func (h *SubscriptionHandler) getBaseURL(
	c *gin.Context,
) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}
