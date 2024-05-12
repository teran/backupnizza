package secretbox

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/teran/secretbox/service"
)

func TestNewSecretboxCommandVariable(t *testing.T) {
	r := require.New(t)

	svcM := service.NewMock()
	defer svcM.AssertExpectations(t)
	svcM.On("CreateToken", "secret_name", "random value").Return(nil).Once()

	rg := func([]rune, uint) string {
		return "random value"
	}

	fn := newSecretboxCommandVariable(svcM, "/usr/local/bin/secretbox-cli", "socket", "secret_name", rg)
	v, err := fn()
	r.NoError(err)
	r.Equal("/usr/local/bin/secretbox-cli -p unix -l socket -s secret_name -t random value", v)
}

func TestNewSecretboxVariable(t *testing.T) {
	r := require.New(t)

	svcM := service.NewMock()
	defer svcM.AssertExpectations(t)
	svcM.On("GetSecretNoAuth", "secret_name").Return("secret data", nil).Once()

	fn := newSecretboxVariable(svcM, "secret_name")
	v, err := fn()
	r.NoError(err)
	r.Equal("secret data", v)
}
