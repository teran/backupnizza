package command

import (
	"encoding/json"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Environment struct {
	Name      string          `json:"name"`
	Value     *string         `json:"value"`
	ValueFrom *string         `json:"value_from"`
	Options   json.RawMessage `json:"options"`
}

type CommandConfig struct {
	Name        string        `json:"name"`
	Environment []Environment `json:"environment"`
	Arguments   []string      `json:"arguments"`
	Binary      string        `json:"binary"`
}

type GroupOptions struct {
	ExecutionLogic ExecutionLogic        `json:"execution_logic"`
	Subtasks       []GroupOptionsSubtask `json:"subtasks"`
}

func (g GroupOptions) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.ExecutionLogic, validation.Required, validation.In(ExecutionLogicSequentialAll, ExecutionLogicSequentialDontCare, ExecutionLogicParallelAll)),
		validation.Field(&g.Subtasks, validation.Required),
	)
}

type GroupOptionsSubtask struct {
	Name        string        `json:"name"`
	Environment []Environment `json:"environment"`
	Arguments   []string      `json:"arguments"`
	Binary      string        `json:"binary"`
}

type ExecutionLogic string

const (
	ExecutionLogicSequentialAll      ExecutionLogic = "sequential_all"
	ExecutionLogicSequentialDontCare ExecutionLogic = "sequential_dont_care"
	ExecutionLogicParallelAll        ExecutionLogic = "parallel_all"
)

func (el *ExecutionLogic) UnmarshalJSON(data []byte) error {
	v, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	switch v {
	case string(ExecutionLogicSequentialAll):
		*el = ExecutionLogicSequentialAll
	case string(ExecutionLogicSequentialDontCare):
		*el = ExecutionLogicSequentialDontCare
	case string(ExecutionLogicParallelAll):
		*el = ExecutionLogicParallelAll
	default:
		*el = ExecutionLogicSequentialAll
	}

	return nil
}
