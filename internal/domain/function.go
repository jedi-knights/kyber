// Package domain holds the pure logic of kyber — the Function and Package
// models, the Metric interface that every code-quality check implements, the
// Registry that holds those metrics, and the Score/Finding/Report types that
// flow back to the reporters.
//
// The domain package depends only on the Go standard library (including
// go/ast, go/parser, and go/token). It must not import any of the adapter
// packages; dependencies flow inward.
package domain

import (
	"go/ast"
	"go/token"
)

// Function is a single parsed function or method. It carries the AST node,
// the FileSet needed to resolve positions, and a back-reference to its
// enclosing Package so metrics can inspect siblings (interfaces, globals)
// without re-parsing the source.
type Function struct {
	Name        string
	Receiver    string
	Package     *Package
	File        string
	FuncDecl    *ast.FuncDecl
	FileSet     *token.FileSet
	SourceLines []string
	Doc         string
}

// Position returns the function's starting source position.
func (f *Function) Position() token.Position {
	return f.FileSet.Position(f.FuncDecl.Pos())
}

// LineCount returns the number of source lines spanned by the function,
// inclusive of the opening and closing braces.
func (f *Function) LineCount() int {
	start := f.FileSet.Position(f.FuncDecl.Pos()).Line
	end := f.FileSet.Position(f.FuncDecl.End()).Line
	return end - start + 1
}

// ParamCount returns the total number of parameters declared on the
// function. Each name in a grouped declaration (e.g. `a, b int`) counts
// individually; an anonymous parameter (no name) counts as one.
func (f *Function) ParamCount() int {
	if f.FuncDecl.Type.Params == nil {
		return 0
	}
	n := 0
	for _, field := range f.FuncDecl.Type.Params.List {
		if len(field.Names) == 0 {
			n++
			continue
		}
		n += len(field.Names)
	}
	return n
}

// Package is the enclosing package, providing the cross-function context
// that interface-aware and global-aware metrics depend on.
type Package struct {
	Name       string
	ImportPath string
	Files      []*ast.File
	FileSet    *token.FileSet
	Interfaces map[string]*ast.InterfaceType
	Globals    map[string]token.Pos
}
