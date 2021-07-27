package treesitter_test

import (
	"github.com/stretchr/testify/assert"
	"glide/lib/treesitter"
	"testing"
)

func TestJSONLanguage(t *testing.T) {
	assert.NotNil(t, treesitter.JSON)
}
