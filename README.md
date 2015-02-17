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
> go-find-references -file /usr/local/go/src/sort/search.go -offset 2254

search.go:59  
func Search(n int, f func(int) bool) int {  
search.go:84  
return Search(len(a), func(i int) bool { return a[i] >= x })  
search.go:93  
return Search(len(a), func(i int) bool { return a[i] >= x })  
search.go:102  
return Search(len(a), func(i int) bool { return a[i] >= x })  
search_test.go:53  
i := Search(e.n, e.f)  
search_test.go:82  
i := Search(n, func(i int) bool { count++; return i >= x })  
search_test.go:155  
i := Search(size, func(i int) bool { return i >= targ })  
