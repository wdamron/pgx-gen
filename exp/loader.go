package exp

import (
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/types"
)

type Program struct {
	loader.Program
}

func Load(importPath string) (*Program, error) {
	cfg := &loader.Config{
		Fset:        nil,
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
	}
	cfg.Import(importPath)
	p, err := cfg.Load()
	return &Program{p}, err
}

func (p *Program) StructTypes() []types.Object {
	var objs []types.Object
	for path, pkgInfo := range p.Imported {
		scope := p.Pkg.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if !isStruct(obj.Type()) {
				continue
			}
			objs = append(objs, obj)
		}
	}
	return objs
}

func isStruct(t types.Type) bool {
	_, ok := t.(*types.Struct)
	return ok
}

func isBasic(t types.Type) bool {
	_, ok := t.(*types.Basic)
	return ok
}
