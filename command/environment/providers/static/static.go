package static

import (
	"github.com/teran/backupnizza/command/environment"
)

func NewStaticVariable(s string) environment.Value {
	return func() (string, error) {
		return s, nil
	}
}
