# go-find-references
Find all References of an Identifier in a Codebase

### Usage

> go-find-references [flags]

###### flags:
`-file [string]`: the complete path of the file that contains the identifier.  
`-offset [int]`: the byte offset of the identifier in the above file.  
`-root [string]`: the directory in which to search for references (optional). If omitted, the directory containing `file` will be used.

For each found reference, the program will output two lines of text:  
The first line contains the path of the file relative to `root` followed by `:lineNumber`.
The second line contains the (trimmed) source text line containing the reference.

### example output:  
> go-find-references -file /usr/lib/go/src/pkg/sort/search.go -offset 2254

/usr/share/go/src/pkg/sort/search.go:59:6
func Search(n int, f func(int) bool) int {
/usr/share/go/src/pkg/sort/search.go:84:9
	return Search(len(a), func(i int) bool { return a[i] >= x })
/usr/share/go/src/pkg/sort/search.go:93:9
	return Search(len(a), func(i int) bool { return a[i] >= x })
/usr/share/go/src/pkg/sort/search.go:102:9
	return Search(len(a), func(i int) bool { return a[i] >= x })

### TODOs:

* Improve error messages  
* Add support for non-expression identifiers (e.g. keys in struct initialization)
