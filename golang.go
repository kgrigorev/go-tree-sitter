package treesitter

// #cgo CFLAGS: -Itree-sitter-go/src
// #include "tree-sitter-go/src/parser.c"
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
