package config

import (
	"encoding/json"

	time "github.com/teran/go-time"
)

type Secret struct {
	Name    string          `json:"name"`
	Source  string          `json:"source"`
	Options json.RawMessage `json:"options"`
}

type Secretbox struct {
	EnableSocket     bool          `json:"enable_socket"`
	Socket           string        `json:"socket"`
	MaxTokenTTL      time.Duration `json:"max_token_ttl"`
	CLIPath          string        `json:"cli_path"`
	TokensGCSchedule string        `json:"tokens_gc_schedule"`
	Secrets          []Secret      `json:"secrets"`
}
