package engineutils

import (
	"fmt"
	"wintercloak/internal/engines"
)

func NewEngine(engineName string, args engines.EngineArgs) (engines.Engine, error) {
	if constructor, found := engines.NameToFactoryMapping[engineName]; found {
		eng, err := constructor(args)
		if err != nil {
			return nil, err
		}

		return eng, nil
	}

	return nil, fmt.Errorf("unknown engine '%s'", engineName)
}
