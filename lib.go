package treesitter

// #cgo CFLAGS: -Ivendor/tree-sitter/lib/include -Ivendor/tree-sitter/lib/src
// #include <lib.c>
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

type Language struct {
	ts *C.TSLanguage
}

type Parser struct {
	ts *C.TSParser
}

func NewParser() (*Parser, error) {
	ts, err := C.ts_parser_new()
	if err != nil {
		return nil, err
	}
	parser := &Parser{ts: ts}
	runtime.SetFinalizer(parser, finalizeParser)
	return parser, nil
}

func finalizeParser(parser *Parser) error {
	if parser == nil {
		return nil
	}
	_, err := C.ts_parser_delete(parser.ts)
	return err
}

func (parser *Parser) Delete() error {
	return finalizeParser(parser)
}

func (parser *Parser) SetLanguage(language *Language) (bool, error) {
	ok := C.ts_parser_set_language(parser.ts, language.ts)
	return bool(ok), nil
}

func (parser *Parser) ParseString(oldTree *Tree, input string) (*Tree, error) {
	cstr := C.CString(input)
	defer C.free(unsafe.Pointer(cstr))

	var ot *C.struct_TSTree
	if oldTree != nil {
		ot = oldTree.ts
	}

	ts, err := C.ts_parser_parse_string(parser.ts, ot, cstr, C.uint(len(input)))
	if err != nil {
		return nil, err
	}
	return wrapTree(ts), nil
}

type Tree struct {
	ts *C.TSTree
}

func wrapTree(ts *C.TSTree) *Tree {
	tree := &Tree{ts: ts}
	runtime.SetFinalizer(tree, finalizeTree)
	return tree
}

func finalizeTree(tree *Tree) {
	_, err := C.ts_tree_delete(tree.ts)
	if err != nil {
		panic(err)
	}
}

func (tree *Tree) Delete() error {
	_, err := C.ts_tree_delete(tree.ts)
	return err
}

func (tree *Tree) RootNode() (Node, error) {
	var (
		root C.TSNode
		err  error
	)
	root, err = C.ts_tree_root_node(tree.ts)
	if err != nil {
		return Node{}, err
	}
	return Node{ts: root}, nil
}

func (tree *Tree) Edit(inputEdit *InputEdit) error {
	_, err := C.ts_tree_edit(tree.ts, inputEdit.toC())
	return err
}

type InputEdit struct {
	startByte   uint32
	oldEndByte  uint32
	newEndByte  uint32
	startPoint  Point
	oldEndPoint Point
	newEndPoint Point
}

func (inputEdit *InputEdit) toC() *C.TSInputEdit {
	return nil
}

type Node struct {
	ts C.TSNode
}

func (node Node) String() (string, error) {
	cstr, err := C.ts_node_string(node.ts)
	if err != nil {
		return "", err
	}
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString(cstr), nil
}

func (node Node) Type() (string, error) {
	cstr, err := C.ts_node_type(node.ts)
	if err != nil {
		return "", err
	}
	return C.GoString(cstr), nil
}

func (node Node) Parent() (Node, error) {
	ts, err := C.ts_node_parent(node.ts)
	if err != nil {
		return Node{}, err
	}
	return Node{ts: ts}, nil
}

func (node Node) Child(index uint32) (Node, error) {
	ts, err := C.ts_node_child(node.ts, C.uint(index))
	if err != nil {
		return Node{}, err
	}
	return Node{ts: ts}, nil
}

func (node Node) ChildCount() (uint32, error) {
	ts, err := C.ts_node_child_count(node.ts)
	if err != nil {
		return 0, err
	}
	return uint32(ts), nil
}

func (node Node) NamedChild(index uint32) (Node, error) {
	ts, err := C.ts_node_named_child(node.ts, C.uint(index))
	if err != nil {
		return Node{}, err
	}
	return Node{ts: ts}, nil
}

func (node Node) NamedChildCount() (uint32, error) {
	ts, err := C.ts_node_named_child_count(node.ts)
	if err != nil {
		return 0, err
	}
	return uint32(ts), nil
}

func (node Node) NextSibling() (Node, error) {
	ts, err := C.ts_node_next_sibling(node.ts)
	if err != nil {
		return Node{}, err
	}

	return Node{ts: ts}, nil
}

func (node Node) IsNull() (bool, error) {
	ts, err := C.ts_node_is_null(node.ts)
	if err != nil {
		return false, err
	}
	return bool(ts), nil
}

func (node Node) StartByte() (uint32, error) {
	ts, err := C.ts_node_start_byte(node.ts)
	if err != nil {
		return 0, err
	}
	return uint32(ts), nil
}

func (node Node) EndByte() (uint32, error) {
	ts, err := C.ts_node_end_byte(node.ts)
	if err != nil {
		return 0, err
	}
	return uint32(ts), nil
}

