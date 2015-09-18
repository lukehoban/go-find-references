package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/lukehoban/ident"
)

var (
	byteOffset = flag.Int("offset", -1, "the byte offset of the identifier in the file")
	filePath   = flag.String("file", "", "the file path containing the identifier")
	searchRoot = flag.String("root", "", "the root directory in which to search for references")
	showIdent  = flag.Bool("ident", false, "whether to show the name and position where the identifier is defined")
)

func main() {
	flag.Parse()

	if *filePath == "" || *byteOffset == -1 {
		flag.Usage()
		return
	}

	if *searchRoot == "" {
		*searchRoot = path.Dir(*filePath)
	}

	def, err := ident.Lookup(*filePath, *byteOffset)
	if err != nil {
		reportError(err)
		return
	}

	if *showIdent {
		fmt.Println(def.Name, def.Position)
	}

	refs, errs := def.FindReferences(*searchRoot, true)
	for {
		select {
		case ref, ok := <-refs:
			if !ok {
				return
			}
			reportReference(ref)
		case err, ok := <-errs:
			if !ok {
				return
			}
			reportError(err)
		}
	}
}

func reportError(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
}

func reportReference(ref ident.Reference) {
	f, err := os.Open(ref.Position.Filename)
	if err != nil {
		reportError(err)
	}

	_, err = f.Seek(int64(ref.Position.Offset-ref.Position.Column+1), 0)
	if err != nil {
		reportError(err)
	}

	line, err := bufio.NewReader(f).ReadString('\n')
	if err != nil {
		reportError(err)
	}

	fmt.Println(ref.Position)
	fmt.Println(line[:len(line)-1])
}
