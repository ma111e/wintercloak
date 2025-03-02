package engines

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"wintercloak/internal/utils"
)

type TextParser struct {
	OriginalFilepath string
	Content          []byte
	patchMarker      string
	skipEmptyString  bool
	Initialized      bool
}

func (p *TextParser) Patch(applyFn func(string) ([]byte, error)) error {
	if !p.Initialized {
		return fmt.Errorf("parser not initialized")
	}

	extracted := p.ExtractMarkedStrings()
	for _, val := range extracted {
		patched, err := applyFn(val)
		if err != nil {
			return err
		}

		// Replace PATCH('string') or PATCH("string") with the generated code
		// Note that all the same strings will be patched with the same generated code
		pattern := fmt.Sprintf(`['"]%s['"]`, regexp.QuoteMeta(val))
		toPatchRegex := regexp.MustCompile(pattern)

		p.Content = []byte(toPatchRegex.ReplaceAllLiteralString(string(p.Content), string(patched)))
	}

	return nil
}

func (p *TextParser) ExtractMarkedStrings() []string {
	sc := bufio.NewScanner(strings.NewReader(string(p.Content)))
	var extracted []string

	for sc.Scan() {
		line := sc.Text()
		if strings.Contains(line, p.patchMarker) {
			sc.Scan()
			nextLine := sc.Text()
			extracted = append(extracted, utils.ExtractAllStrings(nextLine)...)
		}
	}

	return extracted
}

func (p *TextParser) Init(opt RunOptions) error {
	if !opt.IgnoreExtensions && !slices.Contains(opt.TargetExtensions, filepath.Ext(opt.TargetPath)) {
		return fmt.Errorf("failed to load the parser: '%s' is not a supported file (valid extensions: %v)", opt.TargetPath, opt.TargetExtensions)
	}

	p.skipEmptyString = opt.SkipEmpty
	p.patchMarker = opt.PatchMarker
	p.OriginalFilepath = opt.TargetPath
	err := p.Reload()
	if err != nil {
		return err
	}

	p.Initialized = true

	return nil
}

func (p *TextParser) Reload() error {
	text, err := os.ReadFile(p.OriginalFilepath)
	if err != nil {
		return err
	}

	p.Content = text
	return nil
}
