package treesitter

// #cgo CFLAGS: -Ivendor/tree-sitter-go/src
// #include "vendor/tree-sitter-go/src/parser.c"
import "C"

var (
	Go *Language
)

func init() {
	language, err := C.tree_sitter_go()
	if err != nil {
		panic(err)
	}

	Go = &Language{ts: language}
}
