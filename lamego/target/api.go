package target

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// THIS FILE: Core Programming Interface for Code Generation

type FileManager interface {
	RequestFileForCode(relpath string) (CodeCursor, bool)
	FlushAll()
}

type CursorConfig struct {
	// TODO: It may be desirable to add a subIndent,
	//       for example to add a base indentation of
	//       4 spaces, and a sub-indentation of 2 spaces
	//       for things like long boolean expressions.
	IndentToken string
}

type DeFactoFileManager struct {
	RootPath     string
	CursorConfig CursorConfig
	files        map[string]CodeCursor
}

func NewDeFactoFileManager(outputDir string, cc CursorConfig) DeFactoFileManager {
	return DeFactoFileManager{
		RootPath:     outputDir,
		CursorConfig: cc,
		files:        map[string]CodeCursor{},
	}
}

func (fm DeFactoFileManager) RequestFileForCode(relpath string) (CodeCursor, bool) {
	cc, exists := fm.files[relpath]
	if exists {
		return cc, false
	}
	cc = NewStringCodeCursor(fm.CursorConfig)
	fm.files[relpath] = cc
	return cc, true
}

func (fm DeFactoFileManager) FlushAll() {
	for path, cursor := range fm.files {
		realPath := filepath.Join(fm.RootPath, path)
		err := os.MkdirAll(filepath.Dir(realPath), os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(realPath, []byte(cursor.GetString()), 0644)
		// TODO: Need a LaME error type to return
		if err != nil {
			panic(err)
		}
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
	NewSubCursor(name string) CodeCursor
	GetSubCursor(name string) CodeCursor
}

type InnerCursorGettable interface {
	GetString() string
}

type GettableString struct { S *string }
func (s GettableString) GetString() string {
	return *(s.S)
}

type StringCodeCursor struct {
	subCursors map[string]CodeCursor
	gettables []InnerCursorGettable
	code        *string
	lineStarted bool

	config CursorConfig
	indent      int
}

func NewStringCodeCursor(conf CursorConfig) *StringCodeCursor {
	var initialString string
	return &StringCodeCursor{
		config: conf,
		subCursors: map[string]CodeCursor{},
		gettables: []InnerCursorGettable{
			GettableString{S: &initialString},
		},
		code: &initialString,
	}
}

func (cc *StringCodeCursor) writeIndent() {
	for i := 0; i < cc.indent; i++ {
		*cc.code += cc.config.IndentToken
	}
}

func (cc *StringCodeCursor) AddLine(line string) {
	cc.StartLine()
	defer cc.EndLine()
	*cc.code += strings.TrimSpace(line)
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
	*cc.code += "\n"
}

func (cc *StringCodeCursor) AddString(str string) {
	if !cc.lineStarted {
		panic("Invalid use of CodeCursor - must start a line first!")
	}
	*cc.code += str
}

func (cc *StringCodeCursor) IncrIndent() {
	cc.indent++
}

func (cc *StringCodeCursor) DecrIndent() {
	cc.indent--
}

func (cc *StringCodeCursor) GetString() string {
	fullString := ""
	for _, gettable := range cc.gettables {
		fullString += gettable.GetString()
	}
	return fullString
}

func (cc *StringCodeCursor) NewSubCursor(name string) CodeCursor {
	// TODO: locks?
	scc := NewStringCodeCursor(cc.config)
	// Register subcursor to gettables and subCursors map
	cc.subCursors[name] = scc
	cc.gettables = append(cc.gettables, scc)

	// Need to add a new GettableString for data after this cursor
	var newWriteString string
	cc.gettables = append(cc.gettables, GettableString{S: &newWriteString})
	cc.code = &newWriteString

	return scc
}

func (cc *StringCodeCursor) GetSubCursor(name string) CodeCursor {
	return cc.subCursors[name]
}

type ClassGenerator interface {
	WriteClass(cls Class, fm FileManager)
}
