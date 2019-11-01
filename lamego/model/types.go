package model

type Primitive byte

const (
	PrimitiveString Primitive = 's'
	PrimitiveBool             = 'b'
	PrimitiveInt              = 'i'
	PrimitiveFloat            = 'f'
	PrimitiveObject           = 'o'
	PrimitiveVoid             = 'v'
	PrimitiveLaME             = 'l'
)

var primitives = map[string]Primitive{
	"string": PrimitiveString,
	"bool":   PrimitiveBool,
	"int":    PrimitiveInt,
	"float":  PrimitiveFloat,
	"object": PrimitiveObject, // TODO: support for this might be weird
	"void":   PrimitiveVoid,
	"lame":   PrimitiveLaME,
	"LaME":   PrimitiveLaME, // Also allow stylized version, of course
}

type Type struct {
	Primitive  Primitive
	Identifier string
}

var Bool = Type{Primitive: PrimitiveBool}

type Visibility string

const (
	VisibilityPublic  Visibility = ""
	VisibilityPrivate            = "private"
)
