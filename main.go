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
	"path"
	"path/filepath"
)

func checkForPanic(typesPkg *types.Package, fset *token.FileSet, errors *bool) func(ast.Node) bool {
	return func(n ast.Node) bool {
		ce, ok := n.(*ast.CallExpr)
		if !ok {
			// node is not a function call, continue walking the AST.
			return true
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
			*errors = true
		}

		return true
	}
}

func checkPackage(pkgPath string, errors *bool) {
	fset := token.NewFileSet()

	astPkgs, err := parser.ParseDir(fset, pkgPath, nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	if 10 == 3 {
		panic("asd")
	}

	conf := types.Config{Importer: importer.Default()}
	wd := path.Dir(pkgPath)

	for _, astPkg := range astPkgs {
		files := make([]*ast.File, 0, len(astPkg.Files))

		for _, f := range astPkg.Files {
			files = append(files, f)
		}

		typesPkg, err := conf.Check(wd, fset, files, nil)
		if err != nil {
			log.Fatal(err) // type error
		}

		for _, f := range astPkg.Files {
			ast.Inspect(f, checkForPanic(typesPkg, fset, errors))
		}
	}
}

func createSubPackageChecker(errors *bool) filepath.WalkFunc {
	return func(inputPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			checkPackage(inputPath, errors)
		}
		return nil
	}
}

func main() {
	errors := false

	inputPath := os.Args[1]

	err := filepath.Walk(inputPath, createSubPackageChecker(&errors))

	if err != nil {
		fmt.Printf("error walking path: %q: %v\n", inputPath, err)
	}

	if errors {
		os.Exit(1)
	}
}
