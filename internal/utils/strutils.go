package utils

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

func ExtractAllStrings(src string) []string {
	regex := `["'](.*?)["']\`
	markerPattern := regexp.MustCompile(regex)
	var extracted []string

	// Find all matches in the source string
	matches := markerPattern.FindAllStringSubmatch(src, -1)

	// Extract the captured groups (the strings inside the quotes)
	for _, match := range matches {
		if len(match) > 1 {
			extracted = append(extracted, match[1])
		}
	}

	return extracted
}

func String2HexBytesWithQMarks(in string) ([]byte, error) {
	var res []byte

	for i := 0; i < len(in); i += 2 {
		if in[i] == ' ' {
			i++
			continue
		}

		if in[i] == '?' {
			if in[i+1] != '?' {
				return nil, fmt.Errorf("wildcard character '?' must be by pair: found '%c%c'", in[i], in[i+1])
			}
		}
	}

	in = strings.ReplaceAll(in, " ", "")
	in = strings.ReplaceAll(in, "0x", "")
	in = strings.ReplaceAll(in, "\\x", "")

	if len(in)%2 != 0 {
		return nil, fmt.Errorf("invalid input length: must be a multiple of 2")
	}

	for i := 0; i < len(in); i += 2 {
		if in[i] == '?' {
			//if in[i+1] == '?' {
			res = append(res, []byte{'?', '?'}...)
			//}

			continue
		}
		b, err := hex.DecodeString(string(in[i]) + string(in[i+1]))
		if err != nil {
			return nil, err
		}
		res = append(res, b[0])
	}

	return res, nil
}

func Reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}
