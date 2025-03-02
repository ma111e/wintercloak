package engines

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"wintercloak/internal/utils"
)

type GoParser struct {
	OriginalFilepath string
	Content          []byte
	StringNodes      []*ast.BasicLit
	FnCallNodes      []*ast.CallExpr
	skipEmptyString  bool
	AstFile          *ast.File
	Fset             *token.FileSet
	Initialized      bool
	RunOptions       RunOptions
	BuildConstraints []string
	WclkDirectives   []string
}

const (
	DONOTEDIT = "// Code generated by Wintercloak DO NOT EDIT."
)

//func NewGoParser(opt RunOptions) (*GoParser, error) {
//	if filepath.Ext(opt.TargetPath) != ".go" {
//		return nil, fmt.Errorf("failed to load the Go parser: '%s' is not a go file", opt.TargetPath)
//	}
//
//	p := &GoParser{skipEmptyString: opt.SkipEmpty, OriginalFilepath: opt.TargetPath}
//	err := p.Reload()
//	if err != nil {
//		return nil, err
//	}
//
//	p.Initialized = true
//
//	return p, nil
//}

func (p *GoParser) Patch(applyFn func(string) ([]byte, error)) error {
	if !p.Initialized {
		return fmt.Errorf("parser not initialized")
	}

	p.AddDoNotEditMention()

	err := p.StringConstsToVar()
	if err != nil {
		return err
	}

	if p.RunOptions.PatchMarker != "" {
		err = p.ParseMarkedStrings(applyFn)
		if err != nil {
			return err
		}
	} else {
		err = p.ParseStrings(applyFn)
		if err != nil {
			return err
		}
	}

	err = p.RemoveGoGenerate()
	if err != nil {
		return err
	}

	return nil
}

func (p *GoParser) Init(opt RunOptions) error {
	if !opt.IgnoreExtensions && filepath.Ext(opt.TargetPath) != ".go" {
		return fmt.Errorf("failed to load the Go parser: '%s' is not a go file", opt.TargetPath)
	}

	p.OriginalFilepath = opt.TargetPath
	p.RunOptions = opt

	err := p.Reload()
	if err != nil {
		return err
	}

	p.Initialized = true
	return nil
}

func (p *GoParser) RefreshAst() {
	for _, decl := range p.AstFile.Decls {
		ast.Walk(p, decl)
	}
}

func (p *GoParser) Reload() error {
	text, err := os.ReadFile(p.OriginalFilepath)
	if err != nil {
		return err
	}

	// Ugly hack to allow the DONOTEDIT comment to be put before the package keyword using AST
	// Top directives (//xx:yy) are ignored by the AST parser
	// go/scanner could be used, however keeping a single representation of the file throughout its modification
	// is better
	if bytes.HasPrefix(text, []byte("package")) {
		// Add empty comment to trick the ast parser
		text = append([]byte("//\n"), text...)
	}

	fset := token.NewFileSet()
	p.AstFile, err = parser.ParseFile(fset, p.OriginalFilepath, text, parser.ParseComments)
	if err != nil {
		return err
	}

	p.Content = text
	p.Fset = fset

	for _, decl := range p.AstFile.Decls {
		ast.Walk(p, decl)
	}

	p.LoadBuildConstraints()
	p.LoadWintercloakDirectives()

	return nil
}

func (p *GoParser) Visit(n ast.Node) ast.Visitor {
	if lit, ok := n.(*ast.BasicLit); ok {
		if lit.Kind == token.STRING || lit.Kind == token.CONST {
			p.StringNodes = append(p.StringNodes, lit)
		}

		return nil
	} else if callExpr, ok := n.(*ast.CallExpr); ok {
		p.FnCallNodes = append(p.FnCallNodes, callExpr)
	} else if decl, ok := n.(*ast.GenDecl); ok {
		if decl.Tok == token.IMPORT {
			return nil
		}
	} else if _, ok := n.(*ast.StructType); ok {
		// Avoid messing with annotation strings.
		return nil
	}

	return p
}

