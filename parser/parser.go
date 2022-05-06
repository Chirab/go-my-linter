package memcheck

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"unicode"
)

type Parser struct {
	flag     string
	filename string
}

func NewParser(flag, filename string) *Parser {
	return &Parser{
		flag:     flag,
		filename: filename,
	}
}

func (p *Parser) Run() {
	isFile, err := p.checkIFileOrDirectory()
	if err != nil {
		log.Fatal(err)
		return
	}

	if !isFile {
		p.recursiveFileParse()
		return
	}

	if isGoFile := p.checkGoFile(p.filename); !isGoFile {
		fmt.Println("Golang file is required : ", p.filename)
		return
	}
	p.parseSingleFile()
}

func (p *Parser) checkGoFile(f string) bool {
	if !strings.HasSuffix(f, ".go") {
		return false
	}

	return true
}

func (p *Parser) checkIFileOrDirectory() (bool, error) {
	file, err := os.Open(p.filename)
	if err != nil {
		return false, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}

	if fileInfo.IsDir() {
		return false, nil
	}
	defer file.Close()
	return true, nil
}

func (p *Parser) parseSingleFile() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, p.filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	p.inspectNode(node, fset)

}

func (p *Parser) inspectNode(node *ast.File, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch p.flag {
		case "-p":
			if val, ok := n.(*ast.CallExpr); ok {
				if ret, ok := val.Fun.(*ast.SelectorExpr); ok {
					if ret.Sel.Name == "Println" {
						log.Printf("Println call found on line %d at %s in package %s\n", fset.Position(val.Pos()).Line, fset.Position(val.Pos()).Filename, node.Name.Name)
					}
				}
			}
		case "-c":
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
		return true
	})
}

func (p *Parser) recursiveFileParse() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	directory := dir + "/" + p.filename
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, directory, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	if len(pkgs) == 0 {
		return
	}

	for _, v := range pkgs {
		for _, file := range v.Files {
			p.inspectNode(file, fset)
		}
	}

}
