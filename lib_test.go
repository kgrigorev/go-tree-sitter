package treesitter_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	treesitter "glide/lib/go-tree-sitter"
	"io"
	"os"
	"runtime"
	"testing"
)

func TestNode_NamedChild(t *testing.T) {
	var (
		nodeType   string
		childCount uint32
		err        error
	)
	root := parseJson(t, "[1, null]")
	nodeType, err = root.Type()
	assert.NoError(t, err)
	assert.Equal(t, "document", nodeType)
	childCount, err = root.ChildCount()
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), childCount)

	arrayNode, err := root.NamedChild(0)
	assert.NoError(t, err)
	nodeType, err = arrayNode.Type()
	assert.NoError(t, err)
	assert.Equal(t, "array", nodeType)
	childCount, err = arrayNode.ChildCount()
	assert.NoError(t, err)
	assert.Equal(t, uint32(5), childCount)

	numberNode, err := arrayNode.NamedChild(0)
	assert.NoError(t, err)
	nodeType, err = numberNode.Type()
	assert.NoError(t, err)
	assert.Equal(t, "number", nodeType)
	childCount, err = numberNode.ChildCount()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), childCount)
}

func TestTraverse(t *testing.T) {
	_, f, _, _ := runtime.Caller(0)
	r, _ := os.Open(f)
	data, _ := io.ReadAll(r)

	root := parseGo(t, string(data))
	traverse(root)
}

func traverse(node treesitter.Node) {
	if ok, err := node.IsNull(); ok || err != nil {
		return
	}

	fmt.Println(node.Type())
	fmt.Println(node.StartPoint())
	fmt.Println(node.EndPoint())

	count, err := node.ChildCount()
	if err != nil {
		return
	}
	for i := uint32(0); i < count; i++ {
		child, err := node.Child(i)
		if err != nil {
			return
		}
		traverse(child)
	}
}
