package engines

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"wintercloak/internal/utils"
)

type EngineArgs []EngineArg

type EngineArg struct {
	Name              string      `json:"name"`
	Title             string      `json:"title"`
	Description       string      `json:"description"`
	Type              string      `json:"type"`
	DefaultValue      interface{} `json:"defaultValue"`
	Value             interface{} `json:"value"`
	AvailableOptions  []string    `json:"availableOptions"`
	TargetLanguage    string      `json:"targetLanguage"`
	NotTargetLanguage string      `json:"notTargetLanguage"`
	//DependsOn        string      `json:"dependsOn"`
}

func (ea *EngineArgs) InitDefault() {
	for idx, arg := range *ea {
		(*ea)[idx].Value = arg.DefaultValue
	}
}

func (ea *EngineArgs) GetByName(name string) *EngineArg {
	for i := range *ea {
		if (*ea)[i].Name == name {
			return &(*ea)[i]
		}
	}
	return nil
}

// GetStringByName returns the value as string if found and of correct type
func (ea *EngineArgs) GetStringByName(name string) (string, error) {
	arg := ea.GetByName(name)
	if arg == nil {
		return "", fmt.Errorf("argument %s not found", name)
	}

	if val, ok := arg.Value.(string); ok {
		return val, nil
	}
	return "", fmt.Errorf("argument %s is not a string", name)
}

// GetFloat64ByName returns the value as int if found and of correct type
func (ea *EngineArgs) GetFloat64ByName(name string) (float64, error) {
	arg := ea.GetByName(name)
	if arg == nil {
		return 0, fmt.Errorf("argument %s not found", name)
	}

	switch arg.Value.(type) {
	case string:
		if arg.Value == "" {
			return 0, nil
		}
		return strconv.ParseFloat(arg.Value.(string), 64)

	case int:
		return float64(arg.Value.(int)), nil
	case int8:
		return float64(arg.Value.(int8)), nil
	case int16:
		return float64(arg.Value.(int16)), nil
	case int32:
		return float64(arg.Value.(int32)), nil
	case int64:
		return float64(arg.Value.(int64)), nil
	case uint:
		return float64(arg.Value.(uint)), nil
	case uint8:
		return float64(arg.Value.(uint8)), nil
	case uint16:
		return float64(arg.Value.(uint16)), nil
	case uint32:
		return float64(arg.Value.(uint32)), nil
	case uint64:
		return float64(arg.Value.(uint64)), nil
	case float32:
		return float64(arg.Value.(float32)), nil
	case float64:
		return arg.Value.(float64), nil
	}

	return 0, fmt.Errorf("argument '%s' is not a valid number", name)
}

// GetByteByName returns the value as byte if found and of correct type
func (ea *EngineArgs) GetByteByName(name string) (byte, error) {
	arg := ea.GetByName(name)
	if arg == nil {
		return 0, fmt.Errorf("argument %s not found", name)
	}

	switch v := arg.Value.(type) {
	case byte:
		return v, nil
	case string:
		parsed, err := utils.String2HexBytesWithQMarks(v)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse byte representation")
		}

		if bytes.Equal(parsed, []byte("??")) {
			return byte(rand.Intn(256)), err
		}

		if len(parsed) == 1 {
			return parsed[0], nil
		}

		if len(parsed) == 0 {
			return byte(0), nil
		}

		return 0, fmt.Errorf("failed to parse byte representation: malformed input")
	}

	return 0, fmt.Errorf("argument %s is not a valid representation of a byte (\\xXX, 0xXX, ?? or a single ASCII char) ", name)
}

// GetByteSliceByName returns the value as []byte if found and of correct type
func (ea *EngineArgs) GetByteSliceByName(name string) ([]byte, error) {
	arg := ea.GetByName(name)
	if arg == nil {
		return nil, fmt.Errorf("argument %s not found", name)
	}

	if val, ok := arg.Value.([]byte); ok {
		return val, nil
	}
	return nil, fmt.Errorf("argument %s is not a []byte", name)
}
