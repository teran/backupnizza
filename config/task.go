package config

import "encoding/json"

type Task struct {
	Name     string          `json:"name"`
	Kind     string          `json:"kind"`
	Schedule string          `json:"schedule"`
	Options  json.RawMessage `json:"options"`
}