func (p *GoParser) HasGoBuildConstraint(constraint string) bool {
	return slices.Contains(p.BuildConstraints, constraint)
}

func (p *GoParser) StringConstsToVar() error {
	for i := 0; i < len(p.AstFile.Decls); i++ {
		astDecl := p.AstFile.Decls[i]

		switch astDecl.(type) {
		case *ast.GenDecl:
			astDecl := astDecl.(*ast.GenDecl)
			if astDecl.Tok == token.CONST {
				if utils.ConstOnlyHasStrings(astDecl) {
					astDecl.Tok = token.VAR
				}
			}
		}
	}

	return nil
}

//func (p *GoParser) RemovePlusBuildConstraint(constraint string) error {
//	prefix := "// +build"
//	argSep := ","
//
//	for _, commGroup := range p.AstFile.Comments {
//		for commIdx, comm := range commGroup.List {
//			if strings.HasPrefix(comm.Text, prefix) {
//				line := strings.Replace(comm.Text, prefix+" ", "", 1)
//
//				constraintIdx := -1
//				argsChunks := strings.Split(line, argSep)
//				for idx, chunk := range argsChunks {
//					if chunk == constraint {
//						constraintIdx = idx
//					}
//				}
//
//				if constraintIdx == -1 {
//					// constraint not found
//					continue
//				}
//
//				var fixed string
//
//				argsChunks = append(argsChunks[:constraintIdx], argsChunks[constraintIdx+1:]...)
//				args := strings.Join(argsChunks, argSep)
//				if len(argsChunks) > 0 {
//					fixed = fmt.Sprintf("%s %s", prefix, args)
//					commGroup.List[commIdx].Text = fixed
//				} else {
//					commGroup.List = append(commGroup.List[:commIdx], commGroup.List[commIdx+1:]...)
//				}
//
//			}
//		}
//	}
//
//	return nil
//}

func (p *GoParser) RemoveGoBuildConstraint(constraint string) error {
	prefix := "//go:build"
	argSep := " && "

	for _, commGroup := range p.AstFile.Comments {
		for commIdx, comm := range commGroup.List {
			if strings.HasPrefix(comm.Text, prefix) {
				line := strings.Replace(comm.Text, prefix+" ", "", 1)

				constraintIdx := -1
				argsChunks := strings.Split(line, argSep)
				for idx, chunk := range argsChunks {
					if chunk == constraint {
						constraintIdx = idx
					}
				}

				if constraintIdx == -1 {
					// constraint not found
					continue
				}

				var fixed string

				argsChunks = append(argsChunks[:constraintIdx], argsChunks[constraintIdx+1:]...)
				args := strings.Join(argsChunks, argSep)
				if len(argsChunks) > 0 {
					fixed = fmt.Sprintf("%s %s", prefix, args)
					commGroup.List[commIdx].Text = fixed
				} else {
					commGroup.List = append(commGroup.List[:commIdx], commGroup.List[commIdx+1:]...)
				}

			}
		}
	}

	return nil
}

func (p *GoParser) LoadWintercloakDirectives() bool {
	prefix := "wintercloak:"

	for _, commGroup := range p.AstFile.Comments {
		for _, rawComm := range commGroup.List {
			comm := utils.CommentText(rawComm.Text)
			if strings.HasPrefix(comm, prefix) {
				comm = strings.Replace(comm, prefix, "", 1)
				p.WclkDirectives = append(p.WclkDirectives, strings.TrimSpace(comm))
			}
		}
	}

	return false
}

func (p *GoParser) LoadBuildConstraints() bool {
	prefix := "go:build "

	for _, commGroup := range p.AstFile.Comments {
		for _, rawComm := range commGroup.List {
			comm := utils.CommentText(rawComm.Text)
			if strings.HasPrefix(comm, prefix) {
				comm = strings.Replace(comm, prefix, "", 1)
				for _, constraint := range strings.Split(comm, " && ") {
					p.BuildConstraints = append(p.BuildConstraints, strings.TrimSpace(constraint))
				}
			}
		}
	}

	return false
}

