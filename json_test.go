package treesitter_test

import (
	treesitter "github.com/kgrigorev/go-tree-sitter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONLanguage(t *testing.T) {
	assert.NotNil(t, treesitter.JSON)
}
