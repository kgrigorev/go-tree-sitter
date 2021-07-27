package treesitter

// #cgo CFLAGS: -Itree-sitter-json/src
// #include "tree-sitter-json/src/parser.c"
import "C"

var (
	JSON *Language
)

func init() {
	language, err := C.tree_sitter_json()
	if err != nil {
		panic(err)
	}

	JSON = &Language{ts: language}
}
