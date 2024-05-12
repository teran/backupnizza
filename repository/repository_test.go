package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	r := require.New(t)

	repo, err := NewRepository("some_label", "some_password")
	r.NoError(err)
	r.Equal("some_label", repo.Label())
	r.Equal("some_password", repo.password())
}
