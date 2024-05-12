package secretbox

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/teran/backupnizza/command/environment"
	random "github.com/teran/go-random"
	"github.com/teran/secretbox/service"
)

func NewSecretboxCommandVariable(svc service.Service, cliPath string, socket, secretName string) environment.Value {
	return newSecretboxCommandVariable(svc, cliPath, socket, secretName, random.String)
}

func newSecretboxCommandVariable(svc service.Service, cliPath string, socket, secretName string, randomGenerator func([]rune, uint) string) environment.Value {
	return func() (string, error) {
		token := randomGenerator(random.AlphaNumeric, 64)

		err := svc.CreateToken(context.TODO(), secretName, token)
		if err != nil {
			return "", errors.Wrap(err, "error creating access token in secretbox")
		}

		cmd := fmt.Sprintf("%s -p unix -l %s -s %s -t %s", cliPath, socket, secretName, token)
		log.WithFields(log.Fields{
			"component":   "secretbox_variable_generator",
			"secret_name": secretName,
		}).Trace("generated command")

		return cmd, nil
	}
}

func NewSecretboxVariable(svc service.Service, secretName string) environment.Value {
	return newSecretboxVariable(svc, secretName)
}

func newSecretboxVariable(svc service.Service, secretName string) environment.Value {
	return func() (string, error) {
		secret, err := svc.GetSecretNoAuth(context.TODO(), secretName)
		if err != nil {
			return "", errors.Wrap(err, "error requesting token")
		}

		return secret, nil
	}
}
