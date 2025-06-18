package repositories

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
)

type SubscriptionStorage interface {
	Save(sub *models.Subscription) error
	MarkConfirmed(token string) error
	Remove(token string) error
}
