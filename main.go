package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"code.google.com/p/rog-go/exp/go/ast"
	"code.google.com/p/rog-go/exp/go/token"
)

var (
	byteOffset = flag.Int("offset", -1, "the byte offset of the identifier in the file")
	filePath   = flag.String("file", "", "the file path containing the identifier")
	searchRoot = flag.String("root", "", "the root directory in which to search for references")
)

func main() {
	flag.Parse()

	if *filePath == "" || *byteOffset == -1 {
		flag.Usage()
		return
	}

	*filePath = normalizePath(*filePath, false)
	if *searchRoot == "" {
		*searchRoot = getRootPath(*filePath)
	} else {
		*searchRoot = normalizePath(*searchRoot, true)
	}

	var asts []*ast.File
	asts = append(asts, parseAST(*filePath))
	for path := range getFilesRecursive(*searchRoot) {
		if path != *filePath {
			asts = append(asts, parseAST(path))
		}
	}

	ident, declPos := findDecl(asts[0], *byteOffset)

	for _, ast := range asts {
		display(findReferences(ast, declPos, ident))
	}
}

func display(c chan token.Position) {
	for p := range c {
		f, err := os.Open(p.Filename)
		if err != nil {
			panic(err)
		}
		f.Seek(int64(p.Offset-p.Column+1), 0)
		line, err := bufio.NewReader(f).ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s:%d\n%s\n", strings.TrimPrefix(p.Filename, *searchRoot), p.Line, strings.TrimSpace(line))
	}
}
