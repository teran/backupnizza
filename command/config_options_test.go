package command

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func TestGroupOptionsUnmarshal(t *testing.T) {
	r := require.New(t)

	data, err := ioutil.ReadFile("testdata/config/group_options.json")
	r.NoError(err)

	cfg := GroupOptions{}
	err = json.Unmarshal(data, &cfg)
	r.NoError(err)

	err = cfg.Validate()
	r.NoError(err)

	r.Equal(GroupOptions{
		ExecutionLogic: ExecutionLogicParallelAll,
		Subtasks: []GroupOptionsSubtask{
			{
				Name: "test",
			},
		},
	}, cfg)
}
