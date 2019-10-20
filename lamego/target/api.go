package target

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// THIS FILE: Core Programming Interface for Code Generation

type FileManager interface {
	RequestFileForCode(relpath string) CodeCursor
	FlushAll()
}

type CursorConfig struct {
	IndentToken string
}

type DeFactoFileManager struct {
	RootPath     string
	CursorConfig CursorConfig
	files        map[string]CodeCursor
}

func (fm DeFactoFileManager) RequestFileForCode(relpath string) CodeCursor {
	cc, exists := fm.files[relpath]
	if exists {
		return cc
	}
	cc = NewStringCodeCursor(fm.CursorConfig)
	fm.files[relpath] = cc
	return cc
}

func (fm DeFactoFileManager) FlushAll() {
	for path, cursor := range fm.files {
		realPath := filepath.Join(fm.RootPath, path)
		err := ioutil.WriteFile(realPath, []byte(cursor.GetString()), 0644)
		// TODO: Need a LaME error type to return
		panic(err)
	}
}

type CodeCursor interface {
	AddLine(line string)
	AddString(str string)
	StartLine()
	EndLine()
	IncrIndent()
	DecrIndent()
	GetString() string
}

type StringCodeCursor struct {
	code        string
	lineStarted bool

	// TODO: It may be desirable to add a subIndent,
	//       for example to add a base indentation of
	//       4 spaces, and a sub-indentation of 2 spaces
	//       for things like long boolean expressions.
	indent      int
	indentToken string
}

func NewStringCodeCursor(conf CursorConfig) *StringCodeCursor {
	return &StringCodeCursor{
		indentToken: conf.IndentToken,
	}
}

func (cc *StringCodeCursor) writeIndent() {
	for i := 0; i < cc.indent; i++ {
		cc.code += cc.indentToken
	}
}

func (cc *StringCodeCursor) AddLine(line string) {
	cc.StartLine()
	defer cc.EndLine()
	cc.code += strings.TrimSpace(line)
}

func (cc *StringCodeCursor) StartLine() {
	if cc.lineStarted {
		panic("Invalid use of CodeCursor - must finish the previous line first!")
	}
	cc.lineStarted = true
	cc.writeIndent()
}

func (cc *StringCodeCursor) EndLine() {
	if !cc.lineStarted {
		panic("Invalid use of CodeCursor - must start a line first!")
	}
	cc.lineStarted = false
	cc.code += "\n"
}

func (cc *StringCodeCursor) AddString(str string) {
	if !cc.lineStarted {
		panic("Invalid use of CodeCursor - must start a line first!")
	}
	cc.code += str
}

func (cc *StringCodeCursor) IncrIndent() {
	cc.indent++
}

func (cc *StringCodeCursor) DecrIndent() {
	cc.indent--
}

func (cc *StringCodeCursor) GetString() string {
	return cc.code
}

type ClassGenerator interface {
	WriteClass(cls Class, fm FileManager)
}
