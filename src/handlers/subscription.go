package handlers

import (
	"fmt"
	"github.com/Van4ooo/genesis_case_task/src/db"
	"github.com/Van4ooo/genesis_case_task/src/models"
	"github.com/Van4ooo/genesis_case_task/src/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type subscriptionRequest struct {
	Email     string `form:"email" json:"email" binding:"required,email"`
	City      string `form:"city"  json:"city"  binding:"required"`
	Frequency string `form:"frequency" json:"frequency" binding:"required,oneof=hourly daily"`
}

func Subscribe(c *gin.Context) {
	req, err := parseSubscriptionRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := saveSubscription(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already subscribed"})
		return
	}

	if err := sendConfirmation(c, req.Email, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Confirmation email sent"})
}

func Confirm(c *gin.Context) {
	token := c.Param("token")
	if err := setSubscriptionConfirmed(token); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed"})
}

func Unsubscribe(c *gin.Context) {
	token := c.Param("token")
	if err := deleteSubscription(token); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

func parseSubscriptionRequest(c *gin.Context) (*subscriptionRequest, error) {
	var req subscriptionRequest
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") {
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
	if err := validateRequest(req); err != nil {
		return nil, err
	}
	return &req, nil
}

func validateRequest(req subscriptionRequest) error {
	if req.Email == "" || req.City == "" {
		return fmt.Errorf("email and city are required")
	}
	if req.Frequency != "hourly" && req.Frequency != "daily" {
		return fmt.Errorf("frequency must be 'hourly' or 'daily'")
	}
	return nil
}

func saveSubscription(req *subscriptionRequest) (string, error) {
	token := uuid.NewString()
	sub := &models.Subscription{
		Email:     req.Email,
		City:      req.City,
		Frequency: req.Frequency,
		Token:     token,
	}
	if err := db.DB.Create(sub).Error; err != nil {
		return "", err
	}
	return token, nil
}

func sendConfirmation(c *gin.Context, email, token string) error {
	baseURL := getBaseURL(c)
	return services.SendConfirmationEmail(email, baseURL, token)
}

func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

func setSubscriptionConfirmed(token string) error {
	var sub models.Subscription
	if err := db.DB.Where("token = ?", token).First(&sub).Error; err != nil {
		return err
	}
	sub.Confirmed = true
	return db.DB.Save(&sub).Error
}

func deleteSubscription(token string) error {
	return db.DB.Where("token = ?", token).Delete(&models.Subscription{}).Error
}

func RenderSubscribePage(c *gin.Context) {
	c.File("static/subscribe.html")
}
