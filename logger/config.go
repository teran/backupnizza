package logger

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	log "github.com/sirupsen/logrus"
)

type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

type Config struct {
	Level  log.Level `json:"level"`
	Format Format    `json:"format"`
}

func NewConfigFromJSON(in json.RawMessage) (*Config, error) {
	cfg := &Config{
		Level:  log.WarnLevel,
		Format: "text",
	}

	if err := json.Unmarshal(in, cfg); err != nil {
		return nil, err
	}
	return cfg, cfg.Validate()
}

func (cfg *Config) Validate() error {
	return validation.ValidateStruct(cfg,
		validation.Field(&cfg.Level, validation.Required),
		validation.Field(&cfg.Format, validation.Required, validation.In(FormatText, FormatJSON)),
	)
}