func (p *GoParser) RemoveBuildConstraints(constraint string) error {
	//err := p.RemovePlusBuildConstraint(constraint)
	//if err != nil {
	//	return err
	//}

	return p.RemoveGoBuildConstraint(constraint)
}

func (p *GoParser) AddDoNotEditMention() {
	doNotEditComment := ast.Comment{
		Text:  DONOTEDIT,
		Slash: p.AstFile.FileStart, // Position at the start of file
	}

	doNotEditCommentGroup := &ast.CommentGroup{
		List: []*ast.Comment{&doNotEditComment},
	}

	// If there are existing comments, prepend our comment
	if len(p.AstFile.Comments) > 0 {
		p.AstFile.Comments[0].List = append([]*ast.Comment{&doNotEditComment}, p.AstFile.Comments[0].List...)
	} else {
		p.AstFile.Comments = []*ast.CommentGroup{doNotEditCommentGroup}
	}

}

func (p *GoParser) RemoveGoGenerate() error {
	goGenLine := fmt.Sprintf("//go:generate %s", filepath.Base(os.Args[0]))

	for _, commGroup := range p.AstFile.Comments {
		for idx, comm := range commGroup.List {
			if strings.HasPrefix(comm.Text, goGenLine) {
				commGroup.List = append(commGroup.List[:idx], commGroup.List[idx+1:]...)
			}
		}
	}
	return nil
}

func (p *GoParser) UpdateGoEmbed(from string, to string) error {
	for _, commGroup := range p.AstFile.Comments {
		for idx, comm := range commGroup.List {
			if strings.HasPrefix(comm.Text, "//go:embed ") {
				commGroup.List[idx].Text = strings.Replace(commGroup.List[idx].Text, from, to, 1)
			}
		}
	}
	return nil
}

func (p *GoParser) ParseMarkedStrings(apply func(string) ([]byte, error)) error {
	for _, node := range p.FnCallNodes {
		// Check if this is a function call with the name matching the patch marker
		if funcIdent, ok := node.Fun.(*ast.Ident); ok && funcIdent.Name == p.RunOptions.PatchMarker {
			// We found a function call with the patch marker name
			// Only process if there's at least one argument
			if len(node.Args) > 0 {
				// Get the first argument, which should be a string
				if strArg, ok := node.Args[0].(*ast.BasicLit); ok && strArg.Kind == token.STRING {
					// Extract the string value
					unquoted, err := strconv.Unquote(strArg.Value)
					if err != nil {
						return err
					}

					// Apply the transformation function
					patched, err := apply(unquoted)
					if err != nil {
						return err
					}

					// Update the argument with the new value
					strArg.Value = string(patched)
				}
			}
		}
	}

	return nil

}

func (p *GoParser) ParseStrings(apply func(string) ([]byte, error)) error {
	for _, node := range p.StringNodes {

		unquoted, err := strconv.Unquote(node.Value)
		if err != nil {
			return err
		}

		if unquoted == "" && p.skipEmptyString {
			continue
		}

		patched, err := apply(unquoted)
		if err != nil {
			return err
		}

		node.Value = string(patched)
	}

	return nil
}

func (p *GoParser) WritePatchedFile(patched []byte) error {
	ext := filepath.Ext(p.OriginalFilepath)
	base := strings.TrimSuffix(p.OriginalFilepath, ext)
	destFilepath := base + "_gen" + ext

	err := utils.Copy(p.OriginalFilepath, destFilepath)
	if err != nil {
		return err
	}

	err = os.WriteFile(destFilepath, patched, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p *GoParser) Len() int {
	return len(p.StringNodes)
}

func (p *GoParser) Swap(i, j int) {
	p.StringNodes[i], p.StringNodes[j] = p.StringNodes[j], p.StringNodes[i]
}

func (p *GoParser) Less(i, j int) bool {
	return p.StringNodes[i].Pos() < p.StringNodes[j].Pos()
}
