package utils

import (
	"fmt"
	"go/ast"
	"go/token"
)

func StringConstsToVar(astFile *ast.File) error {
	for i := 0; i < len(astFile.Decls); i++ {
		astDecl := astFile.Decls[i]

		switch astDecl.(type) {
		case *ast.GenDecl:
			astDecl := astDecl.(*ast.GenDecl)
			if astDecl.Tok == token.CONST {
				if ConstOnlyHasStrings(astDecl) {
					astDecl.Tok = token.VAR
				}
			}
		}
	}

	return nil
}

func ConstOnlyHasStrings(decl *ast.GenDecl) bool {
	for _, spec := range decl.Specs {
		if cs, ok := spec.(*ast.ValueSpec); ok {
			if !specIsString(cs) {
				return false
			}
		}
	}
	return true
}

func specIsString(v *ast.ValueSpec) bool {
	if v.Type != nil {
		s, ok := v.Type.(fmt.Stringer)
		if ok && s.String() == "string" {
			return true
		}
	}
	if len(v.Values) != 1 {
		return false
	}
	return exprIsString(v.Values[0])
}

func exprIsString(e ast.Expr) bool {
	switch e := e.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			return true
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			return exprIsString(e.X) || exprIsString(e.Y)
		}
		return false
	case *ast.ParenExpr:
		return exprIsString(e.X)
	}
	return false
}
