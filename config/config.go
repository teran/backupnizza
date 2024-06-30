package config

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Common    Common          `json:"common"`
	Secretbox Secretbox       `json:"secretbox"`
	Log       json.RawMessage `json:"log"`
	Tasks     []Task          `json:"tasks"`
}

func NewFromFile(filename string) (*Config, error) {
	switch strings.ToLower(path.Ext(filename)) {
	case ".yaml", ".json":
		return newFromYAMLFile(filename)
	}
	return nil, errors.Errorf("unexpected file format: `%s`", filename)
}

func newFromYAMLFile(filename string) (*Config, error) {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "error opening configuration file")
	}

	data := map[string]any{}
	if err := yaml.Unmarshal(fileData, &data); err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	return &cfg, json.Unmarshal(jsonData, &cfg)
}
