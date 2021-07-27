package treesitter_test

import (
	treesitter "github.com/kgrigorev/go-tree-sitter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewParser(t *testing.T) {
	parser, err := treesitter.NewParser()
	assert.NoError(t, err)
	assert.NotNil(t, parser)
	err = parser.Delete()
	assert.NoError(t, err)
}

func TestParser_ParseString(t *testing.T) {
	parser := newJsonParser(t)
	defer parser.Delete()

	tree, err := parser.ParseString(nil, "[1, null]")
	assert.NoError(t, err)
	assert.NotNil(t, tree)

	root, err := tree.RootNode()
	assert.NoError(t, err)
	assert.NotNil(t, root)
	s, _ := root.String()
	assert.NotEqual(t, "", s)
}

func newJsonParser(t *testing.T) *treesitter.Parser {
	t.Helper()
	parser, err := treesitter.NewParser()
	require.NoError(t, err)

	ok, err := parser.SetLanguage(treesitter.JSON)
	require.NoError(t, err)
	require.True(t, ok)
	return parser
}

func parseJson(t *testing.T, input string) treesitter.Node {
	t.Helper()
	parser := newJsonParser(t)

	tree, err := parser.ParseString(nil, input)
	require.NoError(t, err)
	require.NotNil(t, tree)

	root, err := tree.RootNode()
	require.NoError(t, err)
	require.NotNil(t, root)
	return root
}

func newGoParser(t *testing.T) *treesitter.Parser {
	t.Helper()
	parser, err := treesitter.NewParser()
	require.NoError(t, err)

	ok, err := parser.SetLanguage(treesitter.Go)
	require.NoError(t, err)
	require.True(t, ok)
	return parser
}

func parseGo(t *testing.T, input string) treesitter.Node {
	t.Helper()
	parser := newGoParser(t)

	tree, err := parser.ParseString(nil, input)
	require.NoError(t, err)
	require.NotNil(t, tree)

	root, err := tree.RootNode()
	require.NoError(t, err)
	require.NotNil(t, root)
	return root
}
