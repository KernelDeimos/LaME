package target

var PrimitiveBool = "bool"
var PrimitiveString = "string"
var PrimitiveInt = "int"
var PrimitiveFloat = "float"
var PrimitiveObject = "object"
var PrimitiveVoid = "void"

var Bool Type = Type{
	TypeOfType: PrimitiveType,
	Identifier: PrimitiveBool,
}
