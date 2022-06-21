package memcheck

import (
	"go/ast"
	"go/token"
	"log"
	"strings"
	"unicode"
)

type Inspect struct {
	flag string
	node *ast.File
	fset *token.FileSet
}

func InspectNode(flag string, node *ast.File, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch flag {
		case "-p":
			findPrintln(node, n, fset)
		case "-c":
			findWrongComments(node, n, fset)
		}
		return true
	})
}

func findWrongComments(node *ast.File, n ast.Node, fset *token.FileSet) {
	if val, ok := n.(*ast.FuncDecl); ok {
		if val.Doc.Text() != "" {
			if string(val.Doc.Text()[0]) != "" {
				if len(strings.Split(val.Doc.Text(), "\n")) == 2 {
					if unicode.IsUpper(rune(val.Doc.Text()[0])) == false || strings.HasSuffix(val.Doc.Text(), ".") == true {
						log.Printf("Rewrite your comment at line %d at %s\n", fset.Position(val.Pos()).Line-1, fset.Position(val.Pos()).Filename)
					}
				}
			}
		}

	}
}

func findPrintln(node *ast.File, n ast.Node, fset *token.FileSet) {
	if val, ok := n.(*ast.CallExpr); ok {
		if ret, ok := val.Fun.(*ast.SelectorExpr); ok {
			if ret.Sel.Name == "Println" {
				log.Printf("Println call found on line %d at %s in package %s\n", fset.Position(val.Pos()).Line, fset.Position(val.Pos()).Filename, node.Name.Name)
			}
		}
	}
}
