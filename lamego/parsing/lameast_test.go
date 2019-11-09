package parsing

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/alecthomas/participle"
)

func TestIt(t *testing.T) {
	ast := &Block{}
	parser, err := participle.Build(ast)
	if err != nil {
		t.Error(err)
	}

	err = parser.ParseString(`
	= myVar "Hello, World",
	x this.print [ myVar ],
	if x test [] []
	`, ast)

	b, err := json.Marshal(ast)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(b))
}
