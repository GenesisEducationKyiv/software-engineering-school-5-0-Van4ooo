package subscription

import (
	"github.com/google/uuid"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
)

type Storage interface {
	Save(sub *models.Subscription) error
	MarkConfirmed(token string) error
	Remove(token string) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (svc *Service) RegisterNew(
	req *models.SubscriptionRequest) (*models.Subscription, error) {
	token := uuid.NewString()
	sub := req.ToSubscription(token)

	if err := svc.storage.Save(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (svc *Service) ConfirmByToken(token string) error {
	return svc.storage.MarkConfirmed(token)
}

func (svc *Service) CancelByToken(token string) error {
	return svc.storage.Remove(token)
}
