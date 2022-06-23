package memcheck

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"strings"
	"sync"
	"unicode"
)

type Inspect struct {
	flag string
	node *ast.File
	fset *token.FileSet
}

type basicNamingInfo struct {
	line     int
	filename string
	typeLoop string
}

func InspectNode(flag string, node *ast.File, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch flag {
		case "-p":
			findPrintln(node, n, fset)
		case "-c":
			findWrongComments(n, fset)
		case "-n":
			// Do not use abbreviations to name basic type variable
			findBasicNamingVariableForLoop(node, n, fset)
		case "-l":
			methodsUnderLimitLines(node, n, fset)
		}
		return true
	})
}

func convertNodeString(fset *token.FileSet, ret ast.Node) (string, error) {
	var typeNameBuf bytes.Buffer
	err := printer.Fprint(&typeNameBuf, fset, ret)
	if err != nil {
		log.Fatalf("failed printing %s", err)
		return "", err
	}
	return typeNameBuf.String(), nil
}

func findRangeLoop(wg *sync.WaitGroup, rlp chan basicNamingInfo, n ast.Node, fset *token.FileSet) {
	defer wg.Done()
	bs := basicNamingInfo{}
	if val, ok := n.(*ast.RangeStmt); ok {
		if val.Value != nil {
			if ret, ok := val.Value.(ast.Node); ok {
				varname, err := convertNodeString(fset, ret)
				if err != nil {
					log.Fatal(err)
				}

				if len(varname) == 1 {
					bs.filename = fset.Position(val.Pos()).Filename
					bs.line = fset.Position(val.Pos()).Line
					bs.typeLoop = "range"
					rlp <- bs

				}
			}
		}
	}
}

func findForLoop(wg *sync.WaitGroup, flp chan basicNamingInfo, n ast.Node, fset *token.FileSet) {
	defer wg.Done()
	bs := basicNamingInfo{}
	if val, ok := n.(*ast.ForStmt); ok {
		if ret, ok := val.Init.(*ast.AssignStmt); ok {
			t := ret.Lhs[0]
			if ret, ok := t.(ast.Node); ok {
				varname, err := convertNodeString(fset, ret)
				if err != nil {
					log.Fatal(err)
				}

				if len(varname) == 1 {
					bs.filename = fset.Position(val.Pos()).Filename
					bs.line = fset.Position(val.Pos()).Line
					bs.typeLoop = "for"
					flp <- bs
				}
			}
		}
	}
}

func findBasicNamingVariableForLoop(node *ast.File, n ast.Node, fset *token.FileSet) {

	rlp := make(chan basicNamingInfo, 1)
	flp := make(chan basicNamingInfo, 1)

	var wg sync.WaitGroup
	wg.Add(2)
	go findForLoop(&wg, flp, n, fset)
	go findRangeLoop(&wg, rlp, n, fset)

	wg.Wait()
	close(rlp)
	close(flp)
	for v := range rlp {
		log.Printf("abbreviation to name basic type variables in %s loop at line %d in %s", v.typeLoop, v.line, v.filename)
	}

	for v := range flp {
		log.Printf("abbreviation to name basic type variables in %s loop at line %d in %s", v.typeLoop, v.line, v.filename)
	}
}

func findWrongComments(n ast.Node, fset *token.FileSet) {
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

func methodsUnderLimitLines(node *ast.File, n ast.Node, fset *token.FileSet) {
	if val, ok := n.(*ast.FuncDecl); ok {
		sum := (fset.Position(val.Body.Rbrace).Line - fset.Position(val.Body.Lbrace).Line) - 1
		if sum > 80 {
			log.Printf("methods above 80 lines at %s", fset.Position(val.Body.Pos()))
		}

	}
}
