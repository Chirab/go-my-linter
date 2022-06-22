package memcheck

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Parser struct {
	flag string
	args string
}

func NewParser(flag, args string) *Parser {
	return &Parser{
		flag: flag,
		args: args,
	}
}

func (p *Parser) Run() error {
	fls, err := p.getRecursiveFiles()
	if err != nil {
		log.Fatal(err)
		return err
	}

	asts, fset, err := p.getAstFiles(fls)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err := p.processFiles(asts, fset); err != nil {
		return err
	}

	return nil
}

func (p *Parser) processFiles(astfiles []*ast.File, fs *token.FileSet) error {
	for _, file := range astfiles {
		InspectNode(p.flag, file, fs)
	}
	return nil
}

func (p *Parser) checkGoFile(f string) bool {
	if !strings.HasSuffix(f, ".go") {
		return false
	}

	return true
}

func (p *Parser) getRecursiveFiles() ([]string, error) {
	fileList := make([]string, 0, 1)
	err := filepath.Walk(p.args,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if ok := p.checkGoFile(path); ok {
				fileList = append(fileList, path)
			}

			return nil
		})
	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func (p *Parser) getAstFiles(fileList []string) ([]*ast.File, *token.FileSet, error) {
	astFileList := make([]*ast.File, 0, 1)
	fset := token.NewFileSet()
	for _, filepath := range fileList {
		asf, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
			return nil, nil, err
		}

		if asf == nil {
			return nil, nil, nil
		}

		astFileList = append(astFileList, asf)
	}

	return astFileList, fset, nil
}
