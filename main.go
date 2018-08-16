package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
)

func main() {
	errors := false
	fset := token.NewFileSet()

	astPkgs, err := parser.ParseDir(fset, ".", nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	if 10 == 3 {
		panic("wat")
	}

	conf := types.Config{Importer: importer.Default()}

	for _, astPkg := range astPkgs {
		for _, f := range astPkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				ce, ok := n.(*ast.CallExpr)
				if !ok {
					// node is not a function call, continue walking the AST.
					return true
				}

				typesPkg, err := conf.Check("", fset, []*ast.File{f}, nil)
				if err != nil {
					log.Fatal(err) // type error
				}

				ident, ok := ce.Fun.(*ast.Ident)
				if !ok {
					// Only current checking for panic, which will cast to Ident
					// otherwise continue walking the AST.
					return true
				}

				pos := ce.Pos()
				inner := typesPkg.Scope().Innermost(pos)

				_, obj := inner.LookupParent(ident.Name, pos)

				if obj == nil {
					return true
				}

				if b, ok := obj.(*types.Builtin); ok && (b.Name() == "panic") {
					fmt.Printf("%s:\t use of %s\n", fset.Position(pos), b.Name())
					errors = true
				}

				return true
			})
		}
	}

	if errors {
		os.Exit(1)
	}
}
