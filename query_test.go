package treesitter_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	treesitter "glide/lib/go-tree-sitter"
	"testing"
)

func TestNewQuery(t *testing.T) {
	language := treesitter.JSON
	q, err := treesitter.NewQuery(language, "(")
	assert.NoError(t, err)
	assert.Error(t, q.Error())
	err = q.Delete()
	assert.NoError(t, err)
}

func TestNewQueryCursor(t *testing.T) {
	qc, err := treesitter.NewQueryCursor()
	assert.NoError(t, err)
	err = qc.Delete()
	assert.NoError(t, err)
}

func TestQueryCursor_Exec(t *testing.T) {
	var (
		query       *treesitter.Query
		queryCursor *treesitter.QueryCursor
		err         error
	)
	root := parseGo(t, `
package main
import "fmt"

func main() {
	fmt.Println("hello")
    run(3)
}

func run(a int) error {
	if a > 1 {
    	return fmt.Errorf("a greater than 1")
    }
    return nil
}
`)
	language := treesitter.Go
	//query, err = treesitter.NewQuery(language, "(comment) @c")
	query, err = treesitter.NewQuery(language, "(function_declaration (identifier) @func_id)")
	assert.NoError(t, err)
	assert.Error(t, query.Error())

	queryCursor, err = treesitter.NewQueryCursor()
	assert.NoError(t, err)

	err = queryCursor.Exec(query, root)
	for m, ok, err := queryCursor.NextMatch(); ok && err == nil; m, ok, err = queryCursor.NextMatch() {
		for _, capture := range m.Captures {
			node := capture.Node
			fmt.Println(node.String())
			var (
				point treesitter.Point
				b     uint32
				err   error
			)
			point, err = node.StartPoint()
			fmt.Println("start point", point, err)
			b, err = node.StartByte()
			fmt.Println("start byte", b, err)
			point, err = node.EndPoint()
			fmt.Println("end point", point, err)
			b, err = node.EndByte()
			fmt.Println("end byte", b, err)

			fmt.Println()
		}
	}
}
