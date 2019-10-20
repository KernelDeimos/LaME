package model

type Primitive byte

var primitives = map[string]Primitive{
	"string": 's',
	"bool":   'b',
	"int":    'i',
	"float":  'f',
	"object": 'o', // TODO: support for this might be weird
	"void":   'v',
	"lame":   'l',
	"LaME":   'l', // Also allow stylized version, of course
}

const (
	PrimitiveString = 's'
	PrimitiveBool   = 'b'
	PrimitiveInt    = 'i'
	PrimitiveFloat  = 'f'
	PrimitiveObject = 'o'
	PrimitiveVoid   = 'v'
	PrimitiveLaME   = 'l'
)

type Type struct {
	Primitive  Primitive
	Identifier string
}

var Bool = Type{Primitive: PrimitiveBool}
