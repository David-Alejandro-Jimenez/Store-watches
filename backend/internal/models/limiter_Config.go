package models

type LimiterConfig struct {
	RequestPerSecond float64 `mapstructure:"request_rate"`
	Burst int `mapstructure:"burst"`
}