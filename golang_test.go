package treesitter_test

import (
	"glide/lib/treesitter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoLanguage(t *testing.T) {
	assert.NotNil(t, treesitter.Go)
}
