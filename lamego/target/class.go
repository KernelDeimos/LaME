package target

import (
	"github.com/KernelDeimos/LaME/lamego/model"
)

type Visibility byte

const (
	VisibilityPublic    = 'c'
	VisibilityPrivate   = 'e'
	VisibilityProtected = 'd'
)

type TypeOfType byte

const (
	PrimitiveType TypeOfType = 'p'
	LaMEType                 = 'l'
	TargetType               = 't'
)

type Type struct {
	TypeOfType TypeOfType
	Identifier string
}

type Variable struct {
	Name string
	Type Type
}

type ClassVariable struct {
	Variable
	Visibility Visibility
}

type Method struct {
	Name      string
	Arguments []Variable
	Return    Variable
	Code      model.FakeBlock
}

type ClassReference string

// Class is a universal definition of a class. Sure, some languages
// don't support object-oriented programming, but that doesn't matter.
// C programmers often accidentally implement a concept of classes in
// their code without even thinking about it.
type Class struct {
	Package string
	Name    string

	Meta ClassMeta

	// Extend should be for language-related purposes only. For instance,
	// in Java this is critical in defining the exception hierarchy.
	// Otherwise, real "extending" should happen in the code generator,
	// just like a compiler would have to do.
	Extend ClassReference

	// If a language doesn't support this, the compatability models better
	// find a way.
	Implements []ClassReference
	Variables  []Variable
	Methods    []Method
}
