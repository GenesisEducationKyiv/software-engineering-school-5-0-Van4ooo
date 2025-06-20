package config

import "os"

type EnvProvider interface {
	Get(key string) string
}

type OSProvider struct{}

func (p OSProvider) Get(key string) string {
	return os.Getenv(key)
}
