package services

import (
	"github.com/google/uuid"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/repositories"
)

type SubscriptionService interface {
	RegisterNew(sub *models.SubscriptionRequest) (string, error)
	ConfirmByToken(token string) error
	CancelByToken(token string) error
}

type SubscriptionServiceImpl struct {
	storage repositories.SubscriptionStorage
}

func NewSubscriptionService(
	storage repositories.SubscriptionStorage) *SubscriptionServiceImpl {
	return &SubscriptionServiceImpl{storage: storage}
}

func (svc *SubscriptionServiceImpl) RegisterNew(
	req *models.SubscriptionRequest) (string, error) {
	token := uuid.NewString()
	sub := &models.Subscription{
		Email:     req.Email,
		City:      req.City,
		Frequency: req.Frequency,
		Token:     token,
	}
	if err := svc.storage.Save(sub); err != nil {
		return "", err
	}
	return token, nil
}

func (svc *SubscriptionServiceImpl) ConfirmByToken(token string) error {
	return svc.storage.MarkConfirmed(token)
}

func (svc *SubscriptionServiceImpl) CancelByToken(token string) error {
	return svc.storage.Remove(token)
}
