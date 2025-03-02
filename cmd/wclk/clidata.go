package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CLIDataMapping is a custom type for parsing key=value flags
type CLIDataMapping map[string]string

func (kv *CLIDataMapping) String() string {
	b, _ := json.Marshal(*kv)
	return string(b)
}

func (kv *CLIDataMapping) Set(rawValue string) error {
	parts := strings.SplitN(rawValue, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid key-rawValue pair: %s, expected format key=rawValue", rawValue)
	}

	(*kv)[parts[0]] = parts[1]
	return nil
}

func (kv *CLIDataMapping) Type() string {
	return "CLIDataMapping"
}