func (node Node) StartPoint() (Point, error) {
	ts, err := C.ts_node_start_point(node.ts)
	if err != nil {
		return Point{}, err
	}
	return Point{
		row:    uint32(ts.row),
		column: uint32(ts.column),
	}, nil
}

func (node Node) EndPoint() (Point, error) {
	ts, err := C.ts_node_end_point(node.ts)
	if err != nil {
		return Point{}, err
	}
	return Point{
		row:    uint32(ts.row),
		column: uint32(ts.column),
	}, nil
}

type Point struct {
	row, column uint32
}

type Query struct {
	ts        *C.TSQuery
	input     unsafe.Pointer
	errOffset uint32
	errType   QueryError
}

func NewQuery(language *Language, source string) (*Query, error) {
	input := C.CString(source)
	q := &Query{input: unsafe.Pointer(input)}
	runtime.SetFinalizer(q, finalizeQuery)
	var (
		errOffset C.uint32_t     = 0
		errType   C.TSQueryError = 0
	)
	ts, err := C.ts_query_new(language.ts, input, C.uint(len(source)), &errOffset, &errType)
	if err != nil {
		return nil, err
	}
	if errOffset != 0 {
		q.errOffset = uint32(errOffset)
	}
	if errType != 0 {
		q.errType = QueryError(errType)
	}
	q.ts = ts
	return q, nil
}

func (query *Query) Error() error {
	if query == nil {
		return nil
	}
	return fmt.Errorf("query error at %d of type %s", query.errOffset, query.errType)
}

func finalizeQuery(query *Query) error {
	if query == nil {
		return nil
	}
	if query.input != nil {
		_, err := C.free(query.input)
		if err != nil {
			return err
		}
	}
	if query.ts != nil {
		_, err := C.ts_query_delete(query.ts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (query *Query) Delete() error {
	return finalizeQuery(query)
}

type QueryCursor struct {
	ts *C.TSQueryCursor
}

func NewQueryCursor() (*QueryCursor, error) {
	ts, err := C.ts_query_cursor_new()
	if err != nil {
		return nil, err
	}
	qc := &QueryCursor{ts: ts}
	return qc, nil
}

func finalizeQueryCursor(queryCursor *QueryCursor) error {
	if queryCursor == nil {
		return nil
	}
	_, err := C.ts_query_cursor_delete(queryCursor.ts)
	return err
}

func (queryCursor *QueryCursor) Delete() error {
	return finalizeQueryCursor(queryCursor)
}

func (queryCursor *QueryCursor) Exec(query *Query, node Node) error {
	if queryCursor == nil || query == nil {
		return nil
	}
	_, err := C.ts_query_cursor_exec(queryCursor.ts, query.ts, node.ts)
	return err
}

func (queryCursor *QueryCursor) NextMatch() (QueryMatch, bool, error) {
	if queryCursor == nil {
		return QueryMatch{}, false, nil
	}
	var tsQueryMatch C.TSQueryMatch
	ok, err := C.ts_query_cursor_next_match(queryCursor.ts, &tsQueryMatch)
	if err != nil {
		return QueryMatch{}, bool(ok), err
	}
	qm := QueryMatch{
		Id:           uint32(tsQueryMatch.id),
		PatternIndex: uint16(tsQueryMatch.pattern_index),
		CaptureCount: uint16(tsQueryMatch.capture_count),
	}
	var (
		tsQueryCapture C.TSQueryCapture
	)
	start := unsafe.Pointer(tsQueryMatch.captures)
	size := unsafe.Sizeof(tsQueryCapture)
	for i := uint16(0); start != nil && i < qm.CaptureCount; i++ {
		value := *(*C.TSQueryCapture)(unsafe.Pointer(uintptr(start) + uintptr(i)*size))
		qm.Captures = append(qm.Captures, QueryCapture{
			Node:  Node{ts: value.node},
			Index: uint32(value.index),
		})
	}
	return qm, bool(ok), nil
}

type QueryCapture struct {
	Node  Node
	Index uint32
}

type QueryMatch struct {
	Id           uint32
	PatternIndex uint16
	CaptureCount uint16
	Captures     []QueryCapture
}

type QueryError int

func (qe QueryError) String() string {
	switch qe {
	case QueryErrorNone:
		return "NODE"
	case QueryErrorSyntax:
		return "SYNTAX"
	case QueryErrorNodeType:
		return "NODE_TYPE"
	case QueryErrorField:
		return "FIELD"
	case QueryErrorCapture:
		return "CAPTURE"
	case QueryErrorStructure:
		return "STRUCTURE"
	}
	return ""
}

const (
	QueryErrorNone QueryError = iota
	QueryErrorSyntax
	QueryErrorNodeType
	QueryErrorField
	QueryErrorCapture
	QueryErrorStructure
)
