package command

import (
	"fmt"

	"github.com/teran/backupnizza/command/environment"
)

func mapMapToSlice(in map[string]environment.Value) ([]string, error) {
	out := []string{}
	for k, v := range in {
		envValue, err := v()
		if err != nil {
			return nil, err
		}
		out = append(out, fmt.Sprintf("%s=%s", k, envValue))
	}
	return out, nil
}
