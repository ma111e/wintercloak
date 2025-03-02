package utils

import (
	"go/ast"
	"go/token"
	"strconv"
)

func PatchImport(astFile *ast.File, packagePath string) error {
	var hasImports bool

importsLoop:
	for i := 0; i < len(astFile.Decls); i++ {
		astDecl := astFile.Decls[i]

		switch astDecl.(type) {
		case *ast.GenDecl:
			astDecl := astDecl.(*ast.GenDecl)

			if astDecl.Tok == token.IMPORT {
				for _, spec := range astDecl.Specs {
					// Skip special import "C"
					if spec.(*ast.ImportSpec).Path.Value == strconv.Quote("C") {
						continue importsLoop
					}
				}
				hasImports = true
				// Add the remote import
				newImportSpec := &ast.ImportSpec{Path: &ast.BasicLit{Value: strconv.Quote(packagePath)}}
				astDecl.Specs = append(astDecl.Specs, newImportSpec)
			}
		}
	}

	if !hasImports {
		newImport := &ast.GenDecl{
			Tok:    token.IMPORT,
			Lparen: 1, // Need non-zero Lparen and Rparen so that printer
			Rparen: 1, // treats this as a factored import.
		}

		newImportSpec := &ast.ImportSpec{Path: &ast.BasicLit{Value: strconv.Quote(packagePath)}}
		newImport.Specs = append(newImport.Specs, newImportSpec)
		astFile.Decls = append([]ast.Decl{newImport}, astFile.Decls...)
	}

	return nil
}
