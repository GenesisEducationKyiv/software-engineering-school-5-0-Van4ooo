package config

import (
	"errors"
	"testing"
)

type mockEnv struct {
	vars map[string]string
}

func (m mockEnv) Get(key string) string {
	return m.vars[key]
}

// nolint: dupl
func TestPostgres_LoadValidate(t *testing.T) {
	p := &Postgres{}
	env := mockEnv{vars: map[string]string{}}
	if err := p.Load(env); err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}
	if err := p.Validate(); err == nil {
		t.Errorf("expected Validate error when DATABASE_URL missing, got nil")
	}

	env = mockEnv{vars: map[string]string{"DATABASE_URL": "https://example.com/db"}}
	p = &Postgres{}
	if err := p.Load(env); err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}
	if p.URL != "https://example.com/db" {
		t.Errorf("Load did not set URL correctly: got %q", p.URL)
	}
	if err := p.Validate(); err != nil {
		t.Errorf("unexpected Validate error: %v", err)
	}
}

func TestSMTP_LoadValidate(t *testing.T) {
	s := &SMTP{}
	env := mockEnv{vars: map[string]string{}}
	if err := s.Load(env); err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}
	if err := s.Validate(); err == nil {
		t.Errorf("expected Validate error when SMTP vars missing, got nil")
	}

	env = mockEnv{vars: map[string]string{
		"SMTP_HOST": "smtp.example.com",
		"SMTP_ADDR": "noreply@example.com",
		"SMTP_NAME": "AppName",
		"SMTP_PASS": "secret",
	}}
	s = &SMTP{}
	if err := s.Load(env); err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}
	// nolint: lll
	if s.Host != "smtp.example.com" || s.Addr != "noreply@example.com" || s.Name != "AppName" || s.Pass != "secret" {
		t.Errorf("Load did not set SMTP fields correctly: %v", s)
	}
	if err := s.Validate(); err != nil {
		t.Errorf("unexpected Validate error: %v", err)
	}
}

// nolint: dupl
func TestWeatherAPI_LoadValidate(t *testing.T) {
	w := &WeatherAPI{}
	env := mockEnv{vars: map[string]string{}}
	if err := w.Load(env); err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}
	if err := w.Validate(); err == nil {
		t.Errorf("expected Validate error when WEATHER_API_KEY missing, got nil")
	}

	env = mockEnv{vars: map[string]string{"WEATHER_API_KEY": "apikey123"}}
	w = &WeatherAPI{}
	if err := w.Load(env); err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}
	if w.Key != "apikey123" {
		t.Errorf("Load did not set Key correctly: got %q", w.Key)
	}
	if err := w.Validate(); err != nil {
		t.Errorf("unexpected Validate error: %v", err)
	}
}

func TestSetupConfig_Success(t *testing.T) {
	vars := map[string]string{
		"DATABASE_URL":    "https://db.example.com",
		"SMTP_HOST":       "smtp.example.com",
		"SMTP_ADDR":       "noreply@example.com",
		"SMTP_NAME":       "App",
		"SMTP_PASS":       "pass",
		"WEATHER_API_KEY": "weatherkey",
	}
	env := mockEnv{vars: vars}
	cfg, err := SetupConfig(env)
	if err != nil {
		t.Fatalf("expected SetupConfig success, got error: %v", err)
	}

	if cfg.DB.URL != vars["DATABASE_URL"] {
		t.Errorf("DB.URL = %q; want %q", cfg.DB.URL, vars["DATABASE_URL"])
	}
	if cfg.SMTP.Host != vars["SMTP_HOST"] {
		t.Errorf("SMTP.Host = %q; want %q", cfg.SMTP.Host, vars["SMTP_HOST"])
	}
	if cfg.WeatherAPI.Key != vars["WEATHER_API_KEY"] {
		t.Errorf("WeatherAPI.Key = %q; want %q", cfg.WeatherAPI.Key, vars["WEATHER_API_KEY"])
	}
}

func TestSetupConfig_Failure(t *testing.T) {
	vars := map[string]string{
		"DATABASE_URL": "https://db.example.com",
		"SMTP_HOST":    "smtp.example.com",
		"SMTP_ADDR":    "noreply@example.com",
		"SMTP_NAME":    "App",
		"SMTP_PASS":    "pass",
	}
	env := mockEnv{vars: vars}
	_, err := SetupConfig(env)
	if err == nil {
		t.Fatalf("expected SetupConfig error when WEATHER_API_KEY missing, got nil")
	}

	if !errors.Is(err, errors.New("")) && !contains(err.Error(), "WEATHER_API_KEY") {
		t.Errorf("error = %v; want mention of WEATHER_API_KEY", err)
	}
}

// nolint: lll
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (func() bool { return len(s) == len(substr) || s[:len(substr)] == substr || contains(s[1:], substr) })()
}
