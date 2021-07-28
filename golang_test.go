package treesitter_test

import (
	treesitter "glide/lib/go-tree-sitter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoLanguage(t *testing.T) {
	assert.NotNil(t, treesitter.Go)
}
