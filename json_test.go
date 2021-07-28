package treesitter_test

import (
	"github.com/stretchr/testify/assert"
	treesitter "glide/lib/go-tree-sitter"
	"testing"
)

func TestJSONLanguage(t *testing.T) {
	assert.NotNil(t, treesitter.JSON)
}
