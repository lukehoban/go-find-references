package main

import (
	"os"
	"strings"
)

func getFilesRecursive(root string) chan string {
	out := make(chan string)
	go func() {
		dir, err := os.Open(root)
		if err != nil {
			panic(err)
		}
		infos, err := dir.Readdir(0)
		if err != nil {
			panic(err)
		}
		for _, fi := range infos {
			if strings.HasPrefix(fi.Name(), ".") {
				continue
			}
			path := dir.Name() + fi.Name()
			if fi.IsDir() {
				for file := range getFilesRecursive(path + "/") {
					out <- file
				}
			} else if strings.HasSuffix(fi.Name(), ".go") {
				out <- path
			}
		}
		close(out)
	}()
	return out
}

func normalizePath(path string, isDir bool) string {
	if path == "" {
		return path
	}

	fields := strings.Split(path, "/")
	s := fields[0]
	for _, f := range fields {
		if f != "" {
			s += "/" + f
		}
	}
	if isDir {
		s += "/"
	}
	return s
}

func getRootPath(filepath string) string {
	return filepath[:strings.LastIndex(filepath, "/")+1]
}
