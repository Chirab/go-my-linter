package main

import (
	"fmt"
	"os"

	memcheck "github.com/chirab/go-memcheck/parser"
)

const (
	InfoColor   = "\033[1;34m%s\033[0m"
	NoticeColor = "\033[1;36m%s\033[0m"
)

func help() {
	fmt.Printf(NoticeColor, "Usage:\n")
	fmt.Printf(NoticeColor, "my_program -argument [option_file]\n")
	fmt.Printf(NoticeColor, "-h : display help\n")
	fmt.Printf(NoticeColor, "-v : display version\n")
	fmt.Printf(NoticeColor, "-p : inspect fmt.Println function call\n")
	fmt.Printf(NoticeColor, "-c : inspect comments\n")
	fmt.Printf(NoticeColor, "-ab : display abbreviations to name basic type variables")
}

func version() {
	fmt.Printf(InfoColor, "go-my-linter v.1.0\n")
}

func main() {
	if len(os.Args[1:]) == 0 {
		return
	}

	/*if os.Args[1] == "-h" {
		help()
		return
	}*/

	switch os.Args[1] {
	case "-h":
		help()
	case "-v":
		version()
	}

	if len(os.Args[1:]) != 2 {
		return
	}

	p := memcheck.NewParser(os.Args[1], os.Args[2])
	p.Run()

}
