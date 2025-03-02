package engines

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"wintercloak/internal/utils"
)

type MaskEngine struct {
	Engine
	GoParser
	CustomPattern []byte
}

func init() {
	args := EngineArgs{
		{
			Name:         "pattern",
			Title:        "Pattern",
			Description:  "A pattern to repeat instead of random bytes",
			Type:         "[]byte",
			DefaultValue: "",
		},
	}

	Available.RegisterNewEngine(
		"Mask",
		"Use a XOR mask for each string",
		"",
		"golang",
		"ph:mask-happy-fill",
		args,
		NewMaskEngine,
	)
}

func NewMaskEngine(args EngineArgs) (Engine, error) {
	rawCustomPattern, err := args.GetStringByName("pattern")
	if err != nil {
		return nil, err
	}

	customPattern, err := utils.String2HexBytesWithQMarks(rawCustomPattern)
	if err != nil {
		return nil, err
	}

	return &MaskEngine{
		GoParser:      GoParser{},
		CustomPattern: customPattern,
	}, nil
}

func (e *MaskEngine) Run(runOptions RunOptions) ([]byte, error) {
	err := e.Init(runOptions)
	if err != nil {
		return []byte{}, err
	}

	if e.RunOptions.UseGoBuildTag && !e.HasGoBuildConstraint(e.RunOptions.GoBuildTag) {
		// skip
		return []byte{}, fmt.Errorf("'%s' filtered out because of missing go constraint '%s'", runOptions.TargetPath, runOptions.GoBuildTag)
	}

	err = e.Patch(e.Apply)
	if err != nil {
		return []byte{}, err
	}

	if e.RunOptions.UseGoBuildTag {
		err = e.RemoveBuildConstraints(e.RunOptions.GoBuildTag)
		if err != nil {
			return []byte{}, err
		}
	}

	patched, err := utils.GenerateFile(e.Fset, e.AstFile)
	if err != nil {
		return []byte{}, err
	}

	return patched, nil
}

func (e *MaskEngine) WritePatchedFile(patched []byte) error {
	return e.GoParser.WritePatchedFile(patched)
}

func (e *MaskEngine) Apply(val string) ([]byte, error) {
	var res bytes.Buffer
	res.WriteString("\n(func() string {\n")
	res.WriteString("\tmask := []byte(\"")
	mask := make([]byte, len(val))
	patternLength := len(e.CustomPattern)

	if len(e.CustomPattern) > 0 {
		for i := range mask {
			if e.CustomPattern[i%patternLength] == '?' {
				mask[i] = byte(rand.Intn(256))
			} else {
				mask[i] = e.CustomPattern[i%patternLength]
			}

			res.WriteString(fmt.Sprintf("\\x%02x", mask[i]))
		}
	} else {
		for i := range mask {
			mask[i] = byte(rand.Intn(256))
			res.WriteString(fmt.Sprintf("\\x%02x", mask[i]))
		}
	}

	res.WriteString("\")\n\tmaskedStr := []byte(\"")
	for i, x := range []byte(val) {
		res.WriteString(fmt.Sprintf("\\x%02x", x^mask[i]))
	}
	res.WriteString("\")\n\tres := make([]byte, ")
	res.WriteString(strconv.Itoa(len(mask)))
	res.WriteString(`)
        for i, m := range mask {
            res[i] = m ^ maskedStr[i]
        }
        return string(res)
        }())`)

	return res.Bytes(), nil
}
