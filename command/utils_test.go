package command

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/teran/backupnizza/command/environment"
)

func TestMapMapToSlice(t *testing.T) {
	r := require.New(t)

	v, err := mapMapToSlice(map[string]environment.Value{
		"key1": func() (string, error) { return "value1", nil },
		"key2": func() (string, error) { return "value2", nil },
	})
	r.NoError(err)
	r.ElementsMatch([]string{
		"key1=value1",
		"key2=value2",
	}, v)
}
