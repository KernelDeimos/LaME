package target

import (
	"strings"
)

// THIS FILE: Core Programming Interface for Code Generation

type CodeCursor interface {
	AddLine(line string)
	IncrementIndent()
	DecrementIndent()
}

type StringCodeCursor struct {
	code        string
	indent      int
	indentToken string
}

func (cc StringCodeCursor) writeIndent() {
	for i := 0; i < cc.indent; i++ {
		cc.code += cc.indentToken
	}
}

func (cc StringCodeCursor) AddLine(line string) {
	cc.writeIndent()
	cc.code += strings.TrimSpace(line) + "\n"
}

type ClassGenerator interface {
	WriteClass(cls Class, cc CodeCursor)
}
