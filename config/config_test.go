package config

import (
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	time "github.com/teran/go-time"
)

var configSample = Config{
	Common: Common{
		PackSize: 128,
		CacheDir: "/Volumes/SSD/Temp/Cache/restic",
		TmpDir:   "/Volumes/SSD/Temp/restic",
	},
	Log: json.RawMessage(mustJSON(map[string]any{
		"format": "text",
		"level":  "TRACE",
	})),
	Secretbox: Secretbox{
		EnableSocket:     true,
		Socket:           "/var/run/backup_secretbox.sock",
		CLIPath:          "/usr/local/bin/secretbox-cli",
		MaxTokenTTL:      15 * time.Second,
		TokensGCSchedule: "* * * * *",
		Secrets: []Secret{
			{
				Name:    "test secret",
				Source:  "onepassword",
				Options: json.RawMessage(`{"label":"restic:Vault"}`),
			},
		},
	},
	Tasks: []Task{
		{
			Name:     "task1",
			Kind:     "command",
			Schedule: "*/1 */2 */3 */4 */5",
			Options: json.RawMessage(mustJSON(map[string]any{
				"arguments": []string{"backup", "~/Documents"},
				"binary":    "/usr/local/bin/restic",
				"environment": []map[string]any{
					{
						"name":  "RESTIC_PACK_SIZE",
						"value": "128",
					},
					{
						"name":  "RESTIC_REPOSITORY",
						"value": "/Volumes/Vault/Backup",
					},
					{
						"name":       "RESTIC_PASSWORD",
						"value_from": "secretbox",
						"options": map[string]any{
							"name": "vault",
						},
					},
				},
			})),
		},
	},
}

func TestNewFromYAMLFile(t *testing.T) {
	r := require.New(t)

	cfg, err := newFromYAMLFile("testdata/config.yaml")
	r.NoError(err)
	r.Equal(&configSample, cfg)
}

func TestNewFromFile(t *testing.T) {
	type testCase struct {
		name            string
		inputFilename   string
		expOutputConfig *Config
		expOutputError  error
	}

	tcs := []testCase{
		{
			name:            "JSON config",
			inputFilename:   "testdata/config.json",
			expOutputConfig: &configSample,
		},
		{
			name:            "YAML config",
			inputFilename:   "testdata/config.yaml",
			expOutputConfig: &configSample,
		},
		{
			name:           "unknown config type",
			inputFilename:  "testdata/config.blah",
			expOutputError: errors.Errorf("unexpected file format: `%s`", "testdata/config.blah"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			cfg, err := NewFromFile(tc.inputFilename)
			if tc.expOutputError != nil {
				r.Error(err)
				r.Equal(tc.expOutputError.Error(), err.Error())
			} else {
				r.NoError(err)
				r.Equal(tc.expOutputConfig, cfg)
			}
		})
	}
}

func mustJSON(in map[string]any) []byte {
	data, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}

	return data
}
