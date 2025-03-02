package utils

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
)

func GenerateFile(fset *token.FileSet, astFile *ast.File) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := printer.Fprint(buf, fset, astFile); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Copy(src string, dst string) error {
	rawContent, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, rawContent, 0644)
	if err != nil {
		return err
	}

	return nil
}
