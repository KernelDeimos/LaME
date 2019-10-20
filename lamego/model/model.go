package model

var ModelRegistry map[string]Model

type LangSpec string

const (
	LanguageGo  = "go"
	LanguageES6 = "javascript"
)

type Argument struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type Field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type Method struct {
	Name     string            `yaml:"name"`
	Args     []Argument        `yaml:"args"`
	Gencode  string            `yaml:"gencode"`
	Hardcode map[string]string `yaml:"hardcode"`
}

type Model struct {
	ID      string   `yaml:"id"`
	Type    string   `yaml:"type"`
	Fields  []Field  `yaml:"fields"`
	Methods []Method `yaml:"methods"`
}

func (f Field) GetTypeObject() Type {
	typePrimitive, isPrimitive := primitives[f.Type]
	if isPrimitive {
		return Type{
			Primitive: typePrimitive,
		}
	}
	return Type{
		Primitive:  PrimitiveLaME,
		Identifier: f.Type,
	}
}

func (m Model) generateHardcode(language LangSpec) {
}
