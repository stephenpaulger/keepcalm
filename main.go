package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func getFuncName(ce *ast.CallExpr) (string, error) {
	switch fun := ce.Fun.(type) {
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", fun.X, fun.Sel), nil
	case *ast.Ident:
		return fun.Name, nil
	}
	return "", errors.New("could not determine call expression's identifier")
}

func main() {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, ".", nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				ce, ok := n.(*ast.CallExpr)
				if ok {
					if fn, err := getFuncName(ce); err == nil {
						fmt.Printf("%s:\t%s\n", fset.Position(ce.Pos()), fn)
					}
				}
				return true
			})
		}
	}

	return

}
