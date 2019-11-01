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
	Name       string     `yaml:"name"`
	Type       string     `yaml:"type"`
	Visibility Visibility `yaml:"visibility"`
}

type Method struct {
	Name       string            `yaml:"name"`
	Args       []Argument        `yaml:"args"`
	Return     string            `yaml:"return"`
	Gencode    string            `yaml:"gencode"`
	Hardcode   map[string]string `yaml:"hardcode"`
	Visibility Visibility        `yaml:"visibility"`
}

type ModelMeta struct {
	GencodeSyntaxFrontend string `yaml:"gencodeSyntaxFrontend"`
}

type Model struct {
	ID      string    `yaml:"id"`
	Type    string    `yaml:"type"`
	Meta    ModelMeta `yaml:"meta"`
	Fields  []Field   `yaml:"fields"`
	Methods []Method  `yaml:"methods"`
}

func NewDefaultModel() Model {
	return Model{
		Meta: ModelMeta{
			GencodeSyntaxFrontend: "LisPI-Natural",
		},
	}
}

func (f Field) GetTypeObject() Type {
	return GetTypeObject(f.Type)
}

// TODO: this is silly; just do it properly in engine
func GetTypeObject(typ string) Type {
	typePrimitive, isPrimitive := primitives[typ]
	if isPrimitive {
		return Type{
			Primitive: typePrimitive,
		}
	}
	return Type{
		Primitive:  PrimitiveLaME,
		Identifier: typ,
	}
}

func (m Model) generateHardcode(language LangSpec) {
}
