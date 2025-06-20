package services

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/repositories"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services/email"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services/weather"
)

type JobScheduler interface {
	AddFunc(spec string, cmd func()) (cron.EntryID, error)
	Start()
	GetSchedules() []SchedulerConfig
}

type CronScheduler struct {
	cron      *cron.Cron
	schedules []SchedulerConfig
}

func NewCronScheduler(loc *time.Location, schedules []SchedulerConfig) *CronScheduler {
	return &CronScheduler{
		cron:      cron.New(cron.WithLocation(loc)),
		schedules: schedules,
	}
}

func (s *CronScheduler) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return s.cron.AddFunc(spec, cmd)
}

func (s *CronScheduler) Start() {
	s.cron.Start()
}

func (s *CronScheduler) GetSchedules() []SchedulerConfig {
	return s.schedules
}

func RunScheduler(cfg config.AppSettings, db *gorm.DB) {
	schedules := []SchedulerConfig{
		{Spec: "0 8 * * *", Frequency: "daily"},
		{Spec: "0 * * * *", Frequency: "hourly"},
	}
	scheduler := NewCronScheduler(time.UTC, schedules)

	store := repositories.NewSubscriptionStore(db)
	weatherService := weather.NewService(cfg.GetWeatherAPI())
	emailSender := email.NewSender(cfg.GetSMTP())

	svc := NewSchedulerService(scheduler, store, weatherService, emailSender)

	go svc.Start()
	log.Println("Scheduler started")
}

type SchedulerConfig struct {
	Spec      string
	Frequency string
}

type SubscriptionsFetcher interface {
	FetchByFrequency(freq string) ([]models.Subscription, error)
}

type SchedulerService struct {
	scheduler      JobScheduler
	store          SubscriptionsFetcher
	weatherService WeatherService
	emailSender    EmailSender
}

type EmailSender interface {
	Send(template email.Template) error
}

type WeatherService interface {
	GetByCity(city string) (*models.Weather, error)
}

func NewSchedulerService(
	scheduler JobScheduler,
	store SubscriptionsFetcher,
	service WeatherService,
	sender EmailSender,
) *SchedulerService {
	return &SchedulerService{
		scheduler:      scheduler,
		store:          store,
		weatherService: service,
		emailSender:    sender,
	}
}

func (s *SchedulerService) Start() {
	for _, cfg := range s.scheduler.GetSchedules() {
		if _, err := s.scheduler.AddFunc(cfg.Spec, func(freq string) func() {
			return func() { s.sendWeatherUpdates(freq) }
		}(cfg.Frequency)); err != nil {
			log.Fatalf("Error scheduling %s updates (%s): %v", cfg.Frequency, cfg.Spec, err)
		}
	}
	s.scheduler.Start()
	select {}
}

func (s *SchedulerService) sendWeatherUpdates(freq string) {
	subs, err := s.store.FetchByFrequency(freq)
	if err != nil {
		log.Printf("Error fetching subscriptions: %v", err)
		return
	}

	for _, sub := range subs {
		weather, err := s.weatherService.GetByCity(sub.City)
		if err != nil {
			log.Printf("Error fetching weather for %s: %v", sub.City, err)
			continue
		}

		body := fmt.Sprintf(
			"Current weather in %s: %.1f°C, %s (Humidity %.0f%%)",
			sub.City,
			weather.Temperature,
			weather.Description,
			weather.Humidity,
		)

		tmpl := email.NewSimpleMail(sub.Email, "Weather Update", body)
		if err := s.emailSender.Send(tmpl); err != nil {
			log.Printf("Error sending email to %s: %v", sub.Email, err)
		}
	}
}
