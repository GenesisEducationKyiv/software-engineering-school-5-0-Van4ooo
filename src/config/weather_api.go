package config

import "fmt"

type WeatherSettings interface {
	ConfigSource
	GetKey() string
	GetBaseURL() string
	GenUrl(city string) string
}

type WeatherAPI struct {
	Key     string `validate:"required"`
	BaseURL string `validate:"required"`
}

func (w *WeatherAPI) Load(provider EnvProvider) error {
	w.Key = provider.Get("WEATHER_API_KEY")
	w.BaseURL = provider.Get("WEATHER_API_BASE_URL")
	return nil
}

func (w *WeatherAPI) Validate() error {
	if err := validate.Struct(w); err != nil {
		return mapValidationErrorsToEnvVars(err, "WeatherAPI",
			map[string]string{
				"Key":     "WEATHER_API_KEY",
				"BaseURL": "WEATHER_API_BASE_URL",
			})
	}
	return nil
}

func (w *WeatherAPI) GetKey() string     { return w.Key }
func (w *WeatherAPI) GetBaseURL() string { return w.BaseURL }
func (w *WeatherAPI) GenUrl(city string) string {
	return fmt.Sprintf("%s/current.json?key=%s&q=%s", w.GetBaseURL(), w.GetKey(), city)
}

var _ WeatherSettings = (*WeatherAPI)(nil)
