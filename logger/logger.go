package logger

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func SetupFromConfig(cfg *Config) error {
	log.SetLevel(cfg.Level)

	switch cfg.Format {
	case FormatText:
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat:  time.RFC3339,
			FullTimestamp:    true,
			ForceQuote:       true,
			DisableSorting:   false,
			QuoteEmptyFields: true,
		})
	case FormatJSON:
		log.SetFormatter(&log.JSONFormatter{})
	default:
		return errors.Errorf("unexpected format: %s", cfg.Format)
	}
	return nil
}
